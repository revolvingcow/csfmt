package rules

import (
	"bytes"
	"fmt"
	"testing"
)

func TestClosingParenthesisMustBeSpacedCorrectly(t *testing.T) {
	input := []byte(`if (true ){}
if (true) {}
public void something(int i ) {}
switch (foo ){}
(2)+1`)
	expected := []byte(`if (true) {}
if (true) {}
public void something(int i) {}
switch (foo) {}
(2) +1`)

	actual := applyClosingParenthesisMustBeSpacedCorrectly(input)
	if !bytes.Equal(expected, actual) {
		fmt.Println(string(actual))
		t.Fail()
	}
}
