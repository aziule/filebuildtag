package internal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewTag(t *testing.T) {
	testCases := map[string]struct {
		name        string
		expected    *Tag
		expectedErr error
	}{
		"success": {
			name: "foo",
			expected: &Tag{
				name: "foo",
				len:  3,
			},
		},
		"missing name": {
			name:        "",
			expectedErr: errors.New("empty tag"),
		},
	}
	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			tag, err := NewTag(tC.name)
			if tC.expectedErr != nil {
				assert.Nil(t, tag)
				assert.EqualError(t, err, tC.expectedErr.Error())
			} else {
				assert.Equal(t, tC.expected, tag)
				assert.NoError(t, err)
			}
		})
	}
}
