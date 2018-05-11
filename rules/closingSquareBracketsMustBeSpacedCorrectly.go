package rules

import (
	"bytes"
	"regexp"

	"github.com/revolvingcow/csfmt"
)

var closingSquareBracketsMustBeSpacedCorrectly = &csfmt.Rule{
	Name:        "Closing square brackets must be spaced correctly",
	Enabled:     true,
	Apply:       applyClosingSquareBracketsMustBeSpacedCorrectly,
	Description: ``,
}

func applyClosingSquareBracketsMustBeSpacedCorrectly(source []byte) []byte {
	return scan(source, func(line, literal []byte) []byte {
		re := regexp.MustCompile(`([\S])([\t ]+)([\]])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1$3"))
		}

		re = regexp.MustCompile(`([\]])([\S])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1 $2"))
		}
		re = regexp.MustCompile(`([\]]) ([;])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1$2"))
		}

		return line
	})
}
