package rules

import (
	"bytes"
	"testing"
)

func TestCommasMustBeSpacedCorrectly(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		{
			description: "fix spacing around multiline commas",
			given: []byte(
				"public void FunctionName(string s ,int i\n" +
					"	, object x)\n" +
					"{\n" +
					"	int[] b = new [1,   3,4 ,5];\n" +
					"	var o = new {\n" +
					"		blah = \"stomething\",\n" +
					"		meh = 0,\n" +
					"		dude = true,\n" +
					"	};\n" +
					"}"),
			expected: []byte(
				"public void FunctionName(string s, int i, object x)\n" +
					"{\n" +
					"	int[] b = new [1, 3, 4, 5];\n" +
					"	var o = new {\n" +
					"		blah = \"stomething\",\n" +
					"		meh = 0,\n" +
					"		dude = true,\n" +
					"	};\n" +
					"}"),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := applyCommasMustBeSpacedCorrectly(test.given)
			if !bytes.Equal(test.expected, actual) {
				t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
			}
		})
	}
}
