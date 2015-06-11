package rules

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSymbolsMustBeSpacedCorrectly(t *testing.T) {
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
	as IEnumerable<Namespace.ClassName>`)
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
	for (var i = 1; i > -1; i-- ) {}
	if (i || i) --i;
	if (i && i) i-- * 0;
	if (1 + 1 == 2) i - 1 + (i * 3) / (i / 1)
	a && !b
	if !(true);
	as IEnumerable<Namespace.ClassName>`)

	actual := applySymbolsMustBeSpacedCorrectly(input)
	if !bytes.Equal(expected, actual) {
		fmt.Println(string(actual))
		t.Fail()
	}
}
