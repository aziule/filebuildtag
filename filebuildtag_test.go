package filebuildtag

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_Lint(t *testing.T) {
	testdata := analysistest.TestData()
	testCases := map[string]struct {
		pattern   string
		buildTags []string
	}{
		"buildtag - std lib linter's original test file": {
			pattern:   "buildtag",
			buildTags: []string{"*:foo"},
		},
		"filebuildtag - wildcard match": {
			pattern: "filebuildtag_wildcard",
			buildTags: []string{
				"*tag1_suff.go:tag1",
				"*tag2_suff.go:tag2",
			},
		},
		"filebuildtag - exact match": {
			pattern: "filebuildtag_exact",
			buildTags: []string{
				"pref_tag1_suff.go:tag1",
				"pref_tag2_suff.go:tag2",
			},
		},
		"filebuildtag - no tags": {
			pattern:   "filebuildtag_exact",
			buildTags: []string{},
		},
	}
	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			fileTags = fileTagsFlag{}
			for _, buildTag := range tt.buildTags {
				err := fileTags.Set(buildTag)
				require.NoError(t, err)
			}
			analysistest.Run(t, testdata, Analyzer, tt.pattern)
		})
	}
}
