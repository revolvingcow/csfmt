package rules

import (
	"bufio"
	"bytes"
	"regexp"
	"unicode"

	"github.com/revolvingcow/csfmt"
)

var singleLineCommentsMustBeginWithSingleSpace = &csfmt.Rule{
	Name:        "Single line comments must begin with single space",
	Enabled:     true,
	Apply:       applySingleLineCommentsMustBeginWithSingleSpace,
	Description: ``,
}

func applySingleLineCommentsMustBeginWithSingleSpace(source []byte) []byte {
	reString := regexp.MustCompile(`".*"`)

	lines := []byte{}
	buffer := bytes.NewBuffer(source)
	scanner := bufio.NewScanner(buffer)

	for scanner.Scan() {
		// Add a newline character on each line after the first
		if len(lines) > 0 {
			lines = append(lines, byte('\n'))
		}

		line := scanner.Bytes()

		literal := []byte{}
		if reString.Match(line) {
			literal = reString.Find(line)
		}

		// Handle comments with no space.
		re := regexp.MustCompile(`(\s*)[/]{2}\s{0}(\S+)`)
		for re.Match(line) {
			if !bytes.Contains(literal, re.Find(line)) {
				line = re.ReplaceAll(line, []byte("$1// $2"))
			}
		}

		// Handle comments with more than one space
		re = regexp.MustCompile(`(\s*)[/]{2}\s{2,}(\S+)`)
		for re.Match(line) {
			if !bytes.Contains(literal, re.Find(line)) {
				line = re.ReplaceAll(line, []byte("$1// $2"))
			}
		}

		// Adjust for URIs
		re = regexp.MustCompile(`(\s*)[:]{1}[/]{2}\s+(\S+)`)
		for re.Match(line) {
			if !bytes.Contains(literal, re.Find(line)) {
				line = re.ReplaceAll(line, []byte("$1://$2"))
			}
		}

		// Trim ending
		line = bytes.TrimRightFunc(line, unicode.IsSpace)
		lines = append(lines, line...)
	}

	return lines
}
