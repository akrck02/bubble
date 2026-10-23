package material

import (
	"testing"
)

func TestASwitchFlipsWhenTouched(t *testing.T) {
	on, flips := false, 0
	screen := newScreen(t, func(look Look) {
		if look.Switch("sound", &on) {
			flips++
		}
	})
	screen.frame()
	screen.tap("sound")
	if !on || flips != 1 {
		t.Fatalf("on %v after %d flips; want on after one", on, flips)
	}
}

func TestSegmentedButtonsChooseOne(t *testing.T) {
	selected := 0
	screen := newScreen(t, func(look Look) {
		look.SegmentedButtons("speed", []string{"Slow", "Fast"}, &selected)
	})
	screen.frame()
	screen.tap("speed", "Fast")
	if selected != 1 {
		t.Fatalf("selected %d; want Fast", selected)
	}
}

func TestChoicesAnswerTouches(t *testing.T) {
	ticked, chosen, tab := false, 0, 0
	screen := newScreen(t, func(look Look) {
		look.Checkbox("Lit at night", &ticked)
		for index, option := range []string{"North", "South"} {
			if look.RadioButton(option, index == chosen) {
				chosen = index
			}
		}
		look.Tabs("views", []string{"Map", "List"}, &tab)
	})
	screen.frame()
	screen.tap("check Lit at night")
	screen.tap("radio South")
	screen.tap("views", "List")
	if !ticked || chosen != 1 || tab != 1 {
		t.Fatalf("ticked %v, chosen %d, tab %d: want ticked, South, List", ticked, chosen, tab)
	}
}

func TestAnInputChipTellsItsCrossApart(t *testing.T) {
	removed := false
	screen := newScreen(t, func(look Look) {
		if _, gone := look.InputChip("harbour"); gone {
			removed = true
		}
	})
	screen.frame()
	screen.tap("input harbour", "remove")
	if !removed {
		t.Fatal("touching the cross did not remove the chip")
	}
}
