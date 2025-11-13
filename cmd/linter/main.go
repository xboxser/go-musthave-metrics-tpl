package main

import (
	"metrics/linter_analyzer"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(linter_analyzer.Analyzer)
}
