package loglint

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLowerCaseRule(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "lowercase")
}

func TestEnglishRule(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "english")
}

func TestSpecialCharsRule(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "specials")
}

func TestSensitiveRule(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "sensitive")
}
