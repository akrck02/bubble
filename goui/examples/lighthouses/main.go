// Command lighthouses is the example of docs/5-guides/5-material-components.md: a phone-sized
// window in Bubble's Material, with a top bar and a button that counts the lights.
package main

import (
	"fmt"
	"log"

	"github.com/akrck02/bubble"
	"github.com/akrck02/bubble/goui"
	"github.com/akrck02/bubble/goui/material"
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
	"github.com/akrck02/go-ui/window"
)

func main() {
	accent, err := bubble.ParseColor("#2C5C9E")
	if err != nil {
		log.Fatal(err)
	}
	theme, err := bubble.NewTheme(bubble.Options{Style: bubble.StyleMaterial, Mode: bubble.ModeLight,
		Density: bubble.DensityComfortable, Accent: accent})
	if err != nil {
		log.Fatal(err)
	}
	options := window.Options{Title: "Lighthouses", Width: 412, Height: 892}
	options.Fonts = goui.FontsFor(theme.Options.Style) // Roboto, for the window to draw in
	lights := 0
	err = window.Run(options, func(context *ui.Context) {
		context.Theme = goui.Theme(theme)
		look := material.New(context, theme)
		look.TopAppBar(material.TopBar{Title: "Lighthouses"})
		if look.FilledButton("plus", "New light") {
			lights++
		}
		look.Text(fmt.Sprintf("%d lights", lights), 16, draw2d.Regular, look.Roles.OnSurface)
	})
	if err != nil {
		log.Fatal(err)
	}
}
