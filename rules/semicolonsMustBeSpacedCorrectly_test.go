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
		{description: "when none are found", given: []byte("public void FunctionName(string s, int i)"), expected: []byte("public void FunctionName(string s, int i)")},
		{description: "with inline comment", given: []byte("var i = 0;// blah"), expected: []byte("var i = 0; // blah")},
		{description: "with no trailing space", given: []byte("for (i = 0;i < 4;i++) {"), expected: []byte("for (i = 0; i < 4; i++) {")},
		{description: "with leading space and trailing space", given: []byte("return s + i.ToString() ; "), expected: []byte("return s + i.ToString();")},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := applySemicolonsMustBeSpacedCorrectly(test.given)
			if !bytes.Equal(test.expected, actual) {
				t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
			}
		})
	}
}
