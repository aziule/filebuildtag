package filebuildtag

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

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

// Analyzer used to run the linter.
var Analyzer = &analysis.Analyzer{
	Name:     "filebuildtag",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

type fileTagsFlag []string

func (f *fileTagsFlag) String() string {
	return ""
}

func (f *fileTagsFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

var fileTags fileTagsFlag

func init() {
	Analyzer.Flags.Var(&fileTags, "filetag", `assign a file name pattern to a tag using the form <file-name-pattern>:<tag>, for example "foo/*_integration_test.go:integration"`)
}

func run(pass *analysis.Pass) (interface{}, error) {
	if len(fileTags) == 0 {
		return nil, nil
	}

	expectedTags := map[string]string{}
	for i := range fileTags {
		parts := strings.Split(fileTags[i], ":")
		expectedTags[parts[0]] = parts[1]
	}
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}
	inspector.Preorder(nodeFilter, func(node ast.Node) {
		f := node.(*ast.File)
		fileName := getFileName(pass, f.Pos())
		tags := internal.CheckGoFile(pass, f)

		for pattern, tag := range expectedTags {
			ok, _ := filepath.Match(pattern, fileName)
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
	return nil, nil
}

func getFileName(pass *analysis.Pass, pos token.Pos) string {
	path := pass.Fset.Position(pos).Filename
	_, file := filepath.Split(path)
	return file
}
