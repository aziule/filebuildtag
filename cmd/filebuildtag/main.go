package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/aziule/filebuildtag"
	"golang.org/x/tools/go/analysis/singlechecker"
)

const tagsDoc = `Map file names to build tags using the form "pattern:tag". For example:
- Single file: "*foo.go:tag1"
- Multiple files: "*foo.go:tag1,*foo2.go:tag2"`

func main() {
	filetags := ""
	flag.StringVar(&filetags, "filetags", "", tagsDoc)
	flag.Parse()

	if filetags == "" {
		fmt.Println(`Missing mandatory flag "filetags":`)
		flag.PrintDefaults()
		return
	}
	filetagMap, err := toMap(filetags)
	if err != nil {
		fmt.Println(err)
		flag.PrintDefaults()
		return
	}

	analyzer := filebuildtag.NewAnalyzer(filetagMap)
	singlechecker.Main(analyzer)
}

func toMap(value string) (map[string]string, error) {
	filetags := strings.Split(value, ",")
	m := make(map[string]string, len(filetags))
	for i := 0; i < len(filetags); i++ {
		filetag := strings.TrimSpace(filetags[i])
		parts := strings.Split(filetag, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("Malformed argument: %s", parts)
		}
		m[parts[0]] = parts[1]
	}
	return m, nil
}
