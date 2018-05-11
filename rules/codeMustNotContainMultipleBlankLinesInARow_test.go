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
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := applyCodeMustNotContainMultipleBlankLinesInARow(test.given)
			if !bytes.Equal(test.expected, actual) {
				t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
			}
		})
	}
}
