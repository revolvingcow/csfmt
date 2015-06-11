package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	extensions = [...]string{".cs"}
)

// SourceFile represents a file declared as source code.
type SourceFile struct {
	Path string
}

func (f *SourceFile) Exists() bool {
	if _, err := os.Stat(f.Path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (f *SourceFile) IsDir() bool {
	fi, err := os.Stat(f.Path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func (f *SourceFile) IsDotNet() bool {
	name := strings.ToLower(f.Path)
	for _, extension := range extensions {
		if strings.HasSuffix(name, extension) {
			return true
		}
	}
	return false
}

// Read the file contents.
func (f *SourceFile) Read() ([]byte, error) {
	contents, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

// Write the contents out to a stream.
func (f *SourceFile) Write(contents []byte) error {
	fi, err := os.OpenFile(f.Path, os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer fi.Close()

	_, err = fi.Write(contents)
	if err != nil {
		return err
	}
	return nil
}

// Walk a directory's file structure looking for source files.
func (f *SourceFile) Walk() chan SourceFile {
	c := make(chan SourceFile)

	go func() {
		filepath.Walk(f.Path, func(p string, fi os.FileInfo, e error) error {
			if e != nil {
				return e
			}

			s := SourceFile{
				Path: p,
			}

			if s.IsDotNet() {
				c <- s
			}

			return nil
		})

		defer close(c)
	}()

	return c
}
