package main

import (
	"github.com/aziule/filebuildtag/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
