# fsmock

A golang mock implementation of `fs.FS` and friends for testing.

![CI Status][ci-img-url] 
[![Go Report Card][go-report-card-img-url]][go-report-card-url] 
[![Package Doc][package-doc-img-url]][package-doc-url] 
[![Releases][release-img-url]][release-url]

[ci-img-url]: https://github.com/halimath/fsmock/workflows/CI/badge.svg
[go-report-card-img-url]: https://goreportcard.com/badge/github.com/halimath/fsmock
[go-report-card-url]: https://goreportcard.com/report/github.com/halimath/fsmock
[package-doc-img-url]: https://img.shields.io/badge/GoDoc-Reference-blue.svg
[package-doc-url]: https://pkg.go.dev/github.com/halimath/fsmock
[release-img-url]: https://img.shields.io/github/v/release/halimath/fsmock.svg
[release-url]: https://github.com/halimath/fsmock/releases

`fsmock` implements a mock filesystem satisfying `fs.FS` and other interfaces
to enable easier testing of code that uses `fs.FS` to access file systems.

# Installation

`fsmock` is provided as a go module and requires go >= 1.16.

```shell
go get github.com/halimath/fsmock@main
```

# Usage

`fsmock` provides two basic types `Dir` and `File`. These can be used to build
up a filesystem in plain go. Use the provided functions `NewDir` and `NewFile`
to create them conveniently.

Create a new filesystem by invoking `fsmock.New` providing a root directory.

```go
fsys := fsmock.New(fsmock.NewDir("",
    fsmock.EmptyFile("go.mod"),
    fsmock.EmptyFile("go.sum"),
    fsmock.NewDir("cmd",
        fsmock.TextFile("main.go", "package main"),
    ),
    fsmock.NewDir("internal",
        fsmock.EmptyFile("tool.go"),
        fsmock.EmptyFile("tool_test.go"),
    ),
))
```

You can use the methods defined in `fs.FS` and other interfaces from the `fs`
package to access files and directories:

```go
f, err := fsys.Open("cmd/main.go")
if err != nil {
    panic(err)
}
```

```go
fs.WalkDir(fsys, "", func(path string, d fs.DirEntry, err error) error {
    // ...
    return nil
})
```

## Modifying the filesystem

In addition to the read-only functions defined by the `fs` interfaces, `fsmock`
provides some helper functions to modify the filesystem.

To update a modification timestamp of either a file or directory use the
`Touch` function. This will also create an empty file if the named file does
not exist (just like the Unix `touch` command does):

```go
if err := fsys.Touch("internal/foo/foo.go"); err != nil {
    panic(err)
}
```

To create a directory use the `Mkdir` function which works like the `mkdir`
Unix shell command (_without_ the `-p` option):

```go
if err := fsys.Mkdir("internal/foo"); err != nil {
    panic(err)
}
```

See [`fsmock_test.go`](./fsmock_test.go) for a full-blown example.

# License

Copyright 2022 Alexander Metzner.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
