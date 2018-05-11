package rules

import (
	"bytes"
	"testing"
)

func TestUsingDirectivesMustBeOrderedAlphabeticallyByNamespace(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		{
			description: "sort alphabetically",
			given: []byte(
				"using System;\n" +
					"using System.Collections.Generic;\n" +
					"using System.Linq;\n" +
					"using System.Web.Services;\n" +
					"using System.Web.UI;\n" +
					"using Use.This.Example.Concrete;\n" +
					"using Use.This.Example.Extensions;\n" +
					"\n" +
					"namespace Company.Blah {}"),
			expected: []byte(
				"using System;\n" +
					"using System.Collections.Generic;\n" +
					"using System.Linq;\n" +
					"using System.Web.Services;\n" +
					"using System.Web.UI;\n" +
					"using Use.This.Example.Concrete;\n" +
					"using Use.This.Example.Extensions;\n" +
					"\n" +
					"namespace Company.Blah {}"),
		},
		{
			description: "ignore using blocks",
			given: []byte(
				"using Company;\n" +
					"using CompanyB.System;\n" +
					"using Company.Collections.Generic;\n" +
					"using Company.Linq;\n" +
					"\n" +
					"namespace Company.Blah {}\n" +
					"using (var something = new Something()) {}"),
			expected: []byte(
				"using Company;\n" +
					"using Company.Collections.Generic;\n" +
					"using Company.Linq;\n" +
					"using CompanyB.System;\n" +
					"\n" +
					"namespace Company.Blah {}\n" +
					"using (var something = new Something()) {}"),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := applyUsingDirectivesMustBeOrderedAlphabeticallyByNamespace(test.given)
			if !bytes.Equal(test.expected, actual) {
				t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
			}
		})
	}
}
