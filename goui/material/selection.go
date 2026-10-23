package material

import (
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// Slider chooses a value between two ends by dragging, in steps (0 for any value): go-ui's slider
// in the scheme's colours and a finger's size. It returns true when the value changed.
func (look Look) Slider(name string, value *float32, minimum, maximum, step float32) bool {
	context := look.Context
	saved := context.Theme
	defer func() { context.Theme = saved }()
	context.Theme.ControlLine, context.Theme.Sunken = look.Roles.SecondaryContainer, look.Roles.SecondaryContainer
	if step <= 0 {
		return context.Slider(name, value, minimum, maximum)
	}
	return context.SliderStep(name, value, minimum, maximum, step)
}

// Switch is Material's switch: a 52 × 32 track whose handle grows and moves to the end when on. It
// flips the value and returns true when touched.
func (look Look) Switch(name string, on *bool) bool {
	scheme := look.Roles
	// the handle slides across and grows as the track fills with colour
	progress := look.Context.Animate(name, boolFloat(*on), switchMotion)
	track := ui.Look{Fill: mix(scheme.SurfaceContainerHighest, scheme.Primary, progress),
		Border: stateLayer(scheme.Outline, 1-progress), BorderWidth: 2, Radius: 16}
	handle, handleSize := mix(scheme.Outline, scheme.OnPrimary, progress), 16+8*progress
	travel := 52 - 8 - handleSize
	layout := ui.Layout{Width: ui.Fixed(52), Height: ui.Fixed(32), Padding: ui.Padding{Left: 4 + travel*progress, Right: 4},
		AlignX: ui.Start, AlignY: ui.Center}
	touched := look.Context.Clickable(name, layout, ui.Looks{Normal: track, Hovered: track, Pressed: track}, func() {
		size := ui.Fixed(handleSize)
		look.Context.Frame(ui.Layout{Width: size, Height: size}, ui.Look{Fill: handle, Radius: handleSize / 2}, nil)
	})
	if touched {
		*on = !*on
	}
	look.Context.Describe(ui.Description{Role: ui.SwitchRole, Label: shownName(name), Checked: *on})
	return touched
}

// SegmentedButtons are two to five choices side by side sharing the row's width, one of them
// chosen (a view, a speed, a period). It changes selected and returns true when another is
// touched. name tells the group apart from another one on the same screen.
func (look Look) SegmentedButtons(name string, options []string, selected *int) bool {
	scheme := look.Roles
	changed := false
	look.Context.Push(name)
	defer look.Context.Pop()
	look.Context.Box(ui.Layout{Width: ui.Grow(), Gap: 8}, func() {
		for index, option := range options {
			chosen := index == *selected
			normal, ink := ui.Look{Border: scheme.Outline, BorderWidth: 1, Radius: 20}, scheme.OnSurface
			if chosen {
				normal, ink = ui.Look{Fill: scheme.SecondaryContainer, Border: scheme.Outline, BorderWidth: 1, Radius: 20}, scheme.OnSecondaryContainer
			}
			layout := ui.Layout{Width: ui.Share(1), Height: ui.Fixed(40), Gap: 6, AlignX: ui.Center, AlignY: ui.Center}
			touched := look.Context.Clickable(option, layout, pressedOnly(normal, stateLayer(ink, 0.12)), func() {
				if chosen {
					look.Context.IconSized("check", 18, ink)
				}
				look.Text(option, 14, draw2d.Medium, ink)
			})
			look.Context.Describe(ui.Description{Selected: chosen})
			if touched && !chosen {
				*selected, changed = index, true
			}
		}
	})
	return changed
}

// Checkbox is a square that is ticked or not, with a label after it; touching either flips it. It
// returns true when touched.
func (look Look) Checkbox(label string, on *bool) bool {
	scheme := look.Roles
	layout := ui.Layout{Height: ui.Fixed(48), Padding: ui.Padding{Right: 8}, Gap: 4, AlignY: ui.Center}
	touched := look.Context.Clickable("check "+label, layout, ui.Looks{Pressed: ui.Look{Fill: stateLayer(scheme.OnSurface, 0.08), Radius: 24}}, func() {
		look.Context.Box(ui.Layout{Width: ui.Fixed(48), Height: ui.Fixed(48), AlignX: ui.Center, AlignY: ui.Center}, func() {
			box := ui.Layout{Width: ui.Fixed(18), Height: ui.Fixed(18), AlignX: ui.Center, AlignY: ui.Center}
			if *on {
				look.Context.Frame(box, ui.Look{Fill: scheme.Primary, Radius: 2}, func() {
					look.Context.IconSized("check", 14, scheme.OnPrimary)
				})
				return
			}
			look.Context.Frame(box, ui.Look{Border: scheme.OnSurfaceVariant, BorderWidth: 2, Radius: 2}, nil)
		})
		look.Text(label, 16, draw2d.Regular, scheme.OnSurface)
	})
	if touched {
		*on = !*on
	}
	look.Context.Describe(ui.Description{Role: ui.CheckboxRole, Checked: *on})
	return touched
}

// RadioButton is one of a few choices of which one is chosen, with a label after it. It returns
// true when touched; the program then chooses it.
func (look Look) RadioButton(label string, chosen bool) bool {
	scheme := look.Roles
	layout := ui.Layout{Height: ui.Fixed(48), Padding: ui.Padding{Right: 8}, Gap: 4, AlignY: ui.Center}
	touched := look.Context.Clickable("radio "+label, layout, ui.Looks{Pressed: ui.Look{Fill: stateLayer(scheme.OnSurface, 0.08), Radius: 24}}, func() {
		look.Context.Box(ui.Layout{Width: ui.Fixed(48), Height: ui.Fixed(48), AlignX: ui.Center, AlignY: ui.Center}, func() {
			ring := ui.Layout{Width: ui.Fixed(20), Height: ui.Fixed(20), AlignX: ui.Center, AlignY: ui.Center}
			colour := scheme.OnSurfaceVariant
			if chosen {
				colour = scheme.Primary
			}
			look.Context.Frame(ring, ui.Look{Border: colour, BorderWidth: 2, Radius: 10}, func() {
				if chosen {
					look.Context.Frame(ui.Layout{Width: ui.Fixed(10), Height: ui.Fixed(10)}, ui.Look{Fill: colour, Radius: 5}, nil)
				}
			})
		})
		look.Text(label, 16, draw2d.Regular, scheme.OnSurface)
	})
	look.Context.Describe(ui.Description{Role: ui.RadioRole, Checked: chosen})
	return touched
}
