package rules

import (
	"bytes"
	"testing"
)

func TestTabsMustNotBeUsed(t *testing.T) {
	input := []byte("public void FunctionName(string s, int i)\n{\n\tvar i = 0; // blah\n\tfor (i = 0; i < 4; i++) {\n\t\t// Do something\n\t}\n\treturn s + i.ToString();\n}")
	expected := []byte("public void FunctionName(string s, int i)\n{\n    var i = 0; // blah\n    for (i = 0; i < 4; i++) {\n        // Do something\n    }\n    return s + i.ToString();\n}")

	actual := applyTabsMustNotBeUsed(input)
	if !bytes.Equal(expected, actual) {
		t.Fail()
	}
}
