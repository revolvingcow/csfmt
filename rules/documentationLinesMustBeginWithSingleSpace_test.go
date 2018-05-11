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
		{description: "missing space between comment and opening XML", given: []byte("///<summary>"), expected: []byte("/// <summary>")},
		{description: "missing space between comment and text", given: []byte("///The summary."), expected: []byte("/// The summary.")},
		{description: "missing space between comment and closing XML", given: []byte("///</summary>"), expected: []byte("/// </summary>")},
		{description: "do nothing if okay", given: []byte("/// <param name=\"foo\">The foo.</param>"), expected: []byte("/// <param name=\"foo\">The foo.</param>")},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := applyDocumentationLinesMustBeginWithSingleSpace(test.given)
			if !bytes.Equal(test.expected, actual) {
				t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
			}
		})
	}
}
