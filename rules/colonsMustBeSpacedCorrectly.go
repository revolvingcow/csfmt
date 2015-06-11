package rules

var colonsMustBeSpacedCorrectly = &Rule{
	Name:        "Colons must be spaced correctly",
	Enabled:     false,
	Apply:       applyColonsMustBeSpacedCorrectly,
	Description: ``,
}

func applyColonsMustBeSpacedCorrectly(source []byte) []byte {
	return source
}
