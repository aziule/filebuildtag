package main

import (
	"github.com/aziule/filebuildtag/pkg/filebuildtag"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(filebuildtag.Analyzer)
}
