package rules

// Rule is a style rule to look for and apply within the source code.
type Rule struct {
	Name        string
	Description string
	Enabled     bool
	Apply       func(source []byte) []byte
}
