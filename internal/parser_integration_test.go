// +build integration

package internal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseFileContents(t *testing.T) {
	testCases := map[string]struct {
		fileName    string
		tagName     string
		expected    bool
		expectedErr error
	}{
		"file that does not exist": {
			fileName:    "foo.go",
			tagName:     "bar",
			expected:    false,
			expectedErr: errors.New("could not open file foo.go: open foo.go: no such file or directory"),
		},
		"tag matches": {
			fileName: "./testdata/single_tag.go",
			tagName:  "foo",
			expected: true,
		},
		"tag matches with blank lines": {
			fileName: "./testdata/single_tag_with_blank_lines.go",
			tagName:  "foo",
			expected: true,
		},
		"tag matches in first position": {
			fileName: "./testdata/several_tags.go",
			tagName:  "foo",
			expected: true,
		},
		"tag matches in second position": {
			fileName: "./testdata/several_tags.go",
			tagName:  "bar",
			expected: true,
		},
		"another, shorter tag exists": {
			fileName: "./testdata/other_tag_shorter.go",
			tagName:  "bar",
			expected: false,
		},
		"no tags": {
			fileName: "./testdata/no_tags.go",
			tagName:  "bar",
			expected: false,
		},
	}
	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			p := NewContentsParser()
			tag, err := NewTag(tC.tagName)
			require.NoError(t, err)

			found, err := p.Parse(tC.fileName, tag)
			assert.Equal(t, tC.expected, found)

			if tC.expectedErr != nil {
				assert.EqualError(t, err, tC.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
