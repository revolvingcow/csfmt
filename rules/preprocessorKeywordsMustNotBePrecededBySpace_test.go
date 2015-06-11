package rules

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPreprocessorKeywordsMustNotBePrecededBySpace(t *testing.T) {
	input := []byte(`
# if
	// Something
# else
	// Something
# elif
	// Something
# endif
# define
# undef
# warning
# error
# line
# region
# endregion
# pragma
# pragma warning
# pragma checksum`)
	expected := []byte(`
#if
	// Something
#else
	// Something
#elif
	// Something
#endif
#define
#undef
#warning
#error
#line
#region
#endregion
#pragma
#pragma warning
#pragma checksum`)

	actual := applyPreprocessorKeywordsMustNotBePrecededBySpace(input)
	if !bytes.Equal(expected, actual) {
		fmt.Println(string(actual))
		t.Fail()
	}
}
