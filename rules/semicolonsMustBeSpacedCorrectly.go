package rules

import "regexp"

var semicolonsMustBeSpacedCorrectly = &Rule{
	Name:        "Semicolons must be spaced correctly",
	Enabled:     true,
	Apply:       applySemicolonsMustBeSpacedCorrectly,
	Description: ``,
}

func applySemicolonsMustBeSpacedCorrectly(source []byte) []byte {
	// Look for leading spaces
	re := regexp.MustCompile(`\s;`)
	for re.Match(source) {
		source = re.ReplaceAllLiteral(source, []byte(";"))
	}

	// Add trailing spaces as necessary
	re = regexp.MustCompile(`;(\S)`)
	for re.Match(source) {
		source = re.ReplaceAll(source, []byte("; $1"))
	}
	return source
}
