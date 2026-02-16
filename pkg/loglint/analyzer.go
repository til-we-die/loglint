package loglint

import (
	"encoding/json"
	"flag"
	"go/ast"
	"go/token"
	"os"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type Config struct {
	EnableRules             []string `json:"enable_rules"`
	CustomSensitiveKeywords []string `json:"custom_sensitive_keywords"`
	CustomSensitivePatterns []string `json:"custom_sensitive_patterns"`
	CustomZapSensitiveKeys  []string `json:"custom_zap_sensitive_keys"`

	LowerCase struct {
		Enabled      bool `json:"enabled"`
		AllowNumbers bool `json:"allow_numbers"`
		AllowSymbols bool `json:"allow_symbols"`
	} `json:"lowercase"`

	English struct {
		Enabled         bool     `json:"enabled"`
		AllowedNonLatin []string `json:"allowed_non_latin"`
	} `json:"english"`

	SpecialChars struct {
		Enabled      bool `json:"enabled"`
		AllowEmojis  bool `json:"allow_emojis"`
		AllowSpecial bool `json:"allow_special"`
	} `json:"special_chars"`

	Sensitive struct {
		Enabled       bool     `json:"enabled"`
		ExtraPatterns []string `json:"extra_patterns"`
		StrictMode    bool     `json:"strict_mode"`
	} `json:"sensitive"`
}

func defaultConfig() *Config {
	cfg := &Config{
		EnableRules: []string{
			"lowercase",
			"english",
			"special_chars",
			"sensitive",
		},
		CustomSensitiveKeywords: []string{},
		CustomSensitivePatterns: []string{},
		CustomZapSensitiveKeys:  []string{},
	}

	cfg.LowerCase.Enabled = true
	cfg.LowerCase.AllowNumbers = true
	cfg.LowerCase.AllowSymbols = true

	cfg.English.Enabled = true
	cfg.English.AllowedNonLatin = []string{}

	cfg.SpecialChars.Enabled = true
	cfg.SpecialChars.AllowEmojis = false
	cfg.SpecialChars.AllowSpecial = false

	cfg.Sensitive.Enabled = true
	cfg.Sensitive.ExtraPatterns = []string{}
	cfg.Sensitive.StrictMode = false

	return cfg
}

func LoadConfig(path string) (*Config, error) {
	cfg := defaultConfig()

	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

var Analyzer = &analysis.Analyzer{
	Name:     "loglint",
	Doc:      "Linter for log messages to ensure they follow best practices",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Flags:    flagSet(),
}

var (
	configPath string
	config     *Config
)

func flagSet() flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.StringVar(&configPath, "config", "", "path to configuration file")
	return *fs
}

type logCall struct {
	call    *ast.CallExpr
	message string
	pos     token.Pos
	logger  string // "slog" or "zap"
	level   string // "Info", "Error", etc
}

func run(pass *analysis.Pass) (interface{}, error) {
	var err error
	config, err = LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

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

		checkRules(pass, logInfo, config)
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
