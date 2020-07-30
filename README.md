# filebuildtag
Linter that matches Go file naming patterns to their expected build tags.

---

[![GoDoc](https://godoc.org/github.com/aziule/filebuildtag?status.svg)](https://godoc.org/github.com/aziule/filebuildtag)
[![Go Report Card](https://goreportcard.com/badge/github.com/aziule/filebuildtag)](https://goreportcard.com/report/github.com/aziule/filebuildtag)

[Jump to Installation and usage](#installation-and-usage)

## Benefits

Match file naming patterns to build tags and make sure these files always have the correct build tags in the `// +build` instruction.

## Features

### One-To-One match

Example: files named `foo.go` must include the `bar` build tag.

File: `foo.go`
```go
// +build bar

package foo
```

### Many-To-One match

Example: files ending with `_suffix.go` must include the `bar` build tag.

File: `a_suffix.go`
```go
// +build bar

package foo
```

File: `b_suffix.go`
```go
// +build bar

package foo
```

### Go's `buildtag` linter support

Built on top of the `buildtag` linter, it supports its Go files features.

### And also

* Run it as a standalone command
* Integrate it as a part of a runner using the provided `analysis.Analyzer`

## Real world use case

Let's say you name your integration tests `*_integration_test.go` and you run them as part of your CI pipeline.

You might forget to add the `// +build integration` instruction on a newly created file, or you might remove it inadvertently, 
and it can have some consequences, such as never running during your pipeline.

As a consequence, you think your code works when it doesn't, because it is not tested (but you believe it is).

This linter can help with such issues and let you know when you forgot to add the expected build tags.

## Installation and usage

**Install with Go install**

```shell
GO111MODULE=on go get github.com/aziule/filebuildtag/cmd/filebuildtag
```

**Install and build from source**
1. Clone the repo
2. Build the executable
```shell
make build
```

**Usage**

```shell
// All files named "foo.go" must have the "bar" tag
filebuildtag -filetags "foo.go:bar" .

// All files ending with "_integration_test.go" must have the "integration" tag
filebuildtag -filetags "*_integration_test.go:integration" .

// Both of the above
filebuildtag -filetags "foo.go:bar,*_integration_test.go:integration" .
```

*Note: file name patterns match using Go's `filepath.Match` method and therefore support all of its features.
See [File patterns](#file-patterns) for more information and examples.*

## Using with runners

To facilitate the integration with existing linter runners, you can use the `Analyzer` provided:
```go
cfg := filebuildtag.Config{}
cfg.WithFiletag("foo.go", "tag1").
    WithFiletag("*_suffix.go", "tag2")
analyzer := NewAnalyzer(cfg)
```

## File patterns

### Syntax

From the official `filepath.Match` doc, file patterns syntax is the following:

```
pattern:
	{ term }
term:
	'*'         matches any sequence of non-Separator characters
	'?'         matches any single non-Separator character
	'[' [ '^' ] { character-range } ']'
	            character class (must be non-empty)
	c           matches character c (c != '*', '?', '\\', '[')
	'\\' c      matches character c

character-range:
	c           matches character c (c != '\\', '-', ']')
	'\\' c      matches character c
	lo '-' hi   matches character c for lo <= c <= hi
```

### Examples

| file 👇 pattern 👉 | foo.go | *.go | ba?.go | *_test.go |
|--------------------|--------|------|--------|-----------|
| foo.go             | ✅     | ✅   | 🚫     | 🚫        |
| bar.go             | 🚫     | ✅   | ✅     | 🚫        |
| baz.go             | 🚫     | ✅   | ✅     | 🚫        |
| a_test.go          | 🚫     | ✅   | 🚫     | ✅        |
| something          | 🚫     | 🚫   | 🚫     | 🚫        |

## Roadmap

* Support for folder name matching (`/pkg/**/foo.go`, `/pkg/foo/*.go`, etc.).

## Contributing

A bug to report? A feature to add? Please feel free to open an issue or to propose pull requests!

## License

Some of the code was copied from Go's `buildtag` linter and adapted to match the needs of the `filebuildtag` linter.
Those files have the mandatory copyright header and their license can be found in `LICENSE.google`.

You can also find the original code [here](https://github.com/golang/tools/tree/master/go/analysis/passes/buildtag).
