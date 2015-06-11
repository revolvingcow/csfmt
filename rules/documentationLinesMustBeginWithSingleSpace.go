package rules

import "regexp"

var documentationLinesMustBeginWithSingleSpace = &Rule{
	Name:        "Documentation lines must begin with a single space",
	Enabled:     true,
	Apply:       applyDocumentationLinesMustBeginWithSingleSpace,
	Description: ``,
}

func applyDocumentationLinesMustBeginWithSingleSpace(source []byte) []byte {
	re := regexp.MustCompile(`([/]{3})(\S)`)
	for re.Match(source) {
		source = re.ReplaceAll(source, []byte("$1 $2"))
	}
	return source
}
