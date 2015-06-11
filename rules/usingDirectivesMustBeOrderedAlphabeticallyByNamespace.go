package rules

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var usingDirectivesMustBeOrderedAlphabeticallyByNamespace = &Rule{
	Name:        "Using directives must be ordered alphabetically by namespace",
	Enabled:     true,
	Apply:       applyUsingDirectivesMustBeOrderedAlphabeticallyByNamespace,
	Description: ``,
}

func applyUsingDirectivesMustBeOrderedAlphabeticallyByNamespace(source []byte) []byte {
	reCommentShortBegin := regexp.MustCompile(`\A\s*(/{2,})`)
	reCommentLongBegin := regexp.MustCompile(`\A\s*(/\*)`)
	reCommentLongEnd := regexp.MustCompile(`\*/`)

	short := false
	long := false
	skipNewline := false
	lines := []byte{}
	usings := []string{}

	buffer := bytes.NewBuffer(source)
	scanner := bufio.NewScanner(buffer)
	for scanner.Scan() {
		// Add a newline character on each line after the first
		if len(lines) > 0 && !skipNewline {
			lines = append(lines, byte('\n'))
		}

		skipNewline = false
		line := scanner.Bytes()
		short = reCommentShortBegin.Match(line)
		long = long || reCommentLongBegin.Match(line)

		if !(short || long) {
			// Find usings
			re := regexp.MustCompile(`(.*)(using)([\t| ])([^\(])(.*)([;])`)
			if re.Match(line) {
				using := re.ReplaceAll(line, []byte("$2 $4$5"))
				usings = append(usings, string(using))
				skipNewline = true
				line = []byte{}
			}
		}

		short = false
		long = long && !reCommentLongEnd.Match(line)
		lines = append(lines, line...)
	}

	if len(usings) > 0 {
		// Sort the usings and add them to the top of the file
		s := sort.StringSlice(usings)
		if !sort.IsSorted(s) {
			sort.Sort(s)
		}

		lines = append([]byte(fmt.Sprintf("%s;\n\n", strings.Join(s, ";\n"))), lines...)
	}

	return lines
}
