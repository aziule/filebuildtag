// Copyright 2013 The Go Authors. All rights reserved.

package internal

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

func CheckGoFile(pass *analysis.Pass, f *ast.File) []string {
	tags := []string{}
	pastCutoff := false
	for _, group := range f.Comments {
		// A +build comment is ignored after or adjoining the package declaration.
		if group.End()+1 >= f.Package {
			pastCutoff = true
		}

		// "+build" is ignored within or after a /*...*/ comment.
		if !strings.HasPrefix(group.List[0].Text, "//") {
			pastCutoff = true
			continue
		}

		// Check each line of a //-comment.
		for _, c := range group.List {
			if !strings.Contains(c.Text, "+build") {
				continue
			}
			lineTags, err := checkLine(c.Text, pastCutoff)
			if err != nil {
				pass.Reportf(c.Pos(), "%s", err)
				continue
			}
			tags = append(tags, lineTags...)
		}
	}
	return tags
}

// checkLine checks a line that starts with "//" and contains "+build".
func checkLine(line string, pastCutoff bool) ([]string, error) {
	tags := []string{}
	line = strings.TrimPrefix(line, "//")
	line = strings.TrimSpace(line)

	if strings.HasPrefix(line, "+build") {
		fields := strings.Fields(line)
		if fields[0] != "+build" {
			// Comment is something like +buildasdf not +build.
			return nil, fmt.Errorf("possible malformed +build comment")
		}
		if pastCutoff {
			return nil, fmt.Errorf("+build comment must appear before package clause and be followed by a blank line")
		}
		argTags, err := checkArguments(fields)
		if err != nil {
			return nil, err
		}
		tags = append(tags, argTags...)
	} else {
		// Comment with +build but not at beginning.
		if !pastCutoff {
			return nil, fmt.Errorf("possible malformed +build comment")
		}
	}
	return tags, nil
}

func checkArguments(fields []string) ([]string, error) {
	tags := []string{}
	for _, arg := range fields[1:] {
		for _, elem := range strings.Split(arg, ",") {
			if strings.HasPrefix(elem, "!!") {
				return nil, fmt.Errorf("invalid double negative in build constraint: %s", arg)
			}

			isTagExcluded := strings.HasPrefix(elem, "!")
			if isTagExcluded {
				elem = elem[1:]
			}
			for _, c := range elem {
				if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' && c != '.' {
					return nil, fmt.Errorf("invalid non-alphanumeric build constraint: %s", arg)
				}
			}
			if isTagExcluded {
				continue
			}
			tags = append(tags, elem)
		}
	}
	return tags, nil
}
