package main

import (
	"errors"
	"testing"

	"github.com/aziule/gofilebuildtags/internal"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func Test_NewLinter(t *testing.T) {
	tags := make(map[string]*internal.Tag)
	tags["foo"], _ = internal.NewTag("foo")
	tags["bar"], _ = internal.NewTag("bar")
	tags["baz"], _ = internal.NewTag("baz")

	fileNameParser, err := internal.NewFileNameParser([]*internal.Tag{tags["foo"], tags["bar"], tags["baz"]})
	require.NoError(t, err)

	testCases := map[string]struct {
		tags        []string
		expected    *Linter
		expectedErr error
	}{
		"empty tags list": {
			tags:        []string{},
			expectedErr: ErrEmptyTagsList,
		},
		"empty tag only": {
			tags:        []string{""},
			expectedErr: errors.New("empty tag"),
		},
		"empty tag in the list": {
			tags:        []string{"foo", ""},
			expectedErr: errors.New("empty tag"),
		},
		"several times the same tag": {
			tags: []string{"foo", "bar", "foo", "baz"},
			expected: &Linter{
				fileNameParser: fileNameParser,
				contentsParser: internal.NewContentsParser(),
				tags:           []*internal.Tag{tags["foo"], tags["bar"], tags["baz"]},
			},
		},
	}
	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			linter, err := NewLinter(tC.tags)
			if tC.expectedErr != nil {
				assert.Nil(t, linter)
				assert.EqualError(t, err, tC.expectedErr.Error())
			} else {
				assert.Equal(t, linter, tC.expected)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Check(t *testing.T) {
	fileName := "foo.go"
	anErr := errors.New("an error occurred")
	tags := make(map[string]*internal.Tag)
	tags["foo"], _ = internal.NewTag("foo")
	tags["bar"], _ = internal.NewTag("bar")

	testCases := map[string]struct {
		fileNameParser internal.TagParser
		contentsParser internal.TagParser
		tags           []*internal.Tag
		expected       []*Issue
		expectedErr    error
	}{
		"file name parser error": {
			fileNameParser: newStubTagParser(false, anErr),
			tags:           []*internal.Tag{tags["foo"]},
			expected: []*Issue{
				{
					FileName: fileName,
					Reason:   anErr.Error(),
				},
			},
		},
		"file contents parser error": {
			fileNameParser: newStubTagParser(false, nil),
			contentsParser: newStubTagParser(false, anErr),
			tags:           []*internal.Tag{tags["foo"]},
			expected: []*Issue{
				{
					FileName: fileName,
					Reason:   anErr.Error(),
				},
			},
		},
		"tag in file name only": {
			fileNameParser: newStubTagParser(true, nil),
			contentsParser: newStubTagParser(false, nil),
			tags:           []*internal.Tag{tags["foo"]},
			expected: []*Issue{
				{
					FileName: fileName,
					Reason:   "tag foo found in the file name but missing from the build tags",
				},
			},
		},
		"tag in file contents only": {
			fileNameParser: newStubTagParser(false, nil),
			contentsParser: newStubTagParser(true, nil),
			tags:           []*internal.Tag{tags["foo"]},
			expected: []*Issue{
				{
					FileName: fileName,
					Reason:   "tag foo found in the file build tags but missing from the file name",
				},
			},
		},
		"tag neither in file name or contents": {
			fileNameParser: newStubTagParser(false, nil),
			contentsParser: newStubTagParser(false, nil),
			tags:           []*internal.Tag{tags["foo"]},
		},
		"tag in both file name and contents": {
			fileNameParser: newStubTagParser(true, nil),
			contentsParser: newStubTagParser(true, nil),
			tags:           []*internal.Tag{tags["foo"]},
		},
		"error with several tags": {
			fileNameParser: newStubTagParser(true, nil),
			contentsParser: newStubTagParser(false, nil),
			tags:           []*internal.Tag{tags["foo"], tags["bar"]},
			expected: []*Issue{
				{
					FileName: fileName,
					Reason:   "tag foo found in the file name but missing from the build tags",
				},
				{
					FileName: fileName,
					Reason:   "tag bar found in the file name but missing from the build tags",
				},
			},
		},
	}
	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			linter := &Linter{
				fileNameParser: tC.fileNameParser,
				contentsParser: tC.contentsParser,
				tags:           tC.tags,
			}
			found, err := linter.Check(fileName)
			if tC.expectedErr != nil {
				assert.Nil(t, found)
				assert.EqualError(t, err, tC.expectedErr.Error())
			} else {
				assert.Equal(t, tC.expected, found)
				assert.NoError(t, err)
			}
		})
	}
}

type StubTagParser struct {
	parseResult bool
	parseErr    error
}

func newStubTagParser(parseResult bool, parseErr error) *StubTagParser {
	return &StubTagParser{
		parseResult: parseResult,
		parseErr:    parseErr,
	}
}

func (p *StubTagParser) Parse(s string, t *internal.Tag) (bool, error) {
	return p.parseResult, p.parseErr
}
