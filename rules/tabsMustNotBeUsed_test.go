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
		{
			description: "remove tabs and replace with spaces",
			given:       []byte("public void FunctionName(string s, int i)\n{\n\tvar i = 0; // blah\n\tfor (i = 0; i < 4; i++) {\n\t\t// Do something\n\t}\n\treturn s + i.ToString();\n}"),
			expected:    []byte("public void FunctionName(string s, int i)\n{\n    var i = 0; // blah\n    for (i = 0; i < 4; i++) {\n        // Do something\n    }\n    return s + i.ToString();\n}"),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := applyTabsMustNotBeUsed(test.given)
			if !bytes.Equal(test.expected, actual) {
				t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
			}
		})
	}
}
