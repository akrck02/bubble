package material

import (
	"fmt"
	"unicode/utf8"

	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// Field is what a text field says around what is typed.
type Field struct {
	Label       string // above the box
	Placeholder string // in the box, quiet, while it is empty
	Supporting  string // under the box: a hint
	Error       string // under the box instead of the hint, in red, with the box red too
	MaxLength   int    // the most characters it takes; 0 = no limit
	// Counter shows how many characters there are of MaxLength, under the box at its end.
	Counter bool
	// Filled draws Material's filled field, a tinted box with a line under it, rather than the
	// outlined one: for a form on a plain screen, where boxes drawn by their border look empty.
	Filled bool
	// Prefix and Suffix are written in the box before and after what is typed: a currency, a unit.
	Prefix, Suffix string
	// Secret shows a dot for each character (a password). With Reveal, an eye at the box's end
	// shows the characters while *Reveal is true, and touching it turns *Reveal over.
	Secret bool
	Reveal *bool
}

// TextField is Material's text field: a label, the box typed in (outlined, or filled), and a hint
// or an error and a count of characters under it. On a phone the on-screen keyboard comes up while
// it is typed in. It returns true when the text changed. name tells fields apart.
func (look Look) TextField(name string, text *string, field Field) bool {
	scheme := look.Roles
	words := look.words()
	look.Context.SetEditingWords(ui.EditingWords{Cut: words.Cut, Copy: words.Copy, Paste: words.Paste, SelectAll: words.SelectAll})
	changed := false
	focused := look.Context.FieldFocused(name)
	look.Context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Gap: 4}, func() {
		if field.Label != "" {
			look.Text(field.Label, 12, draw2d.Medium, look.fieldColour(field, focused, scheme.OnSurfaceVariant))
			look.Context.Describe(ui.Description{Hidden: true}) // read with the field, which it names
		}
		look.fieldBox(field, focused, func() {
			if field.Prefix != "" {
				look.Text(field.Prefix, 16, draw2d.Regular, scheme.OnSurfaceVariant)
			}
			hidden := field.Secret && (field.Reveal == nil || !*field.Reveal)
			changed = look.Context.TextFieldWith(name, text, ui.TextFieldOptions{Placeholder: field.Placeholder,
				MaxLength: field.MaxLength, Plain: true, Secret: hidden})
			// the label over the box names the field; a screen reader reads it with what is typed
			look.Context.Describe(ui.Description{Label: field.Label})
			if field.Suffix != "" {
				look.Text(field.Suffix, 16, draw2d.Regular, scheme.OnSurfaceVariant)
			}
			if field.Secret && field.Reveal != nil {
				icon := map[bool]string{false: "eye", true: "eye-off"}[*field.Reveal]
				word := map[bool]string{false: look.words().ShowPassword, true: look.words().HidePassword}[*field.Reveal]
				if look.NamedIconButton(icon, word) {
					*field.Reveal = !*field.Reveal
				}
			}
		})
		look.fieldFoot(field, *text)
	})
	return changed
}

// fieldColour is a field's colour for its label and its border: red with an error, the main
// colour while typed in, else the one given.
func (look Look) fieldColour(field Field, focused bool, otherwise draw2d.Color) draw2d.Color {
	switch {
	case field.Error != "":
		return look.Roles.Error
	case focused:
		return look.Roles.Primary
	}
	return otherwise
}

// fieldBox draws the box a field is typed in, 56 pixels tall: outlined by a border that thickens
// while typed in, or filled with a line under it that does.
func (look Look) fieldBox(field Field, focused bool, content func()) {
	scheme := look.Roles
	thickness := float32(1)
	if focused || field.Error != "" {
		thickness = 2
	}
	right := float32(16)
	if field.Secret && field.Reveal != nil {
		right = 4 // the eye's own target keeps it off the edge
	}
	row := ui.Layout{Width: ui.Grow(), Height: ui.Grow(), Padding: ui.Padding{Left: 16, Right: right}, Gap: 8, AlignY: ui.Center}
	if !field.Filled {
		box := ui.Look{Fill: scheme.Surface, Border: look.fieldColour(field, focused, scheme.Outline), BorderWidth: thickness, Radius: 4}
		look.Context.Frame(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(56)}, box, func() { look.Context.Box(row, content) })
		return
	}
	box := ui.Look{Fill: scheme.SurfaceContainerHighest, Radius: 4, Clip: true}
	look.Context.Frame(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Fixed(56)}, box, func() {
		look.Context.Box(row, content)
		line := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(thickness)}
		look.Context.Frame(line, ui.Look{Fill: look.fieldColour(field, focused, scheme.OnSurfaceVariant)}, nil)
	})
}

// fieldFoot writes the hint or the error under a field, and the count of its characters at the end.
func (look Look) fieldFoot(field Field, text string) {
	scheme := look.Roles
	note, colour := field.Supporting, scheme.OnSurfaceVariant
	if field.Error != "" {
		note, colour = field.Error, scheme.Error
	}
	counted := field.Counter && field.MaxLength > 0
	if note == "" && !counted {
		return
	}
	look.Context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.PadXY(16, 0), Gap: 16}, func() {
		look.Context.Box(ui.Layout{Width: ui.Grow()}, func() {
			if note != "" {
				look.Paragraph(note, 12, draw2d.Regular, colour)
			}
		})
		if counted {
			look.Text(fmt.Sprintf("%d/%d", utf8.RuneCountInString(text), field.MaxLength), 12, draw2d.Regular, colour)
		}
	})
}

// SearchBar is a pill to search in, at the top of a screen. It returns true when the text changed.
func (look Look) SearchBar(name string, text *string, placeholder string) bool {
	scheme := look.Roles
	changed := false
	bar := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(56), Padding: ui.Padding{Left: 16, Right: 8}, Gap: 8, AlignY: ui.Center}
	look.Context.Frame(bar, ui.Look{Fill: scheme.SurfaceContainerHigh, Radius: 28}, func() {
		look.Context.IconSized("search", 24, scheme.OnSurfaceVariant)
		// the pill is the field's box: the field draws none of its own, nor a focus ring inside it
		changed = look.Context.TextFieldWith(name, text, ui.TextFieldOptions{Placeholder: placeholder, Plain: true})
		if *text != "" && look.NamedIconButton("x", look.words().Clear) {
			*text, changed = "", true
		}
	})
	return changed
}

// Dropdown is what an exposed dropdown keeps: the choice, and whether its choices show.
type Dropdown struct {
	Selected int
	open     bool
}

// Dropdown is Material's exposed dropdown menu: a box like a text field's showing the choice, which
// opens the choices under it when touched. It returns true when another one is chosen. name tells
// dropdowns apart.
func (look Look) Dropdown(name, label string, choices []string, dropdown *Dropdown) bool {
	scheme := look.Roles
	changed := false
	field := Field{Label: label}
	look.Context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Gap: 4}, func() {
		look.Text(label, 12, draw2d.Medium, look.fieldColour(field, dropdown.open, scheme.OnSurfaceVariant))
		box := ui.Look{Fill: scheme.Surface, Border: look.fieldColour(field, dropdown.open, scheme.Outline), BorderWidth: 1, Radius: 4}
		if dropdown.open {
			box.BorderWidth = 2
		}
		layout := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(56), Padding: ui.Padding{Left: 16, Right: 12}, Gap: 8, AlignY: ui.Center}
		if look.Context.Clickable(name, layout, ui.Looks{Normal: box, Hovered: box, Pressed: box}, func() {
			look.Context.Box(ui.Layout{Width: ui.Grow()}, func() {
				if dropdown.Selected >= 0 && dropdown.Selected < len(choices) {
					look.Text(choices[dropdown.Selected], 16, draw2d.Regular, scheme.OnSurface)
				}
			})
			look.Context.IconSized(map[bool]string{false: "chevron-down", true: "chevron-up"}[dropdown.open], 24, scheme.OnSurfaceVariant)
		}) {
			dropdown.open = !dropdown.open
		}
		menu := ui.Look{Fill: scheme.SurfaceContainer, Radius: 4, Shadow: true}
		look.Context.Popover(name+" choices", &dropdown.open, menu, func() {
			for index, choice := range choices {
				if look.dropdownChoice(choice, index == dropdown.Selected) {
					changed = index != dropdown.Selected
					dropdown.Selected, dropdown.open = index, false
				}
			}
		})
	})
	return changed
}

// dropdownChoice is one of a dropdown's choices, 48 pixels tall, the chosen one tinted and ticked.
func (look Look) dropdownChoice(choice string, chosen bool) bool {
	scheme := look.Roles
	normal := ui.Look{}
	if chosen {
		normal.Fill = scheme.SecondaryContainer
	}
	layout := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(48), Padding: ui.PadXY(12, 0), Gap: 12, AlignY: ui.Center}
	touched := look.Context.Clickable("choice "+choice, layout, pressedOnly(normal, stateLayer(scheme.OnSurface, 0.1)), func() {
		look.Context.Box(ui.Layout{Width: ui.Grow()}, func() { look.Text(choice, 16, draw2d.Regular, scheme.OnSurface) })
		if chosen {
			look.Context.IconSized("check", 20, scheme.OnSurface)
		}
	})
	look.Context.Describe(ui.Description{Selected: chosen})
	return touched
}

// Search is what a search view keeps: what is typed, and whether it is open over the screen.
type Search struct {
	Query  string
	Active bool
}

// SearchView is Material's search: a pill at the top of a screen which, touched, opens over the
// whole screen ready to type in, with a way back and a cross to empty it; content shows under it
// what the program finds (the last searches, suggestions, results), and closes the search by
// setting search.Active to false when one is chosen. The system's back closes it too. It returns
// true when the query changed. name tells searches apart.
func (look Look) SearchView(name, placeholder string, search *Search, content func()) bool {
	if !search.Active {
		return look.searchPill(name, placeholder, search)
	}
	scheme := look.Roles
	insets := look.Context.Insets()
	changed := false
	arrived := look.Context.Appear("search "+name, arrivingMotion)
	options := ui.LayerOptions{Width: ui.Grow(), Height: ui.Grow(), Hidden: 1 - arrived,
		Leave: ui.Leave{Duration: leavingMotion * 3 / 4}}
	dismissed := look.Context.Modal(name+" open", options, func() {
		view := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow(), Padding: ui.Padding{Top: insets.Top}}
		look.Context.Frame(view, ui.Look{Fill: scheme.SurfaceContainerHigh}, func() {
			bar := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(72), Padding: ui.Padding{Left: 4 + insets.Left, Right: 4 + insets.Right}, Gap: 4, AlignY: ui.Center}
			look.Context.Box(bar, func() {
				if look.NamedIconButton("arrow-left", look.words().Back) {
					search.Active = false
				}
				if arrived == 0 {
					look.Context.FocusField("query") // open ready to type in
				}
				changed = look.Context.TextFieldWith("query", &search.Query, ui.TextFieldOptions{Placeholder: placeholder, Plain: true})
				if search.Query != "" && look.NamedIconButton("x", look.words().Clear) {
					search.Query, changed = "", true
				}
			})
			look.Divider()
			look.Context.ScrollArea("found", ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow()}, content)
		})
	})
	if dismissed {
		search.Active = false
	}
	return changed
}

// searchPill is the closed search: a pill like the search bar's, which opens the search when
// touched.
func (look Look) searchPill(name, placeholder string, search *Search) bool {
	scheme := look.Roles
	pill := ui.Look{Fill: scheme.SurfaceContainerHigh, Radius: 28}
	layout := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(56), Padding: ui.Padding{Left: 16, Right: 16}, Gap: 12, AlignY: ui.Center}
	if look.Context.Clickable(name, layout, pressedOnly(pill, scheme.SurfaceContainerHighest), func() {
		look.Context.IconSized("search", 24, scheme.OnSurfaceVariant)
		shown, colour := placeholder, scheme.OnSurfaceVariant
		if search.Query != "" {
			shown, colour = search.Query, scheme.OnSurface
		}
		look.Text(shown, 16, draw2d.Regular, colour)
	}) {
		search.Active = true
	}
	return false
}
