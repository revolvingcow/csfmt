package rules

import "regexp"

var preprocessorKeywordsMustNotBePrecededBySpace = &Rule{
	Name:        "Single line comments must begin with single space",
	Enabled:     true,
	Apply:       applyPreprocessorKeywordsMustNotBePrecededBySpace,
	Description: ``,
}

func applyPreprocessorKeywordsMustNotBePrecededBySpace(source []byte) []byte {
	keywords := `(if|else|elif|endif|define|undef|warning|error|line|region|endregion|pragma|pragma warning|pragma checksum)`
	re := regexp.MustCompile(`([#])(\t| )+` + keywords)
	for re.Match(source) {
		source = re.ReplaceAll(source, []byte("$1$3"))
	}
	return source
}
