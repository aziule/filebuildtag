package filebuildtag

import (
	"errors"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_Lint(t *testing.T) {
	testdata := analysistest.TestData()
	testCases := map[string]struct {
		pattern string
		flags   string
	}{
		"match files with a wildcard": {
			pattern: "filebuildtag_wildcard",
			flags:   "*tag1_suff.go:tag1,*tag2_suff.go:tag2",
		},
		"match exact file names": {
			pattern: "filebuildtag_exact",
			flags:   "pref_tag1_suff.go:tag1,pref_tag2_suff.go:tag2",
		},
		"match exact file name without tags": {
			pattern: "filebuildtag_exact",
			flags:   "",
		},
		"the std lib linter's original test file must have the foo tag": {
			pattern: "buildtag",
			flags:   "*:foo",
		},
	}
	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			analyzer := Analyzer
			flags := newFlagSet(t, tt.flags)
			analyzer.Flags = flags
			analysistest.Run(t, testdata, analyzer, tt.pattern)
		})
	}
}

func Test_parseFlags(t *testing.T) {
	emptyFiletags := map[string]string{}
	testCases := map[string]struct {
		flags       flag.FlagSet
		expected    map[string]string
		expectedErr error
	}{
		"no flags": {
			flags:       flag.FlagSet{},
			expected:    emptyFiletags,
			expectedErr: nil,
		},
		"empty tag": {
			flags:       newFlagSet(t, "  "),
			expected:    emptyFiletags,
			expectedErr: nil,
		},
		"malformed flag": {
			flags:       newFlagSet(t, "foo"),
			expectedErr: errors.New(`malformed argument: "foo", must be of the form "pattern:tag"`),
		},
		"empty file pattern": {
			flags:       newFlagSet(t, ":foo"),
			expectedErr: errors.New(`malformed argument: ":foo", must be of the form "pattern:tag"`),
		},
		"empty build tag": {
			flags:       newFlagSet(t, "foo:"),
			expectedErr: errors.New(`malformed argument: "foo:", must be of the form "pattern:tag"`),
		},
		"single file pattern": {
			flags: newFlagSet(t, "*:foo"),
			expected: map[string]string{
				"*": "foo",
			},
		},
		"several file patterns": {
			flags: newFlagSet(t, "foo:bar,bar:baz"),
			expected: map[string]string{
				"foo": "bar",
				"bar": "baz",
			},
		},
	}
	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			found, err := parseFlags(tt.flags)
			require.Equal(t, tt.expected, found)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func newFlagSet(t *testing.T, args string) flag.FlagSet {
	fs := flags()
	err := fs.Set(FlagFiletagsName, args)
	require.NoError(t, err)
	return fs
}
