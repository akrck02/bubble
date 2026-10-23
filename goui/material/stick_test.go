package material

import (
	"testing"

	"github.com/akrck02/go-renderer/input"
)

func TestTheStickFollowsItsFingerAndNoOther(t *testing.T) {
	var stick Stick
	screen := newScreen(t, func(look Look) { look.Stick("walk", &stick, 120) })
	screen.frame()
	area, _ := screen.context.Area("walk")
	middleX, middleY := float64(area.X+area.Width/2), float64(area.Y+area.Height/2)

	screen.state.Feed().SetTouch(7, 300, 300) // a finger elsewhere does not take it
	screen.frame()
	if stick.Held {
		t.Fatal("a finger away from the stick took it")
	}
	screen.state.Feed().SetTouch(8, middleX, middleY) // on it
	screen.frame()
	if !stick.Held || stick.Finger != 8 || stick.Forward != 0 {
		t.Fatalf("held %v by %d, forward %v; want held by the finger on it, in the middle", stick.Held, stick.Finger, stick.Forward)
	}
	screen.state.Feed().SetTouch(8, middleX, middleY-200) // dragged far up
	screen.frame()
	if stick.Forward < 0.99 || stick.Sideways != 0 {
		t.Fatalf("forward %v, sideways %v; want all the way ahead", stick.Forward, stick.Sideways)
	}
	if !stick.Touching(input.Touch{Finger: 8}) || stick.Touching(input.Touch{Finger: 7}) {
		t.Fatal("the stick does not tell its finger from the other one")
	}
	screen.state.Feed().LiftTouch(8)
	screen.frame()
	if stick.Held || stick.Forward != 0 {
		t.Fatal("letting go does not bring the stick back")
	}
}

func TestAStickInsideSwipingPagesTakesItsFingerFromThem(t *testing.T) {
	var stick Stick
	var pager Pager
	screen := newScreen(t, func(look Look) {
		look.TabPager("views", []string{"Bay", "List"}, &pager, func(index int) {
			look.Stick("sail", &stick, 120)
		})
	})
	screen.clocked()
	screen.frame()
	screen.frame()
	area, found := screen.context.Area("views page 0", "sail")
	if !found {
		t.Fatal("the stick was not laid out inside the first page")
	}
	middleX, middleY := float64(area.X+area.Width/2), float64(area.Y+area.Height/2)

	screen.drag(middleX, middleY, 200, 0, false) // far to the right, held
	if !stick.Held || stick.Sideways < 0.99 {
		t.Fatalf("held %v, sideways %v; want the stick held all the way right", stick.Held, stick.Sideways)
	}
	if pager.shown != 0 {
		t.Fatalf("pages shown at %v while steering; want them still", pager.shown)
	}
	screen.drag(middleX+200, middleY, 0, 0, true)
	if stick.Held || pager.Selected != 0 {
		t.Fatalf("held %v, page %d after letting go; want the stick let go on the first page", stick.Held, pager.Selected)
	}
}
