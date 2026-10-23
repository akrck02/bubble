package main

import (
	"fmt"

	"github.com/akrck02/bubble/goui/material"
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
	"github.com/akrck02/go-ui/window"
)

// feedbackState is what the feedback screen remembers.
type feedbackState struct {
	progress float32
	loaded   bool
}

// draw shows what tells the person something is happening: loading indicators and progress,
// skeletons, and what goes over the screen (dialogs, sheets, notes and notifications).
func (feedback *feedbackState) draw(program *gallery, look material.Look) {
	context, roles := look.Context, look.Roles
	overlays := &program.overlays
	program.largeTopBar(look, "Feedback", "Loading, progress, dialogs and notes", "feedback")
	page(context, "feedback", func() {
		look.Heading("Loading")
		row(context, func() {
			look.LoadingIndicator(48)
			look.ContainedLoadingIndicator(48)
			look.CircularProgress(material.Indeterminate, 48, false)
			look.CircularProgress(material.Indeterminate, 48, true)
			look.CircularProgress(feedback.progress, 48, true)
		})
		look.WavyProgress(material.Indeterminate)
		look.WavyProgress(feedback.progress)
		look.LinearProgress(feedback.progress)
		look.Slider("progress", &feedback.progress, 0, 1, 0)
		look.Heading("Skeletons")
		look.ListItem("loaded", material.Item{Headline: "Loaded", Trailing: func() { look.Switch("loaded switch", &feedback.loaded) }})
		for index := range 3 {
			if feedback.loaded {
				look.ListItem(fmt.Sprint("story ", index), material.Item{Icon: "newspaper", Headline: fmt.Sprint("Story ", index+1), Supporting: "Arrived"})
				continue
			}
			context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(72), Padding: ui.PadXY(16, 0), Gap: 16, AlignY: ui.Center}, func() {
				look.Skeleton(ui.Fixed(40), ui.Fixed(40), 20)
				context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Gap: 8}, func() {
					look.Skeleton(ui.Percent(0.6), ui.Fixed(14), 7)
					look.Skeleton(ui.Percent(0.4), ui.Fixed(12), 6)
				})
			})
		}
		look.Heading("Over the screen")
		row(context, func() {
			if look.TonalButton("message-square", "Dialog") {
				overlays.dialogShown = true
			}
			if look.TonalButton("panel-bottom", "Sheet") {
				overlays.sheetShown = true
			}
			if look.TonalButton("panel-right", "Side sheet") {
				overlays.sideShown = true
			}
			if look.TonalButton("map", "Standard sheet") {
				overlays.placesShown, overlays.places.Expanded = true, false
			}
			if look.TonalButton("info", "Snackbar") {
				overlays.snackbar.Show("Harbour Light switched off", "Undo", after(6))
			}
			if look.TonalButton("bell", "Notification") {
				overlays.notify()
			}
		})
		look.Paragraph("Notifications show in the phone's shade, even with the gallery in the background.", 14, draw2d.Regular, roles.OnSurfaceVariant)
	})
}

// notify sends a notification, asking first where the phone wants it allowed.
func (overlays *overlaysState) notify() {
	if !window.NotificationsAllowed() {
		window.AskForNotifications()
		overlays.snackbar.Show("Allow notifications, then touch again", "", after(4))
		return
	}
	if window.Notify(1, "Harbour Light", "Lit for the night, reaching 18 miles") {
		overlays.snackbar.Show("Notification sent: look at the shade", "", after(4))
		return
	}
	overlays.snackbar.Show("Notifications are not shown here", "", after(4))
}
