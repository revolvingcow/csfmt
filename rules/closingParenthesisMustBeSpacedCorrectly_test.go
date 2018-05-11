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
		{description: "with leading space and no trailing space", given: []byte("if (true ){}"), expected: []byte("if (true) {}")},
		{description: "with proper spacing", given: []byte("if (true) {}"), expected: []byte("if (true) {}")},
		{description: "function leading space only", given: []byte("public void something(int i ) {}"), expected: []byte("public void something(int i) {}")},
		{description: "switch statement leading space and no trailing", given: []byte("switch (foo ){}"), expected: []byte("switch (foo) {}")},
		{description: "in arithmetic", given: []byte("(2)+1"), expected: []byte("(2) +1")},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := applyClosingParenthesisMustBeSpacedCorrectly(test.given)
			if !bytes.Equal(test.expected, actual) {
				t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
			}
		})
	}
}
