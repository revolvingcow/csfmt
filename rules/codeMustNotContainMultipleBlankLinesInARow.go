package rules

import (
	"regexp"

	"github.com/revolvingcow/csfmt"
)

var codeMustNotContainMultipleBlankLinesInARow = &csfmt.Rule{
	Name:        "Code must not contain multiple blank lines in a row",
	Enabled:     true,
	Apply:       applyCodeMustNotContainMultipleBlankLinesInARow,
	Description: ``,
}

func applyCodeMustNotContainMultipleBlankLinesInARow(source []byte) []byte {
	re := regexp.MustCompile("\n{3,}")
	for re.Match(source) {
		source = re.ReplaceAllLiteral(source, []byte("\n"))
	}
	return source
}
