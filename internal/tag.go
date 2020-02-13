package internal

import "errors"

type Tag struct {
	name string
	len  int
}

func NewTag(name string) (*Tag, error) {
	if name == "" {
		return nil, errors.New("empty tag")
	}

	return &Tag{
		name: name,
		len:  len(name),
	}, nil
}
