package rules

import (
	"bytes"
	"testing"
)

func TestCodeMustNotContainMultipleWhitespaceInARow(t *testing.T) {
	input := []byte(`if  (i  == 0)  {  }`)
	expected := []byte(`if (i == 0) { }`)

	actual := applyCodeMustNotContainMultipleWhitespaceInARow(input)
	if !bytes.Equal(expected, actual) {
		t.Fail()
	}
}
