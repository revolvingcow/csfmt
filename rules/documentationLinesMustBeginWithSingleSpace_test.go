package rules

import (
	"bytes"
	"testing"
)

func TestDocumentationLineMustBeginWithSingleSpace(t *testing.T) {
	input := []byte(`///<summary>
///The summary.
///</summary>
/// <param name="foo">The foo.</param>
/// <returns>The bar.</returns>`)
	expected := []byte(`/// <summary>
/// The summary.
/// </summary>
/// <param name="foo">The foo.</param>
/// <returns>The bar.</returns>`)

	actual := applyDocumentationLinesMustBeginWithSingleSpace(input)
	if !bytes.Equal(expected, actual) {
		t.Fail()
	}
}
