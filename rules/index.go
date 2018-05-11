package rules

var Library = []*Rule{
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
	colonsMustBeSpacedCorrectly,
	codeMustNotContainMultipleWhitespaceInARow,
	tabsMustNotBeUsed,
}

func Enabled() []*Rule {
	enabled := []*Rule{}
	for _, rule := range Library {
		if rule.Enabled {
			enabled = append(enabled, rule)
		}
	}
	return enabled
}
