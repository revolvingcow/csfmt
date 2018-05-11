package rules

import (
	"regexp"
)

var tabsMustNotBeUsed = &Rule{
	Name:        "Tabs must not be used",
	Enabled:     true,
	Apply:       applyTabsMustNotBeUsed,
	Description: `A violation of this rule occurs whenver the code contains a tab character.`,
}

func applyTabsMustNotBeUsed(source []byte) []byte {
	re := regexp.MustCompile(`\t`)
	for re.Match(source) {
		source = re.ReplaceAllLiteral(source, []byte("    "))
	}
	return source
}
