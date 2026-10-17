package bubble

// Code is an error with a stable name that programs check and people can search for.
// A code is made once, so returning or comparing one costs no allocation.
type Code struct {
	name    string
	meaning string
}

// Error gives the code's name and what it means.
func (code *Code) Error() string {
	return code.name + ": " + code.meaning
}

// Name is the code's stable name, such as "bubble.invalid-colour".
func (code *Code) Name() string {
	return code.name
}

func define(name string, meaning string) *Code {
	return &Code{name: name, meaning: meaning}
}

// InvalidColour means a colour is not written as #RRGGBB or #RRGGBBAA.
var InvalidColour = define("bubble.invalid-colour", "a colour must be written as #RRGGBB or #RRGGBBAA")

// UnknownStyle means the style is neither paper nor material.
var UnknownStyle = define("bubble.unknown-style", "the style must be paper or material")

// UnknownMode means the mode is neither light nor dark.
var UnknownMode = define("bubble.unknown-mode", "the mode must be light or dark")

// UnknownDensity means the density is neither comfortable nor compact.
var UnknownDensity = define("bubble.unknown-density", "the density must be comfortable or compact")
