package material

import (
	"testing"

	"github.com/akrck02/go-renderer/draw2d"
)

func TestTheFloatingMenuOpensChoosesAndCloses(t *testing.T) {
	menu := &FloatingMenu{Icon: "plus", Actions: []Destination{{Icon: "file", Label: "Note"}, {Icon: "image", Label: "Picture"}}}
	chosen := -1
	screen := newScreen(t, func(look Look) {
		if index := look.FloatingMenu("create", menu); index >= 0 {
			chosen = index
		}
	})
	screen.frame()
	screen.tap("create", "toggle")
	if !menu.Open {
		t.Fatal("touching the button does not open the menu")
	}
	screen.frame() // the actions rise in
	screen.tap("create", "action Picture")
	if chosen != 1 || menu.Open {
		t.Fatalf("chosen %d, open %v: want Picture, and the menu closed", chosen, menu.Open)
	}
	screen.tap("create", "toggle")
	screen.state.Feed().PressBack()
	screen.frame()
	if menu.Open {
		t.Fatal("the back does not close the menu")
	}
}

func TestGroupedAndSplitButtonsAnswerTheirParts(t *testing.T) {
	selected, action, more, tool := 0, false, false, -1
	screen := newScreen(t, func(look Look) {
		look.ButtonGroup("size", []string{"S", "M", "L"}, &selected)
		if main, menu := look.SplitButton("send", "Send"); main || menu {
			action, more = action || main, more || menu
		}
		if index := look.FloatingToolbar("tools", []Destination{{Icon: "pencil", Label: "Edit"}, {Icon: "share-2", Label: "Share"}, {Icon: "trash-2", Label: "Delete"}}, 0); index >= 0 {
			tool = index
		}
	})
	screen.frame()
	screen.tap("size", "L")
	screen.tap("split Send", "more")
	screen.tap("tools", "tool trash-2")
	if selected != 2 || action || !more || tool != 2 {
		t.Fatalf("selected %d, action %v, more %v, tool %d: want L, the menu half alone, and the bin", selected, action, more, tool)
	}
}

// drawn tells whether a text is in the last frame's drawing.
func (screen *screen) drawn(text string) bool {
	for _, command := range screen.context.Drawing().Commands {
		if command.Kind == draw2d.TextKind && command.Text == text {
			return true
		}
	}
	return false
}

func TestAClosedFloatingMenusActionsFadeAwayUntouchable(t *testing.T) {
	menu := &FloatingMenu{Icon: "plus", Actions: []Destination{{Icon: "file", Label: "Note"}}}
	chosen := -1
	screen := newScreen(t, func(look Look) {
		if index := look.FloatingMenu("create", menu); index >= 0 {
			chosen = index
		}
	})
	screen.clocked()
	screen.frame()
	screen.tap("create", "toggle")
	for range 20 {
		screen.frame()
	}
	screen.tap("create", "toggle") // it closes in the frame the finger lifts
	screen.frame()
	if !screen.drawn("Note") {
		t.Fatal("the action went at once when the menu closed")
	}
	if _, found := screen.context.Area("create", "action Note"); found {
		t.Fatal("an action on its way out can still be touched")
	}
	for range 20 {
		screen.frame()
	}
	if screen.drawn("Note") || chosen != -1 {
		t.Fatalf("drawn %v, chosen %d after the menu had time to close; want it gone, nothing chosen", screen.drawn("Note"), chosen)
	}
}
