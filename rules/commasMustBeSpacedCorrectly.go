package rules

import (
	"bytes"
	"regexp"
)

var commasMustBeSpacedCorrectly = &Rule{
	Name:        "Commas must be spaced correctly",
	Enabled:     true,
	Apply:       applyCommasMustBeSpacedCorrectly,
	Description: ``,
}

func applyCommasMustBeSpacedCorrectly(source []byte) []byte {
	// Look for leading spaces
	re := regexp.MustCompile(`(\n)*[\s]+,`)
	for re.Match(source) {
		source = re.ReplaceAllLiteral(source, []byte(","))
	}

	return scan(source, func(line, literal []byte) []byte {
		// Add trailing spaces as necessary
		re = regexp.MustCompile(`(\S),(\w|\d)`)
		for re.Match(line) {
			if !bytes.Contains(literal, re.Find(line)) {
				line = re.ReplaceAll(line, []byte("$1, $2"))
			}
		}

		// Look for too many trailing spaces
		re = regexp.MustCompile(`\,  `)
		for re.Match(line) {
			if !bytes.Contains(literal, re.Find(line)) {
				line = re.ReplaceAll(line, []byte(", "))
			}
		}
		return line
	})
}
