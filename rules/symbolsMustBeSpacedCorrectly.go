package rules

import (
	"bytes"
	"regexp"

	"github.com/revolvingcow/csfmt"
)

var symbolsMustBeSpacedCorrectly = &csfmt.Rule{
	Name:        "Symbols must be spaced correctly",
	Enabled:     false,
	Apply:       applySymbolsMustBeSpacedCorrectly,
	Description: ``,
}

func applySymbolsMustBeSpacedCorrectly(source []byte) []byte {
	return scan(source, func(line, literal []byte) []byte {
		// Look for pairings
		re := regexp.MustCompile(`([\w\)])([<>!\+\-\*\^%/\^=&\|\?]?[=\|&\?]|[<>\?\:])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1 $2"))
		}
		re = regexp.MustCompile(`([<>!\+\-\*\^%/\^=&\|\?]?[=\|&\?]|[<>\?\:])([\w!])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1 $2"))
		}

		// Incrementors and decrementors
		re = regexp.MustCompile(`([^\(])([\W])(\+\+|\-\-)(\w)`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1$2 $3$4"))
		}
		re = regexp.MustCompile(`(\w)(\+\+|\-\-)([^\)])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1$2 $3$4"))
		}

		// Unary operators
		re = regexp.MustCompile(`([\w])([!])([\w|\(])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1 $2$3"))
		}

		// Singlets
		re = regexp.MustCompile(`([\w\)])([\*/])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1 $2"))
		}
		re = regexp.MustCompile(`([\*/])([\w\(])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1 $2"))
		}

		re = regexp.MustCompile(`([^\+])([\+])([^\+=])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1 $2 $3"))
		}
		re = regexp.MustCompile(`([^\-])([\-])([^\-=])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1 $2 $3"))
		}

		// Fix negatives
		re = regexp.MustCompile(`([\+=<>\?])( *)([\-])([ ]+)([\d])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1 $3$5"))
		}

		// Fix generics
		re = regexp.MustCompile(`( < )(.*)( >\s*)\(`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("<$2>("))
		}
		re = regexp.MustCompile(`( < )(.*)( >\s*)(\w*)`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("<$2> $4"))
		}

		return line
	})
}
