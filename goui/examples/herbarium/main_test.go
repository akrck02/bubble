package main

import (
	"testing"

	"github.com/akrck02/bubble"
	"github.com/akrck02/go-ui/window"
)

func TestTheHerbariumOpensInEveryLook(test *testing.T) {
	for _, look := range []struct {
		style bubble.Style
		glass bool
		dark  bool
	}{
		{bubble.StylePaper, false, false}, {bubble.StylePaper, false, true},
		{bubble.StyleMaterial, false, false}, {bubble.StylePaper, true, false},
	} {
		garden, err := newHerbarium(look.style, look.glass, window.Options{Width: 1280, Height: 760, Dark: look.dark})
		if err != nil {
			test.Fatalf("%+v: %v", look, err)
		}
		if look.glass != (garden.ground != nil) {
			test.Errorf("%+v: glass's ground should be there only on glass", look)
		}
	}
}

func TestTheSearchFindsAPlantWhateverTheCase(test *testing.T) {
	if !matches("Rosemary", "  rose ") {
		test.Error("rose should find Rosemary")
	}
	if matches("Lavender", "rose") {
		test.Error("rose should not find Lavender")
	}
}
