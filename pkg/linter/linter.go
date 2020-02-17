package main

import (
	"errors"
	"fmt"

	"github.com/aziule/gofilebuildtags/internal"
)

// ErrEmptyTagsList is an error returned when no tags are provided.
var ErrEmptyTagsList = errors.New("empty tags list")

// Linter is the main linter object, used to parse file names or file contents.
type Linter struct {
	fileNameParser internal.TagParser
	contentsParser internal.TagParser
	tags           []*internal.Tag
}

// NewLinter creates a new Linter.
func NewLinter(tags []string) (*Linter, error) {
	nbTags := len(tags)

	if nbTags == 0 {
		return nil, ErrEmptyTagsList
	}

	seenTags := make(map[string]bool)
	var tagObjs []*internal.Tag

	for _, tag := range tags {
		if _, ok := seenTags[tag]; ok {
			continue
		}

		seenTags[tag] = true

		t, err := internal.NewTag(tag)
		if err != nil {
			return nil, err
		}

		tagObjs = append(tagObjs, t)
	}

	if len(tagObjs) == 0 {
		return nil, ErrEmptyTagsList
	}

	fileNameParser, err := internal.NewFileNameParser(tagObjs)
	if err != nil {
		return nil, err
	}

	return &Linter{
		fileNameParser: fileNameParser,
		contentsParser: internal.NewContentsParser(),
		tags:           tagObjs,
	}, nil
}

// Check tries to parse a file using both its name and its contents and checks if the file should implement
// a tag among the list of provided tags when creating the linter.
func (l *Linter) Check(fileName string) ([]*Issue, error) {
	var issues []*Issue

	for _, tag := range l.tags {
		inName, err := l.fileNameParser.Parse(fileName, tag)
		if err != nil {
			issues = append(issues, &Issue{
				FileName: fileName,
				Reason:   err.Error(),
			})
			continue
		}

		inContents, err := l.contentsParser.Parse(fileName, tag)
		if err != nil {
			issues = append(issues, &Issue{
				FileName: fileName,
				Reason:   err.Error(),
			})
			continue
		}

		if inName && !inContents {
			issues = append(issues, &Issue{
				FileName: fileName,
				Reason:   fmt.Sprintf("tag %s found in the file name but missing from the build tags", tag.Name()),
			})
			continue
		}

		if !inName && inContents {
			issues = append(issues, &Issue{
				FileName: fileName,
				Reason:   fmt.Sprintf("tag %s found in the file build tags but missing from the file name", tag.Name()),
			})
			continue
		}
	}

	return issues, nil
}

// Issue represents an issue that occurred when parsing a file.
type Issue struct {
	FileName string
	Reason   string
}
