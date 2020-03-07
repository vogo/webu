// +build ignore

// Command makeStatic reads a set of files and writes a Go source file to "static.go"
// that declares a map of string constants containing contents of the input files.
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/vogo/vogo/vio/vioutil"
)

func main() {
	pk := os.Args[1]
	srcDir := os.Args[2]
	targetDir := os.Args[3]

	if err := makeStatic(pk, srcDir, targetDir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func makeStatic(pk, srcDir, targetDir string) error {
	staticFile := filepath.Join(targetDir, "static.go")
	f, err := os.Create(staticFile)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%v\n\n%v\n\npackage %s\n\n", license, warning, pk)
	fmt.Fprintf(buf, "var StaticFiles = map[string]string{\n")

	files, err := vioutil.ListFileNames(srcDir, "", ".html")
	if err != nil {
		return err
	}

	for _, file := range files {
		b, err := ioutil.ReadFile(filepath.Join(srcDir, file))
		if err != nil {
			return err
		}
		fmt.Fprintf(buf, "\t%q: ", file)
		if utf8.Valid(b) {
			fmt.Fprintf(buf, "`%s`", sanitize(b))
		} else {
			fmt.Fprintf(buf, "%q", b)
		}
		fmt.Fprintln(buf, ",\n")
	}
	fmt.Fprintln(buf, "}")
	fmtBuf, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	return ioutil.WriteFile(staticFile, fmtBuf, 0666)
}

// sanitize prepares a valid UTF-8 string as a raw string constant.
func sanitize(b []byte) []byte {
	// Replace ` with `+"`"+`
	b = bytes.Replace(b, []byte("`"), []byte("`+\"`\"+`"), -1)

	// Replace BOM with `+"\xEF\xBB\xBF"+`
	// (A BOM is valid UTF-8 but not permitted in Go source files.
	// I wouldn't bother handling this, but for some insane reason
	// jquery.js has a BOM somewhere in the middle.)
	return bytes.Replace(b, []byte("\xEF\xBB\xBF"), []byte("`+\"\\xEF\\xBB\\xBF\"+`"), -1)
}

const warning = `// Code generated by "makestatic.go"; DO NOT EDIT.`

var license = `// Copyright 2019 wongoo. All rights reserved.`
