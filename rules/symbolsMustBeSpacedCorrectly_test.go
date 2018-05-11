package rules

import (
	"bytes"
	"testing"
)

func TestSymbolsMustBeSpacedCorrectly(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		{description: "documentation comment", given: []byte(`/// <summary>`), expected: []byte(`/// <summary>`)},
		{description: "assignment", given: []byte(`int i=0;`), expected: []byte(`int i = 0;`)},
		{description: "equality and increment by", given: []byte(`if (i==0) i+=2;`), expected: []byte(`if (i == 0) i += 2;`)},
		{description: "not equal and decrement by", given: []byte(`if (i!=0) i-=2;`), expected: []byte(`if (i != 0) i -= 2;`)},
		{description: "greater than and divide by", given: []byte(`if (i>0) i/=1;`), expected: []byte(`if (i > 0) i /= 1;`)},
		{description: "less than and multiply by", given: []byte(`if (i<0) i*=1;`), expected: []byte(`if (i < 0) i *= 1;`)},
		{description: "single comment", given: []byte(`// This is a comment`), expected: []byte(`// This is a comment`)},
		{description: "less than or equal to, incrementor, and multiplication", given: []byte(`if (i<=0) i++*0;`), expected: []byte(`if (i <= 0) i++ * 0;`)},
		{description: "inline comment with asterisk", given: []byte(`/* This is another comment */`), expected: []byte(`/* This is another comment */`)},
		{description: "greater than or equal to pre-incrementor", given: []byte(`if (i>=0)++i;`), expected: []byte(`if (i >= 0) ++i;`)},
		{description: "comment with asterisk and XML tag", given: []byte("/* <This is a comment>\n *\n */"), expected: []byte("/* <This is a comment>\n *\n */")},
		{description: "common for loop", given: []byte(`for (var i=1; i>-1; i--) {}`), expected: []byte(`for (var i = 1; i > -1; i--) {}`)},
		{description: "OR comparison", given: []byte(`if (i||i)--i;`), expected: []byte(`if (i || i) --i;`)},
		{description: "AND comparison", given: []byte(`if (i&&i) i--*0;`), expected: []byte(`if (i && i) i-- * 0;`)},
		{description: "arithmetic including parens", given: []byte(`if (1+1==2) i-1+(i*3)/(i/1)`), expected: []byte(`if (1 + 1 == 2) i - 1 + (i * 3) / (i / 1)`)},
		{description: "AND comparison with negation", given: []byte(`a&&!b`), expected: []byte(`a && !b`)},
		{description: "IF statement with negation of parens", given: []byte(`if!(true);`), expected: []byte(`if !(true);`)},
		{description: "generics", given: []byte(`as IEnumerable<Namespace.ClassName>`), expected: []byte(`as IEnumerable<Namespace.ClassName>`)},
		{description: "coalesce", given: []byte(`null??string.Empty;`), expected: []byte(`null ?? string.Empty;`)},
		{description: "ternary", given: []byte(`(true)?true:false;`), expected: []byte(`(true) ? true : false;`)},
		{description: "IF statement with multiline conditionals", given: []byte("if (true\n&& false)"), expected: []byte("if (true\n&& false)")},
		{description: "attributes", given: []byte(`[Route("/api/[controller]")]`), expected: []byte(`[Route("/api/[controller]")]`)},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := applySymbolsMustBeSpacedCorrectly(test.given)
			if !bytes.Equal(test.expected, actual) {
				t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
			}
		})
	}
}
