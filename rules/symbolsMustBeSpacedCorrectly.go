package rules

import (
	"bufio"
	"bytes"
	"regexp"
)

var symbolsMustBeSpacedCorrectly = &Rule{
	Name:        "Symbols must be spaced correctly",
	Enabled:     true,
	Apply:       applySymbolsMustBeSpacedCorrectly,
	Description: ``,
}

func applySymbolsMustBeSpacedCorrectly(source []byte) []byte {
	reCommentShortBegin := regexp.MustCompile(`\A\s*(/{2,})`)
	reCommentLongBegin := regexp.MustCompile(`\A\s*(/\*)`)
	reCommentLongEnd := regexp.MustCompile(`\*/`)

	short := false
	long := false
	lines := []byte{}
	buffer := bytes.NewBuffer(source)
	scanner := bufio.NewScanner(buffer)
	for scanner.Scan() {
		// Add a newline character on each line after the first
		if len(lines) > 0 {
			lines = append(lines, byte('\n'))
		}

		line := scanner.Bytes()
		short = reCommentShortBegin.Match(line)
		long = long || reCommentLongBegin.Match(line)

		if !(short || long) {
			// Look for pairings
			re := regexp.MustCompile(`([\w\)])([<>!\+\-\*\^%/\^=&\|\?]?[=\|&\?]|[<>\?\:])`)
			line = re.ReplaceAll(line, []byte("$1 $2"))
			re = regexp.MustCompile(`([<>!\+\-\*\^%/\^=&\|\?]?[=\|&\?]|[<>\?\:])([\w!])`)
			line = re.ReplaceAll(line, []byte("$1 $2"))

			// Incrementors and decrementors
			re = regexp.MustCompile(`([^\(])([\W])(\+\+|\-\-)(\w)`)
			line = re.ReplaceAll(line, []byte("$1$2 $3$4"))
			re = regexp.MustCompile(`(\w)(\+\+|\-\-)([^\)])`)
			line = re.ReplaceAll(line, []byte("$1$2 $3$4"))

			// Unary operators
			re = regexp.MustCompile(`([\w])([!])([\w|\(])`)
			line = re.ReplaceAll(line, []byte("$1 $2$3"))

			// Singlets
			re = regexp.MustCompile(`([\w\)])([\*/])`)
			line = re.ReplaceAll(line, []byte("$1 $2"))
			re = regexp.MustCompile(`([\*/])([\w\(])`)
			line = re.ReplaceAll(line, []byte("$1 $2"))
			re = regexp.MustCompile(`([^\+])([\+])([^\+=])`)
			line = re.ReplaceAll(line, []byte("$1 $2 $3"))
			re = regexp.MustCompile(`([^\-])([\-])([^\-=])`)
			line = re.ReplaceAll(line, []byte("$1 $2 $3"))

			// Fix negatives
			re = regexp.MustCompile(`([\+=<>\?])([ ])([\-])([ ])([\d])`)
			line = re.ReplaceAll(line, []byte("$1 $3$5"))

			// Fix generics
			re = regexp.MustCompile(`( < )(.*)( >\s*)`)
			line = re.ReplaceAll(line, []byte("<$2>"))
		}

		short = false
		long = long && !reCommentLongEnd.Match(line)
		lines = append(lines, line...)
	}

	return lines
}

// func applySymbolsMustBeSpacedCorrectly(source []byte) []byte {
// 	alphaNumeric := `(\S)`
// 	symbols := `([|=|\+|\-|\*|/|\||&|\^|<|>|!|])`
// 	unary := `([!])`
//
// 	reCommentShortBegin := regexp.MustCompile(`\A\s*(/{2,})`)
// 	reCommentLongBegin := regexp.MustCompile(`\A\s*(/\*)`)
// 	reCommentLongEnd := regexp.MustCompile(`\*/`)
//
// 	short := false
// 	long := false
// 	lines := []byte{}
// 	buffer := bytes.NewBuffer(source)
// 	scanner := bufio.NewScanner(buffer)
// 	for scanner.Scan() {
// 		// Add a newline character on each line after the first
// 		if len(lines) > 0 {
// 			lines = append(lines, byte('\n'))
// 		}
//
// 		line := scanner.Bytes()
//
// 		short = reCommentShortBegin.Match(line)
// 		long = long || reCommentLongBegin.Match(line)
//
// 		if !(short || long) {
// 			// Break apart all the symbols
// 			re := regexp.MustCompile(alphaNumeric + symbols)
// 			for re.Match(line) {
// 				line = re.ReplaceAll(line, []byte("$1 $2"))
// 			}
// 			re = regexp.MustCompile(symbols + alphaNumeric)
// 			for re.Match(line) {
// 				line = re.ReplaceAll(line, []byte("$1 $2"))
// 			}
//
// 			// Fix basic spacing
// 			re = regexp.MustCompile(symbols + `([\s])` + symbols)
// 			for re.Match(line) {
// 				line = re.ReplaceAll(line, []byte("$1$3"))
// 			}
//
// 			// Fix unary operators
// 			re = regexp.MustCompile(`([\S])([\s]*)` + unary + `([\s])([\S])`)
// 			for re.Match(line) {
// 				line = re.ReplaceAll(line, []byte("$1 $3$5"))
// 			}
//
// 			// Fix generics
// 			re = regexp.MustCompile(`( < )(.*)( >\s*)`)
// 			for re.Match(line) {
// 				line = re.ReplaceAll(line, []byte("<$2>"))
// 			}
//
// 			// Fix incrementer/decrementer
// 			re = regexp.MustCompile(`(--|\+\+)` + symbols)
// 			for re.Match(line) {
// 				line = re.ReplaceAll(line, []byte("$1 $2"))
// 			}
// 			re = regexp.MustCompile(symbols + `(--|\+\+)`)
// 			for re.Match(line) {
// 				line = re.ReplaceAll(line, []byte("$1 $2"))
// 			}
// 			re = regexp.MustCompile(`([A-z]|[\d]|[\(])([\s])(--|\+\+)`)
// 			for re.Match(line) {
// 				line = re.ReplaceAll(line, []byte("$1$3"))
// 			}
// 			re = regexp.MustCompile(`(--|\+\+)([\s])([A-z]|[\d]|[\)])`)
// 			for re.Match(line) {
// 				line = re.ReplaceAll(line, []byte("$1$3"))
// 			}
//
// 			lines = append(lines, line...)
// 		} else {
// 			lines = append(lines, line...)
// 		}
//
// 		short = false
// 		long = long && !reCommentLongEnd.Match(line)
// 	}
//
// 	return lines
// }
