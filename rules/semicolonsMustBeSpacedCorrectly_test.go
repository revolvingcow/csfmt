package rules

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSemicolonsMustBeSpacedCorrectly(t *testing.T) {
	input := []byte(`public void FunctionName(string s, int i)
{
	var i = 0; // blah
	for (i = 0;i < 4;i++) {
		// Do something
	}
	return s + i.ToString() ;
}`)
	expected := []byte(`public void FunctionName(string s, int i)
{
	var i = 0; // blah
	for (i = 0; i < 4; i++) {
		// Do something
	}
	return s + i.ToString();
}`)

	actual := applySemicolonsMustBeSpacedCorrectly(input)
	if !bytes.Equal(expected, actual) {
		fmt.Println(string(actual))
		t.Fail()
	}
}
