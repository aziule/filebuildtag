package gofilebuildtags

import (
	"errors"
	"path/filepath"
)

var ErrEmptyFileName = errors.New("empty file name")

type file struct {
	name string
	len  int
	ext  string
}

func newFile(name string) (*file, error) {
	if name == "" {
		return nil, ErrEmptyFileName
	}

	return &file{
		name: name,
		len:  len(name),
		ext:  filepath.Ext(name),
	}, nil
}

func (f *file) hasTagInName(tag *tag) bool {
	if f.len-len(f.ext) < tag.len {
		return false
	}

	return f.name[f.len-tag.len-len(f.ext):f.len-len(f.ext)] == tag.name
}

func (f *file) hasTagInFile(tag *tag) bool {
	return true
}
