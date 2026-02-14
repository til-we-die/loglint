package loglint

import (
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

	firstChar := []rune(log.message)[0]

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

	for i, r := range log.message {
		if isEmoji(r) || isForbiddenPunctuation(r) {
			pass.Reportf(log.pos, "log message should not contain special characters or emojis (found %q at position %d)", r, i+1)
			return
		}
	}

	trimmed := strings.TrimRightFunc(log.message, func(r rune) bool {
		return r == '!' || r == '?' || r == '.'
	})

	if len(trimmed) < len(log.message) && len(log.message)-len(trimmed) > 1 {
		pass.Reportf(log.pos, "log message should not contain multiple punctuation marks at the end")
	}
}

// Проверка на чувствительные данные
func checkSensitive(pass *analysis.Pass, log *logCall) {
	if log.message == "" {
		return
	}

	sensitiveWords := []string{
		"password", "pass", "pwd",
		"token", "secret", "key",
		"api_key", "apikey",
		"credit", "card", "ssn",
		"auth", "credentials",
	}

	messageLower := strings.ToLower(log.message)

	for _, word := range sensitiveWords {
		if strings.Contains(messageLower, word) {
			pass.Reportf(log.pos, "log message may contain sensitive data: %q", word)
			return
		}
	}

	// добавить проверку аргументов
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
	// Диапазоны эмодзи в Unicode
	return (r >= 0x1F600 && r <= 0x1F64F) || // Эмотиконы
		(r >= 0x1F300 && r <= 0x1F5FF) || // Символы и пиктограммы
		(r >= 0x1F680 && r <= 0x1F6FF) || // Транспорт и символы
		(r >= 0x2600 && r <= 0x26FF) || // Разные символы
		(r >= 0x2700 && r <= 0x27BF) || // Символы Dingbats
		(r >= 0xFE00 && r <= 0xFE0F) || // Варианты селекторов
		(r >= 0x1F900 && r <= 0x1F9FF) // Дополнительные символы
}

// добавить регулярные выражения
var (
	emailRegex      = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	creditCardRegex = regexp.MustCompile(`\b(?:\d[ -]*?){13,16}\b`)
)
