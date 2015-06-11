package rules

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCodeMustNotContainMultipleWhitespaceInARow(t *testing.T) {
	input := []byte(`
	if  (i  == 0)  {  }`)
	expected := []byte(`
	if (i == 0) { }`)

	actual := applyCodeMustNotContainMultipleWhitespaceInARow(input)
	if !bytes.Equal(expected, actual) {
		fmt.Println(string(actual))
		t.Fail()
	}
}
