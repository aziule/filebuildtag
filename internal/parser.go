package internal

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
)

const (
	// TestFileSuffix represents the suffix expected in Go test files
	TestFileSuffix string = "_test"

	regexFmt       string = `(?:[\w_])*%s%s\.go`
	contentsPrefix string = "// +build "
)

var (
	ErrEmptyFileName = errors.New("empty file name")
)

// TagParser is the interface used to parse a tag from a string.
type TagParser interface {
	Parse(string, *Tag) (bool, error)
}

// FileNameParser implements the TagParser interface and parses a tag from a file name.
type FileNameParser struct {
	regexs map[string]*regexp.Regexp
}

// NewFileNameParser creates a new FileNameParser.
func NewFileNameParser(tags []*Tag) (*FileNameParser, error) {
	p := &FileNameParser{
		regexs: make(map[string]*regexp.Regexp),
	}

	for _, tag := range tags {
		regex, err := regexp.Compile(fmt.Sprintf(regexFmt, tag.name, TestFileSuffix))
		if err != nil {
			return nil, fmt.Errorf("could not compile regex for tag %s: %v", tag.name, err)
		}
		p.regexs[tag.name] = regex
	}

	return p, nil
}

// Parse implementation.
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

// ContentsParser implements the TagParser interface and parses the content of a file, given its name.
type ContentsParser struct {
	expected map[string]string
}

// NewContentsParser creates a new ContentsParser.
func NewContentsParser() *ContentsParser {
	return &ContentsParser{
		expected: make(map[string]string),
	}
}

// Parse implementation.
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

		// Allow blank lines, as per Go's specs regarding build tags
		if len(line) == 0 {
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
