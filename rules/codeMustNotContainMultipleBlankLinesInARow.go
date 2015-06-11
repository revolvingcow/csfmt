package rules

import "regexp"

var codeMustNotContainMultipleBlankLinesInARow = &Rule{
	Name:        "Code must not contain multiple blank lines in a row",
	Enabled:     true,
	Apply:       applyCodeMustNotContainMultipleBlankLinesInARow,
	Description: ``,
}

func applyCodeMustNotContainMultipleBlankLinesInARow(source []byte) []byte {
	re := regexp.MustCompile("(\n\n)+")
	for re.Match(source) {
		source = re.ReplaceAllLiteral(source, []byte("\n"))
	}
	return source
}
