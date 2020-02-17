package internal

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
)

const (
	TestFileSuffix string = "_test"

	regexFmt       string = `(?:[\w_])*%s%s\.go`
	contentsPrefix string = "// +build "
)

var (
	ErrEmptyFileName = errors.New("empty file name")
)

type TagParser interface {
	Parse(string, *Tag) (bool, error)
}

type FileNameParser struct {
	regexs map[string]*regexp.Regexp
}

func NewFileNameParser(suffix string, tags []*Tag) (*FileNameParser, error) {
	p := &FileNameParser{
		regexs: make(map[string]*regexp.Regexp),
	}

	for _, tag := range tags {
		regex, err := regexp.Compile(fmt.Sprintf(regexFmt, tag.name, suffix))
		if err != nil {
			return nil, fmt.Errorf("could not compile regex for tag %s: %v", tag.name, err)
		}
		p.regexs[tag.name] = regex
	}

	return p, nil
}

func (p *FileNameParser) Parse(fileName string, tag *Tag) (bool, error) {
	if fileName == "" {
		return false, ErrEmptyFileName
	}

	regex, ok := p.regexs[tag.name]
	if !ok {
		return false, errors.New("tag mismatch")
	}

	return regex.Match([]byte(fileName)), nil
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
		expected = contentsPrefix + tag.name
		p.expected[tag.name] = expected
	}

	buf := bufio.NewReader(f)
	hasTag := false
	done := false

	for !done {
		line, _, err := buf.ReadLine()
		if err != nil {
			done = true
			continue
		}

		if len(line) < len(contentsPrefix)+tag.len {
			done = true
			continue
		}

		if string(line) == expected {
			hasTag = true
			done = true
			continue
		}

		if string(line)[:len(contentsPrefix)] == contentsPrefix {
			continue
		}
	}

	return hasTag, nil
}
