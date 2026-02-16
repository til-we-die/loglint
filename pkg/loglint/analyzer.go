package loglint

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglint",
	Doc:      "Linter for log messages to ensure they follow best practices",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

type logCall struct {
	call    *ast.CallExpr
	message string
	pos     token.Pos
	logger  string // "slog" или "zap"
	level   string // "Info", "Error", etc
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		logger, level, ok := isLogCall(call)
		if !ok {
			return
		}

		message := extractLogMessage(call)
		if message == "" {
			return
		}

		logInfo := &logCall{
			call:    call,
			message: message,
			pos:     call.Pos(),
			logger:  logger,
			level:   level,
		}

		checkRules(pass, logInfo)
	})

	return nil, nil
}

func isLogCall(call *ast.CallExpr) (string, string, bool) {
	fun, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", "", false
	}

	methodName := fun.Sel.Name

	switch methodName {
	case "Debug", "Info", "Warn", "Error":
		if ident, ok := fun.X.(*ast.Ident); ok {
			if ident.Name == "slog" {
				return "slog", methodName, true
			}
		}

		return "zap", methodName, true
	}

	return "", "", false
}

func extractLogMessage(call *ast.CallExpr) string {
	if len(call.Args) == 0 {
		return ""
	}

	return extractStringFromExpr(call.Args[0])
}

func extractStringFromExpr(expr ast.Expr) string {
	switch v := expr.(type) {

	case *ast.BasicLit:
		if v.Kind == token.STRING && len(v.Value) >= 2 {
			return v.Value[1 : len(v.Value)-1]
		}

	case *ast.Ident:
		if obj := v.Obj; obj != nil && obj.Kind == ast.Con {
			if vs, ok := obj.Decl.(*ast.ValueSpec); ok && len(vs.Values) > 0 {
				if lit, ok := vs.Values[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
					return lit.Value[1 : len(lit.Value)-1]
				}
			}
		}

	case *ast.CallExpr:
		if sel, ok := v.Fun.(*ast.SelectorExpr); ok {
			if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "fmt" && sel.Sel.Name == "Sprintf" {
				if len(v.Args) > 0 {
					if lit, ok := v.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
						return lit.Value[1 : len(lit.Value)-1]
					}
				}
			}
		}
	}
	return ""
}

func checkRules(pass *analysis.Pass, log *logCall) {
	checkLowerCase(pass, log)
	checkEnglish(pass, log)
	checkSpecialChars(pass, log)
	checkSensitive(pass, log)
	checkZapFields(pass, log)
}
