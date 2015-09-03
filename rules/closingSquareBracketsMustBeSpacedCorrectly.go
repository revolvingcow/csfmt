package rules

import "regexp"

var closingSquareBracketsMustBeSpacedCorrectly = &Rule{
	Name:        "Closing square brackets must be spaced correctly",
	Enabled:     true,
	Apply:       applyClosingSquareBracketsMustBeSpacedCorrectly,
	Description: ``,
}

func applyClosingSquareBracketsMustBeSpacedCorrectly(source []byte) []byte {
	re := regexp.MustCompile(`([\S])([\t ]+)([\]])`)
	source = re.ReplaceAll(source, []byte("$1$3"))

	// re = regexp.MustCompile(`([\]])([\S])`)
	// source = re.ReplaceAll(source, []byte("$1 $2"))
	// re = regexp.MustCompile(`([\]]) ([;])`)
	// source = re.ReplaceAll(source, []byte("$1$2"))

	return source
}
