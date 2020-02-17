package internal

import "errors"

// Tag represents a build tag.
type Tag struct {
	name string
	len  int
}

// NewTag creates a new Tag.
func NewTag(name string) (*Tag, error) {
	if name == "" {
		return nil, errors.New("empty tag")
	}

	return &Tag{
		name: name,
		len:  len(name),
	}, nil
}

// Name returns the tag's name
func (t *Tag) Name() string {
	return t.name
}
