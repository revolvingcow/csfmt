package rules

import "regexp"

var singleLineCommentsMustBeginWithSingleSpace = &Rule{
	Name:        "Single line comments must begin with single space",
	Enabled:     true,
	Apply:       applySingleLineCommentsMustBeginWithSingleSpace,
	Description: ``,
}

func applySingleLineCommentsMustBeginWithSingleSpace(source []byte) []byte {
	re := regexp.MustCompile(`([/]{2})([\S])`)
	for re.Match(source) {
		source = re.ReplaceAll(source, []byte("$1 $2"))
	}
	re = regexp.MustCompile(`([/]{2})([\s]{2,})([\S])`)
	for re.Match(source) {
		source = re.ReplaceAll(source, []byte("$1 $3"))
	}
	re = regexp.MustCompile(`([/]{2})([\s])([/]{1,})([\s]*)`)
	for re.Match(source) {
		source = re.ReplaceAll(source, []byte("$1$3"))
	}
	return source
}
