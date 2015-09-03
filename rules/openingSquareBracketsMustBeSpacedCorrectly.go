package rules

import "regexp"

var openingSquareBracketsMustBeSpacedCorrectly = &Rule{
	Name:        "Opening square brackets must be spaced correctly",
	Enabled:     true,
	Apply:       applyOpeningSquareBracketsMustBeSpacedCorrectly,
	Description: ``,
}

func applyOpeningSquareBracketsMustBeSpacedCorrectly(source []byte) []byte {
	re := regexp.MustCompile(`([\S])([\t ]+)([\[])`)
	source = re.ReplaceAll(source, []byte("$1$3"))
	re = regexp.MustCompile(`([\[])([\t ]+)([\S])`)
	source = re.ReplaceAll(source, []byte("$1$3"))

	return source
}
