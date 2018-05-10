package rules

import (
	"bytes"
	"testing"
)

func TestSingleLineCommentsMustBeingWithSingleSpace(t *testing.T) {
	input := []byte(`//This is a comment with no space
//  This is a comment with too many spaces
////int i = 0;
//////int i = 0;
////return i;
//http://www.example.com`)
	expected := []byte(`// This is a comment with no space
// This is a comment with too many spaces
// // int i = 0;
// // // int i = 0;
// // return i;
// http://www.example.com`)

	actual := applySingleLineCommentsMustBeginWithSingleSpace(input)
	if !bytes.Equal(expected, actual) {
		t.Fail()
	}
}
