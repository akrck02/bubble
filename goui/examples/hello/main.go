// Command hello is the Go program of Bubble's getting started (docs/2-getting-started.md): a
// heading, a field and a button in Go UI, wearing Bubble with the herbarium's green.
package main

import (
	"flag"
	"log"

	"github.com/akrck02/bubble"
	"github.com/akrck02/bubble/goui"
	ui "github.com/akrck02/go-ui"
	"github.com/akrck02/go-ui/window"
)

func main() {
	accent, err := bubble.ParseColor("#4C9A5B")
	if err != nil {
		log.Fatal(err)
	}
	options := bubble.DefaultOptions()
	options.Accent = accent
	theme, err := bubble.NewTheme(options)
	if err != nil {
		log.Fatal(err)
	}
	windowOptions := window.Options{Title: "Hello, Bubble"}
	windowOptions.RegisterFlags(flag.CommandLine) // -dark, -capture and the rest
	flag.Parse()
	name := "Rosemary"
	err = window.Run(windowOptions, func(context *ui.Context) {
		context.Theme = goui.Theme(theme)
		context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.Pad(24), Gap: 12}, func() {
			context.Heading("Hello, Bubble")
			context.TextField("name", &name)
			if context.PrimaryButton("Save") {
				log.Println("saved", name)
			}
		})
	})
	if err != nil {
		log.Fatal(err)
	}
}
