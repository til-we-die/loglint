package main

import (
	"os"

	"github.com/til-we-die/loglint/pkg/loglint"
	"golang.org/x/tools/go/analysis"
)

type GolangCILintConfig struct {
	ConfigPath string `json:"config"`
}

func New(conf any) ([]*analysis.Analyzer, error) {
	if conf != nil {
		if configMap, ok := conf.(map[string]interface{}); ok {
			if configPath, ok := configMap["config"].(string); ok {
				os.Args = append(os.Args, "-config="+configPath)
			}
		}
	}

	return []*analysis.Analyzer{loglint.Analyzer}, nil
}

// main нужен только для сборки как плагина
func main() {}
