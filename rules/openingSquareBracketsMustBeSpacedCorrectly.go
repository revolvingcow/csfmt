package rules

import (
	"bytes"
	"regexp"
)

var openingSquareBracketsMustBeSpacedCorrectly = &Rule{
	Name:        "Opening square brackets must be spaced correctly",
	Enabled:     true,
	Apply:       applyOpeningSquareBracketsMustBeSpacedCorrectly,
	Description: ``,
}

func applyOpeningSquareBracketsMustBeSpacedCorrectly(source []byte) []byte {
	return scan(source, func(line, literal []byte) []byte {
		re := regexp.MustCompile(`([\S])([\t ]+)([\[])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1$3"))
		}

		re = regexp.MustCompile(`([\[])([\t ]+)([\S])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1$3"))
		}

		return line
	})
}
