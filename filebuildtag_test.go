package filebuildtag

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_Lint(t *testing.T) {
	testdata := analysistest.TestData()
	testCases := map[string]struct {
		pattern string
		cfg     Config
	}{
		"buildtag - std lib linter's original test file": {
			pattern: "buildtag",
			cfg:     Config{filetags: map[string]string{"*": "foo"}},
		},
		"filebuildtag - wildcard match": {
			pattern: "filebuildtag_wildcard",
			cfg: Config{filetags: map[string]string{
				"*tag1_suff.go": "tag1",
				"*tag2_suff.go": "tag2",
			}},
		},
		"filebuildtag - exact match": {
			pattern: "filebuildtag_exact",
			cfg: Config{filetags: map[string]string{
				"pref_tag1_suff.go": "tag1",
				"pref_tag2_suff.go": "tag2",
			}},
		},
		"filebuildtag - no tags": {
			pattern: "filebuildtag_exact",
			cfg:     Config{filetags: map[string]string{}},
		},
	}
	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			analyzer := NewAnalyzer(tt.cfg)
			analysistest.Run(t, testdata, analyzer, tt.pattern)
		})
	}
}
