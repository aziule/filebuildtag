package gofilebuildtags

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

var ErrEmptyTagsList = errors.New("empty tags list")

type Linter struct {
	tags []*tag
}

func NewLinter(tags []string) (*Linter, error) {
	nbTags := len(tags)

	if nbTags == 0 {
		return nil, ErrEmptyTagsList
	}

	seenTags := make(map[string]bool)
	var tagObjs []*tag

	for _, tag := range tags {
		if _, ok := seenTags[tag]; ok {
			continue
		}

		seenTags[tag] = true

		t, err := newTag(tag)
		if err != nil {
			return nil, err
		}

		tagObjs = append(tagObjs, t)
	}

	if len(tagObjs) == 0 {
		return nil, ErrEmptyTagsList
	}

	return &Linter{
		tags: tagObjs,
	}, nil
}

func (l *Linter) Check(fileName string) (bool, error) {
	f, err := newFile(fileName)
	if err != nil {
		return false, err
	}

	for _, tag := range l.tags {
		hasTagInName := f.hasTagInName(tag)
		hasTagInFile := f.hasTagInFile(tag)

		fmt.Println(hasTagInName, hasTagInFile)
	}

	//requiredTags, err := requiredTagsBasedOnFileName(fileName)
	//if err != nil {
	//	issues = append(issues, &Issue{
	//		FileName: fileName,
	//		Reason:   err.Error(),
	//	})
	//}
	//actualTags, err := parseTags(fileName)
	//if err != nil {
	//	issues = append(issues, &Issue{
	//		FileName: fileName,
	//		Reason:   err.Error(),
	//	})
	//}
	//fmt.Println(requiredTags, actualTags)
	return true, nil
}

type Issue struct {
	FileName string
	Reason   string
	Line     int
}

//func requiredTagsBasedOnFileName(fileName string) ([]string, error) {
//	var requiredTags []string
//
//	for pattern, tag := range fileNamesTags {
//		m, err := regexp.Match(pattern, []byte(fileName))
//		if err != nil {
//			return nil, err
//		}
//		if m == false {
//			continue
//		}
//		requiredTags = append(requiredTags, tag)
//	}
//
//	return requiredTags, nil
//}

func parseTags(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("couldn't read: %v", scanner.Err())
		}
		//fmt.Println(">", string(scanner.Bytes()))
	}
	return nil, nil
}
