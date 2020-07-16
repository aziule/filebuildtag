package analyzer

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "gofilebuildtags",
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

type FileTag struct {
	FilePattern string
	BuildTag    string
}

var expectedFileTags []FileTag

func run(pass *analysis.Pass) (interface{}, error) {
	if len(buildTags) == 0 {
		return nil, nil
	}

	if expectedFileTags == nil {
		expectedFileTags = make([]FileTag, 0, len(buildTags))
		for _, buildTag := range buildTags {
			parts := strings.Split(buildTag, ":")
			expectedFileTags = append(expectedFileTags, FileTag{
				FilePattern: parts[0],
				BuildTag:    parts[1],
			})
		}
	}

	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}
	inspector.Preorder(nodeFilter, func(node ast.Node) {
		f := node.(*ast.File)
		fileName := getFileName(pass, f.Pos())

		for i := range expectedFileTags {
			ok, _ := filepath.Match(expectedFileTags[i].FilePattern, fileName)
			if !ok {
				continue
			}

			if len(f.Comments) == 0 {
				pass.Reportf(f.Pos(), `missing expected build tag: "%s"`, expectedFileTags[i].BuildTag)
				return
			}

			hasBuildTagsAnnotation := false

			for _, commentGroup := range f.Comments {
				// Stop when we reach a comment after the pkg keyword
				if commentGroup.End() >= f.Package {
					break
				}

				hasBuildTagsAnnotation = true

				foundTag := false
				for _, comment := range commentGroup.List {
					if hasTag(expectedFileTags[i].BuildTag, comment) {
						foundTag = true
					}
				}
				if !foundTag {
					pass.Reportf(f.Pos(), `missing expected build tag: "%s"`, expectedFileTags[i].BuildTag)
				}
			}

			if !hasBuildTagsAnnotation {
				pass.Reportf(f.Comments[0].Pos(), "missing build tags annotation")
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

func hasTag(tag string, comment *ast.Comment) bool {
	tags := getTags(comment.Text)
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

func getTags(str string) []string {
	if !strings.HasPrefix(str, "// +build ") {
		return nil
	}
	str = strings.ReplaceAll(str, ",", " ")
	return strings.Split(str[10:], " ")
}
