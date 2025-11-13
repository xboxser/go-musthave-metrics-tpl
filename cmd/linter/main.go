package main

import (
	"metrics/linteranalyzer"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(linteranalyzer.Analyzer)
}
