package main

import (
	"bytes"
	"log"
	"os"

	"code.revolvingcow.com/revolvingcow/csfmt/rules"
)

func main() {
	sourceFiles := []SourceFile{}

	// Determine what files to format
	argc := len(os.Args)
	if argc < 2 {
		return
	} else if argc == 2 && os.Args[1] == "..." {
		// Walk the file structure from the current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return
		}

		s := SourceFile{
			Path: cwd,
		}

		for sourceFile := range s.Walk() {
			sourceFiles = append(sourceFiles, sourceFile)
		}
	} else {
		// Assuming multiple files were given
		// for _, a := range os.Args[1:] {
		a := os.Args[1]
		s := SourceFile{
			Path: a,
		}

		if s.Exists() {
			if s.IsDir() {
				for sourceFile := range s.Walk() {
					sourceFiles = append(sourceFiles, sourceFile)
				}
			} else if s.IsDotNet() {
				sourceFiles = append(sourceFiles, s)
			}
		}
		// }
	}

	count := len(sourceFiles)
	modified := 0
	queuedRules := rules.EnabledRules()

	for _, s := range sourceFiles {
		contents, err := s.Read()
		if err != nil {
			log.Fatalln(err)
		}
		original := contents

		for _, rule := range queuedRules {
			contents = rule.Apply(contents)
		}

		if bytes.Compare(original, contents) != 0 {
			modified++
		}

		s.Write(contents)
		// fmt.Println(string(contents))
	}

	log.Printf("Modified %d of %d files using %d rules\n", modified, count, len(queuedRules))
}
