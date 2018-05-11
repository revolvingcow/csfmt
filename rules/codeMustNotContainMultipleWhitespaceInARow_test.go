package rules

import (
	"bytes"
	"testing"
)

func TestCodeMustNotContainMultipleWhitespaceInARow(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		{description: "remove multiple spaces", given: []byte("if  (i  == 0)  {  }"), expected: []byte("if (i == 0) { }")},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := applyCodeMustNotContainMultipleWhitespaceInARow(test.given)
			if !bytes.Equal(test.expected, actual) {
				t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
			}
		})
	}
}
