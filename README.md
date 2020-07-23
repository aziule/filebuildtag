# filebuildtag

Check that Go files have the expected build tags in the "// +build" instruction

Define file patterns and assign them to build tags, for instance:
* File "bar.go" must have the "baz" build tag
* Files "*_integration_test.go" must have the "integration" build tag