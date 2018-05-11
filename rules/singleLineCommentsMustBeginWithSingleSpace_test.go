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
		{description: "missing space between comment and text", given: []byte("//This is a comment with no space"), expected: []byte("// This is a comment with no space")},
		{description: "too many spaces between comment and text", given: []byte("//  This is a comment with too many spaces"), expected: []byte("// This is a comment with too many spaces")},
		{description: "double commented code", given: []byte("////int i = 0;"), expected: []byte("// // int i = 0;")},
		{description: "triple commented code", given: []byte("//////int i = 0;"), expected: []byte("// // // int i = 0;")},
		{description: "comment with URI", given: []byte("//http://www.example.com"), expected: []byte("// http://www.example.com")},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := applySingleLineCommentsMustBeginWithSingleSpace(test.given)
			if !bytes.Equal(test.expected, actual) {
				t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
			}
		})
	}
}
