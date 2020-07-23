# filebuildtag

Check that Go files have the expected build tags in the `// +build` instruction

Define file patterns and assign them to build tags, for instance:
* File "bar.go" must have the "baz" build tag
* Files "*_integration_test.go" must have the "integration" build tag

# Example - Unit VS integration tests

Let's say we want to enforce the "integration" build tag on our integration test files.

Given the following, sample folder structure:
```
project
│
└───pkg
│   │   foo.go
│   │   bar.go
│
└───test
    │   foo.go
    │   bar_integration_test.go
    │   baz_integration_test.go
```

Test file `foo_test.go` is a unit test file.

Test files `bar_integration_test.go` and `baz_integration_test.go` should only run when the `integration` build 
tag is present.

File: foo.go
```go
package test

func Test_Foo(t *testing.T){ /* ... */ }
```

File: bar_integration_test.go
```go
// +build integration

package test

func Test_Bar(t *testing.T){ /* ... */ }
```

File: baz_integration_test.go
```go
// +build integration

package test

func Test_Baz(t *testing.T){ /* ... */ }
```

To make sure the integration test files always have the `// +build integration` instruction, use the following arguments
to the linter: `-filetag="*_integration_test.go:integration"`.