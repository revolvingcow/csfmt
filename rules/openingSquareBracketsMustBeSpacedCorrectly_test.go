package rules

import (
	"bytes"
	"fmt"
	"testing"
)

func TestOpeningSquareBracketsMustBeSpacedCorrectly(t *testing.T) {
	input := []byte(`new string [] {};
new int [ 1];
new string [
	"blah",
	"meh",
	"bleh"
]`)
	expected := []byte(`new string[] {};
new int[1];
new string[
	"blah",
	"meh",
	"bleh"
]`)

	actual := applyOpeningSquareBracketsMustBeSpacedCorrectly(input)
	if !bytes.Equal(expected, actual) {
		fmt.Println(string(actual))
		t.Fail()
	}
}
