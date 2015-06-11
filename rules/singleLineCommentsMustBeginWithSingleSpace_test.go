package rules

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSingleLineCommentsMustBeingWithSingleSpace(t *testing.T) {
	input := []byte(`//This is a comment with no space
//  This is a comment with too many spaces
////int i = 0;
//////int i = 0;
////return i;`)
	expected := []byte(`// This is a comment with no space
// This is a comment with too many spaces
////int i = 0;
//////int i = 0;
////return i;`)

	actual := applySingleLineCommentsMustBeginWithSingleSpace(input)
	if !bytes.Equal(expected, actual) {
		fmt.Println(string(actual))
		t.Fail()
	}
}
