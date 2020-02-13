package gofilebuildtags

import "errors"

type tag struct {
	name string
	len  int
}

func newTag(name string) (*tag, error) {
	if name == "" {
		return nil, errors.New("empty tag")
	}

	return &tag{
		name: name,
		len:  len(name),
	}, nil
}
