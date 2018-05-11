package rules

import (
	"bytes"
	"regexp"

	"github.com/revolvingcow/csfmt"
)

var semicolonsMustBeSpacedCorrectly = &csfmt.Rule{
	Name:        "Semicolons must be spaced correctly",
	Enabled:     true,
	Apply:       applySemicolonsMustBeSpacedCorrectly,
	Description: ``,
}

func applySemicolonsMustBeSpacedCorrectly(source []byte) []byte {
	return scan(source, func(line, literal []byte) []byte {
		// Look for leading spaces
		re := regexp.MustCompile(`\s;`)
		for re.Match(line) {
			if !bytes.Contains(literal, re.Find(line)) {
				line = re.ReplaceAllLiteral(line, []byte(";"))
			}
		}

		// Add trailing spaces as necessary
		re = regexp.MustCompile(`;(\S)`)
		for re.Match(line) {
			if !bytes.Contains(literal, re.Find(line)) {
				line = re.ReplaceAll(line, []byte("; $1"))
			}
		}
		return line
	})
}
