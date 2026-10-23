package main

import (
	"testing"
	"time"
	"unicode/utf8"

	"github.com/akrck02/go-renderer/draw2d"
	"github.com/akrck02/go-renderer/input"
	ui "github.com/akrck02/go-ui"
)

type measurer struct{}

func (measurer) MeasureText(text string, size float32, _ draw2d.Weight, _ float32) draw2d.TextMetrics {
	return draw2d.TextMetrics{Width: float32(utf8.RuneCountInString(text)) * size / 2, Ascent: size * 0.8, Descent: size * 0.2}
}

func (measurer) HasIcon(string) bool { return true }

// Every destination can be reached from every other, by its place in the navigation bar.
func TestEveryScreenLeadsToEveryOther(t *testing.T) {
	program := &gallery{inputs: newInputs(), lists: newLists(), tabs: newTabs()}
	context, state := ui.NewWithMeasurer(measurer{}), input.NewState()
	state.WindowWidth, state.WindowHeight = 412, 892
	state.FramebufferWidth, state.FramebufferHeight = 412, 892
	var finger input.TouchPointer
	frame := func() {
		context.Begin(state)
		program.build(context)
		context.End()
		state.Feed().EndFrame()
	}
	tap := func(label string) {
		area, found := context.Area("navigate " + label)
		if !found {
			t.Fatalf("no destination %s on screen %d", label, program.screen)
		}
		// as a phone's finger does: down, a frame, up, a frame
		x, y := area.X+area.Width/2, area.Y+area.Height/2
		finger.Begin(state, 1, x, y, 1)
		frame()
		finger.End(state, 1, x, y, 1)
		frame()
		frame()
	}
	labels := []string{"Actions", "Inputs", "Lists", "Feedback", "Tabs"}
	frame()
	for from := range labels {
		for to := range labels {
			tap(labels[from])
			tap(labels[to])
			if program.screen != to {
				t.Fatalf("from %s, touching %s left the gallery on screen %d", labels[from], labels[to], program.screen)
			}
		}
	}
}

// A lighthouse held is chosen; the back lets it go, and the bin removes what is chosen.
func TestALighthouseHeldIsChosenAndTheBinRemovesIt(t *testing.T) {
	program := &gallery{screen: listsScreen, inputs: newInputs(), lists: newLists(), tabs: newTabs()}
	context, state := ui.NewWithMeasurer(measurer{}), input.NewState()
	state.WindowWidth, state.WindowHeight = 412, 892
	state.FramebufferWidth, state.FramebufferHeight = 412, 892
	now := time.Unix(1000, 0)
	state.Now = func() time.Time { return now }
	var finger input.TouchPointer
	frame := func() {
		now = now.Add(16 * time.Millisecond)
		context.Begin(state)
		program.build(context)
		context.End()
		state.Feed().EndFrame()
	}
	touch := func(frames int, path ...string) {
		area, found := context.Area(path...)
		if !found {
			t.Fatalf("nothing called %v on the lists screen", path)
		}
		x, y := area.X+area.Width/2, area.Y+area.Height/2
		finger.Begin(state, 1, x, y, 1)
		for range frames {
			frame()
		}
		finger.End(state, 1, x, y, 1)
		frame()
		frame()
	}
	for range 30 { // the screen rising into place
		frame()
	}
	touch(40, "lists screen ", "Seal Rock", "row")
	if !program.lists.choosing() || program.lists.chosenCount() != 1 {
		t.Fatal("holding Seal Rock did not choose it")
	}
	state.Feed().PressBack()
	frame()
	if program.lists.choosing() || program.screen != listsScreen {
		t.Fatalf("after the back, choosing %v on screen %d; want nothing chosen, still on the lists", program.lists.choosing(), program.screen)
	}
	touch(40, "lists screen ", "Seal Rock", "row")
	touch(1, "lists screen ", "icon trash-2")
	for _, light := range program.lists.lighthouses {
		if light.name == "Seal Rock" {
			t.Fatal("the bin left Seal Rock in the list")
		}
	}
}

// A lighthouse touched opens its details; the back returns to the list, and the next back leaves
// the lists.
func TestALighthouseTouchedOpensItsDetails(t *testing.T) {
	program := &gallery{screen: listsScreen, inputs: newInputs(), lists: newLists(), tabs: newTabs()}
	context, state := ui.NewWithMeasurer(measurer{}), input.NewState()
	state.WindowWidth, state.WindowHeight = 412, 892
	state.FramebufferWidth, state.FramebufferHeight = 412, 892
	now := time.Unix(1000, 0)
	state.Now = func() time.Time { return now }
	var finger input.TouchPointer
	frame := func() {
		now = now.Add(16 * time.Millisecond)
		context.Begin(state)
		program.build(context)
		context.End()
		state.Feed().EndFrame()
	}
	settle := func() {
		for range 30 {
			frame()
		}
	}
	settle()
	area, found := context.Area("lists screen ", "Seal Rock", "row")
	if !found {
		t.Fatal("no Seal Rock row")
	}
	x, y := area.X+area.Width/2, area.Y+area.Height/2
	finger.Begin(state, 1, x, y, 1)
	frame()
	finger.End(state, 1, x, y, 1)
	settle()
	if program.lists.navigation.Current() != "Seal Rock" {
		t.Fatalf("on %q after touching Seal Rock; want its details", program.lists.navigation.Current())
	}
	state.Feed().PressBack()
	settle()
	if program.lists.navigation.Current() != "" || program.screen != listsScreen {
		t.Fatalf("on %q of screen %d after the back; want the list", program.lists.navigation.Current(), program.screen)
	}
	state.Feed().PressBack()
	settle()
	if program.screen != actionsScreen {
		t.Fatal("the next back did not leave the lists")
	}
}

// With a screen reader, every destination is reached by double taps alone, and every control on it
// has a name to read.
func TestAScreenReaderReachesAndNamesEveryScreen(t *testing.T) {
	program := &gallery{inputs: newInputs(), lists: newLists(), tabs: newTabs()}
	context, state := ui.NewWithMeasurer(measurer{}), input.NewState()
	state.WindowWidth, state.WindowHeight = 412, 892
	state.FramebufferWidth, state.FramebufferHeight = 412, 892
	context.SetScreenReader(true)
	frame := func() {
		context.Begin(state)
		program.build(context)
		context.End()
		state.Feed().EndFrame()
	}
	find := func(label string) (ui.AccessibleNode, bool) {
		for _, node := range context.ScreenReaderNodes() {
			if node.Label == label && node.Pressable {
				return node, true
			}
		}
		return ui.AccessibleNode{}, false
	}
	frame()
	for index, label := range []string{"Actions", "Inputs", "Lists", "Feedback", "Tabs"} {
		node, found := find(label)
		if !found {
			for _, each := range context.ScreenReaderNodes() {
				t.Logf("%q role %d pressable %v area %+v parent %d", each.Label, each.Role, each.Pressable, each.Area, each.Parent)
			}
			t.Fatalf("no destination %s to double tap on screen %d", label, program.screen)
		}
		context.ScreenReaderAct(node.ID, ui.ReaderPress)
		frame()
		frame()
		if program.screen != index {
			t.Fatalf("a double tap on %s left the gallery on screen %d", label, program.screen)
		}
		for _, node := range context.ScreenReaderNodes() {
			if node.Pressable && node.Label == "" {
				t.Errorf("on %s, a control has no name: %+v", label, node)
			}
		}
	}
}
