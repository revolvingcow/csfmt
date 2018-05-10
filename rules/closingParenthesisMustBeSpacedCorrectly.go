package rules

import (
	"bytes"
	"regexp"
)

var closingParenthesisMustBeSpacedCorrectly = &Rule{
	Name:        "Closing parenthesis must be spaced correctly",
	Enabled:     true,
	Apply:       applyClosingParenthesisMustBeSpacedCorrectly,
	Description: ``,
}

func applyClosingParenthesisMustBeSpacedCorrectly(source []byte) []byte {
	return scan(source, func(line, literal []byte) []byte {
		spaceBetween := `([A-z]|=|\+|\-|\*|/|&|\||\^|\{)`

		// Remove leading spaces
		re := regexp.MustCompile(`([\S])(\t| )([\)])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1$3"))
		}

		// Remove trailing spaces
		re = regexp.MustCompile(`([\)])(\t| )([\S])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1$3"))
		}

		// Add space between operators and keywords
		re = regexp.MustCompile(`([\)])` + spaceBetween)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1 $2"))
		}

		return line
	})
}
