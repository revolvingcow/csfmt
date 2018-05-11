package rules

import (
	"github.com/revolvingcow/csfmt"
)

var Library = []*csfmt.Rule{
	codeMustNotContainMultipleBlankLinesInARow,
	usingDirectivesMustBeOrderedAlphabeticallyByNamespace,
	symbolsMustBeSpacedCorrectly,
	commasMustBeSpacedCorrectly,
	semicolonsMustBeSpacedCorrectly,
	singleLineCommentsMustBeginWithSingleSpace,
	documentationLinesMustBeginWithSingleSpace,
	preprocessorKeywordsMustNotBePrecededBySpace,
	openingParenthesisMustBeSpacedCorrectly,
	closingParenthesisMustBeSpacedCorrectly,
	openingSquareBracketsMustBeSpacedCorrectly,
	closingSquareBracketsMustBeSpacedCorrectly,
	codeMustNotContainMultipleWhitespaceInARow,
	tabsMustNotBeUsed,
}

func Enabled() []*csfmt.Rule {
	enabled := []*csfmt.Rule{}
	for _, rule := range Library {
		if rule.Enabled {
			enabled = append(enabled, rule)
		}
	}
	return enabled
}
