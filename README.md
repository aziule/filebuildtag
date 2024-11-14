# filebuildtag
Linter enforcing files to contain expected build tags (`// +build` instruction), based on the file name.

---

[![GoDoc](https://godoc.org/github.com/aziule/filebuildtag?status.svg)](https://godoc.org/github.com/aziule/filebuildtag)
[![Go Report Card](https://goreportcard.com/badge/github.com/aziule/filebuildtag)](https://goreportcard.com/report/github.com/aziule/filebuildtag)

[Jump to Installation and usage](#installation-and-usage)

## Real world use case

Let's say you put integration tests in files named `*_integration_test.go` and you run them as part of your CI pipeline.

You might forget to add the `// +build integration` instruction on a newly created file, or you might remove it inadvertently,
and it can have some consequences, such as never running during the CI pipeline.

As a consequence, you think your code works when it doesn't, because it is not tested, but you believe it is.

This linter can help with such issues and let you know when you forgot to add the expected build tags.

## Features

### Exact match

Example: files named `foo.go` must include the `bar` build tag.

File: `foo.go`
```go
// +build bar

package foo
```

### Wildcard match

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

`filebuildtag` is built on top of the `buildtag` linter, hence it supports its features.

### And also

* Run it as a standalone command
* Integrate it as a part of a runner using the provided `analysis.Analyzer`

## Installation and usage

**Install with Go install**

```shell
go get github.com/aziule/filebuildtag/cmd/filebuildtag
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
filebuildtag --filetags foo.go:bar ./...

// All files ending with "_integration_test.go" must have the "integration" tag
filebuildtag --filetags "*_integration_test.go:integration" ./...

// Both of the above
filebuildtag --filetags "foo.go:bar,*_integration_test.go:integration" ./...

// Only check that the `// +build` instructions are correct (no args to pass) 
filebuildtag ./...
```

*Note: files naming patterns are matched using Go's `filepath.Match` method. Therefore, you can use any of its supported patterns.
See [File patterns](#file-patterns) for more information and examples.*

Head to the [test scenarios](./filebuildtag_test.go) for more examples.

## Using with linters runners

This linter exposes an `Analyzer` (accessible via `filebuildtag.Analyzer`), which is defined as 
a `golang.org/x/tools/go/analysis/analysis.Analyzer` struct. 

Most of the linters runners expect linters to be defined like so, therefore you should not have much trouble integrating it
following the linters runner's doc.

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

## Development



**Run tests**

```shell
make test
```

**Lint**

```shell
make lint
```

## Roadmap

* Support for folder name matching (`/pkg/**/foo.go`, `/pkg/foo/*.go`, etc.).

## Contributing

A bug to report? A feature to add? Please feel free to open an issue or to propose pull requests!

## License

Some of the code was copied from Go's `buildtag` linter and adapted to match the needs of the `filebuildtag` linter.
Those files have the mandatory copyright header and their license can be found in `LICENSE.google`.

You can also find the original code [here](https://github.com/golang/tools/tree/master/go/analysis/passes/buildtag).
