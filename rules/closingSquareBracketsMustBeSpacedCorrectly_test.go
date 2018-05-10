package rules

import (
	"bytes"
	"testing"
)

func TestClosingSquareBracketsMustBeSpacedCorrectly(t *testing.T) {
	input := []byte(`new string[ ]{};
new int[1 ] ;
new string[
	"blah",
	"[meh].[bleh]",
]`)
	expected := []byte(`new string[] {};
new int[1];
new string[
	"blah",
	"[meh].[bleh]",
]`)

	actual := applyClosingSquareBracketsMustBeSpacedCorrectly(input)
	if !bytes.Equal(expected, actual) {
		t.Fail()
	}
}
