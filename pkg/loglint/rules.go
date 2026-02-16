package loglint

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// Проверка на строчную букву
func checkLowerCase(pass *analysis.Pass, log *logCall) {
	if log.message == "" {
		return
	}

	runes := []rune(log.message)
	if len(runes) == 0 {
		return
	}

	firstChar := runes[0]

	if unicode.IsLetter(firstChar) && !unicode.IsLower(firstChar) {
		pass.Reportf(log.pos, "log message should start with a lowercase letter (found %q)", firstChar)
	}
}

// Проверка на английский язык
func checkEnglish(pass *analysis.Pass, log *logCall) {
	if log.message == "" {
		return
	}

	for _, r := range log.message {
		if unicode.IsSpace(r) || unicode.IsDigit(r) || isAllowedPunctuation(r) {
			continue
		}

		if unicode.IsLetter(r) && !isLatinLetter(r) {
			pass.Reportf(log.pos, "log message should be in English (non-Latin character detected: %q)", r)
			return
		}
	}
}

// Проверка на спецсимволы и эмодзи
func checkSpecialChars(pass *analysis.Pass, log *logCall) {
	if log.message == "" {
		return
	}

	runes := []rune(log.message)

	// Эмодзи и запрещённые символы
	for _, r := range runes {
		if isEmoji(r) || isForbiddenPunctuation(r) {
			pass.Reportf(log.pos, "special characters or emojis")
			return
		}
	}

	last := runes[len(runes)-1]

	if last == '!' || last == '.' {
		pass.Reportf(log.pos, "multiple punctuation marks")
		return
	}

	if last == '?' {
		if len(runes) > 1 {
			prev := runes[len(runes)-2]
			if prev == '?' || prev == '!' || prev == '.' {
				pass.Reportf(log.pos, "multiple punctuation marks")
				return
			}
		}
	}
}

// Проверка на чувствительные данные
func checkSensitive(pass *analysis.Pass, log *logCall) {
	if log.message == "" {
		return
	}

	msg := strings.ToLower(log.message)

	keywords := []string{
		"password",
		"pwd",
		"token",
		"secret",
		"api_key",
		"apikey",
		"ssn",
		"credentials",
	}

	for _, kw := range keywords {
		if strings.Contains(msg, kw) {
			pass.Reportf(log.pos, "sensitive data")
			return
		}
	}

	if strings.Contains(msg, "credit card") {
		pass.Reportf(log.pos, "sensitive data")
		return
	}

	words := strings.FieldsFunc(msg, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_'
	})

	for _, w := range words {
		if strings.HasPrefix(w, "password") ||
			strings.HasPrefix(w, "secret") ||
			strings.HasPrefix(w, "token") {
			pass.Reportf(log.pos, "sensitive data")
			return
		}
	}
}

func checkZapFields(pass *analysis.Pass, log *logCall) {
	if log.logger != "zap" {
		return
	}

	sensitiveKeys := []string{
		"password",
		"pwd",
		"pass",
		"token",
		"secret",
		"api_key",
		"apikey",
		"credentials",
	}

	for _, arg := range log.call.Args[1:] {
		call, ok := arg.(*ast.CallExpr)
		if !ok {
			continue
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		funcName := sel.Sel.Name
		if funcName == "" {
			continue
		}

		if len(call.Args) == 0 {
			continue
		}

		keyLit, ok := call.Args[0].(*ast.BasicLit)
		if !ok || keyLit.Kind != token.STRING {
			continue
		}

		key := strings.ToLower(keyLit.Value[1 : len(keyLit.Value)-1])

		for _, sensitive := range sensitiveKeys {
			if strings.Contains(key, sensitive) {
				pass.Reportf(log.pos, "sensitive data in zap field")
				return
			}
		}
	}
}

func isLatinLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isAllowedPunctuation(r rune) bool {
	allowed := []rune{'.', ',', ':', ';', '-', '_', '/', '\\', '@', '=', '+', '*', '&', '%', '$', '#', '!', '?'}
	for _, a := range allowed {
		if r == a {
			return true
		}
	}
	return false
}

func isForbiddenPunctuation(r rune) bool {
	forbidden := []rune{'‼', '⁉', '‽', '…'}
	for _, f := range forbidden {
		if r == f {
			return true
		}
	}
	return false
}

func isEmoji(r rune) bool {
	return (r >= 0x1F600 && r <= 0x1F64F) ||
		(r >= 0x1F300 && r <= 0x1F5FF) ||
		(r >= 0x1F680 && r <= 0x1F6FF) ||
		(r >= 0x2600 && r <= 0x26FF) ||
		(r >= 0x2700 && r <= 0x27BF) ||
		(r >= 0xFE00 && r <= 0xFE0F) ||
		(r >= 0x1F900 && r <= 0x1F9FF)
}

var (
	emailRegex      = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	creditCardRegex = regexp.MustCompile(`\b(?:\d[ -]*?){13,16}\b`)
)
