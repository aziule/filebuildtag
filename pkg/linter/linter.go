package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/aziule/gofilebuildtags/internal"
)

var ErrEmptyTagsList = errors.New("empty tags list")

type Linter struct {
	parser internal.Parser
	tags   []*internal.Tag
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

	// TODO: make it configurable
	parser := internal.NewTagParser("_test")
	return &Linter{
		parser: parser,
		tags:   tagObjs,
	}, nil
}

func (l *Linter) Check(fileName string) (bool, error) {
	for _, tag := range l.tags {
		_, err := l.parser.ParseFileName(tag, fileName)
		if err != nil {
			return false, err
		}

		f, err := os.Open(fileName)
		if err != nil {
			return false, fmt.Errorf("could not open file %s: %v", fileName, err)
		}
		defer f.Close()

		_, err = l.parser.ParseContents(tag, f)
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
