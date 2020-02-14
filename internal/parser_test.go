package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseFileName(t *testing.T) {
	parser := NewFileNameParser(TestFileSuffix, GoFileExt)

	testCases := map[string]struct {
		fileName    string
		tagName     string
		expected    bool
		expectedErr error
	}{
		"placeholder_tag_suffix.ext": {
			fileName: "foo_integration_test.go",
			tagName:  "integration",
			expected: true,
		},
		"tag_suffix.ext": {
			fileName: "integration_test.go",
			tagName:  "integration",
			expected: true,
		},
		"tag.ext": {
			fileName: "integration.go",
			tagName:  "integration",
			expected: false,
		},
		"placeholder_suffix.ext": {
			fileName: "foo_test.go",
			tagName:  "integration",
			expected: false,
		},
		"tag_placeholder_suffix.ext": {
			fileName: "integration_foo_test.go",
			tagName:  "integration",
			expected: false,
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
		"": {
			fileName: "./testdata/generic_with_build_tags.go",
			tagName:  "integration",
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
