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

var filetags string

func main() {
	flag.StringVar(&filetags, "filetags", "", tagsDoc)
	flag.Parse()

	if filetags == "" {
		fmt.Println(`missing mandatory flag "filetags"`)
		flag.PrintDefaults()
		return
	}

	cfg, err := parseFlags(filetags)
	if err != nil {
		fmt.Println(err)
		flag.PrintDefaults()
		return
	}

	analyzer := filebuildtag.NewAnalyzer(*cfg)
	singlechecker.Main(analyzer)
}

func parseFlags(filetags string) (*filebuildtag.Config, error) {
	cfg := &filebuildtag.Config{}
	args := strings.Split(filetags, ",")
	for i := 0; i < len(args); i++ {
		filetag := strings.TrimSpace(args[i])
		parts := strings.Split(filetag, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed argument: %s", parts)
		}

		parts[0] = strings.TrimSpace(parts[0])
		parts[1] = strings.TrimSpace(parts[1])
		if parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("empty argument: %s", parts)
		}
		cfg.WithFiletag(parts[0], parts[1])
	}
	return cfg, nil
}
