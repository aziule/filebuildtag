package internal

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

const goExt string = ".go"

var (
	ErrEmptyFileName = errors.New("empty file name")
)

type TagParser interface {
	Parse(string, *Tag) (bool, error)
}

type FileNameParser struct {
	suffix    string
	suffixLen int
}

func NewFileNameParser(suffix string) *FileNameParser {
	return &FileNameParser{
		suffix:    suffix + goExt,
		suffixLen: len(suffix) + len(goExt),
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

	if ext != goExt {
		return false, nil
	}

	namelen := len(fileName)

	if namelen-p.suffixLen < tag.len {
		return false, nil
	}

	return fileName[namelen-tag.len-p.suffixLen:namelen-p.suffixLen] == tag.name, nil
}

type ContentsParser struct {
	regexTemplate string
	regex         map[string]*regexp.Regexp
}

func NewContentsParser() *ContentsParser {
	return &ContentsParser{
		regexTemplate: "^/",
		regex:         make(map[string]*regexp.Regexp),
	}
}

func (p *ContentsParser) Parse(fileName string, tag *Tag) (bool, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return false, fmt.Errorf("could not open file %s: %v", fileName, err)
	}
	defer f.Close()

	reg, ok := p.regex[tag.name]
	if !ok {
		reg = regexp.MustCompile("")
		p.regex[tag.name] = reg
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

		if line[:9] != "// +build" {
			done = true
			continue
		}
		hasTag = true
	}

	return hasTag, nil
}
