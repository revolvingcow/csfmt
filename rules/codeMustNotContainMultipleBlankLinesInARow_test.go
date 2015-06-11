package rules

import (
	"bytes"
	"testing"
)

func TestCodeMustNotContainMultipleBlankLinesInARow(t *testing.T) {
	input := []byte(`public void FunctionName(string s, int i)
{



	return s + i.ToString();
}`)
	expected := []byte(`public void FunctionName(string s, int i)
{
	return s + i.ToString();
}`)

	actual := applyCodeMustNotContainMultipleBlankLinesInARow(input)
	if !bytes.Equal(expected, actual) {
		t.Fail()
	}
}
