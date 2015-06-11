package rules

import "regexp"

var codeMustNotContainMultipleWhitespaceInARow = &Rule{
	Name:        "Code must not contain multiple whitespaces in a row",
	Enabled:     true,
	Apply:       applyCodeMustNotContainMultipleWhitespaceInARow,
	Description: `A violation of this rule occurs whenver the code contains a tab character.`,
}

func applyCodeMustNotContainMultipleWhitespaceInARow(source []byte) []byte {
	re := regexp.MustCompile(`(\S)[ ]{2,}(\S)`)
	for re.Match(source) {
		source = re.ReplaceAll(source, []byte("$1 $2"))
	}
	return source
}
