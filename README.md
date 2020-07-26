# filebuildtag
Linter to check that Go files have the expected `// +build <tag>` instruction

---

[Jump to the installation and usage.](#installation-and-usage)

## Benefits

Match file name patterns to build tags and make sure these files always have the correct build tags in the `// +build` instruction.

Built on top of Go's `buildtag` linter, it supports all of the features related to Go files.

## Features

### One-To-One match

Example: files named `foo.go` must include the `bar` build tag.

File: `foo.go`
```go
// +build foo

package foo
```

### Many-To-One match

Example: file ending with `_suffix.go` must include the `bar` build tag.

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

Supports features from the linter on Go files.

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

*Note: file name patterns match using Go's `filepath.Match` method and therefore support all of its features.*

## Using with runners

To facilitate the integration with existing linter runners, you can use the `Analyzer` provided:
```go
cfg := filebuildtag.Config{}
cfg.WithFiletag("foo", "tag1").
   .WithFileTag("bar", "tag2")
analyzer := NewAnalyzer(cfg)
```

## Contributing

A bug to report? A feature to add? Please feel free to open an issue or to propose pull requests!

## License

Some of the code was copied from Go's `buildtag` linter and adapted to match the needs of the `filebuildtag` linter.
Those files have the mandatory copyright header and their license can be found in `LICENSE.google`.

You can also find the original code [here](https://github.com/golang/tools/tree/master/go/analysis/passes/buildtag).