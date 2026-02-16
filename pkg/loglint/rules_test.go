package loglint

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzerWithConfig(t *testing.T) {
	testdata := analysistest.TestData()

	configPath = "testdata/config.json"

	analysistest.Run(t, testdata, Analyzer, "simple")
}

func TestLowerCaseRule(t *testing.T) {
	testdata := analysistest.TestData()
	configPath = "testdata/lowercase_config.json"
	analysistest.Run(t, testdata, Analyzer, "lowercase")
}

func TestEnglishRule(t *testing.T) {
	testdata := analysistest.TestData()
	configPath = "testdata/english_config.json"
	analysistest.Run(t, testdata, Analyzer, "english")
}

func TestSpecialCharsRule(t *testing.T) {
	testdata := analysistest.TestData()
	configPath = "testdata/specials_config.json"
	analysistest.Run(t, testdata, Analyzer, "specials")
}

func TestSensitiveRule(t *testing.T) {
	testdata := analysistest.TestData()
	configPath = "testdata/sensitive_config.json"
	analysistest.Run(t, testdata, Analyzer, "sensitive")
}
