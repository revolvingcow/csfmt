package rules

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/revolvingcow/csfmt"
)

var usingDirectivesMustBeOrderedAlphabeticallyByNamespace = &csfmt.Rule{
	Name:        "Using directives must be ordered alphabetically by namespace",
	Enabled:     true,
	Apply:       applyUsingDirectivesMustBeOrderedAlphabeticallyByNamespace,
	Description: ``,
}

func applyUsingDirectivesMustBeOrderedAlphabeticallyByNamespace(source []byte) []byte {
	usings := []string{}
	source = scan(source, func(line, literal []byte) []byte {
		// Find usings
		re := regexp.MustCompile(`(.*)(using)([\t| ])([^\(])(.*)([;])`)
		if re.Match(line) {
			using := re.ReplaceAll(line, []byte("$2 $4$5"))
			usings = append(usings, string(using))
			line = []byte{}
		}

		return line
	})

	if len(usings) > 0 {
		// Sort the usings and add them to the top of the file
		s := sort.StringSlice(usings)
		if !sort.IsSorted(s) {
			sort.Sort(s)
		}

		source = append([]byte(fmt.Sprintf("%s;\n\n", strings.Join(s, ";\n"))), source...)
	}

	return source
}
