package analyser_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/trust-me-im-an-engineer/logcheck/analyser"
)

func TestAll(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyser.Analyzer, "p")
}
