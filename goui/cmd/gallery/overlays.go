package main

import (
	"github.com/akrck02/bubble/goui/material"
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// overlaysState is what goes over every screen: a dialog, sheets and a note at the bottom.
type overlaysState struct {
	dialogShown bool
	sheetShown  bool
	sideShown   bool
	placesShown bool
	places      material.Sheet
	onlyLit     bool
	onlyKept    bool
	snackbar    material.Snackbar
}

// draw shows what is over the screen, after the screen, so it is in front and answers the back
// first.
func (overlays *overlaysState) draw(program *gallery, look material.Look) {
	lights := program.lists.lighthouses
	if overlays.dialogShown {
		answer := look.Dialog("switch off", material.Dialog{Icon: "lightbulb-off", Title: "Switch every light off?",
			Text: "Ships will not see the coast until they are lit again.", Confirm: "Switch off", Dismiss: "Cancel"})
		if answer == material.Confirmed {
			for index := range lights {
				lights[index].lit = false
			}
			overlays.snackbar.Show("Every light is off", "Undo", after(6))
		}
		overlays.dialogShown = answer == material.Open
	}
	if overlays.sheetShown && look.BottomSheet("actions", func() {
		look.ListItem("share", material.Item{Icon: "share-2", Headline: "Share the chart"})
		look.ListItem("copy", material.Item{Icon: "copy", Headline: "Copy the position"})
		if look.ListItem("close", material.Item{Icon: "x", Headline: "Close"}) {
			overlays.sheetShown = false
		}
	}) {
		overlays.sheetShown = false
	}
	if overlays.sideShown && look.SideSheet("filters", "Filters", func() {
		look.Checkbox("Only lights that are lit", &overlays.onlyLit)
		look.Checkbox("Only lights with a keeper", &overlays.onlyKept)
		look.Context.Spacer()
		if look.FilledButton("", "Apply") {
			overlays.sideShown = false
		}
	}) {
		overlays.sideShown = false
	}
	if overlays.placesShown {
		overlays.places.Above = aboveBars(look)
		look.StandardSheet("places", &overlays.places, 220, func() {
			look.Context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.Padding{Left: 16, Right: 4}, AlignY: ui.Center}, func() {
				look.Context.Box(ui.Layout{Width: ui.Grow()}, func() {
					look.Text("Lighthouses near you", 20, draw2d.Regular, look.Roles.OnSurface)
				})
				if look.IconButton("x") {
					overlays.placesShown = false
				}
			})
			look.Context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.Padding{Left: 16, Right: 16, Bottom: 8}}, func() {
				look.Paragraph("Drag the handle up for all of them; the screen under the sheet still works.", 13, draw2d.Regular, look.Roles.OnSurfaceVariant)
			})
			for index, light := range lights {
				look.ListItem(light.name, material.Item{Picture: coast()[index%len(coast())], Headline: light.name, Supporting: light.place})
			}
		})
	}
	overlays.snackbar.Above = aboveBars(look)
	if look.Snackbar(&overlays.snackbar) {
		for index := range lights {
			lights[index].lit = true
		}
	}
}
