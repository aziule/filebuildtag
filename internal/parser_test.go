package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseFileName(t *testing.T) {
	tag, err := NewTag("foo")
	require.NoError(t, err)
	tags := []*Tag{tag}
	parser, err := NewFileNameParser(TestFileSuffix, GoFileExt, tags)
	require.NoError(t, err)

	testCases := map[string]struct {
		fileName    string
		tagName     string
		expected    bool
		expectedErr error
	}{
		"placeholder_tag_suffix.ext": {
			fileName: "something_foo_test.go",
			tagName:  "foo",
			expected: true,
		},
		"tag_suffix.ext": {
			fileName: "foo_test.go",
			tagName:  "foo",
			expected: true,
		},
		"tag.ext": {
			fileName: "foo.go",
			tagName:  "foo",
			expected: false,
		},
		"placeholder_suffix.ext": {
			fileName: "something_test.go",
			tagName:  "foo",
			expected: false,
		},
		"tag_placeholder_suffix.ext": {
			fileName: "foo_something_test.go",
			tagName:  "foo",
			expected: false,
		},
		"_tag_suffix.ext": {
			fileName: "_foo_test.go",
			tagName:  "foo",
			expected: true,
		},
	}
	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			tag, err := NewTag(tC.tagName)
			require.NoError(t, err)

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

func Test_ParseFileContents(t *testing.T) {
	testCases := map[string]struct {
		fileName    string
		tagName     string
		expected    bool
		expectedErr error
	}{
		"tag matches in first position": {
			fileName: "./testdata/generic_with_build_tags.go",
			tagName:  "foo",
			expected: true,
		},
		"tag matches in second position": {
			fileName: "./testdata/generic_with_build_tags.go",
			tagName:  "bar",
			expected: true,
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
