package internal

import (
	"bufio"
	"errors"
	"io"
	"path/filepath"
)

const goExt string = ".go"

var (
	ErrEmptyFileName = errors.New("empty file name")
)

type Parser interface {
	ParseFileName(tag *Tag, fileName string) (bool, error)
	ParseContents(tag *Tag, r io.Reader) (bool, error)
}

type TagParser struct {
	suffix    string
	suffixLen int
}

func NewTagParser(suffix string) *TagParser {
	return &TagParser{
		suffix:    suffix + goExt,
		suffixLen: len(suffix) + len(goExt),
	}
}

func (p *TagParser) ParseFileName(tag *Tag, fileName string) (bool, error) {
	if fileName == "" {
		return false, ErrEmptyFileName
	}

	ext := filepath.Ext(fileName)

	if ext == "" {
		return false, nil
	}

	if ext != goExt {
		return false, nil
	}

	namelen := len(fileName)

	if namelen-p.suffixLen < tag.len {
		return false, nil
	}

	return fileName[namelen-tag.len-p.suffixLen:namelen-p.suffixLen] == tag.name, nil
}

func (p *TagParser) ParseContents(tag *Tag, r io.Reader) (bool, error) {
	buf := bufio.NewReader(r)
	hasTag := false
	done := false

	for !done {
		line, err := buf.ReadString('\n')
		if err != nil {
			done = true
			continue
		}

		if len(line) <= 9+tag.len {
			done = true
			continue
		}

		if line[:9] != "// +build" {
			done = true
			continue
		}
		hasTag = true
	}

	return hasTag, nil
}
