package material

import (
	"time"

	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// Dialog is what a dialog says and offers.
type Dialog struct {
	Icon    string // over the title, centred; none when empty
	Title   string
	Text    string // under the title, broken into lines
	Content func() // more under the text (fields, choices); none when nil
	Confirm string // the main action, at the end
	Dismiss string // the other action, before it; none when empty
}

// Answer is what the person did with a dialog in a frame.
type Answer int

const (
	Open      Answer = iota // nothing yet: keep showing it
	Confirmed               // the main action
	Dismissed               // the other action, a touch outside, Escape
)

// Dialog asks something over the screen, which is dimmed and does not react meanwhile: a card in
// the middle with a title, a text and one or two actions. The program shows it with an if while it
// wants an answer, and stops when the answer is not Open. name tells dialogs apart.
func (look Look) Dialog(name string, dialog Dialog) Answer {
	roles := look.Roles
	width, _ := look.Context.Size()
	answer := Open
	// it rises into place as the screen behind darkens
	arrived := look.Context.Appear("dialog "+name, arrivingMotion)
	padding := look.safePadding(24)
	padding.Top += (1 - arrived) * 48
	options := ui.LayerOptions{AlignX: ui.Center, AlignY: ui.Center, Padding: padding,
		Width: ui.Fixed(min(560, max(width-48, 200))), Hidden: 1 - arrived,
		Leave: ui.Leave{Duration: leavingMotion * 3 / 4}} // it fades out where it is
	card := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.Pad(24), Gap: 16}
	dismissed := look.Context.Modal(name, options, func() {
		look.Context.Frame(card, ui.Look{Fill: roles.SurfaceContainerHigh, Radius: 28, Shadow: true}, func() {
			alignX := ui.Start
			if dialog.Icon != "" {
				alignX = ui.Center
				look.Context.Box(ui.Layout{Width: ui.Grow(), AlignX: ui.Center}, func() {
					look.Context.IconSized(dialog.Icon, 24, roles.Primary)
				})
			}
			look.Context.Box(ui.Layout{Width: ui.Grow(), AlignX: alignX}, func() {
				look.Text(dialog.Title, 24, draw2d.Regular, roles.OnSurface)
				look.Context.Describe(ui.Description{Role: ui.HeadingRole})
			})
			if dialog.Text != "" {
				look.Paragraph(dialog.Text, 14, draw2d.Regular, roles.OnSurfaceVariant)
			}
			if dialog.Content != nil {
				dialog.Content()
			}
			look.Context.Box(ui.Layout{Width: ui.Grow(), Gap: 8, AlignX: ui.End}, func() {
				if dialog.Dismiss != "" && look.TextButton("", dialog.Dismiss) {
					answer = Dismissed
				}
				if dialog.Confirm != "" && look.TextButton("", dialog.Confirm) {
					answer = Confirmed
				}
			})
		})
	})
	if dismissed && answer == Open {
		answer = Dismissed
	}
	return answer
}

// BottomSheet shows content in a sheet rising from the bottom of the screen, which is dimmed
// meanwhile: choices for what was touched, a short form. It returns true when the person touches
// outside it (or presses Escape); the program then stops showing it. name tells sheets apart.
func (look Look) BottomSheet(name string, content func()) (dismissed bool) {
	roles := look.Roles
	insets := look.Context.Insets()
	width, _ := look.Context.Size()
	// it slides up from below the screen as the screen darkens
	arrived := look.Context.Appear("sheet "+name, arrivingMotion)
	options := ui.LayerOptions{AlignX: ui.Center, AlignY: ui.End, Padding: ui.Padding{Top: insets.Top + 72, Bottom: -(1 - arrived) * 400},
		Width: ui.Fixed(min(640, width)), Hidden: 1 - arrived,
		Leave: ui.Leave{Duration: leavingMotion, Y: 400}} // it slides back down
	sheet := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(),
		Padding: ui.Padding{Left: 16 + insets.Left, Top: 16, Right: 16 + insets.Right, Bottom: 16 + insets.Bottom}, Gap: 8}
	return look.Context.Modal(name, options, func() {
		look.Context.Frame(sheet, ui.Look{Fill: roles.SurfaceContainerLow, Radius: 28}, func() {
			handle := ui.Layout{Width: ui.Fixed(32), Height: ui.Fixed(4)}
			look.Context.Box(ui.Layout{Width: ui.Grow(), AlignX: ui.Center, Padding: ui.Padding{Bottom: 8}}, func() {
				look.Context.Frame(handle, ui.Look{Fill: stateLayer(roles.OnSurfaceVariant, 0.4), Radius: 2}, nil)
			})
			content()
		})
	})
}

// SideSheet slides in from the end of the screen over the dimmed screen, headed by a title and a
// cross: filters for a list, the details of what is chosen, room for a form on a wide screen. It
// returns true when it should close: the cross touched, a touch outside, the system's back. name
// tells side sheets apart.
func (look Look) SideSheet(name, title string, content func()) (closed bool) {
	roles := look.Roles
	insets := look.Context.Insets()
	width, _ := look.Context.Size()
	sheetWidth := min(400, width-56)
	arrived := look.Context.Appear("side sheet "+name, arrivingMotion)
	options := ui.LayerOptions{AlignX: ui.End, AlignY: ui.Start, Padding: ui.Padding{Right: -(1 - arrived) * sheetWidth},
		Width: ui.Fixed(sheetWidth), Height: ui.Grow(), Hidden: 1 - arrived,
		Leave: ui.Leave{Duration: leavingMotion, X: sheetWidth}} // back off the side it came from
	dismissed := look.Context.Modal(name, options, func() {
		sheet := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow(),
			Padding: ui.Padding{Left: 24, Top: 12 + insets.Top, Right: 12 + insets.Right, Bottom: 24 + insets.Bottom}, Gap: 8}
		look.Context.Frame(sheet, ui.Look{Fill: roles.SurfaceContainerLow, Radius: 16, Shadow: true}, func() {
			look.Context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(56), AlignY: ui.Center}, func() {
				look.Context.Box(ui.Layout{Width: ui.Grow()}, func() { look.Text(title, 22, draw2d.Regular, roles.OnSurfaceVariant) })
				closed = look.NamedIconButton("x", look.words().Close)
			})
			content()
		})
	})
	return closed || dismissed
}

// Snackbar is a short note at the bottom of the screen for a few seconds, with an action that
// undoes or follows what was just done. A program keeps one, calls Show when something happens and
// draws it with Look.Snackbar in every frame.
type Snackbar struct {
	Text   string
	Action string // a button at its end; none when empty
	// Above is how far over the bottom of the window it sits: a navigation bar's height
	// (NavigationBarHeight) when there is one.
	Above float32
	until time.Time
}

// Show puts a note up for a while (Material's are 4 to 10 seconds).
func (snackbar *Snackbar) Show(text, action string, duration time.Duration) {
	snackbar.Text, snackbar.Action, snackbar.until = text, action, time.Now().Add(duration)
}

// Hide takes the note away before its time.
func (snackbar *Snackbar) Hide() { snackbar.until = time.Time{} }

// Shown tells whether the note is up.
func (snackbar *Snackbar) Shown() bool { return time.Now().Before(snackbar.until) }

// Snackbar draws the note while it is up, over everything, and asks for a frame when it is due to
// go. It returns true when its action is touched, and the note goes then.
func (look Look) Snackbar(snackbar *Snackbar) (actionTouched bool) {
	left := time.Until(snackbar.until)
	if left <= 0 {
		return false
	}
	look.Context.RequestFrame(left)
	roles := look.Roles
	insets := look.Context.Insets()
	width, _ := look.Context.Size()
	arrived := look.Context.Appear("snackbar", arrivingMotion) // it slides up from the bottom edge
	padding := ui.Padding{Left: 16 + insets.Left, Right: 16 + insets.Right,
		Bottom: 16 + insets.Bottom + snackbar.Above - (1-arrived)*(80+insets.Bottom+snackbar.Above)}
	options := ui.LayerOptions{AlignX: ui.Center, AlignY: ui.End, Padding: padding,
		Width: ui.Fixed(min(600, width-padding.Left-padding.Right)), Leave: ui.Leave{Duration: leavingMotion}}
	look.Context.Overlay("snackbar", options, func() {
		bar := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(48), Padding: ui.Padding{Left: 16, Right: 8}, Gap: 8, AlignY: ui.Center}
		look.Context.Frame(bar, ui.Look{Fill: roles.InverseSurface, Radius: 4, Shadow: true}, func() {
			look.Context.Box(ui.Layout{Width: ui.Grow()}, func() {
				look.Text(snackbar.Text, 14, draw2d.Regular, roles.InverseOnSurface)
			})
			if snackbar.Action == "" {
				return
			}
			layout := ui.Layout{Height: ui.Fixed(40), Padding: ui.PadXY(12, 0), AlignY: ui.Center}
			normal := ui.Look{Radius: 20}
			actionTouched = look.Context.Clickable("snackbar action", layout, pressedOnly(normal, stateLayer(roles.InversePrimary, 0.12)), func() {
				look.Text(snackbar.Action, 14, draw2d.Medium, roles.InversePrimary)
			})
		})
	})
	if actionTouched {
		snackbar.Hide()
	}
	return actionTouched
}

// Menu is a button with an icon (three dots, usually) that opens a list of choices under it; a
// screen reader reads the button as "More options". It returns the index of the one chosen in this
// frame, or -1.
func (look Look) Menu(icon string, choices []string) int {
	chosen := -1
	look.Context.MenuButton(icon, "", func() {
		for index, choice := range choices {
			if look.Context.MenuItem("", choice, "") {
				chosen = index
			}
		}
	})
	look.Context.Describe(ui.Description{Label: look.words().MoreOptions})
	return chosen
}

// Sheet is what a standard bottom sheet keeps: whether it is expanded, and how far a finger has
// dragged it from where it rests.
type Sheet struct {
	Expanded bool
	// Above is how far over the bottom of the window it sits: a navigation bar's height
	// (NavigationBarHeight) when there is one.
	Above   float32
	dragged float32 // upward, in pixels, while a finger holds the handle
}

// StandardSheet is a bottom sheet that stays over the screen without dimming it, which goes on
// working above it (a map under a list of its places): it rests showing peek pixels, or expanded
// to most of the screen, and a finger drags it by its handle between the two; let go, it settles
// where it was heading. Its content scrolls inside it. name tells sheets apart.
func (look Look) StandardSheet(name string, sheet *Sheet, peek float32, content func()) {
	roles := look.Roles
	insets := look.Context.Insets()
	width, height := look.Context.Size()
	full := max(height-insets.Top-56-sheet.Above, peek)
	rest := peek
	if sheet.Expanded {
		rest = full
	}
	motion := arrivingMotion
	if sheet.dragged != 0 {
		motion = 0 // the sheet follows the finger
	}
	shown := look.Context.Animate(name+" height", clampFloat(rest+sheet.dragged, peek*0.75, full), motion)
	// at the window's foot the sheet reaches its edge and keeps its content clear of the gesture
	// bar; over a navigation bar it sits on the bar, which is clear of the gesture bar already
	lift, inside := float32(0), insets.Bottom
	if sheet.Above > 0 {
		lift, inside = sheet.Above+insets.Bottom, 0
	}
	const hiddenCorners = 28 // the panel's lower corners, cut off: a sheet has round corners on top only
	options := ui.LayerOptions{AlignX: ui.Center, AlignY: ui.End, Padding: ui.Padding{Bottom: lift},
		Width: ui.Fixed(min(640, width)), Height: ui.Fixed(shown)}
	look.Context.Overlay(name, options, func() {
		cut := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Fixed(shown)}
		look.Context.Frame(cut, ui.Look{Clip: true}, func() {
			panel := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Fixed(shown + hiddenCorners),
				Padding: ui.Padding{Bottom: hiddenCorners + inside}}
			look.Context.Frame(panel, ui.Look{Fill: roles.SurfaceContainerLow, Radius: 28}, func() {
				look.sheetHandle(name, sheet, full-peek)
				look.Context.ScrollArea(name+" content", ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow()}, content)
			})
		})
	})
}

// sheetHandle is the bar a finger drags a standard sheet by: let go, the sheet expands past a
// quarter of the way up or on a flick upwards, and folds down past a quarter of the way down or on
// a flick downwards.
func (look Look) sheetHandle(name string, sheet *Sheet, travel float32) {
	look.Context.SwipeBox(name+" handle", ui.Layout{Width: ui.Grow(), Height: ui.Fixed(36), AlignX: ui.Center, AlignY: ui.Center}, ui.Vertical, func(swipe ui.Swipe) {
		look.Context.Frame(ui.Layout{Width: ui.Fixed(32), Height: ui.Fixed(4)}, ui.Look{Fill: stateLayer(look.Roles.OnSurfaceVariant, 0.4), Radius: 2}, nil)
		sheet.dragged = 0
		if swipe.Held {
			sheet.dragged = -swipe.Offset
		}
		if !swipe.Released {
			return
		}
		moved := -swipe.Offset
		switch {
		case moved > travel/4 || swipe.Velocity < -600:
			sheet.Expanded = true
		case moved < -travel/4 || swipe.Velocity > 600:
			sheet.Expanded = false
		}
	})
}
