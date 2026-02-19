package analyser_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"logcheck/analyser"
)

func TestAll(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyser.Analyzer, "p")
}
