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
		{description: "array declaration", given: []byte("new string[ ]{};"), expected: []byte("new string[] {};")},
		{description: "array declaration with size", given: []byte("new int[1 ] ;"), expected: []byte("new int[1];")},
		{description: "multiline array declaration opening", given: []byte("new string["), expected: []byte("new string[")},
		{description: "multiline array declaration closing", given: []byte("]"), expected: []byte("]")},
		{description: "ignore in string", given: []byte("\"[meh].[bleh ]\","), expected: []byte("\"[meh].[bleh ]\",")},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := applyClosingSquareBracketsMustBeSpacedCorrectly(test.given)
			if !bytes.Equal(test.expected, actual) {
				t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
			}
		})
	}
}
