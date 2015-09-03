package rules

// The basic ruleset comes from [StyleCop][1] with them toggled on or off
// based solely on the authors supreme opinion(s).
//
// [1]: http://www.stylecop.com/docs/StyleCop%20Rules.html

// Rule is a style rule to look for and apply within the source code.
type Rule struct {
	Name        string
	Description string
	Enabled     bool
	Apply       func(source []byte) []byte
}

// Documentation rules
var Documentation = []*Rule{
// elementsMustBeDocument,
// partialElementsMustBeDocumented,
// enumerationItemsMustBeDocumented,
// documentationMustContainValidXml,
// elementDocumentationMustHaveSummary,
// elementDocumentationMustHaveSummaryText,
// elementDocumenationMustNotHaveDefaultSummary,
// elementDocumentationMustBeSpelledCorrectly,
// partialElementDocumentationMustHaveSummary,
// partialElementDocumentationMustHaveSummaryText,
// propertyDocumentationMustHaveValue,
// propertyDocumentationMustHaveValueText,
// elementParametersMustBeDocumented,
// elementParameterDocumentationMustMatchElementParameters,
// elementParameterDocumentationMustDeclareParameterName,
// elementParameterDocumentationMustHaveText,
// elementReturnValueMustBeDocumented,
// elementReturnValueDocumentationMustHaveValue,
// voidReturnValueMustNotBeDocumented,
// genericTypeParametersMustBeDocumented,
// genericTypeParametersMustBeDocumentedPartialClass,
// genericTypeParameterDocumentationMustMatchTypeParameters,
// genericTypeParameterDocumentationMustDeclareParameterName,
// genericTypeParameterDocumentationMustHaveText,
// propertySummaryDocumentationMustMatchAccessors,
// propertySummaryDocumentationMustOmitSetAccessorWithRestrictedAccess,
// elementDocumentationMustNotBeCopiedAndPasted,
// singleLineCommentsMustNotUseDocumentationStyleSlashes,
// documentationTextMustNotBeEmpty,
// documentationTextMustBeginWithACapitalLetter,
// documentationTextMustEndWithAPeriod,
// documentationTextMustContainWhitespace,
// documentationTextMustMeetCharacterPercentage,
// documentationTextMustMeetMinimumCharacterLenght,
// fileMustHaveHeader,
// fileHeaderMustShowCopyright,
// fileHeaderMustHaveCopyrightText,
// fileHeaderCopyrightTextMustMatch,
// fileHeaderMustContainFileName,
// fileHeaderFileNameDocumentationMustMatchFileName,
// fileHeaderMustHaveSummary,
// fileHeaderMustHaveValidCompanyText,
// fileHeaderCompanyNameTextMustMatch,
// fileHeaderFileNameDocumentationMustMatchTypeName,
// constructorSummaryDocumentationMustBeginWithStandardText,
// deconstructorSummaryDocumentationMustBeginWithStandardText,
// documentationHeadersMustNotContainBlankLines,
// includedDocumentationFileDoesNotExist,
// includedDocumentationXPathDoesNotExist,
// includedNodeDoesNotContainValidFileAndPath,
// inheritDocMustBeUsedWithInheritingClass,
}

// Layout rules
var Layout = []*Rule{
	// curlyBracketsForMultiLineStatementsMustNotShareLine,
	// statementMustNotBeOnSingleLine,
	// elementMustNotBeOnSingleLine,
	// curlyBracketsMustNotBeOmitted,
	// allAccessorMustBeMultiLineOrSingleLine,
	// openingCurlyBracketsMustNotBeFollowedByBlankLine,
	// elementDocumentationHeadersMustNotBeFollowedByBlankLine,
	codeMustNotContainMultipleBlankLinesInARow,
	// closingCurlyBracketsMustNotBePrecededByBlankLine,
	// openingCurlyBracketsMustNotBePrecededByBlankLine,
	// chainedStatementBlocksMustNotBePrecededByBlankLine,
	// whileDoFooterMustNotBePrecededByBlankLine,
	// singleLineCommentsMustNotBeFollowedByBlankLine,
	// closingCurlyBracketMustBeFollowedByBlankLine,
	// elementDocumentationHeaderMustBePrecededByBlankLine,
	// singleLineCommentMustBePrecededByBlankLine,
	// elementsMustBeSeparatedByBlankLine,
	// codeMustNotContainBlankLinesAtStartOfFile,
	// codeMustNotContainBlankLinesAtEndOfFile,
}

// Maintainability Rules
var Maintainability = []*Rule{
// statementMustNotUseUnnecessaryParenthesis,
// accessModifierMustBeDeclared,
// fieldsMustBePrivate,
// fileMayOnlyContainASingleClass,
// fileMayOnlyContainASingleNamespace,
// codeAnalysisSuppressionMustHaveJustification,
// debugAssertMustProvideMessageText,
// debugFailMustProvideMessageText,
// arithmeticExpressionsMustDeclarePrecedence,
// conditionalExpressionsMustDeclarePrecendence,
// removeUnnecessaryCode,
// removeDelegateParenthesisWhenPossible,
// attributeConstructorMustNotUseUnnecessaryParenthesis,
}

// Naming Rules
var Naming = []*Rule{
// elementMustBeginWithUpperCaseLetter,
// elementMustBeginWithLowerCaseLetter,
// interfaceNamesMustBeginWithI,
// constFieldNamesMustBeginWithUpperCaseLetter,
// nonPrivateReadonlyFieldsMustBeginWithUpperCaseLetter,
// fieldNamesMustNotUseHungarianNotation,
// fieldNamesMustBeginWithLowerCaseLetter,
// accessibleFieldsMustBeginWithUpperCaseLetter,
// variableNamesMustNotBePrefixed,
// fieldNamesMustNotBeginWithUnderscore,
// fieldNamesMustNotContainUnderscore,
// staticReadonlyFieldsMustBeginWithUpperCaseLetter,
}

// Ordering Rules
var Ordering = []*Rule{
	// usingDirectivesMustBePlacedWithinNamespace,
	// elementsMustAppearInTheCorrectOrder,
	// elementsMustBeOrderedByAccess,
	// constantsMustAppearBeforeFields,
	// staticElementsMustAppearBeforeInstanceElements,
	// partialElementsMustDeclareAccess,
	// declarationKeywordsMustFollowOrder,
	// protectedMustComeBeforeInternal,
	// systemUsingDirectivesMustBePlacedBeforeOtherUsingDirectives,
	// usingAliasDirectivesMustBePlacedAfterOtherUsingDirectives,
	usingDirectivesMustBeOrderedAlphabeticallyByNamespace,
	// usingAliasDirectivesMustBeOrderedAlphabeticallyByAliasName,
	// propertyAccessorsMustFollowOrder,
	// eventAccessorsMustFollowOrder,
	// staticReadonlyElementsMustAppearBeforeStaticNonReadonlyElements,
	// instanceReadonlyElementsMustAppearBeforeInstanceNonReadonlyElements,
}

// Readability Rules
var Readability = []*Rule{
// doNotPrefixCallsWithBaseUnlessLocalImplementationExists,
// prefixLocalCallsWithThis,
// queryClauseMustFollowPreviousClause,
// queryClausesMustBeOnSeparateLinesOrAllOnOneLine,
// queryClauseMustBeginOnNewLineWhenPreviousClauseSpansMultipleLines,
// queryClausesSpanningMultipleLinesMustBeginOnOwnLine,
// codeMustNotContainEmptyStatements,
// codeMustNotContainMultipleStatementsOnOneLine,
// blockStatementsMustNotContainEmbeddedComments,
// blockStatementsMustNotContainEmbeddedRegions,
// openingParenthesisMustBeOnDeclarationLine,
// closingParenthesisMustBeOnLineOfOpeningParenthesis,
// commaMustBeOnSameLineAsPreviousParameter,
// parameterListMustFollowDeclaration,
// parameterMustFollowComma,
// splitParametersMustStartOnLineAfterDeclaration,
// parametersMustBeOnSameLineOrSeparateLines,
// parametersMustNotSpanMultipleLines,
// commentsMustContainText,
// useBuiltInTypeAlias,
// useStringEmptyForEmptyStrings,
// doNotPlaceRegionsWithinElements,
// doNotUseRegions,
// useShorthandForNullableTypes,
// prefixCallsCorrectly,
}

// Spacing Rules
var Spacing = []*Rule{
	symbolsMustBeSpacedCorrectly,
	// keywordsMustBeSpacedCorrectly,
	commasMustBeSpacedCorrectly,
	semicolonsMustBeSpacedCorrectly,
	singleLineCommentsMustBeginWithSingleSpace,
	documentationLinesMustBeginWithSingleSpace,
	preprocessorKeywordsMustNotBePrecededBySpace,
	// operatorKeywordMustBeFollowedBySpace,
	openingParenthesisMustBeSpacedCorrectly,
	closingParenthesisMustBeSpacedCorrectly,
	openingSquareBracketsMustBeSpacedCorrectly,
	closingSquareBracketsMustBeSpacedCorrectly,
	// openingCurlyBracketsMustBeSpacedCorrectly,
	// closingCurlyBracketsMustBeSpacedCorrectly,
	// openingGenericBracketsMustBeSpacedCorrectly,
	// closingGenericBracketsMustBeSpacedCorrectly,
	// openingAttributeBracketsMustBeSpacedCorrectly,
	// closingAttributeBracketsMustBeSpacedCorrectly,
	// nullableTypeSymbolsMustNotBePrecededBySpace,
	// memberAccessSymbolsMustBeSpacedCorrectly,
	// incrementDecrementSymbolsMustBeSpacedCorrectly,
	// negativeSignsMustBeSpacedCorrectly,
	// positiveSignsMustBeSpacedCorrectly,
	// dereferenceAndAccessOfSymbolsMustBeSpacedCorrectly,
	colonsMustBeSpacedCorrectly,
	codeMustNotContainMultipleWhitespaceInARow,
	// codeMustNotContainSpaceAfterNewKeywordInImplicitlyTypedArrayAllocation,
	tabsMustNotBeUsed,
}

func EnabledRules() []*Rule {
	enabled := []*Rule{}

	for _, rule := range Documentation {
		if rule.Enabled {
			enabled = append(enabled, rule)
		}
	}
	for _, rule := range Layout {
		if rule.Enabled {
			enabled = append(enabled, rule)
		}
	}
	for _, rule := range Maintainability {
		if rule.Enabled {
			enabled = append(enabled, rule)
		}
	}
	for _, rule := range Ordering {
		if rule.Enabled {
			enabled = append(enabled, rule)
		}
	}
	for _, rule := range Readability {
		if rule.Enabled {
			enabled = append(enabled, rule)
		}
	}
	for _, rule := range Spacing {
		if rule.Enabled {
			enabled = append(enabled, rule)
		}
	}

	return enabled
}
