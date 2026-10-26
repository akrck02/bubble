package goui

import (
	"strings"

	"github.com/akrck02/bubble"
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// How far a tonal button's fill moves towards its ink under the pointer, and twice that while it
// is held, as Bubble's CSS draws .button.tonal:hover and :active: 15 % on paper, 12 % in Material.
const (
	tonalDeeper         = 0.15
	materialTonalDeeper = 0.12
)

// Tonal is how a tonal button is drawn: its fill and the ink of its label and icon.
type Tonal struct {
	Fill, Ink draw2d.Color
}

// TonalLook is the tonal button of a Bubble theme: the accent's pale tint with its label in the
// accent's ink, as Bubble's CSS draws .button.tonal. Material gives it its own container.
func TonalLook(theme bubble.Theme) Tonal {
	colours := theme.Colors
	if theme.Options.Style == bubble.StyleMaterial {
		return Tonal{Fill: colour(colours.AccentTint), Ink: colour(colours.OnAccentTint)}
	}
	return Tonal{Fill: colour(colours.AccentTint), Ink: colour(colours.AccentInk)}
}

// TonalButton is a button for an action repeated in a list or a grid (sell one, add a tag):
// quieter than the main button and louder than an outline, with no line. It is a finger tall and
// reports a click like go-ui's buttons. Its label follows go-ui's: "Sell one##3" shows "Sell one"
// and tells the third row's button from the others.
func TonalButton(context *ui.Context, theme bubble.Theme, label string) bool {
	return TonalIconButton(context, theme, "", label)
}

// TonalIconButton is TonalButton with an icon (a Lucide name) before its label. A label with no
// shown words, "##Add a tag", gives a square with the icon alone, read aloud by the hidden name.
func TonalIconButton(context *ui.Context, theme bubble.Theme, icon, label string) bool {
	tonal, settings := TonalLook(theme), context.Theme
	height := max(settings.ControlHeight, float32(theme.Space.TouchTarget))
	layout := ui.Layout{Height: ui.Fixed(height), Padding: ui.PadXY(settings.ButtonPadding, 0), Gap: settings.ButtonGap, AlignX: ui.Center, AlignY: ui.Center}
	shown, hidden, _ := strings.Cut(label, "##")
	if shown == "" {
		layout.Width, layout.Padding = ui.Fixed(height), ui.Padding{}
	}
	normal := ui.Look{Fill: tonal.Fill, Radius: settings.ButtonRadius}
	deeper := float32(tonalDeeper)
	if theme.Options.Style == bubble.StyleMaterial {
		deeper = materialTonalDeeper
	}
	hovered, pressed := normal, normal
	hovered.Fill = tonal.Fill.Mix(tonal.Ink, deeper)
	pressed.Fill = tonal.Fill.Mix(tonal.Ink, 2*deeper)
	name := label
	if name == "" {
		name = icon
	}
	clicked := context.Clickable(name, layout, ui.Looks{Normal: normal, Hovered: hovered, Pressed: pressed}, func() {
		if icon != "" {
			context.IconSized(icon, settings.IconSize, tonal.Ink)
		}
		if shown != "" {
			context.Text(shown, ui.TextStyle{Weight: settings.ButtonWeight, Color: tonal.Ink})
		}
	})
	if shown == "" && hidden != "" {
		context.Describe(ui.Description{Role: ui.ButtonRole, Label: hidden})
	}
	return clicked
}
