package material

import (
	"fmt"
	"slices"
	"testing"

	"github.com/akrck02/go-renderer/input"
)

func TestAnExpandableRowOpensAndCloses(t *testing.T) {
	open, shownMore := false, false
	screen := newScreen(t, func(look Look) {
		shownMore = false
		look.ExpandableItem("faq", Item{Headline: "When is it lit?"}, &open, func() { shownMore = true })
	})
	screen.clocked()
	screen.frame()
	screen.tap("faq")
	if !open || !shownMore {
		t.Fatal("touching the row did not open it")
	}
	screen.tap("faq")
	if open || !shownMore {
		t.Fatalf("open %v, more shown %v after touching it again; want it closed, folding away", open, shownMore)
	}
	for range 20 {
		screen.frame()
	}
	if shownMore {
		t.Fatal("what was under the row is still there after it had time to fold")
	}
}

func TestInitialsAreTheFirstLettersOfTwoWords(t *testing.T) {
	for name, want := range map[string]string{"ana maría garcía": "AM", "Ñandú": "Ñ", "": ""} {
		if got := initialsOf(name); got != want {
			t.Fatalf("initials of %q = %q; want %q", name, got, want)
		}
	}
}

// hold keeps a finger still on a piece laid out in the last frame, for frames of 16 ms, and lifts it.
func (screen *screen) hold(frames int, path ...string) {
	screen.t.Helper()
	area, found := screen.context.Area(path...)
	if !found {
		screen.t.Fatalf("nothing called %v was laid out", path)
	}
	x, y := float64(area.X+area.Width/2), float64(area.Y+area.Height/2)
	screen.state.Feed().SetTouch(1, x, y)
	screen.state.Feed().MoveCursor(x, y)
	screen.state.Feed().SetButton(input.MouseLeft, true)
	for range frames {
		screen.frame()
	}
	screen.state.Feed().SetButton(input.MouseLeft, false)
	screen.state.Feed().LiftTouch(1)
	screen.frame()
}

func TestARowHeldIsChosenAndThenATouchChoosesAnother(t *testing.T) {
	chosen := make([]bool, 3)
	touches := 0
	screen := newScreen(t, func(look Look) {
		choosing := slices.Contains(chosen, true)
		for index := range chosen {
			look.Context.Push(fmt.Sprint("light ", index))
			touched := look.ListItem("row", Item{Icon: "lightbulb", Headline: fmt.Sprint("Light ", index), Selected: chosen[index]})
			switch {
			case look.LongPressed():
				chosen[index] = true
			case touched && choosing:
				chosen[index] = !chosen[index]
			case touched:
				touches++
			}
			look.Context.Pop()
		}
	})
	screen.clocked()
	screen.frame()
	screen.hold(40, "light 0", "row") // 640 ms
	if !chosen[0] || touches != 0 {
		t.Fatalf("chosen %v, %d touches after holding the first row; want it chosen and no touch", chosen, touches)
	}
	screen.tap("light 2", "row")
	if !chosen[2] || chosen[1] {
		t.Fatalf("chosen %v after touching the third row while choosing; want it chosen too", chosen)
	}
}
