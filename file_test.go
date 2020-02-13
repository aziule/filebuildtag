package gofilebuildtags

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_hasTagInName(t *testing.T) {
	testCases := map[string]struct {
		fileName string
		tagName  string
		expected bool
	}{
		"success": {
			fileName: "foo_integration.go",
			tagName:  "integration",
			expected: true,
		},
		"file named after the tag": {
			fileName: "integration.go",
			tagName:  "integration",
			expected: true,
		},
		"missing tag": {
			fileName: "foo.go",
			tagName:  "integration",
			expected: false,
		},
	}
	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			testTag, err := newTag(tC.tagName)
			require.NoError(t, err)

			f, err := newFile(tC.fileName)
			require.NoError(t, err)
			assert.Equal(t, tC.expected, f.hasTagInName(testTag))
		})
	}
}
