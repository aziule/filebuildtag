# filebuildtag
Linter to check that Go files have the expected `// +build <tag>` instruction

---

[Jump to the installation and usage.](#installation-and-usage)

## Benefits

Match file name patterns to build tags and make sure these files always have the correct build tags in the `// +build` instruction.

Built on top of Go's `buildtag` linter, it supports all of the features related to Go files.

## Features

**1-to-1 match**
Every file named `foo.go` must include the `bar` build tag.

**many-to-1 match**
Every file named `*foo.go` must include the `bar` build tag.

**Go's `buildtag` linter support**
Supports features from the linter on Go files.

**And also**
* Run it as a standalone command using `cmd/filebuildtag`.
* Integrate it as a part of a runner using the provided `analysis.Analyzer`.

## Use cases

* Avoid running CI with tests you thought should run because you forgot to add the expected build tag to your files.
For example, never forget to add `// +build integration` to your integration test files named `*_integration_test.go`.

## Installation and usage

**With Go install**

```shell
GO111MODULE=on go get github.com/aziule/filebuildtag/cmd/filebuildtag
```

**Build from source**
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

*Note: patterns are matched using Go's `filepath.Match` method and support all of its features.*

## Using with runners

To facilitate the integration with existing linter runners, you can use the `Analyzer` provided:
```go
cfg := filebuildtag.Config{}
cfg.WithFiletag("foo", "tag1").
   .WithFileTag("bar", "tag2")
analyzer := NewAnalyzer(cfg)
```

## License

Some of the code was copied from Go's `buildtag` linter and adapted to match the needs of the `filebuildtag` linter.
Those files have the mandatory copyright header and their license can be found in `LICENSE.google`.

You can also find the original code [here](https://github.com/golang/tools/tree/master/go/analysis/passes/buildtag).