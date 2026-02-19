package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"logcheck/analyser"
)

func main() {
	singlechecker.Main(analyser.Analyzer)
}
