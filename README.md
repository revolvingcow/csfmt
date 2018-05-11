# csfmt

A C# formatting tool from the command line.

This program was converted to **literate programming** with the help of [lmt](https://github.com/driusan/lmt) by **Dave MacFarlane**.

## How to get it

Go ahead and pull the latest from Github

``` shell
go get -u github.com/revolvingcow/csfmt
```

To generate the latest code from this file run

``` shell
lmt README.md
```

and then build normally with Go

``` shell
go build
```

## What problem is being solved?

## Finding files to apply rule sets

From a high level view we are looking for this basic workflow:

 1. Parse command line arguments driving where to look for files
 2. Have some basic statistics to give a gist of what was found/applied
 3. Gather all rules we need to apply
 4. Loop through all found files and apply the rule set

``` go main.go
package main

import (
	"flag"
	<<<main.go imports>>>
)

var (
	<<<main.go vars>>>
)

func main() {
	flag.Parse()
	sourceFiles := []SourceFile{}

	<<<handle command arguments>>>
	<<<setup statistics>>>
	<<<get rules>>>
	<<<apply rules to source files>>>
	<<<output statistics>>>
}
```

### What is a source file?

A source file will be considered any file container source code which would apply to the rule sets. In
our case this typically this will be files with the C# extension of `.cs`.

``` go source.go
package main

import (
	<<<source file imports>>>
)

var (
	<<<source file vars>>>
)

// SourceFile represents a file declared as source code.
type SourceFile struct {
	<<<source file fields>>>
}

<<<source file methods>>>
```

Let us declare the file extensions we are actually looking for a enumerable variable.

``` go "source file vars"
extensions = []string{".cs"}
```

Now we know each source file will have a file path associated to it so we can read them in to memory when needed.
We will declare this string and name it `Path`.

``` go "source file fields"
Path string
```

There are some common properties we will have to check for across all source files. For these we will create methods
off of the `SourceFile` structure. Due to sanity we need to check if the thing even exists. This can return a simple Boolean
value without any need to pass back an error if one is raised.

``` go "source file methods"
func (f *SourceFile) Exists() bool {
	if _, err := os.Stat(f.Path); os.IsNotExist(err) {
		return false
	}
	return true
}
```

and thus need to include the `os` package

``` go "source file imports"
"os"
```

We will also need to know if the file path is actually a directory. This will allow us to

 1. Not create a source file based on a directory
 2. Recursively find source files within the directory if needed
 
``` go "source file methods" +=
func (f *SourceFile) IsDir() bool {
	fi, err := os.Stat(f.Path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}
```

Assuming we have a file now we need to ensure it has an acceptable extension

``` go "source file methods" +=
func (f *SourceFile) IsDotNet() bool {
	name := strings.ToLower(f.Path)
	for _, extension := range extensions {
		if strings.HasSuffix(name, extension) {
			return true
		}
	}
	return false
}
```

which brings in the `strings` package

``` go "source file imports" +=
"strings"
```

To make life a bit easier we'll create some common file I/O methods to handle the file contents.

``` go "source file methods" +=
// Read the file contents.
func (f *SourceFile) Read() ([]byte, error) {
	contents, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

// Write the contents out to a stream.
func (f *SourceFile) Write(contents []byte) error {
	fi, err := os.OpenFile(f.Path, os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer fi.Close()

	_, err = fi.Write(contents)
	if err != nil {
		return err
	}
	return nil
}
```

and also include the `ioutil` package

``` go "source file imports" +=
"io/ioutil"
```
 
To walk directories we need to bring in another package

``` go "source file imports" +=
"path/filepath"
```

Now we can recursively go through directory structures on the file system and only return
the file paths meeting the approved extensions.

``` go "source file methods" +=
// Walk a directory's file structure looking for source files.
func (f *SourceFile) Walk() chan SourceFile {
	c := make(chan SourceFile)

	go func() {
		filepath.Walk(f.Path, func(p string, fi os.FileInfo, e error) error {
			if e != nil {
				return e
			}

			s := SourceFile{
				Path: p,
			}

			if s.IsDotNet() {
				c <- s
			}

			return nil
		})

		defer close(c)
	}()

	return c
}
```

### What is a rule?

``` go rules/rules.go
package rules

// Rule is a style rule to look for and apply within the source code.
type Rule struct {
	Name        string
	Description string
	Enabled     bool
	Apply       func(source []byte) []byte
}
```

### Apply the basic structures to our workflow

``` go "handle command arguments"
// Determine what files to format
argc := len(os.Args)
if argc < 2 {
	return
} else if argc == 2 && os.Args[1] == "..." {
	// Walk the file structure from the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	s := SourceFile{
		Path: cwd,
	}

	for sourceFile := range s.Walk() {
		sourceFiles = append(sourceFiles, sourceFile)
	}
} else {
	// Assuming multiple files were given
	for _, a := range os.Args[1:] {
		s := SourceFile{
			Path: a,
		}

		if s.Exists() {
			if s.IsDir() {
				for sourceFile := range s.Walk() {
					sourceFiles = append(sourceFiles, sourceFile)
				}
			} else if s.IsDotNet() {
				sourceFiles = append(sourceFiles, s)
			}
		}
	}
}
```

and now we bring in the `os` standard package

``` go "main.go imports" +=
"os"
```

We don't need anything too fancy for statistics so maybe just some totals.

``` go "setup statistics"
count := len(sourceFiles)
modified := 0
```

We only want to apply rules which have explicitly been enabled. This allows us
to turn off specific rules when there are known issues or they are still being
actively hacked on.

``` go "get rules"
queuedRules := rules.Enabled()
```

and our rules are found in a sub-package just for organization purposes

``` go "main.go imports" +=
"github.com/revolvingcow/csfmt/rules"
```

``` go "apply rules to source files"
for _, s := range sourceFiles {
	contents, err := s.Read()
	if err != nil {
		log.Fatalln(err)
	}
	original := contents

	for _, rule := range queuedRules {
		contents = rule.Apply(contents)
	}

	if bytes.Compare(original, contents) != 0 {
		modified++
		<<<main.go file has been modified>>>
	}

	<<<main.go final processing of file>>>
}
```

looks like we'll need to import another package since we need to compare to byte arrays

``` go "main.go imports" +=
"bytes"
```

Just a simple message will do detailing how many files were processed, how many were modified,
and the number of rules applied.

``` go "output statistics"
log.Printf("Modified %d of %d files using %d rules\n", modified, count, len(queuedRules))
```

and now include the required package

``` go "main.go imports" +=
"log"
```

### Allow writing changes to source file(s)

Provide a command line flag so the user may determine when it is appropriate to
make potentially destructive changes. This will default to **false** to ensure
a conscience decision has been made to possibly overwrite file contents.

``` go "main.go vars"
flagWrite = flag.Bool("w", false, "write changes to file")
```

We'll probably need to output something so we'll prepare for it

``` go "main.go imports" +=
"fmt"
```

If we have been told to write any modifications to the file then we will only
write to the file when a modification has been detected

``` go "main.go file has been modified"
if *flagWrite {
	s.Write(contents)
}
```

However, if we are not writing to the file then we'll output the file contents
regardless of modification to standard output

``` go "main.go final processing of file"
if !*flagWrite {
	fmt.Println(string(contents))
}
```

## Rules

The basic rule set comes from [StyleCop]([h](https://github.com/StyleCop/StyleCop/tree/master/Project/Docs/Rules/StyleCop%20Rules.html)ttp://www.stylecop.com/docs/StyleCop%20Rules.html) with them toggled on or off
based solely on the author(s) supreme opinion(s).

### Index

The basic **index** should only consist of rules currently finished or being worked
on. For tracking of rules yet to be applied we can keep track of them within this
document thus reducing dead code.

``` go rules/index.go
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
```

Now we need a way to filter out any rules which are not enabled.

``` go rules/index.go +=
func Enabled() []*Rule {
	enabled := []*Rule{}
	for _, rule := range Library {
		if rule.Enabled {
			enabled = append(enabled, rule)
		}
	}
	return enabled
}
```

### Scanning files line-by-line

Several of the rules will need to go line-by-line through the file while checking for
discrepancies. Common caveats may be:

 1. Ignore commented lines
 2. Ignore literal strings
 
In order to apply some form of the DRY principle we can declare a function to make
this a bit less painful:

``` go rules/scan.go
package rules

import (
	<<<scan imports>>>
)

<<<scan methods>>>
```

First, let us go ahead and bring in some packages we know we be used.

 - `bytes` we will be needed to create a buffer from file byte array at a minimum
 - `bufio` will allow us to create a new scanner from a buffer
 - `regexp` because more than likely some pattern matching will be necessary

``` go "scan imports"
"bufio"
"bytes"
"regexp"
```

The `scan` method is intended only as a helper function within the `rules` package so
we will leave it not exported. The main intentions of the function are:

 1. Take in the byte array of the source file
 2. Go through each line of the file
 3. Only allow modification of lines which are not commented
 4. Call an anonymous function to apply further processing
 5. Return the contents as a byte array

``` go "scan methods"
func scan(source []byte, applyFunc func(line, literal []byte) []byte) []byte {
	reCommentShortBegin := regexp.MustCompile(`\A\s*(/{2,})`)
	reCommentLongBegin := regexp.MustCompile(`\A\s*(/\*)`)
	reCommentLongEnd := regexp.MustCompile(`\*/`)
	reString := regexp.MustCompile(`".*"`)

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

		literal := []byte{}
		if reString.Match(line) {
			literal = reString.Find(line)
		}

		if !(short || long) {
			line = applyFunc(line, literal)

			// Trim ending
			line = bytes.TrimRightFunc(line, unicode.IsSpace)
		}

		short = false
		long = long && !reCommentLongEnd.Match(line)
		lines = append(lines, line...)
	}

	return lines
}
```

Oops, looks like we missed an import for a package!

``` go "scan imports" +=
"unicode"
```

### Checklist

#### Documentation

 - [ ] SA1600: ElementsMustBeDocumented
 - [ ] SA1601: PartialElementsMustBeDocumented
 - [ ] SA1602: EnumerationItemsMustBeDocumented
 - [ ] SA1603: DocumentationMustContainValidXml
 - [ ] SA1604: ElementDocumentationMustHaveSummary
 - [ ] SA1605: PartialElementDocumentationMustHaveSummary
 - [ ] SA1606: ElementDocumentationMustHaveSummaryText
 - [ ] SA1607: PartialElementDocumentationMustHaveSummaryText
 - [ ] SA1608: ElementDocumentationMustNotHaveDefaultSummary
 - [ ] SA1609: PropertyDocumentationMustHaveValue
 - [ ] SA1610: PropertyDocumentationMustHaveValueText
 - [ ] SA1611: ElementParametersMustBeDocumented
 - [ ] SA1612: ElementParameterDocumentationMustMatchElementParameters
 - [ ] SA1613: ElementParameterDocumentationMustDeclareParameterName
 - [ ] SA1614: ElementParameterDocumentationMustHaveText
 - [ ] SA1615: ElementReturnValueMustBeDocumented
 - [ ] SA1616: ElementReturnValueDocumentationMustHaveValue
 - [ ] SA1617: VoidReturnValueMustNotBeDocumented
 - [ ] SA1618: GenericTypeParametersMustBeDocumented
 - [ ] SA1619: GenericTypeParametersMustBeDocumentedPartialClass
 - [ ] SA1620: GenericTypeParameterDocumentationMustMatchTypeParameters
 - [ ] SA1621: GenericTypeParameterDocumentationMustDeclareParameterName
 - [ ] SA1622: GenericTypeParameterDocumentationMustHaveText
 - [ ] SA1623: PropertySummaryDocumentationMustMatchAccessors
 - [ ] SA1624: PropertySummaryDocumentationMustOmitSetAccessorWithRestricedAccess
 - [ ] SA1625: ElementDocumentationMustNotBeCopiedAndPasted
 - [ ] SA1626: SingleLineCommentsMustNotUseDocumentationStyleSlashes
 - [ ] SA1627: DocumentationTextMustNotBeEmpty
 - [ ] SA1628: DocumentationTextMustBeginWithACapitalLetter
 - [ ] SA1629: DocumentationTextMustEndWithAPeriod
 - [ ] SA1630: DocumentationTextMustContainWhitespace
 - [ ] SA1631: DocumentationTextMustMeetCharacterPercentage
 - [ ] SA1632: DocumentationTextMustMeetMinimumCharacterLength
 - [ ] SA1633: FileMustHaveHeader
 - [ ] SA1634: FileHeaderMustShowCopyright
 - [ ] SA1635: FileHeaderMustHaveCopyrightText
 - [ ] SA1636: FileHeaderCopyrightTextMustMatch
 - [ ] SA1637: FileHeaderMustContainFileName
 - [ ] SA1638: FileHeaderFileNameDocumentationMustMatchFileName
 - [ ] SA1639: FileHeaderMustHaveSummary
 - [ ] SA1640: FileHeaderMustHaveValidCompanyText
 - [ ] SA1641: FileHeaderCompanyNameTextMustMatch
 - [ ] SA1642: ConstructorSummaryDocumentationMustBeginWithStandardText
 - [ ] SA1643: DestructorSummaryDocumentationMustBeginWithStandardText
 - [ ] SA1644: DocumentationHeadersMustNotContainBlankLines
 - [ ] SA1645: IncludedDocumentationFileDoesNotExist
 - [ ] SA1646: IncludedDocumentationXPathDoesNotExist
 - [ ] SA1647: IncludeNodeDoesNotContainValidFileAndPath
 - [ ] SA1648: InheritDocMustBeUsedWithInheritingClass
 - [ ] SA1649: FileHeaderFileNameDocumentationMustMatchTypeName
 - [ ] SA1650: ElementDocumentationMustBeSpelledCorrectly

#### Layout

 - [ ] SA1500: CurlyBracketsForMultiLineStatementsMustNotShareLine</A></P>
 - [ ] SA1501: StatementMustNotBeOnSingleLine</A></P>
 - [ ] SA1502: ElementMustNotBeOnSingleLine</A></P>
 - [ ] SA1503: CurlyBracketsMustNotBeOmitted</A></P>
 - [ ] SA1504: AllAccessorMustBeMultiLineOrSingleLine</A></P>
 - [ ] SA1505: OpeningCurlyBracketsMustNotBeFollowedByBlankLine</A></P>
 - [ ] SA1506: ElementDocumentationHeadersMustNotBeFollowedByBlankLine</A></P>
 - [x] SA1507: CodeMustNotContainMultipleBlankLinesInARow</A></P>
 - [ ] SA1508: ClosingCurlyBracketsMustNotBePrecededByBlankLine</A></P>
 - [ ] SA1509: OpeningCurlyBracketsMustNotBePrecedededByBlankLine</A></P>
 - [ ] SA1510: ChainedStatementBlocksMustNotBePrecededByBlankLine</A></P>
 - [ ] SA1511: WhileDoFooterMustNotBePrecededByBlankLine</A></P>
 - [ ] SA1512: SingleLineCommentsMustNotBeFollowedByBlankLine</A></P>
 - [ ] SA1513: ClosingCurlyBracketMustBeFollowedByBlankLine</A></P>
 - [ ] SA1514: ElementDocumentationHeaderMustBePrecededByBlankLine</A></P>
 - [ ] SA1515: SingleLineCommentMustBePrecededByBlankLine</A></P>
 - [ ] SA1516: ElementsMustBeSeparatedByBlankLine</A></P>
 - [ ] SA1517: CodeMustNotContainBlankLinesAtStartOfFile</A></P>
 - [ ] SA1518: CodeMustNotContainBlankLinesAtEndOfFile</A></P>

#### Maintainability

 - [ ] SA1119: StatementMustNotUseUnnecessaryParenthesis
 - [ ] SA1400: AccessModifierMustBeDeclared
 - [ ] SA1401: FieldsMustBePrivate
 - [ ] SA1402: FileMayOnlyContainASingleClass
 - [ ] SA1403: FileMayOnlyContainASingleNamespace
 - [ ] SA1404: CodeAnalysisSuppressionMustHaveJustification
 - [ ] SA1405: DebugAssertMustProvideMessageText
 - [ ] SA1406: DebugFailMustProvideMessageText
 - [ ] SA1407: ArithmeticExpressionsMustDeclarePrecedence
 - [ ] SA1408: ConditionalExpressionsMustDeclarePrecendence
 - [ ] SA1409: RemoveUnnecessaryCode
 - [ ] SA1410: RemoveDelegateParenthesisWhenPossible
 - [ ] SA1411: AttributeConstructorMustNotUseUnnecessaryParenthesis
 
#### Naming

 - [ ] SA1300: ElementMustBeginWithUpperCaseLetter
 - [ ] SA1301: ElementMustBeginWithLowerCaseLetter
 - [ ] SA1302: InterfaceNamesMustBeginWithI
 - [ ] SA1303: ConstFieldNamesMustBeginWithUpperCaseLetter
 - [ ] SA1304: NonPrivateReadonlyFieldsMustBeginWithUpperCaseLetter
 - [ ] SA1305: FieldNamesMustNotUseHungarianNotation
 - [ ] SA1306: FieldNamesMustBeginWithLowerCaseLetter
 - [ ] SA1307: AccessibleFieldsMustBeginWithUpperCaseLetter
 - [ ] SA1308: VariableNamesMustNotBePrefixed
 - [ ] SA1309: FieldNamesMustNotBeginWithUnderscore
 - [ ] SA1310: FieldNamesMustNotContainUnderscore
 - [ ] SA1311: StaticReadonlyFieldsMustBeginWithUpperCaseLetter
 
#### Ordering
 
 - [ ] SA1200: UsingDirectivesMustBePlacedWithinNamespace
 - [ ] SA1201: ElementsMustAppearInTheCorrectOrder
 - [ ] SA1202: ElementsMustBeOrderedByAccess
 - [ ] SA1203: ConstantsMustAppearBeforeFields
 - [ ] SA1204: StaticElementsMustAppearBeforeInstanceElements
 - [ ] SA1205: PartialElementsMustDeclareAccess
 - [ ] SA1206: DeclarationKeywordsMustFollowOrder
 - [ ] SA1207: ProtectedMustComeBeforeInternal
 - [ ] SA1208: SystemUsingDirectivesMustBePlacedBeforeOtherUsingDirectives
 - [ ] SA1209: UsingAliasDirectivesMustBePlacedAfterOtherUsingDirectives
 - [x] SA1210: UsingDirectivesMustBeOrderedAlphabeticallyByNamespace
 - [ ] SA1211: UsingAliasDirectivesMustBeOrderedAlphabeticallyByAliasName
 - [ ] SA1212: PropertyAccessorsMustFollowOrder
 - [ ] SA1213: EventAccessorsMustFollowOrder
 - [ ] SA1214: StaticReadonlyElementsMustAppearBeforeStaticNonReadonlyElements
 - [ ] SA1215: InstanceReadonlyElementsMustAppearBeforeInstanceNonReadonlyElements
 - [ ] SA1216: NoValueFirstComparison

#### Readability

 - [ ] SA1100: DoNotPrefixCallsWithBaseUnlessLocalImplementationExists
 - [ ] SA1101: PrefixLocalCallsWithThis
 - [ ] SA1102: QueryClauseMustFollowPreviousClause
 - [ ] SA1103: QueryClausesMustBeOnSeparateLinesOrAllOnOneLine
 - [ ] SA1104: QueryClauseMustBeginOnNewLineWhenPreviousClauseSpansMultipleLines
 - [ ] SA1105: QueryClausesSpanningMultipleLinesMustBeginOnOwnLine
 - [ ] SA1106: CodeMustNotContainEmptyStatements
 - [ ] SA1107: CodeMustNotContainMultipleStatementsOnOneLine
 - [ ] SA1108: BlockStatementsMustNotContainEmbeddedComments
 - [ ] SA1109: BlockStatementsMustNotContainEmbeddedRegions
 - [ ] SA1110: OpeningParenthesisMustBeOnDeclarationLine
 - [ ] SA1111: ClosingParenthesisMustBeOnLineOfOpeningParenthesis
 - [ ] SA1112: ClosingParenthesisMustBeOnLineOfOpeningParenthesis
 - [ ] SA1113: CommaMustBeOnSameLineAsPreviousParameter
 - [ ] SA1114: ParameterListMustFollowDeclaration
 - [ ] SA1115: ParameterMustFollowComma
 - [ ] SA1116: SplitParametersMustStartOnLineAfterDeclaration
 - [ ] SA1117: ParametersMustBeOnSameLineOrSeparateLines
 - [ ] SA1118: ParameterMustNotSpanMultipleLines
 - [ ] SA1120: CommentsMustContainText
 - [ ] SA1121: UseBuiltInTypeAlias
 - [ ] SA1122: UseStringEmptyForEmptyStrings
 - [ ] SA1123: DoNotPlaceRegionsWithinElements
 - [ ] SA1124: DoNotUseRegions
 - [ ] SA1125: UseShorthandForNullableTypes
 - [ ] SA1126: PrefixCallsCorrectly

#### Spacing

 - [ ] SA1000: KeywordsMustBeSpacedCorrectly
 - [ ] SA1001: CommasMustBeSpacedCorrectly
 - [x] SA1002: SemicolonsMustBeSpacedCorrectly
 - [x] SA1003: SymbolsMustBeSpacedCorrectly
 - [x] SA1004: DocumentationLinesMustBeginWithSingleSpace
 - [x] SA1005: SingleLineCommentsMustBeginWithSingeSpace
 - [x] SA1006: PreprocessorKeywordsMustNotBePrecededBySpace
 - [ ] SA1007: OperatorKeywordMustBeFollowedBySpace
 - [x] SA1008: OpeningParenthesisMustBeSpacedCorrectly
 - [x] SA1009: ClosingParenthesisMustBeSpacedCorrectly
 - [x] SA1010: OpeningSquareBracketsMustBeSpacedCorrectly
 - [x] SA1011: ClosingSquareBracketsMustBeSpacedCorrectly
 - [ ] SA1012: OpeningCurlyBracketsMustBeSpacedCorrectly
 - [ ] SA1013: ClosingCurlyBracketsMustBeSpacedCorrectly
 - [ ] SA1014: OpeningGenericBracketsMustBeSpacedCorrectly
 - [ ] SA1015: ClosingGenericBracketsMustBeSpacedCorrectly
 - [ ] SA1016: OpeningAttributeBracketsMustBeSpacedCorrectly
 - [ ] SA1017: ClosingAttributeBracketsMustBeSpacedCorrectly
 - [ ] SA1018: NullableTypeSymbolsMustNotBePrecededBySpace
 - [ ] SA1019: MemberAccessSymbolsMustBeSpacedCorrectly
 - [ ] SA1020: IncrementDecrementSymbolsMustBeSpacedCorrectly
 - [ ] SA1021: NegativeSignsMustBeSpacedCorrectly
 - [ ] SA1022: PositiveSignsMustBeSpacedCorrectly
 - [ ] SA1023: DereferenceAndAccessOfSymbolsMustBeSpacedCorrectly
 - [x] SA1024: ColonsMustBeSpacedCorrectly
 - [x] SA1025: CodeMustNotContainMultipleWhitespaceInARow
 - [ ] SA1026: CodeMustNotContainSpaceAfterNewKeywordInImplicitlyTypedArrayAllocation
 - [x] SA1027: TabsMustNotBeUsed

### SA1002: Semicolons must be spaced correctly

First the template

``` go rules/semicolonsMustBeSpacedCorrectly.go
package rules

import (
	<<<sa1002 imports>>>
)

<<<sa1002 rule>>>
<<<sa1002 application>>>
```

``` go "sa1002 application"
func applySemicolonsMustBeSpacedCorrectly(source []byte) []byte {
	return scan(source, func(line, literal []byte) []byte {
		// Look for leading spaces
		re := regexp.MustCompile(`\s;`)
		for re.Match(line) {
			if !bytes.Contains(literal, re.Find(line)) {
				line = re.ReplaceAllLiteral(line, []byte(";"))
			}
		}

		// Add trailing spaces as necessary
		re = regexp.MustCompile(`;(\S)`)
		for re.Match(line) {
			if !bytes.Contains(literal, re.Find(line)) {
				line = re.ReplaceAll(line, []byte("; $1"))
			}
		}
		return line
	})
}
```

Bring in used packages

``` go "sa1002 imports"
"bytes"
"regexp"
```

Now the logic has been worked out we'll apply create the rule.

``` go "sa1002 rule"
var semicolonsMustBeSpacedCorrectly = &Rule{
	Name:        "Semicolons must be spaced correctly",
	Enabled:     true,
	Apply:       applySemicolonsMustBeSpacedCorrectly,
	Description: ``,
}
```

Setup the test harness

``` go rules/semicolonsMustBeSpacedCorrectly_test.go
package rules

import (
	"bytes"
	"testing"
)

func TestSemicolonsMustBeSpacedCorrectly(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		<<<sa1002 tests>>>
	}

	for _, test := range tests {
		t.Run(test.description, func (t *testing.T) {
				actual := applySemicolonsMustBeSpacedCorrectly(test.given)
				if !bytes.Equal(test.expected, actual) {
					t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
				}
		})
	}
}
```

and declare our expectations

``` go "sa1002 tests"
{description: "when none are found", given: []byte("public void FunctionName(string s, int i)"), expected: []byte("public void FunctionName(string s, int i)")},
{description: "with inline comment", given: []byte("var i = 0;// blah"), expected: []byte("var i = 0; // blah")},
{description: "with no trailing space", given: []byte("for (i = 0;i < 4;i++) {"), expected: []byte("for (i = 0; i < 4; i++) {")},
{description: "with leading space and trailing space", given: []byte("return s + i.ToString() ; "), expected: []byte("return s + i.ToString();")},
```

### SA1003: Symbols must be spaced correctly

First the template

``` go rules/symbolsMustBeSpacedCorrectly.go
package rules

import (
	<<<sa1003 imports>>>
)

<<<sa1003 rule>>>
<<<sa1003 application>>>
```

``` go "sa1003 application"
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
```

Bring in used packages

``` go "sa1003 imports"
"bytes"
"regexp"
```

Now the logic has been worked out we'll apply create the rule.

``` go "sa1003 rule"
var symbolsMustBeSpacedCorrectly = &Rule{
	Name:        "Symbols must be spaced correctly",
	Enabled:     false,
	Apply:       applySymbolsMustBeSpacedCorrectly,
	Description: ``,
}
```

Setup the test harness

``` go rules/symbolsMustBeSpacedCorrectly_test.go
package rules

import (
	"bytes"
	"testing"
)

func TestSymbolsMustBeSpacedCorrectly(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		<<<sa1003 tests>>>
	}

	for _, test := range tests {
		t.Run(test.description, func (t *testing.T) {
				actual := applySymbolsMustBeSpacedCorrectly(test.given)
				if !bytes.Equal(test.expected, actual) {
					t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
				}
		})
	}
}
```

and declare our expectations

``` go "sa1003 tests"
{description: "documentation comment", given: []byte(`/// <summary>`), expected: []byte(`/// <summary>`)},
{description: "assignment", given: []byte(`int i=0;`), expected: []byte(`int i = 0;`)},
{description: "equality and increment by", given: []byte(`if (i==0) i+=2;`), expected: []byte(`if (i == 0) i += 2;`)},
{description: "not equal and decrement by", given: []byte(`if (i!=0) i-=2;`), expected: []byte(`if (i != 0) i -= 2;`)},
{description: "greater than and divide by", given: []byte(`if (i>0) i/=1;`), expected: []byte(`if (i > 0) i /= 1;`)},
{description: "less than and multiply by", given: []byte(`if (i<0) i*=1;`), expected: []byte(`if (i < 0) i *= 1;`)},
{description: "single comment", given: []byte(`// This is a comment`), expected: []byte(`// This is a comment`)},
{description: "less than or equal to, incrementor, and multiplication", given: []byte(`if (i<=0) i++*0;`), expected: []byte(`if (i <= 0) i++ * 0;`)},
{description: "inline comment with asterisk", given: []byte(`/* This is another comment */`), expected: []byte(`/* This is another comment */`)},
{description: "greater than or equal to pre-incrementor", given: []byte(`if (i>=0)++i;`), expected: []byte(`if (i >= 0) ++i;`)},
{description: "comment with asterisk and XML tag", given: []byte("/* <This is a comment>\n *\n */"), expected: []byte("/* <This is a comment>\n *\n */")},
{description: "common for loop", given: []byte(`for (var i=1; i>-1; i--) {}`), expected: []byte(`for (var i = 1; i > -1; i--) {}`)},
{description: "OR comparison", given: []byte(`if (i||i)--i;`), expected: []byte(`if (i || i) --i;`)},
{description: "AND comparison", given: []byte(`if (i&&i) i--*0;`), expected: []byte(`if (i && i) i-- * 0;`)},
{description: "arithmetic including parens", given: []byte(`if (1+1==2) i-1+(i*3)/(i/1)`), expected: []byte(`if (1 + 1 == 2) i - 1 + (i * 3) / (i / 1)`)},
{description: "AND comparison with negation", given: []byte(`a&&!b`), expected: []byte(`a && !b`)},
{description: "IF statement with negation of parens", given: []byte(`if!(true);`), expected: []byte(`if !(true);`)},
{description: "generics", given: []byte(`as IEnumerable<Namespace.ClassName>`), expected: []byte(`as IEnumerable<Namespace.ClassName>`)},
{description: "coalesce", given: []byte(`null??string.Empty;`), expected: []byte(`null ?? string.Empty;`)},
{description: "ternary", given: []byte(`(true)?true:false;`), expected: []byte(`(true) ? true : false;`)},
{description: "IF statement with multiline conditionals", given: []byte("if (true\n&& false)"), expected: []byte("if (true\n&& false)")},
{description: "attributes", given: []byte(`[Route("/api/[controller]")]`), expected: []byte(`[Route("/api/[controller]")]`)},
```

### SA1004: Documentation lines must begin with a single space

First the template

``` go rules/documentationLinesMustBeginWithSingleSpace.go
package rules

import (
	<<<sa1004 imports>>>
)

<<<sa1004 rule>>>
<<<sa1004 application>>>
```

``` go "sa1004 application"
func applyDocumentationLinesMustBeginWithSingleSpace(source []byte) []byte {
	re := regexp.MustCompile(`([/]{3})(\S)`)
	for re.Match(source) {
		source = re.ReplaceAll(source, []byte("$1 $2"))
	}
	return source
}
```

Bring in used packages

``` go "sa1004 imports"
"regexp"
```

Now the logic has been worked out we'll apply create the rule.

``` go "sa1004 rule"
var documentationLinesMustBeginWithSingleSpace = &Rule{
	Name:        "Documentation lines must begin with a single space",
	Enabled:     true,
	Apply:       applyDocumentationLinesMustBeginWithSingleSpace,
	Description: ``,
}
```

Setup the test harness

``` go rules/documentationLinesMustBeginWithSingleSpace_test.go
package rules

import (
	"bytes"
	"testing"
)

func TestDocumentationLinesMustBeginWithSingleSpace(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		<<<sa1004 tests>>>
	}

	for _, test := range tests {
		t.Run(test.description, func (t *testing.T) {
				actual := applyDocumentationLinesMustBeginWithSingleSpace(test.given)
				if !bytes.Equal(test.expected, actual) {
					t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
				}
		})
	}
}
```

and declare our expectations

``` go "sa1004 tests"
{description: "missing space between comment and opening XML", given: []byte("///<summary>"), expected: []byte("/// <summary>")},
{description: "missing space between comment and text", given: []byte("///The summary."), expected: []byte("/// The summary.")},
{description: "missing space between comment and closing XML", given: []byte("///</summary>"), expected: []byte("/// </summary>")},
{description: "do nothing if okay", given: []byte("/// <param name=\"foo\">The foo.</param>"), expected: []byte("/// <param name=\"foo\">The foo.</param>")},
```

### SA1005: Single line comments must begin with a single space

First the template

``` go rules/singleLineCommentsMustBeginWithSingleSpace.go
package rules

import (
	<<<sa1005 imports>>>
)

<<<sa1005 rule>>>
<<<sa1005 application>>>
```

``` go "sa1005 application"
func applySingleLineCommentsMustBeginWithSingleSpace(source []byte) []byte {
	reString := regexp.MustCompile(`".*"`)

	lines := []byte{}
	buffer := bytes.NewBuffer(source)
	scanner := bufio.NewScanner(buffer)

	for scanner.Scan() {
		// Add a newline character on each line after the first
		if len(lines) > 0 {
			lines = append(lines, byte('\n'))
		}

		line := scanner.Bytes()

		literal := []byte{}
		if reString.Match(line) {
			literal = reString.Find(line)
		}

		// Handle comments with no space.
		re := regexp.MustCompile(`(\s*)[/]{2}\s{0}(\S+)`)
		for re.Match(line) {
			if !bytes.Contains(literal, re.Find(line)) {
				line = re.ReplaceAll(line, []byte("$1// $2"))
			}
		}

		// Handle comments with more than one space
		re = regexp.MustCompile(`(\s*)[/]{2}\s{2,}(\S+)`)
		for re.Match(line) {
			if !bytes.Contains(literal, re.Find(line)) {
				line = re.ReplaceAll(line, []byte("$1// $2"))
			}
		}

		// Adjust for URIs
		re = regexp.MustCompile(`(\s*)[:]{1}[/]{2}\s+(\S+)`)
		for re.Match(line) {
			if !bytes.Contains(literal, re.Find(line)) {
				line = re.ReplaceAll(line, []byte("$1://$2"))
			}
		}

		// Trim ending
		line = bytes.TrimRightFunc(line, unicode.IsSpace)
		lines = append(lines, line...)
	}

	return lines
}
```

Bring in used packages

``` go "sa1005 imports"
"bufio"
"bytes"
"regexp"
"unicode"
```

Now the logic has been worked out we'll apply create the rule.

``` go "sa1005 rule"
var singleLineCommentsMustBeginWithSingleSpace = &Rule{
	Name:        "Single line comments must begin with single space",
	Enabled:     true,
	Apply:       applySingleLineCommentsMustBeginWithSingleSpace,
	Description: ``,
}
```

Setup the test harness

``` go rules/singleLineCommentsMustBeginWithSingleSpace_test.go
package rules

import (
	"bytes"
	"testing"
)

func TestSingleLineCommentsMustBeginWithSingleSpace(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		<<<sa1005 tests>>>
	}

	for _, test := range tests {
		t.Run(test.description, func (t *testing.T) {
				actual := applySingleLineCommentsMustBeginWithSingleSpace(test.given)
				if !bytes.Equal(test.expected, actual) {
					t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
				}
		})
	}
}
```

and declare our expectations

``` go "sa1005 tests"
{description: "missing space between comment and text", given: []byte("//This is a comment with no space"), expected: []byte("// This is a comment with no space")},
{description: "too many spaces between comment and text", given: []byte("//  This is a comment with too many spaces"), expected: []byte("// This is a comment with too many spaces")},
{description: "double commented code", given: []byte("////int i = 0;"), expected: []byte("// // int i = 0;")},
{description: "triple commented code", given: []byte("//////int i = 0;"), expected: []byte("// // // int i = 0;")},
{description: "comment with URI", given: []byte("//http://www.example.com"), expected: []byte("// http://www.example.com")},
```

### SA1006: Preprocessor keywords must not be preceded by space

First the template

``` go rules/preprocessorKeywordsMustNotBePrecededBySpace.go
package rules

import (
	<<<sa1006 imports>>>
)

<<<sa1006 rule>>>
<<<sa1006 application>>>
```

``` go "sa1006 application"
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
```

Bring in used packages

``` go "sa1006 imports"
"bytes"
"regexp"
```

Now the logic has been worked out we'll apply create the rule.

``` go "sa1006 rule"
var preprocessorKeywordsMustNotBePrecededBySpace = &Rule{
	Name:        "Single line comments must begin with single space",
	Enabled:     true,
	Apply:       applyPreprocessorKeywordsMustNotBePrecededBySpace,
	Description: ``,
}
```

Setup the test harness

``` go rules/preprocessorKeywordsMustNotBePrecededBySpace_test.go
package rules

import (
	"bytes"
	"testing"
)

func TestPreprocessorKeywordsMustNotBePrecededBySpace(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		<<<sa1006 tests>>>
	}

	for _, test := range tests {
		t.Run(test.description, func (t *testing.T) {
				actual := applyPreprocessorKeywordsMustNotBePrecededBySpace(test.given)
				if !bytes.Equal(test.expected, actual) {
					t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
				}
		})
	}
}
```

and declare our expectations

``` go "sa1006 tests"
{description: "if should not have space before hash", given: []byte("# if"), expected: []byte("#if")},
{description: "else should not have space before hash", given: []byte("# else"), expected: []byte("#else")},
{description: "elif should not have space before hash", given: []byte("# elif"), expected: []byte("#elif")},
{description: "endif should not have space before hash", given: []byte("# endif"), expected: []byte("#endif")},
{description: "define should not have space before hash", given: []byte("# define"), expected: []byte("#define")},
{description: "undef should not have space before hash", given: []byte("# undef"), expected: []byte("#undef")},
{description: "warning should not have space before hash", given: []byte("# warning"), expected: []byte("#warning")},
{description: "error should not have space before hash", given: []byte("# error"), expected: []byte("#error")},
{description: "line should not have space before hash", given: []byte("# line"), expected: []byte("#line")},
{description: "region should not have space before hash", given: []byte("# region"), expected: []byte("#region")},
{description: "endregion should not have space before hash", given: []byte("# endregion"), expected: []byte("#endregion")},
{description: "pragma should not have space before hash", given: []byte("# pragma"), expected: []byte("#pragma")},
{description: "pragma warning should not have space before hash", given: []byte("# pragma warning"), expected: []byte("#pragma warning")},
{description: "pragma checksum should not have space before hash", given: []byte("# pragma checksum"), expected: []byte("#pragma checksum")},
```

### SA1008: Opening parenthesis must be spaced correctly

First the template

``` go rules/openingParenthesisMustBeSpacedCorrectly.go
package rules

import (
	<<<sa1008 imports>>>
)

<<<sa1008 rule>>>
<<<sa1008 application>>>
```

``` go "sa1008 application"
func applyOpeningParenthesisMustBeSpacedCorrectly(source []byte) []byte {
	spaceBetween := `(if|while|for|switch|foreach|using|\+|\-|\*|/|&|\||\^|=)`

	return scan(source, func(line, literal []byte) []byte {
		// Remove leading spaces
		re := regexp.MustCompile(`([\S])(\t| )([\(])`)
		for re.Match(line) {
			line = re.ReplaceAll(line, []byte("$1$3"))
		}

		// Remove trailing spaces
		re = regexp.MustCompile(`([\(])(\t| )([\S])`)
		for re.Match(line) {
			line = re.ReplaceAll(line, []byte("$1$3"))
		}

		// Add space between operators and keywords
		re = regexp.MustCompile(spaceBetween + `([\(])`)
		for re.Match(line) {
			line = re.ReplaceAll(line, []byte("$1 $2"))
		}

		return line
	})
}
```

Bring in used packages

``` go "sa1008 imports"
"regexp"
```

Now the logic has been worked out we'll apply create the rule.

``` go "sa1008 rule"
var openingParenthesisMustBeSpacedCorrectly = &Rule{
	Name:        "Opening parenthesis must be spaced correctly",
	Enabled:     true,
	Apply:       applyOpeningParenthesisMustBeSpacedCorrectly,
	Description: ``,
}
```

Setup the test harness

``` go rules/openingParenthesisMustBeSpacedCorrectly_test.go
package rules

import (
	"bytes"
	"testing"
)

func TestOpeningParenthesisMustBeSpacedCorrectly(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		<<<sa1008 tests>>>
	}

	for _, test := range tests {
		t.Run(test.description, func (t *testing.T) {
				actual := applyOpeningParenthesisMustBeSpacedCorrectly(test.given)
				if !bytes.Equal(test.expected, actual) {
					t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
				}
		})
	}
}
```

and declare our expectations

``` go "sa1008 tests"
{description: "IF statement missing leading space", given: []byte("if(true) {}"), expected: []byte("if (true) {}")},
{description: "do nothing if okay", given: []byte("if (true) {}"), expected: []byte("if (true) {}")},
{description: "remove trailing space", given: []byte("public void something ( int i) {}"), expected: []byte("public void something(int i) {}")},
{description: "SWITCH statement missing leading space", given: []byte("switch(foo) {}"), expected: []byte("switch (foo) {}")},
{description: "in arithmetic", given: []byte("1+(2)"), expected: []byte("1+ (2)")},
{description: "multiline IF statment", given: []byte("if ("), expected: []byte("if (")},
```

### SA1009: Closing parenthesis must be spaced correctly

First the template

``` go rules/closingParenthesisMustBeSpacedCorrectly.go
package rules

import (
	<<<sa1009 imports>>>
)

<<<sa1009 rule>>>
<<<sa1009 application>>>
```

``` go "sa1009 application"
func applyClosingParenthesisMustBeSpacedCorrectly(source []byte) []byte {
	return scan(source, func(line, literal []byte) []byte {
		spaceBetween := `([A-z]|=|\+|\-|\*|/|&|\||\^|\{)`

		// Remove leading spaces
		re := regexp.MustCompile(`([\S])(\t| )([\)])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1$3"))
		}

		// Remove trailing spaces
		re = regexp.MustCompile(`([\)])(\t| )([\S])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1$3"))
		}

		// Add space between operators and keywords
		re = regexp.MustCompile(`([\)])` + spaceBetween)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1 $2"))
		}

		return line
	})
}
```

Bring in used packages

``` go "sa1009 imports"
"bytes"
"regexp"
```

Now the logic has been worked out we'll apply create the rule.

``` go "sa1009 rule"
var closingParenthesisMustBeSpacedCorrectly = &Rule{
	Name:        "Closing parenthesis must be spaced correctly",
	Enabled:     true,
	Apply:       applyClosingParenthesisMustBeSpacedCorrectly,
	Description: ``,
}
```

Setup the test harness

``` go rules/closingParenthesisMustBeSpacedCorrectly_test.go
package rules

import (
	"bytes"
	"testing"
)

func TestClosingParenthesisMustBeSpacedCorrectly(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		<<<sa1009 tests>>>
	}

	for _, test := range tests {
		t.Run(test.description, func (t *testing.T) {
				actual := applyClosingParenthesisMustBeSpacedCorrectly(test.given)
				if !bytes.Equal(test.expected, actual) {
					t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
				}
		})
	}
}
```

and declare our expectations

``` go "sa1009 tests"
{description: "with leading space and no trailing space", given: []byte("if (true ){}"), expected: []byte("if (true) {}")},
{description: "with proper spacing", given: []byte("if (true) {}"), expected: []byte("if (true) {}")},
{description: "function leading space only", given: []byte("public void something(int i ) {}"), expected: []byte("public void something(int i) {}")},
{description: "switch statement leading space and no trailing", given: []byte("switch (foo ){}"), expected: []byte("switch (foo) {}")},
{description: "in arithmetic", given: []byte("(2)+1"), expected: []byte("(2) +1")},
```

### SA1010: Opening square brackets must be spaced correctly

First the template

``` go rules/openingSquareBracketsMustBeSpacedCorrectly.go
package rules

import (
	<<<sa1010 imports>>>
)

<<<sa1010 rule>>>
<<<sa1010 application>>>
```

``` go "sa1010 application"
func applyOpeningSquareBracketsMustBeSpacedCorrectly(source []byte) []byte {
	return scan(source, func(line, literal []byte) []byte {
		re := regexp.MustCompile(`([\S])([\t ]+)([\[])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1$3"))
		}

		re = regexp.MustCompile(`([\[])([\t ]+)([\S])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1$3"))
		}

		return line
	})
}
```

Bring in used packages

``` go "sa1010 imports"
"bytes"
"regexp"
```

Now the logic has been worked out we'll apply create the rule.

``` go "sa1010 rule"
var openingSquareBracketsMustBeSpacedCorrectly = &Rule{
	Name:        "Opening square brackets must be spaced correctly",
	Enabled:     true,
	Apply:       applyOpeningSquareBracketsMustBeSpacedCorrectly,
	Description: ``,
}
```

Setup the test harness

``` go rules/openingSquareBracketsMustBeSpacedCorrectly_test.go
package rules

import (
	"bytes"
	"testing"
)

func TestOpeningSquareBracketsMustBeSpacedCorrectly(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		<<<sa1010 tests>>>
	}

	for _, test := range tests {
		t.Run(test.description, func (t *testing.T) {
				actual := applyOpeningSquareBracketsMustBeSpacedCorrectly(test.given)
				if !bytes.Equal(test.expected, actual) {
					t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
				}
		})
	}
}
```

and declare our expectations

``` go "sa1010 tests"
{description: "array declaration", given: []byte("new string [] {};"), expected: []byte("new string[] {};")},
{description: "array declaration with size", given: []byte("new int [ 1];"), expected: []byte("new int[1];")},
{description: "multiline array declaration opening", given: []byte("new string ["), expected: []byte("new string[")},
{description: "multiline array declaration closing", given: []byte("]"), expected: []byte("]")},
{description: "ignore in string", given: []byte("\"[ meh].[bleh]\","), expected: []byte("\"[ meh].[bleh]\",")},
```

### SA1011: Closing square brackets must be spaced correctly

First the template

``` go rules/closingSquareBracketsMustBeSpacedCorrectly.go
package rules

import (
	<<<sa1011 imports>>>
)

<<<sa1011 rule>>>
<<<sa1011 application>>>
```

``` go "sa1011 application"
func applyClosingSquareBracketsMustBeSpacedCorrectly(source []byte) []byte {
	return scan(source, func(line, literal []byte) []byte {
		re := regexp.MustCompile(`([\S])([\t ]+)([\]])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1$3"))
		}

		re = regexp.MustCompile(`([\]])([\S])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1 $2"))
		}
		re = regexp.MustCompile(`([\]]) ([;])`)
		if !bytes.Contains(literal, re.Find(line)) {
			line = re.ReplaceAll(line, []byte("$1$2"))
		}

		return line
	})
}
```

Bring in used packages

``` go "sa1011 imports"
"bytes"
"regexp"
```

Now the logic has been worked out we'll apply create the rule.

``` go "sa1011 rule"
var closingSquareBracketsMustBeSpacedCorrectly = &Rule{
	Name:        "Closing square brackets must be spaced correctly",
	Enabled:     true,
	Apply:       applyClosingSquareBracketsMustBeSpacedCorrectly,
	Description: ``,
}
```

Setup the test harness

``` go rules/closingSquareBracketsMustBeSpacedCorrectly_test.go
package rules

import (
	"bytes"
	"testing"
)

func TestClosingSquareBracketsMustBeSpacedCorrectly(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		<<<sa1011 tests>>>
	}

	for _, test := range tests {
		t.Run(test.description, func (t *testing.T) {
				actual := applyClosingSquareBracketsMustBeSpacedCorrectly(test.given)
				if !bytes.Equal(test.expected, actual) {
					t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
				}
		})
	}
}
```

and declare our expectations

``` go "sa1011 tests"
{description: "array declaration", given: []byte("new string[ ]{};"), expected: []byte("new string[] {};")},
{description: "array declaration with size", given: []byte("new int[1 ] ;"), expected: []byte("new int[1];")},
{description: "multiline array declaration opening", given: []byte("new string["), expected: []byte("new string[")},
{description: "multiline array declaration closing", given: []byte("]"), expected: []byte("]")},
{description: "ignore in string", given: []byte("\"[meh].[bleh ]\","), expected: []byte("\"[meh].[bleh ]\",")},
```

### SA1025: Code must not contain multiple whitespace in a row

First the template

``` go rules/codeMustNotContainMultipleWhitespaceInARow.go
package rules

import (
	<<<sa1025 imports>>>
)

<<<sa1025 rule>>>
<<<sa1025 application>>>
```

``` go "sa1025 application"
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
```

Bring in used packages

``` go "sa1025 imports"
"bytes"
"regexp"
```

Now the logic has been worked out we'll apply create the rule.

``` go "sa1025 rule"
var codeMustNotContainMultipleWhitespaceInARow = &Rule{
	Name:        "Code must not contain multiple whitespaces in a row",
	Enabled:     true,
	Apply:       applyCodeMustNotContainMultipleWhitespaceInARow,
	Description: `A violation of this rule occurs whenver the code contains a tab character.`,
}
```

Setup the test harness

``` go rules/codeMustNotContainMultipleWhitespaceInARow_test.go
package rules

import (
	"bytes"
	"testing"
)

func TestCodeMustNotContainMultipleWhitespaceInARow(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		<<<sa1025 tests>>>
	}

	for _, test := range tests {
		t.Run(test.description, func (t *testing.T) {
				actual := applyCodeMustNotContainMultipleWhitespaceInARow(test.given)
				if !bytes.Equal(test.expected, actual) {
					t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
				}
		})
	}
}
```

and declare our expectations

``` go "sa1025 tests"
{description: "remove multiple spaces", given: []byte("if  (i  == 0)  {  }"), expected: []byte("if (i == 0) { }")},
```

### SA1027: Tabs must not be used

First the template

``` go rules/tabsMustNotBeUsed.go
package rules

import (
	<<<sa1027 imports>>>
)

<<<sa1027 rule>>>
<<<sa1027 application>>>
```

``` go "sa1027 application"
func applyTabsMustNotBeUsed(source []byte) []byte {
	re := regexp.MustCompile(`\t`)
	for re.Match(source) {
		source = re.ReplaceAllLiteral(source, []byte("    "))
	}
	return source
}
```

Bring in used packages

``` go "sa1027 imports"
"regexp"
```

Now the logic has been worked out we'll apply create the rule.

``` go "sa1027 rule"
var tabsMustNotBeUsed = &Rule{
	Name:        "Tabs must not be used",
	Enabled:     true,
	Apply:       applyTabsMustNotBeUsed,
	Description: `A violation of this rule occurs whenver the code contains a tab character.`,
}
```

Setup the test harness

``` go rules/tabsMustNotBeUsed_test.go
package rules

import (
	"bytes"
	"testing"
)

func TestTabsMustNotBeUsed(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		<<<sa1027 tests>>>
	}

	for _, test := range tests {
		t.Run(test.description, func (t *testing.T) {
				actual := applyTabsMustNotBeUsed(test.given)
				if !bytes.Equal(test.expected, actual) {
					t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
				}
		})
	}
}
```

and declare our expectations

``` go "sa1027 tests"
{
	description: "remove tabs and replace with spaces",
	given: []byte("public void FunctionName(string s, int i)\n{\n\tvar i = 0; // blah\n\tfor (i = 0; i < 4; i++) {\n\t\t// Do something\n\t}\n\treturn s + i.ToString();\n}"),
	expected: []byte("public void FunctionName(string s, int i)\n{\n    var i = 0; // blah\n    for (i = 0; i < 4; i++) {\n        // Do something\n    }\n    return s + i.ToString();\n}"),
},
```

### SA1210: Using directives must be ordered alphabetically by namespace

First the template

``` go rules/usingDirectivesMustBeOrderedAlphabeticallyByNamespace.go
package rules

import (
	<<<sa1210 imports>>>
)

<<<sa1210 rule>>>
<<<sa1210 application>>>
```

``` go "sa1210 application"
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
```

Bring in used packages

``` go "sa1210 imports"
"fmt"
"regexp"
"sort"
"strings"
```

Now the logic has been worked out we'll apply create the rule.

``` go "sa1210 rule"
var usingDirectivesMustBeOrderedAlphabeticallyByNamespace = &Rule{
	Name:        "Using directives must be ordered alphabetically by namespace",
	Enabled:     true,
	Apply:       applyUsingDirectivesMustBeOrderedAlphabeticallyByNamespace,
	Description: ``,
}
```

Setup the test harness

``` go rules/usingDirectivesMustBeOrderedAlphabeticallyByNamespace_test.go
package rules

import (
	"bytes"
	"testing"
)

func TestUsingDirectivesMustBeOrderedAlphabeticallyByNamespace(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		<<<sa1210 tests>>>
	}

	for _, test := range tests {
		t.Run(test.description, func (t *testing.T) {
				actual := applyUsingDirectivesMustBeOrderedAlphabeticallyByNamespace(test.given)
				if !bytes.Equal(test.expected, actual) {
					t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
				}
		})
	}
}
```

and declare our expectations

``` go "sa1210 tests"
{
	description: "sort alphabetically",
	given: []byte(
		"using System;\n" +
		"using System.Collections.Generic;\n" +
		"using System.Linq;\n" +
		"using System.Web.Services;\n" +
		"using System.Web.UI;\n" +
		"using Use.This.Example.Concrete;\n" +
		"using Use.This.Example.Extensions;\n" +
		"\n" +
		"namespace Company.Blah {}"),
	expected: []byte(
		"using System;\n" +
		"using System.Collections.Generic;\n" +
		"using System.Linq;\n" +
		"using System.Web.Services;\n" +
		"using System.Web.UI;\n" +
		"using Use.This.Example.Concrete;\n" +
		"using Use.This.Example.Extensions;\n" +
		"\n" +
		"namespace Company.Blah {}"),
},
{
	description: "ignore using blocks",
	given: []byte(
		"using Company;\n" +
		"using CompanyB.System;\n" +
		"using Company.Collections.Generic;\n" +
		"using Company.Linq;\n" +
		"\n" +
		"namespace Company.Blah {}\n" +
		"using (var something = new Something()) {}"),
	expected: []byte(
		"using Company;\n" +
		"using Company.Collections.Generic;\n" +
		"using Company.Linq;\n" +
		"using CompanyB.System;\n" +
		"\n" +
		"namespace Company.Blah {}\n" +
		"using (var something = new Something()) {}"),
},
```

### SA1507: Code must not contain multiple blank lines in a row

First the template

``` go rules/codeMustNotContainMultipleBlankLinesInARow.go
package rules

import (
	<<<sa1507 imports>>>
)

<<<sa1507 rule>>>
<<<sa1507 application>>>
```

``` go "sa1507 application"
func applyCodeMustNotContainMultipleBlankLinesInARow(source []byte) []byte {
	re := regexp.MustCompile("\n{3,}")
	for re.Match(source) {
		source = re.ReplaceAllLiteral(source, []byte("\n"))
	}
	return source
}
```

Bring in used packages

``` go "sa1507 imports"
"regexp"
```

Now the logic has been worked out we'll apply create the rule.

``` go "sa1507 rule"
var codeMustNotContainMultipleBlankLinesInARow = &Rule{
	Name:        "Code must not contain multiple blank lines in a row",
	Enabled:     true,
	Apply:       applyCodeMustNotContainMultipleBlankLinesInARow,
	Description: ``,
}
```

Setup the test harness

``` go rules/codeMustNotContainMultipleBlankLinesInARow_test.go
package rules

import (
	"bytes"
	"testing"
)

func TestCodeMustNotContainMultipleBlankLinesInARow(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		<<<sa1507 tests>>>
	}

	for _, test := range tests {
		t.Run(test.description, func (t *testing.T) {
				actual := applyCodeMustNotContainMultipleBlankLinesInARow(test.given)
				if !bytes.Equal(test.expected, actual) {
					t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
				}
		})
	}
}
```

and declare our expectations

``` go "sa1507 tests"
{
	description: "remove multiple blank lines",
	given: []byte(
		"public void FunctionName(string s, int i)\n" +
		"{\n" +
		"\n" +
		"\n" +
		"\n" +
		"    return s + i.ToString();\n" +
		"}\n"),
	expected: []byte(
		"public void FunctionName(string s, int i)\n" +
		"{\n" +
		"    return s + i.ToString();\n" +
		"}\n"),
},
```
