package internal

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	TestFileSuffix string = "_test"
	GoFileExt             = ".go"
)

var (
	ErrEmptyFileName = errors.New("empty file name")
)

type TagParser interface {
	Parse(string, *Tag) (bool, error)
}

type FileNameParser struct {
	suffix    string
	suffixLen int
	ext       string
	extLen    int
}

func NewFileNameParser(suffix, ext string) *FileNameParser {
	return &FileNameParser{
		suffix:    suffix,
		suffixLen: len(suffix),
		ext:       ext,
		extLen:    len(ext),
	}
}

func (p *FileNameParser) Parse(fileName string, tag *Tag) (bool, error) {
	if fileName == "" {
		return false, ErrEmptyFileName
	}

	ext := filepath.Ext(fileName)

	if ext == "" {
		return false, nil
	}

	if ext != p.ext {
		return false, nil
	}

	namelen := len(fileName)

	if namelen-p.suffixLen-p.extLen < tag.len {
		return false, nil
	}

	return fileName[namelen-tag.len-p.suffixLen-p.extLen:namelen-p.suffixLen-p.extLen] == tag.name, nil
}

type ContentsParser struct {
	expected map[string]string
}

func NewContentsParser() *ContentsParser {
	return &ContentsParser{
		expected: make(map[string]string),
	}
}

func (p *ContentsParser) Parse(fileName string, tag *Tag) (bool, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return false, fmt.Errorf("could not open file %s: %v", fileName, err)
	}
	defer f.Close()

	expected, ok := p.expected[tag.name]
	if !ok {
		expected = "// +build " + tag.name
		p.expected[tag.name] = expected
	}

	buf := bufio.NewReader(f)
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

		if line == expected {
			hasTag = true
		}
	}

	return hasTag, nil
}
