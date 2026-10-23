package material

import (
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// AssistChip suggests an action related to what is on screen (add to the calendar, share), with
// an icon before its label. It returns true when touched.
func (look Look) AssistChip(icon, label string) bool {
	scheme := look.Roles
	normal := ui.Look{Border: scheme.OutlineVariant, BorderWidth: 1, Radius: 8}
	layout := ui.Layout{Height: ui.Fixed(32), Padding: ui.Padding{Left: 8, Right: 16}, Gap: 8, AlignY: ui.Center}
	return look.Context.Clickable("assist "+label, layout, pressedOnly(normal, stateLayer(scheme.OnSurface, 0.1)), func() {
		if icon != "" {
			look.Context.IconSized(icon, 18, scheme.Primary)
		}
		look.Text(label, 14, draw2d.Medium, scheme.OnSurface)
	})
}

// InputChip is something the person entered or chose (a recipient, a tag), with a cross to take it
// out. It tells whether the chip was touched and whether its cross was.
func (look Look) InputChip(label string) (touched, removed bool) {
	scheme := look.Roles
	look.Context.Push("input " + label)
	defer look.Context.Pop()
	frame := ui.Layout{Height: ui.Fixed(32), Padding: ui.Padding{Left: 12, Right: 4}, Gap: 4, AlignY: ui.Center}
	look.Context.Frame(frame, ui.Look{Border: scheme.OutlineVariant, BorderWidth: 1, Radius: 8}, func() {
		layout := ui.Layout{Height: ui.Grow(), AlignY: ui.Center}
		touched = look.Context.Clickable("label", layout, ui.Looks{}, func() {
			look.Text(label, 14, draw2d.Medium, scheme.OnSurfaceVariant)
		})
		cross := ui.Layout{Width: ui.Fixed(24), Height: ui.Fixed(24), AlignX: ui.Center, AlignY: ui.Center}
		removed = look.Context.Clickable("remove", cross, ui.Looks{Pressed: ui.Look{Fill: stateLayer(scheme.OnSurface, 0.12), Radius: 12}}, func() {
			look.Context.IconSized("x", 16, scheme.OnSurfaceVariant)
		})
	})
	return touched, removed
}

// FilterChip is a choice that is on or off, among others in a row: outlined, or filled with a tick
// when selected. It returns true when touched; the program flips its own value.
func (look Look) FilterChip(label string, selected bool) bool {
	scheme := look.Roles
	normal, ink := ui.Look{Border: scheme.OutlineVariant, BorderWidth: 1, Radius: 8}, scheme.OnSurfaceVariant
	if selected {
		normal, ink = ui.Look{Fill: scheme.SecondaryContainer, Radius: 8}, scheme.OnSecondaryContainer
	}
	layout := ui.Layout{Height: ui.Fixed(32), Padding: ui.PadXY(12, 0), Gap: 8, AlignY: ui.Center}
	return look.Context.Clickable("chip "+label, layout, pressedOnly(normal, stateLayer(ink, 0.12)), func() {
		if selected {
			look.Context.IconSized("check", 18, ink)
		}
		look.Text(label, 14, draw2d.Medium, ink)
	})
}

// Tag is a label that cannot be touched, in a pill of its own colours: a state ("3 running"), a
// count, a kind. ink is the text's colour and tint the pill's.
func (look Look) Tag(text string, ink, tint draw2d.Color) {
	look.Context.Frame(ui.Layout{Height: ui.Fixed(26), Padding: ui.PadXY(10, 0), AlignY: ui.Center},
		ui.Look{Fill: tint, Radius: 13}, func() { look.Text(text, 12, draw2d.Bold, ink) })
}
