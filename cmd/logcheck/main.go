package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/trust-me-im-an-engineer/logcheck/analyser"
)

func main() {
	singlechecker.Main(analyser.Analyzer)
}
