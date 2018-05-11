package rules

import (
	"bytes"
	"testing"
)

func TestPreprocessorKeywordsMustNotBePrecededBySpace(t *testing.T) {
	tests := []struct {
		description string
		given       []byte
		expected    []byte
	}{
		{description: "if should not have space before hash", given: []byte("# if"), expected: []byte("#if")},
		{description: "else should not have space before hash", given: []byte("# else"), expected: []byte("#else")},
		{description: "elif should not have space before hash", given: []byte("# elif"), expected: []byte("#elif")},
		{description: "endif should not have space before hash", given: []byte("# endif"), expected: []byte("#endif")},
		{description: "define should not have space before hash", given: []byte("# define"), expected: []byte("#define")},
		{description: "undef should not have space before hash", given: []byte("# undef"), expected: []byte("#undef")},
		{description: "warning should not have space before hash", given: []byte("# warning"), expected: []byte("#warning")},
		{description: "error should not have space before hash", given: []byte("# error"), expected: []byte("#error")},
		{description: "line should not have space before hash", given: []byte("# line"), expected: []byte("#line")},
		{description: "region should not have space before hash", given: []byte("# region"), expected: []byte("#region")},
		{description: "endregion should not have space before hash", given: []byte("# endregion"), expected: []byte("#endregion")},
		{description: "pragma should not have space before hash", given: []byte("# pragma"), expected: []byte("#pragma")},
		{description: "pragma warning should not have space before hash", given: []byte("# pragma warning"), expected: []byte("#pragma warning")},
		{description: "pragma checksum should not have space before hash", given: []byte("# pragma checksum"), expected: []byte("#pragma checksum")},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := applyPreprocessorKeywordsMustNotBePrecededBySpace(test.given)
			if !bytes.Equal(test.expected, actual) {
				t.Errorf("Got `%s` but wanted `%s`", string(actual), string(test.expected))
			}
		})
	}
}
