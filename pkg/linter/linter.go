package main

import (
	"errors"

	"github.com/aziule/gofilebuildtags/internal"
)

var ErrEmptyTagsList = errors.New("empty tags list")

type Linter struct {
	fileNameParser internal.TagParser
	contentsParser internal.TagParser
	tags           []*internal.Tag
}

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

	return &Linter{
		fileNameParser: internal.NewFileNameParser(internal.TestFileSuffix, internal.GoFileExt),
		contentsParser: internal.NewContentsParser(),
		tags:           tagObjs,
	}, nil
}

func (l *Linter) Check(fileName string) (bool, error) {
	for _, tag := range l.tags {
		_, err := l.fileNameParser.Parse(fileName, tag)
		if err != nil {
			return false, err
		}

		_, err = l.contentsParser.Parse(fileName, tag)
		if err != nil {
			return false, err
		}

		// fmt.Println(fileName, tag, inName, inFile)
	}

	return true, nil
}

type Issue struct {
	FileName string
	Reason   string
	Line     int
}
