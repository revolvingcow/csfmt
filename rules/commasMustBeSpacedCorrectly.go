package rules

import "regexp"

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

	// Add trailing spaces as necessary
	re = regexp.MustCompile(`(\S),(\w|\d)`)
	for re.Match(source) {
		source = re.ReplaceAll(source, []byte("$1, $2"))
	}

	// Look for too many trailing spaces
	re = regexp.MustCompile(`\,  `)
	for re.Match(source) {
		source = re.ReplaceAll(source, []byte(", "))
	}
	return source
}
