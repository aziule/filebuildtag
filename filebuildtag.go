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

// Analyzer used to run the linter.
var Analyzer = &analysis.Analyzer{
	Name:     "filebuildtag",
	Doc:      "Check that files with specific naming have the expected build tags.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

type buildTagsFlag []string

func (f *buildTagsFlag) String() string {
	return "foo"
}

func (f *buildTagsFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

var buildTags buildTagsFlag

func init() {
	Analyzer.Flags.Var(&buildTags, "buildtags", "")
}

func run(pass *analysis.Pass) (interface{}, error) {
	if len(buildTags) == 0 {
		return nil, nil
	}

	expectedTags := map[string]string{}
	for _, buildTag := range buildTags {
		parts := strings.Split(buildTag, ":")
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
