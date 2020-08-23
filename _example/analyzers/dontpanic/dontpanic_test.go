package dontpanic_test

import (
	"testing"

	"github.com/gostaticanalysis/astquery/_example/analyzers/dontpanic"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, dontpanic.Analyzer, "a")
}

