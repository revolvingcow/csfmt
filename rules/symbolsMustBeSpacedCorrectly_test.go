package rules

import (
	"bytes"
	"testing"
)

func TestSymbolsMustBeSpacedCorrectlySimple(t *testing.T) {
	tests := []struct {
		given    []byte
		expected []byte
	}{
		{given: []byte(`/// <summary>`), expected: []byte(`/// <summary>`)},
		{given: []byte(`int i=0;`), expected: []byte(`int i = 0;`)},
		{given: []byte(`if (i==0) i+=2;`), expected: []byte(`if (i == 0) i += 2;`)},
		{given: []byte(`if (i!=0) i-=2;`), expected: []byte(`if (i != 0) i -= 2;`)},
		{given: []byte(`if (i>0) i/=1;`), expected: []byte(`if (i > 0) i /= 1;`)},
		{given: []byte(`if (i<0) i*=1;`), expected: []byte(`if (i < 0) i *= 1;`)},
		{given: []byte(`// This is a comment`), expected: []byte(`// This is a comment`)},
		{given: []byte(`if (i<=0) i++*0;`), expected: []byte(`if (i <= 0) i++ * 0;`)},
		{given: []byte(`/* This is another comment */`), expected: []byte(`/* This is another comment */`)},
		{given: []byte(`if (i>=0)++i;`), expected: []byte(`if (i >= 0) ++i;`)},
		{given: []byte("/* <This is a comment>\n *\n */"), expected: []byte("/* <This is a comment>\n *\n */")},
		{given: []byte(`for (var i=1; i>-1; i--) {}`), expected: []byte(`for (var i = 1; i > -1; i--) {}`)},
		{given: []byte(`if (i||i)--i;`), expected: []byte(`if (i || i) --i;`)},
		{given: []byte(`if (i&&i) i--*0;`), expected: []byte(`if (i && i) i-- * 0;`)},
		{given: []byte(`if (1+1==2) i-1+(i*3)/(i/1)`), expected: []byte(`if (1 + 1 == 2) i - 1 + (i * 3) / (i / 1)`)},
		{given: []byte(`a&&!b`), expected: []byte(`a && !b`)},
		{given: []byte(`if!(true);`), expected: []byte(`if !(true);`)},
		{given: []byte(`as IEnumerable<Namespace.ClassName>`), expected: []byte(`as IEnumerable<Namespace.ClassName>`)},
		{given: []byte(`null??string.Empty;`), expected: []byte(`null ?? string.Empty;`)},
		{given: []byte(`(true)?true:false;`), expected: []byte(`(true) ? true : false;`)},
		{given: []byte(`if (true`), expected: []byte(`if (true`)},
		{given: []byte(`&& false)`), expected: []byte(`&& false)`)},
		{given: []byte(`[Route("/api/[controller]")]`), expected: []byte(`[Route("/api/[controller]")]`)},
	}

	for _, test := range tests {
		actual := applySymbolsMustBeSpacedCorrectly(test.given)
		if !bytes.EqualFold(test.expected, actual) {
			t.Fatalf("Expected `%s` but received `%s`", string(test.expected), string(actual))
		}
	}
}

func TestSymbolsMustBeSpacedCorrectlyComplex(t *testing.T) {
	input := []byte(`/// <summary>
int i=0;
if (i==0) i+=2;
if (i!=0) i-=2;
if (i>0) i/=1;
if (i<0) i*=1;
// This is a comment
if (i<=0) i++*0;
/* This is another comment */
if (i>=0)++i;
/* <This is a comment>
 *
 */
for (var i=1; i>-1; i--) {}
if (i||i)--i;
if (i&&i) i--*0;
if (1+1==2) i-1+(i*3)/(i/1)
a&&!b
if!(true);
as IEnumerable<Namespace.ClassName>
null??string.Empty;
(true)?true:false;
if (true
	&& false)
[Route("/api/[controller]")]`)
	expected := []byte(`/// <summary>
int i = 0;
if (i == 0) i += 2;
if (i != 0) i -= 2;
if (i > 0) i /= 1;
if (i < 0) i *= 1;
// This is a comment
if (i <= 0) i++ * 0;
/* This is another comment */
if (i >= 0) ++i;
/* <This is a comment>
 *
 */
for (var i = 1; i > -1; i--) {}
if (i || i) --i;
if (i && i) i-- * 0;
if (1 + 1 == 2) i - 1 + (i * 3) / (i / 1)
a && !b
if !(true);
as IEnumerable<Namespace.ClassName>
null ?? string.Empty;
(true) ? true : false;
if (true
	&& false)
[Route("/api/[controller]")]`)

	actual := applySymbolsMustBeSpacedCorrectly(input)
	if !bytes.Equal(expected, actual) {
		t.Fail()
	}
}
