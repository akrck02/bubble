package material

import (
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
	"strings"
	"unicode"
)

// Card is a filled container for one subject: a few lines, a picture, its actions.
func (look Look) Card(content func()) {
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.Pad(16), Gap: 12}
	look.Context.Frame(layout, ui.Look{Fill: look.Roles.SurfaceContainerHighest, Radius: 12}, content)
}

// ElevatedCard is a card raised over what is beneath it: over a map or a scene.
func (look Look) ElevatedCard(content func()) {
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.Pad(16), Gap: 12}
	look.Context.Frame(layout, ui.Look{Fill: look.Roles.SurfaceContainerLow, Radius: 12, Shadow: true}, content)
}

// OutlinedCard is a card drawn by its border, among others on the surface.
func (look Look) OutlinedCard(content func()) {
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.Pad(16), Gap: 12}
	look.Context.Frame(layout, ui.Look{Fill: look.Roles.Surface, Border: look.Roles.OutlineVariant, BorderWidth: 1, Radius: 12}, content)
}

// Item is what a list item shows.
type Item struct {
	Icon string       // in a round avatar at the start; none when empty
	Mark draw2d.Color // the avatar's colour; the primary container when zero
	// Picture is a photo in the round avatar instead of the icon (a person, a place).
	Picture    *draw2d.Image
	Headline   string
	Supporting string // the line under the headline; none when empty
	// Trailing draws what goes at the end (a Switch, an IconButton, a Tag); none when nil.
	Trailing func()
	// Selected marks the row chosen, in a list whose rows are chosen with a long press: it is tinted,
	// and a tick takes the avatar's place.
	Selected bool
}

// ListItem is a row of a list, 72 pixels tall with a supporting line and 56 without. It returns
// true when the row is touched. name tells it apart from the other rows.
func (look Look) ListItem(name string, item Item) bool {
	roles := look.Roles
	height := float32(56)
	if item.Supporting != "" {
		height = 72
	}
	layout := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(height), Padding: ui.PadXY(16, 0), Gap: 16, AlignY: ui.Center}
	looks := ui.Looks{Pressed: ui.Look{Fill: stateLayer(roles.OnSurface, 0.1)}}
	if item.Selected {
		chosen := ui.Look{Fill: roles.SecondaryContainer}
		looks = ui.Looks{Normal: chosen, Hovered: chosen, Pressed: chosen}
	}
	return look.Context.Clickable(name, layout, looks, func() {
		switch {
		case item.Selected:
			look.tick()
		case item.Picture != nil:
			look.PictureAvatar(item.Picture, 40)
		case item.Icon != "":
			look.avatar(item)
		}
		look.Context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Gap: 2}, func() {
			look.Text(item.Headline, 16, draw2d.Regular, roles.OnSurface)
			if item.Supporting != "" {
				look.Text(item.Supporting, 14, draw2d.Regular, roles.OnSurfaceVariant)
			}
		})
		if item.Trailing != nil {
			item.Trailing()
		}
	})
}

func (look Look) avatar(item Item) {
	fill, ink := look.Roles.PrimaryContainer, look.Roles.OnPrimaryContainer
	if item.Mark != (draw2d.Color{}) {
		fill, ink = item.Mark, readableOn(item.Mark)
	}
	layout := ui.Layout{Width: ui.Fixed(40), Height: ui.Fixed(40), AlignX: ui.Center, AlignY: ui.Center}
	look.Context.Frame(layout, ui.Look{Fill: fill, Radius: 20}, func() { look.Context.IconSized(item.Icon, 20, ink) })
}

// tick stands in a chosen row's avatar: a tick in a circle of the primary colour.
func (look Look) tick() {
	layout := ui.Layout{Width: ui.Fixed(40), Height: ui.Fixed(40), AlignX: ui.Center, AlignY: ui.Center}
	look.Context.Frame(layout, ui.Look{Fill: look.Roles.Primary, Radius: 20}, func() {
		look.Context.IconSized("check", 20, look.Roles.OnPrimary)
	})
}

// readableOn is black or white, whichever reads better on a colour.
func readableOn(colour draw2d.Color) draw2d.Color {
	lightness := 0.299*colour[0] + 0.587*colour[1] + 0.114*colour[2]
	if lightness > 0.6 {
		return draw2d.Color{0, 0, 0, 0.87}
	}
	return draw2d.Color{1, 1, 1, 1}
}

// SectionTitle is a small heading over a group of rows, in the primary colour.
func (look Look) SectionTitle(text string) {
	look.Context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.Padding{Left: 16, Top: 16, Right: 16, Bottom: 8}}, func() {
		look.Text(text, 14, draw2d.Medium, look.Roles.Primary)
	})
}

// Heading is a section's small title in the primary colour, as SectionTitle, for content that
// keeps its own margins (a page of cards and buttons) rather than a list across the screen.
func (look Look) Heading(text string) {
	look.Context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.Padding{Top: 12, Bottom: 4}}, func() {
		look.Text(text, 14, draw2d.Medium, look.Roles.Primary)
	})
}

// Divider is a thin line across, between groups.
func (look Look) Divider() {
	look.Context.Frame(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(1)}, ui.Look{Fill: look.Roles.OutlineVariant}, nil)
}

// ExpandableItem is a list row that opens to show more under it (an answer under its question, the
// details of an order), its chevron turning. Touching the row opens or closes it; it returns true
// when touched. name tells rows apart.
func (look Look) ExpandableItem(name string, item Item, open *bool, more func()) bool {
	roles := look.Roles
	icon := "chevron-down"
	if *open {
		icon = "chevron-up"
	}
	item.Trailing = func() { look.Context.IconSized(icon, 24, roles.OnSurfaceVariant) }
	touched := look.ListItem(name, item)
	if touched {
		*open = !*open
	}
	// it unfolds under the row, and folds back when closed
	duration := arrivingMotion
	if !*open {
		duration = leavingMotion
	}
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.Padding{Left: 16, Right: 16, Bottom: 16}, Gap: 8}
	look.Context.Fold(name+" more", *open, duration, layout, more)
	return touched
}

// Avatar is a person or a thing shown in a circle: their initials (up to two letters of name) on a
// colour, or an icon when name is empty. size is its side, 40 in lists.
func (look Look) Avatar(name, icon string, tint draw2d.Color, size float32) {
	if tint == (draw2d.Color{}) {
		tint = look.Roles.PrimaryContainer
	}
	ink := readableOn(tint)
	layout := ui.Layout{Width: ui.Fixed(size), Height: ui.Fixed(size), AlignX: ui.Center, AlignY: ui.Center}
	look.Context.Frame(layout, ui.Look{Fill: tint, Radius: size / 2}, func() {
		if initials := initialsOf(name); initials != "" {
			look.Text(initials, size*0.4, draw2d.Medium, ink)
			return
		}
		look.Context.IconSized(icon, size*0.5, ink)
	})
}

// PictureAvatar is a person or a place shown by a photo in a circle. size is its side, 40 in lists.
func (look Look) PictureAvatar(picture *draw2d.Image, size float32) {
	look.Context.Image(picture, ui.Layout{Width: ui.Fixed(size), Height: ui.Fixed(size)}, size/2)
}

// initialsOf is the first letters of a name's first two words, in capitals.
func initialsOf(name string) string {
	var initials []rune
	for _, word := range strings.Fields(name) {
		for _, letter := range word {
			initials = append(initials, unicode.ToUpper(letter))
			break
		}
		if len(initials) == 2 {
			break
		}
	}
	return string(initials)
}

// Banner is a message across the top of the content that stays until acted on (no connection, an
// update to install), with an icon and one or two actions under its text. It returns the index of
// the action touched in this frame, or -1.
func (look Look) Banner(icon, text string, actions ...string) int {
	roles := look.Roles
	touched := -1
	look.Context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow()}, func() {
		look.Context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.Padding{Left: 16, Top: 24, Right: 16, Bottom: 8}, Gap: 16}, func() {
			if icon != "" {
				look.Avatar("", icon, roles.Primary, 40)
			}
			look.Context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.Padding{Top: 10}}, func() {
				look.Paragraph(text, 14, draw2d.Regular, roles.OnSurface)
			})
		})
		look.Context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.Padding{Right: 8, Bottom: 8}, Gap: 8, AlignX: ui.End}, func() {
			for index, action := range actions {
				if look.TextButton("", action) {
					touched = index
				}
			}
		})
		look.Divider()
	})
	return touched
}
