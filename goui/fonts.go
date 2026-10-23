package goui

import (
	"embed"

	"github.com/akrck02/bubble"
	"github.com/akrck02/go-renderer/opengl/overlay"
)

//go:embed fonts/Roboto-Regular.ttf fonts/Roboto-Medium.ttf fonts/Roboto-Bold.ttf
var roboto embed.FS

// Roboto is Android's own typeface (SIL Open Font License, in fonts/), the material style's: give
// it to the window (window.Options.Fonts) and to the context (ui.NewInFonts), which measure and
// draw text in it.
func Roboto() overlay.Fonts {
	regular, _ := roboto.ReadFile("fonts/Roboto-Regular.ttf") // embedded: they are always there
	medium, _ := roboto.ReadFile("fonts/Roboto-Medium.ttf")
	bold, _ := roboto.ReadFile("fonts/Roboto-Bold.ttf")
	return overlay.Fonts{Regular: regular, Medium: medium, Bold: bold}
}

// FontsFor gives the typeface a theme is written in: Roboto for the material style, go-renderer's
// own otherwise.
func FontsFor(style bubble.Style) overlay.Fonts {
	if style == bubble.StyleMaterial {
		return Roboto()
	}
	return overlay.Fonts{}
}
