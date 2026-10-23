package material

import (
	"testing"
	"time"

	"github.com/akrck02/go-renderer/draw2d"
	"github.com/akrck02/go-renderer/input"
)

func TestADialogSaysWhatWasChosen(t *testing.T) {
	answer := Open
	screen := newScreen(t, func(look Look) {
		if got := look.Dialog("delete", Dialog{Title: "Delete the light?", Text: "It cannot be undone.", Confirm: "Delete", Dismiss: "Cancel"}); got != Open {
			answer = got
		}
	})
	screen.frame()
	screen.tap("delete", "button Delete")
	if answer != Confirmed {
		t.Fatalf("answer %v after Delete; want Confirmed", answer)
	}
	answer = Open
	screen.tap("delete", "button Cancel")
	if answer != Dismissed {
		t.Fatalf("answer %v after Cancel; want Dismissed", answer)
	}
	answer = Open
	screen.state.Feed().MoveCursor(5, 5) // the dimmed screen around the card
	screen.state.Feed().SetButton(input.MouseLeft, true)
	screen.frame()
	screen.state.Feed().SetButton(input.MouseLeft, false)
	screen.frame()
	if answer != Dismissed {
		t.Fatalf("answer %v after a touch outside; want Dismissed", answer)
	}
}

func TestASnackbarShowsForItsTimeAndItsActionAnswers(t *testing.T) {
	var snackbar Snackbar
	undone := false
	screen := newScreen(t, func(look Look) {
		if look.Snackbar(&snackbar) {
			undone = true
		}
	})
	snackbar.Show("Light removed", "Undo", time.Minute)
	screen.frame()
	if _, found := screen.context.Area("snackbar"); !found || !snackbar.Shown() {
		t.Fatal("the note is not up")
	}
	screen.tap("snackbar", "snackbar action")
	if !undone || snackbar.Shown() {
		t.Fatalf("undone %v, still shown %v: the action answers and takes the note away", undone, snackbar.Shown())
	}
	snackbar.Show("Saved", "", time.Millisecond)
	time.Sleep(2 * time.Millisecond)
	screen.frame()
	if snackbar.Shown() {
		t.Fatal("the note outlived its time")
	}
}

func TestADialogAnsweredFadesOutInsteadOfVanishing(t *testing.T) {
	asking := true
	screen := newScreen(t, func(look Look) {
		if asking && look.Dialog("light", Dialog{Title: "Light the north lamp?", Confirm: "Light"}) != Open {
			asking = false
		}
	})
	screen.clocked()
	for range 20 {
		screen.frame()
	}
	screen.tap("light", "button Light")
	if asking {
		t.Fatal("touching Light did not answer the dialog")
	}
	screen.frame()
	if !screen.drawn("Light the north lamp?") {
		t.Fatal("the dialog vanished at once")
	}
	for range 20 {
		screen.frame()
	}
	if screen.drawn("Light the north lamp?") {
		t.Fatal("the dialog is still there after it had time to fade")
	}
}

func TestASideSheetClosesWithItsCross(t *testing.T) {
	shown := true
	screen := newScreen(t, func(look Look) {
		if shown && look.SideSheet("filters", "Filters", func() { look.Text("Only lit lights", 14, draw2d.Regular, look.Roles.OnSurface) }) {
			shown = false
		}
	})
	screen.clocked()
	for range 20 {
		screen.frame()
	}
	if !screen.drawn("Filters") {
		t.Fatal("the side sheet's title is not drawn")
	}
	screen.tap("filters", "icon x")
	if shown {
		t.Fatal("touching the cross did not close the side sheet")
	}
}

// sheetTop is where the standard sheet's panel starts, from its drawing.
func sheetTop(screen *screen) float32 {
	top := float32(1e9)
	for _, command := range screen.context.Drawing().Commands {
		if command.Kind == draw2d.RectangleKind && command.Radius == 28 && command.Area.Width > 300 {
			top = min(top, command.Area.Y)
		}
	}
	return top
}

func TestAStandardSheetIsDraggedUpAndSettles(t *testing.T) {
	var sheet Sheet
	screen := newScreen(t, func(look Look) {
		look.StandardSheet("places", &sheet, 200, func() { look.Text("North Cape", 16, draw2d.Regular, look.Roles.OnSurface) })
	})
	screen.clocked()
	for range 3 {
		screen.frame()
	}
	resting := sheetTop(screen)
	if resting < 590 || resting > 610 {
		t.Fatalf("the sheet starts at %v; want 200 pixels up from the 800-pixel screen's foot", resting)
	}
	handle, found := screen.context.Area("places", "places handle")
	if !found {
		t.Fatal("the sheet has no handle")
	}
	screen.drag(200, float64(handle.Y+handle.Height/2), 0, -260, false)
	if held := sheetTop(screen); held > resting-200 {
		t.Fatalf("the sheet is at %v while dragged 260 pixels up; want it following the finger", held)
	}
	for range 15 {
		screen.frame()
	}
	screen.state.Feed().SetButton(input.MouseLeft, false)
	screen.state.Feed().LiftTouch(1)
	for range 30 {
		screen.frame()
	}
	if !sheet.Expanded || sheetTop(screen) > 100 {
		t.Fatalf("expanded %v, top at %v after letting go high; want it expanded to most of the screen", sheet.Expanded, sheetTop(screen))
	}
}
