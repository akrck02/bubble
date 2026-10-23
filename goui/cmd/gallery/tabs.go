package main

import (
	"time"

	"github.com/akrck02/bubble/goui/material"
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// tabsState is what the tabs screen remembers.
type tabsState struct {
	pager    material.Pager
	carousel material.Carousel
	coast    material.Carousel
	stick    material.Stick
	boatX    float32 // the boat's place in the bay, from 0 to 1 across and down
	boatY    float32
	lastStep time.Time
}

func newTabs() tabsState { return tabsState{boatX: 0.5, boatY: 0.5, lastStep: time.Now()} }

// draw shows tabs over pages a finger swipes from side to side: a carousel, cards, and a stick
// that sails a boat.
func (tabs *tabsState) draw(program *gallery, look material.Look) {
	program.topBar(look, "Tabs", "Swipe from page to page")
	look.TabPager("pages", []string{"Places", "Cards", "Stick"}, &tabs.pager, func(index int) {
		switch index {
		case 0:
			tabs.places(look)
		case 1:
			tabs.cards(look)
		default:
			tabs.bay(look)
		}
	})
}

// places is a carousel of the island's places, over a list that still scrolls up and down.
func (tabs *tabsState) places(look material.Look) {
	context, roles := look.Context, look.Roles
	page(context, "places", func() {
		look.Heading("Places to see")
		tints := []draw2d.Color{roles.PrimaryContainer, roles.SecondaryContainer, roles.ErrorContainer, roles.SurfaceContainerHighest}
		names := []string{"North Cape", "The harbour", "Seal Rock", "Old Tower", "Sand Bar", "East Point"}
		icons := []string{"mountain", "anchor", "fish", "castle", "waves", "sunrise"}
		look.Carousel("places", &tabs.carousel, len(names), 160, func(index int) {
			look.CarouselItem(200, tints[index%len(tints)], icons[index], names[index])
		})
		look.Paragraph("Slide the row sideways: it glides on and settles at a place. Up and down, the page scrolls.",
			14, draw2d.Regular, roles.OnSurfaceVariant)
		look.Heading("Browse the coast")
		look.MultiBrowseCarousel("coast", &tabs.coast, len(names), 200, 180, func(index int, width float32) {
			title := names[index]
			if width < 120 {
				title = "" // a narrow item shows its photo only
			}
			look.CarouselPicture(coast()[index%len(coast())], 180, title)
		})
	})
}

// cards shows the three kinds of cards and tags.
func (tabs *tabsState) cards(look material.Look) {
	context, roles := look.Context, look.Roles
	page(context, "cards", func() {
		look.OutlinedCard(func() {
			context.Image(coast()[3], ui.Layout{Width: ui.Grow(), Height: ui.Fixed(160)}, 12)
			look.Text("The north light at night", 18, draw2d.Medium, roles.OnSurface)
			look.Paragraph("Its lamp turns every ten seconds, seen eighteen miles out.", 14, draw2d.Regular, roles.OnSurfaceVariant)
		})
		look.Card(func() {
			look.Text("Tonight's watch", 18, draw2d.Medium, roles.OnSurface)
			look.Paragraph("Two lighthouses lit, two dark. The wind turns west after midnight.", 14, draw2d.Regular, roles.OnSurfaceVariant)
			context.Box(ui.Layout{Gap: 8}, func() {
				look.Tag("2 lit", roles.OnPrimaryContainer, roles.PrimaryContainer)
				look.Tag("2 dark", roles.OnErrorContainer, roles.ErrorContainer)
			})
		})
		look.ElevatedCard(func() { look.Text("An elevated card", 16, draw2d.Medium, roles.OnSurface) })
		look.OutlinedCard(func() { look.Text("An outlined card", 16, draw2d.Medium, roles.OnSurface) })
	})
}

// bay sails a boat with the stick: the boat keeps moving while the stick is held.
func (tabs *tabsState) bay(look material.Look) {
	context, roles := look.Context, look.Roles
	now := time.Now()
	step := float32(now.Sub(tabs.lastStep).Seconds())
	tabs.lastStep = now
	tabs.boatX = clamp(tabs.boatX+tabs.stick.Sideways*step*0.4, 0, 1)
	tabs.boatY = clamp(tabs.boatY-tabs.stick.Forward*step*0.4, 0, 1)
	wide := look.Wide()
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow(), Padding: ui.Pad(16), Gap: 16}
	if wide {
		layout.Direction, layout.AlignY = ui.Horizontal, ui.End
	}
	context.Box(layout, func() {
		width, height := context.Size()
		size := max(min(width-32, height-360), 64)
		boat := float32(28)
		padding := ui.Padding{Left: tabs.boatX * (size - boat), Top: tabs.boatY * (size - boat)}
		context.Frame(ui.Layout{Width: ui.Fixed(size), Height: ui.Fixed(size), Padding: padding},
			ui.Look{Fill: roles.PrimaryContainer, Radius: 16}, func() {
				context.Frame(ui.Layout{Width: ui.Fixed(boat), Height: ui.Fixed(boat), AlignX: ui.Center, AlignY: ui.Center},
					ui.Look{Fill: roles.Primary, Radius: boat / 2}, func() { context.IconSized("sailboat", 18, roles.OnPrimary) })
			})
		context.Spacer()
		look.Stick("stick", &tabs.stick, 128)
	})
}

func clamp(value, minimum, maximum float32) float32 { return min(max(value, minimum), maximum) }

// uiRow is a row whose children sit side by side with a gap, centred up and down.
func uiRow(gap float32) ui.Layout {
	return ui.Layout{Width: ui.Grow(), Gap: gap, AlignY: ui.Center}
}
