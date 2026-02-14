package main

import (
	"github.com/til-we-die/loglint/pkg/loglint"
	"golang.org/x/tools/go/analysis"
)

// для совместимости с golangci-lint
func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{loglint.Analyzer}, nil
}

// main нужен только для сборки как плагина
func main() {
}
