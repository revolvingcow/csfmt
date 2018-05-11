package rules

import (
	"bytes"
	"regexp"

	"github.com/revolvingcow/csfmt"
)

var preprocessorKeywordsMustNotBePrecededBySpace = &csfmt.Rule{
	Name:        "Single line comments must begin with single space",
	Enabled:     true,
	Apply:       applyPreprocessorKeywordsMustNotBePrecededBySpace,
	Description: ``,
}

func applyPreprocessorKeywordsMustNotBePrecededBySpace(source []byte) []byte {
	keywords := `(if|else|elif|endif|define|undef|warning|error|line|region|endregion|pragma|pragma warning|pragma checksum)`
	re := regexp.MustCompile(`([#])(\t| )+` + keywords)

	return scan(source, func(line, literal []byte) []byte {
		for re.Match(line) {
			if !bytes.Contains(literal, re.Find(line)) {
				line = re.ReplaceAll(line, []byte("$1$3"))
			}
		}
		return line
	})
}
