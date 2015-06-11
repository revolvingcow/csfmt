package rules

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCommasMustBeSpacedCorrectly(t *testing.T) {
	input := []byte(`public void FunctionName(string s ,int i
	, object x)
{
	int[] b = new [1,   3,4 ,5];
	var o = new {
		blah = "stomething",
		meh = 0,
		dude = true,
	};
}`)
	expected := []byte(`public void FunctionName(string s, int i, object x)
{
	int[] b = new [1, 3, 4, 5];
	var o = new {
		blah = "stomething",
		meh = 0,
		dude = true,
	};
}`)

	actual := applyCommasMustBeSpacedCorrectly(input)
	if !bytes.Equal(expected, actual) {
		fmt.Println(string(actual))
		t.Fail()
	}
}
