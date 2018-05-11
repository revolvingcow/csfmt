package rules

import (
	"bufio"
	"bytes"
	"regexp"
	"unicode"
)

func scan(source []byte, applyFunc func(line, literal []byte) []byte) []byte {
	reCommentShortBegin := regexp.MustCompile(`\A\s*(/{2,})`)
	reCommentLongBegin := regexp.MustCompile(`\A\s*(/\*)`)
	reCommentLongEnd := regexp.MustCompile(`\*/`)
	reString := regexp.MustCompile(`".*"`)

	short := false
	long := false
	lines := []byte{}
	buffer := bytes.NewBuffer(source)
	scanner := bufio.NewScanner(buffer)

	for scanner.Scan() {
		// Add a newline character on each line after the first
		if len(lines) > 0 {
			lines = append(lines, byte('\n'))
		}

		line := scanner.Bytes()
		short = reCommentShortBegin.Match(line)
		long = long || reCommentLongBegin.Match(line)

		literal := []byte{}
		if reString.Match(line) {
			literal = reString.Find(line)
		}

		if !(short || long) {
			line = applyFunc(line, literal)

			// Trim ending
			line = bytes.TrimRightFunc(line, unicode.IsSpace)
		}

		short = false
		long = long && !reCommentLongEnd.Match(line)
		lines = append(lines, line...)
	}

	return lines
}
