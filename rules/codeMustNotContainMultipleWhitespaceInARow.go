package rules

import (
	"bytes"
	"regexp"
)

var codeMustNotContainMultipleWhitespaceInARow = &Rule{
	Name:        "Code must not contain multiple whitespaces in a row",
	Enabled:     true,
	Apply:       applyCodeMustNotContainMultipleWhitespaceInARow,
	Description: `A violation of this rule occurs whenver the code contains a tab character.`,
}

func applyCodeMustNotContainMultipleWhitespaceInARow(source []byte) []byte {
	return scan(source, func(line, literal []byte) []byte {
		re := regexp.MustCompile(`(\S)[ ]{2,}(\S)`)
		if !bytes.Contains(literal, re.Find(line)) {
			for re.Match(line) {
				line = re.ReplaceAll(line, []byte("$1 $2"))
			}
		}
		return line
	})
}
