package internal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseFileName(t *testing.T) {
	tag, err := NewTag("foo")
	require.NoError(t, err)
	tags := []*Tag{tag}
	parser, err := NewFileNameParser(TestFileSuffix, tags)
	require.NoError(t, err)

	testCases := map[string]struct {
		fileName    string
		tagName     string
		expected    bool
		expectedErr error
	}{
		"placeholder_tag_suffix.go": {
			fileName: "something_foo_test.go",
			tagName:  "foo",
			expected: true,
		},
		"tag_suffix.go": {
			fileName: "foo_test.go",
			tagName:  "foo",
			expected: true,
		},
		"tag.go": {
			fileName: "foo.go",
			tagName:  "foo",
			expected: false,
		},
		"placeholder_suffix.go": {
			fileName: "something_test.go",
			tagName:  "foo",
			expected: false,
		},
		"tag_placeholder_suffix.go": {
			fileName: "foo_something_test.go",
			tagName:  "foo",
			expected: false,
		},
		"_tag_suffix.go": {
			fileName: "_foo_test.go",
			tagName:  "foo",
			expected: true,
		},
		"placeholder_tag_suffix": {
			fileName: "something_foo_test",
			tagName:  "foo",
			expected: false,
		},
		"tag_suffix": {
			fileName: "foo_test",
			tagName:  "foo",
			expected: false,
		},
		"tag": {
			fileName: "foo",
			tagName:  "foo",
			expected: false,
		},
		"placeholder_suffix": {
			fileName: "something_test",
			tagName:  "foo",
			expected: false,
		},
		"tag_placeholder_suffix": {
			fileName: "foo_something_test",
			tagName:  "foo",
			expected: false,
		},
		"_tag_suffix": {
			fileName: "_foo_test",
			tagName:  "foo",
			expected: false,
		},
		"tag_suffix.ext": {
			fileName: "foo_test.ext",
			tagName:  "foo",
			expected: false,
		},
		"empty file name": {
			fileName:    "",
			tagName:     "",
			expected:    false,
			expectedErr: errors.New("empty file name"),
		},
		"tag mismatch": {
			fileName:    "something.go",
			tagName:     "bar",
			expected:    false,
			expectedErr: errors.New("tag mismatch"),
		},
	}
	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			tag, _ := NewTag(tC.tagName)

			got, err := parser.Parse(tC.fileName, tag)
			assert.Equal(t, tC.expected, got)

			if tC.expectedErr != nil {
				assert.EqualError(t, err, tC.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
