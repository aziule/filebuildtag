package gofilebuildtags

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Check(t *testing.T) {
	testDir := "./testdata"
	files, err := ioutil.ReadDir(testDir)
	require.NoError(t, err)

	linter, err := NewLinter([]string{"integration", "component"})
	require.NoError(t, err)
	for _, file := range files {
		linter.Check(fmt.Sprintf("%s/%s", testDir, file.Name()))
	}
}
