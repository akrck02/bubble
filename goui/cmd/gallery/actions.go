package main

import (
	"github.com/akrck02/bubble/goui/material"
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// actionsState is what the actions screen remembers.
type actionsState struct {
	size   int
	create material.FloatingMenu
}

// draw shows the buttons, from the most to the least important, Material 3 Expressive's split
// button and button group, and a toolbar and a floating menu over the screen.
func (actions *actionsState) draw(program *gallery, look material.Look) {
	context, roles := look.Context, look.Roles
	notes := &program.overlays.snackbar
	program.largeTopBar(look, "Actions", "Buttons, floating buttons and toolbars", "actions")
	page(context, "actions", func() {
		look.Heading("Buttons")
		row(context, func() {
			look.FilledButton("plus", "Filled")
			look.TonalButton("", "Tonal")
			look.OutlinedButton("", "Outlined")
			look.TextButton("", "Text")
		})
		row(context, func() {
			look.NamedIconButton("heart", "Like")
			look.NamedIconButton("bookmark", "Keep")
			look.RoundButton("map", 48, false)
			look.Context.Describe(ui.Description{Label: "Show the map"}) // what a screen reader calls it
			look.RoundButton("locate-fixed", 48, true)
			look.Context.Describe(ui.Description{Label: "Where am I"})
		})
		look.Heading("Split button and group")
		if action, more := look.SplitButton("send", "Send the report"); action {
			notes.Show("Report sent", "", after(3))
		} else if more {
			notes.Show("Later: schedule it, or send it to someone else", "", after(3))
		}
		look.ButtonGroup("size", []string{"Small", "Medium", "Large"}, &actions.size)
		look.Heading("Assist chips")
		row(context, func() {
			if look.AssistChip("calendar-plus", "Add to calendar") {
				notes.Show("Added to the calendar", "", after(2))
			}
			look.AssistChip("share-2", "Share")
			look.AssistChip("navigation", "Directions")
		})
		look.Heading("Over the screen")
		look.Paragraph("A toolbar floats at the foot for what is done to the content, and the floating button opens into a few actions.",
			14, draw2d.Regular, roles.OnSurfaceVariant)
	})
	above := aboveBars(look)
	tools := []material.Destination{{Icon: "pencil", Label: "Edit"}, {Icon: "share-2", Label: "Share"}, {Icon: "trash-2", Label: "Delete"}}
	switch look.FloatingToolbar("tools", tools, above) {
	case 0:
		notes.Show("Edit", "", after(2))
	case 1:
		notes.Show("Share", "", after(2))
	case 2:
		notes.Show("Deleted", "Undo", after(4))
	}
	actions.create.Icon, actions.create.Label, actions.create.Above = "plus", "Create", above
	actions.create.Actions = []material.Destination{{Icon: "file", Label: "Note"}, {Icon: "image", Label: "Picture"}, {Icon: "clock", Label: "Reminder"}}
	if chosen := look.FloatingMenu("create", &actions.create); chosen >= 0 {
		notes.Show("New "+actions.create.Actions[chosen].Label, "", after(2))
	}
}
