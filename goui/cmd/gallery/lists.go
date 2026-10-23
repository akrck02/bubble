package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/akrck02/bubble/goui/material"
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
	"github.com/akrck02/go-ui/window"
)

// listsState is what the lists screen remembers.
type listsState struct {
	lighthouses []lighthouse
	added       int
	refreshing  int // frames left of a pretend refresh
	answers     [3]bool
	bannerShown bool
	navigation  material.Navigation // the list, then a lighthouse's details
}

func newLists() listsState {
	return listsState{bannerShown: true, lighthouses: []lighthouse{
		{name: "North Cape", place: "Cliffs, north shore", lit: true},
		{name: "Harbour Light", place: "End of the pier", lit: true},
		{name: "Seal Rock", place: "Two miles out"},
		{name: "Old Tower", place: "The hill above the town"},
		{name: "East Point", place: "The lighthouse keeper's cottage", lit: true},
		{name: "Sand Bar", place: "On the bar at low tide"}}}
}

// darkCount is the navigation's badge: how many lighthouses are dark, none when all are lit.
func (lists *listsState) darkCount() string {
	dark := 0
	for _, light := range lists.lighthouses {
		if !light.lit {
			dark++
		}
	}
	if dark == 0 {
		return ""
	}
	return fmt.Sprint(dark)
}

// choosing tells whether lighthouses are being chosen: one was held, and touches choose more.
func (lists *listsState) choosing() bool { return lists.chosenCount() > 0 }

func (lists *listsState) chosenCount() int {
	count := 0
	for _, light := range lists.lighthouses {
		if light.chosen {
			count++
		}
	}
	return count
}

func (lists *listsState) chooseNone() {
	for index := range lists.lighthouses {
		lists.lighthouses[index].chosen = false
	}
}

// removeChosen takes the chosen lighthouses out of the list, and says how many went.
func (lists *listsState) removeChosen() int {
	before := len(lists.lighthouses)
	lists.lighthouses = slices.DeleteFunc(lists.lighthouses, func(light lighthouse) bool { return light.chosen })
	return before - len(lists.lighthouses)
}

// draw shows the list of lighthouses, or one's details, which slide in when its row is touched.
func (lists *listsState) draw(program *gallery, look material.Look) {
	look.Screens("lists", &lists.navigation, func(key string) {
		if key == "" {
			lists.drawList(program, look)
			return
		}
		lists.drawLighthouse(program, look, key)
	})
}

// drawList shows a list that refreshes when pulled, rows swiped away or held to choose them, a
// banner and rows that open.
func (lists *listsState) drawList(program *gallery, look material.Look) {
	context, scheme := look.Context, look.Roles
	notes := &program.overlays.snackbar
	choosing := lists.choosing()
	if choosing {
		lists.choiceBar(look, notes)
	} else {
		program.topBar(look, "Lists", "Pull to refresh, hold a row to choose it")
	}
	if lists.bannerShown {
		switch look.Banner("wifi-off", "The lights' positions were last brought up to date a week ago.", "Dismiss", "Refresh now") {
		case 0:
			lists.bannerShown = false
		case 1:
			lists.bannerShown, lists.refreshing = false, 90
		}
	}
	if lists.refreshing > 0 {
		lists.refreshing--
		context.RequestFrame(0)
		if lists.refreshing == 0 {
			notes.Show("Positions brought up to date", "", after(2))
		}
	}
	layout := ui.Layout{Width: ui.Grow(), Height: ui.Grow()}
	if look.PullToRefresh("lights", lists.refreshing > 0, layout, func() {
		for index := 0; index < len(lists.lighthouses); index++ {
			light := &lists.lighthouses[index]
			context.Push(light.name)
			picture := coast()[index%len(coast())]
			if look.SwipeToDismiss("swipe", func() { lists.row(look, light, picture, choosing) }) {
				notes.Show(light.name+" removed", "", after(3))
				lists.lighthouses = append(lists.lighthouses[:index], lists.lighthouses[index+1:]...)
				index--
			}
			context.Pop()
		}
		look.SectionTitle("People")
		for _, person := range []string{"Ana Lighthouse-Keeper", "Tom Boatman", "Lena Harbour"} {
			context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(56), Padding: ui.PadXY(16, 0), Gap: 16, AlignY: ui.Center}, func() {
				look.Avatar(person, "", draw2d.Color{}, 40)
				look.Text(person, 16, draw2d.Regular, scheme.OnSurface)
			})
		}
		look.SectionTitle("Questions")
		questions := []string{"When are the lights lit?", "Who keeps them?", "Can they be visited?"}
		answers := []string{"From sunset to sunrise, and in fog.", "Keepers live by three of them; the others run on their own.",
			"The Old Tower opens on summer weekends."}
		for index, question := range questions {
			look.ExpandableItem(fmt.Sprint("question ", index), material.Item{Headline: question}, &lists.answers[index], func() {
				look.Paragraph(answers[index], 14, draw2d.Regular, scheme.OnSurfaceVariant)
			})
		}
		context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(96)}, nil)
	}) {
		lists.refreshing = 90 // about a second and a half at sixty frames
	}
	context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.Pad(16), AlignX: ui.End}, func() {
		if look.ExtendedFloatingButton("plus", "New light") {
			lists.added++
			lists.lighthouses = append(lists.lighthouses, lighthouse{name: fmt.Sprint("New light ", lists.added), place: "Not placed yet"})
		}
	})
}

// row is a lighthouse's row: held, it is chosen; touched while choosing, it is chosen or let go.
func (lists *listsState) row(look material.Look, light *lighthouse, picture *draw2d.Image, choosing bool) {
	touched := look.ListItem("row", material.Item{Picture: picture, Headline: light.name, Supporting: light.place,
		Selected: light.chosen, Trailing: func() { look.Switch("lit", &light.lit) }})
	switch {
	case look.LongPressed():
		light.chosen = true
	case touched && choosing:
		light.chosen = !light.chosen
	case touched:
		lists.navigation.Go(light.name)
	}
}

// choiceBar takes the top bar's place while lighthouses are chosen: how many, a cross to let them
// go, and a bin to remove them.
func (lists *listsState) choiceBar(look material.Look, notes *material.Snackbar) {
	count := lists.chosenCount()
	bar := material.TopBar{Title: fmt.Sprint(count, " chosen"), Back: true, Close: true, Actions: func() {
		if look.NamedIconButton("share-2", "Share") && !window.Share("Lighthouses", lists.chosenNames()) {
			look.Context.Copy(lists.chosenNames()) // where there is no sheet to share through
			notes.Show("Copied: paste them where you like", "", after(3))
		}
		if look.NamedIconButton("copy", "Copy") {
			look.Context.Copy(lists.chosenNames())
			notes.Show("Copied to the clipboard", "", after(3))
		}
		if look.NamedIconButton("trash-2", "Remove") {
			notes.Show(lightsRemoved(lists.removeChosen()), "", after(3))
		}
	}}
	if look.TopAppBar(bar) {
		lists.chooseNone()
	}
}

// chosenNames is the chosen lighthouses, one a line with where they are, to share or copy.
func (lists *listsState) chosenNames() string {
	var lines []string
	for _, light := range lists.lighthouses {
		if light.chosen {
			lines = append(lines, light.name+" — "+light.place)
		}
	}
	return strings.Join(lines, "\n")
}

func lightsRemoved(count int) string {
	if count == 1 {
		return "1 lighthouse removed"
	}
	return fmt.Sprint(count, " lighthouses removed")
}

// drawLighthouse shows one lighthouse: its photo, where it is, whether it is lit, and ways to pass
// it on.
func (lists *listsState) drawLighthouse(program *gallery, look material.Look, name string) {
	context, scheme := look.Context, look.Roles
	index := slices.IndexFunc(lists.lighthouses, func(light lighthouse) bool { return light.name == name })
	if look.TopAppBar(material.TopBar{Title: name, Back: true}) {
		lists.navigation.Back()
	}
	if index < 0 {
		look.Paragraph("This lighthouse was taken out of the list.", 16, draw2d.Regular, scheme.OnSurfaceVariant)
		return
	}
	light := &lists.lighthouses[index]
	page(context, "light "+name, func() {
		context.Image(coast()[index%len(coast())], ui.Layout{Width: ui.Grow(), Height: ui.Fixed(220)}, 16)
		look.Text(light.place, 16, draw2d.Regular, scheme.OnSurfaceVariant)
		look.ListItem("lit", material.Item{Icon: "lightbulb", Headline: "Lit tonight", Supporting: "Its keeper switches it on at sunset",
			Trailing: func() { look.Switch("lit switch", &light.lit) }})
		look.Paragraph("Touch the arrow or go back with the phone's back to return to the list; this screen slid in over it, and slides back out.",
			14, draw2d.Regular, scheme.OnSurfaceVariant)
		row(context, func() {
			if look.FilledButton("share-2", "Share") && !window.Share(light.name, light.name+" — "+light.place) {
				context.Copy(light.name + " — " + light.place)
				program.overlays.snackbar.Show("Copied to the clipboard", "", after(3))
			}
			if look.OutlinedButton("copy", "Copy") {
				context.Copy(light.name + " — " + light.place)
				program.overlays.snackbar.Show("Copied to the clipboard", "", after(3))
			}
		})
	})
}
