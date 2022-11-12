package fsmock_test

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/halimath/fsmock"
)

func Example() {
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

	f, err := fsys.Open("cmd/main.go")
	if err != nil {
		panic(err)
	}

	c, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(c))

	fmt.Println("---")

	fs.WalkDir(fsys, "", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			fmt.Println(path)
		}

		return nil
	})

	fmt.Println("---")

	cmdTests, err := fs.Glob(fsys, "cmd/*_test.go")
	if err != nil {
		panic(err)
	}
	fmt.Println(len(cmdTests))

	if err := fsys.Touch("cmd/main_test.go"); err != nil {
		panic(err)
	}

	_, err = fsys.ReadFile("cmd/main_test.go")
	if err != nil {
		panic(err)
	}

	cmdTests, err = fs.Glob(fsys, "cmd/*_test.go")
	if err != nil {
		panic(err)
	}
	fmt.Println(len(cmdTests))

	fmt.Println("---")

	if err := fsys.Mkdir("internal/foo"); err != nil {
		panic(err)
	}

	if err := fsys.Touch("internal/foo/foo.go"); err != nil {
		panic(err)
	}

	fsub, err := fsys.Sub("internal/foo")
	if err != nil {
		panic(err)
	}

	if _, err := fsub.Open("foo.go"); err != nil {
		panic(err)
	}

	fmt.Println("---")

	fsys.Rm("internal")

	fs.WalkDir(fsys, "", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fmt.Println(path)

		return nil
	})

	// Output: package main
	// ---
	// go.mod
	// go.sum
	// cmd/main.go
	// internal/tool.go
	// internal/tool_test.go
	// ---
	// 0
	// 1
	// ---
	// ---
	//
	// go.mod
	// go.sum
	// cmd
	// cmd/main.go
	// cmd/main_test.go
}
