package rules

import (
	"bytes"
	"fmt"
	"testing"
)

func TestOpeningParenthesisMustBeSpacedCorrectly(t *testing.T) {
	input := []byte(`if(true) {}
if (true) {}
public void something ( int i) {}
switch(foo) {}
1+(2)
if (`)
	expected := []byte(`if (true) {}
if (true) {}
public void something(int i) {}
switch (foo) {}
1+ (2)
if (`)

	actual := applyOpeningParenthesisMustBeSpacedCorrectly(input)
	if !bytes.Equal(expected, actual) {
		fmt.Println(string(actual))
		t.Fail()
	}
}
