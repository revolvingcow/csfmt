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
		{description: "IF statement missing leading space", given: []byte("if(true) {}"), expected: []byte("if (true) {}")},
		{description: "do nothing if okay", given: []byte("if (true) {}"), expected: []byte("if (true) {}")},
		{description: "remove trailing space", given: []byte("public void something ( int i) {}"), expected: []byte("public void something(int i) {}")},
		{description: "SWITCH statement missing leading space", given: []byte("switch(foo) {}"), expected: []byte("switch (foo) {}")},
		{description: "in arithmetic", given: []byte("1+(2)"), expected: []byte("1+ (2)")},
		{description: "multiline IF statment", given: []byte("if ("), expected: []byte("if (")},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := applyOpeningParenthesisMustBeSpacedCorrectly(test.given)
			if !bytes.Equal(test.expected, actual) {
				t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
			}
		})
	}
}
