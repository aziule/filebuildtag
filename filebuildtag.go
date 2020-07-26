package filebuildtag

import (
	"go/ast"
	"path/filepath"

	"github.com/aziule/filebuildtag/internal"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Doc of the linter.
const Doc = `check that Go files have the expected build tags in the "// +build" instruction

Define file patterns and assign them to build tags, for instance:
	File "bar.go" must have the "baz" build tag
	Files "*_integration_test.go" must have the "integration" build tag`

var analyzer = &analysis.Analyzer{
	Name:     "filebuildtag",
	Doc:      Doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// NewAnalyzer creates an analysis.Analyzer with config params, ready to be used.
func NewAnalyzer(filetags map[string]string) *analysis.Analyzer {
	analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
		run(filetags, pass)
		return nil, nil
	}
	return analyzer
}

func run(filetags map[string]string, pass *analysis.Pass) {
	if len(filetags) == 0 || pass == nil {
		return
	}
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}
	inspector.Preorder(nodeFilter, func(node ast.Node) {
		f := node.(*ast.File)
		filename := getFilename(pass, f)
		tags := internal.CheckGoFile(pass, f)
		for pattern, tag := range filetags {
			ok, _ := filepath.Match(pattern, filename)
			if !ok {
				continue
			}

			foundTag := false
			for i := range tags {
				if tags[i] == tag {
					foundTag = true
					break
				}
			}

			if !foundTag {
				pass.Reportf(f.Pos(), `missing expected build tag: "%s"`, tag)
			}
		}
	})
}

func getFilename(pass *analysis.Pass, file *ast.File) string {
	path := pass.Fset.Position(file.Pos()).Filename
	_, filename := filepath.Split(path)
	return filename
}
