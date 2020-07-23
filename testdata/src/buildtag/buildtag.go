// Package buildtag contains the original test file from the Go tool's linter: "buildtag".
//
// The code was adapted to match the needs and requirements of the filebuildtag linter.
//
// Copyright 2013 The Go Authors. All rights reserved.
// Original code: https://github.com/golang/tools/tree/master/go/analysis/passes/buildtag

// This file contains tests for the buildtag checker.

// +builder // want `possible malformed \+build comment`
// +build !testfix foo

// Mention +build // want `possible malformed \+build comment`

// +build nospace // want "build comment must appear before package clause and be followed by a blank line"
package buildtag

// +build toolate // want "build comment must appear before package clause and be followed by a blank line$"

var _ = 3

var _ = `
// +build notacomment
`
