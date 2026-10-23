package material

import (
	"testing"

	ui "github.com/akrck02/go-ui"
)

// read finds the node a screen reader reads with a label, or fails.
func (screen *screen) read(label string) ui.AccessibleNode {
	screen.t.Helper()
	var labels []string
	for _, node := range screen.context.ScreenReaderNodes() {
		if node.Label == label {
			return node
		}
		labels = append(labels, node.Label)
	}
	screen.t.Fatalf("no node reads %q in %q", label, labels)
	return ui.AccessibleNode{}
}

func TestSelectionPiecesSayWhatTheyAreAndWhetherTheyAreOn(t *testing.T) {
	lamp, fog := true, false
	segment := 1
	screen := newScreen(t, func(look Look) {
		look.Switch("Lamp", &lamp)
		look.Checkbox("Fog horn", &fog)
		look.RadioButton("North", true)
		look.RadioButton("South", false)
		look.SegmentedButtons("view", []string{"Day", "Night"}, &segment)
	})
	screen.context.SetScreenReader(true)
	screen.frame()
	if node := screen.read("Lamp"); node.Role != ui.SwitchRole || !node.Checked {
		t.Fatalf("a switch, on: %+v", node)
	}
	if node := screen.read("Fog horn"); node.Role != ui.CheckboxRole || node.Checked {
		t.Fatalf("a checkbox, not ticked: %+v", node)
	}
	if north, south := screen.read("North"), screen.read("South"); north.Role != ui.RadioRole || !north.Checked || south.Checked {
		t.Fatalf("radio buttons, the chosen one checked: %+v %+v", north, south)
	}
	if night, day := screen.read("Night"), screen.read("Day"); !night.Selected || day.Selected {
		t.Fatalf("the chosen segment is selected: %+v %+v", night, day)
	}
	screen.context.ScreenReaderAct(screen.read("Lamp").ID, ui.ReaderPress)
	screen.frame()
	if lamp {
		t.Fatal("a double tap flips the switch, as a touch does")
	}
}

func TestNavigationIsReadWithWhereThePersonIs(t *testing.T) {
	selected := 0
	places := []Destination{{Icon: "home", Label: "Home"}, {Icon: "map", Label: "Map"}}
	tab := 1
	screen := newScreen(t, func(look Look) {
		look.TopAppBar(TopBar{Title: "Lighthouses", Back: true})
		look.Tabs("views", []string{"Near", "All"}, &tab)
		look.NavigationBar(places, &selected)
	})
	screen.context.SetScreenReader(true)
	screen.frame()
	if node := screen.read("Lighthouses"); node.Role != ui.HeadingRole {
		t.Fatalf("the bar's title is a heading: %+v", node)
	}
	if node := screen.read("Back"); node.Role != ui.ButtonRole || !node.Pressable {
		t.Fatalf("the arrow says what it does, not its icon: %+v", node)
	}
	if node := screen.read("All"); node.Role != ui.TabRole || !node.Selected {
		t.Fatalf("the tab shown: %+v", node)
	}
	if home, place := screen.read("Home"), screen.read("Map"); !home.Selected || place.Selected || home.Role != ui.TabRole {
		t.Fatalf("the destination shown is selected: %+v %+v", home, place)
	}
	screen.context.ScreenReaderAct(screen.read("Map").ID, ui.ReaderPress)
	screen.frame()
	if selected != 1 {
		t.Fatal("a double tap goes to the destination")
	}
}

func TestAFieldIsReadByItsLabelAndTypedInByTheReader(t *testing.T) {
	name := "Keeper"
	screen := newScreen(t, func(look Look) {
		look.TextField("name", &name, Field{Label: "Name", Placeholder: "Who keeps it"})
	})
	screen.context.SetScreenReader(true)
	screen.frame()
	if node := screen.read("Name"); node.Role != ui.FieldRole || node.Value != "Keeper" {
		t.Fatalf("the label over the box names the field: %+v", node)
	}
	screen.context.ScreenReaderSetText(screen.read("Name").ID, "Lamp keeper")
	screen.frame()
	if name != "Lamp keeper" {
		t.Fatalf("the screen reader types in the field: %q", name)
	}
}

func TestADialogHidesTheScreenBehindIt(t *testing.T) {
	screen := newScreen(t, func(look Look) {
		look.FilledButton("plus", "Add one")
		look.Dialog("confirm", Dialog{Title: "Put the lamp out?", Confirm: "Put it out", Dismiss: "Keep it"})
	})
	screen.context.SetScreenReader(true)
	screen.frame()
	screen.frame() // the dialog arrives
	for _, node := range screen.context.ScreenReaderNodes() {
		if node.Label == "Add one" {
			t.Fatal("a button behind the dialog is not read")
		}
	}
	if node := screen.read("Put the lamp out?"); node.Role != ui.HeadingRole {
		t.Fatalf("the dialog's title is a heading: %+v", node)
	}
	screen.read("Put it out")
}

func TestButtonsWithAnIconAloneSayWhatTheyDo(t *testing.T) {
	menu := &FloatingMenu{Icon: "plus", Label: "Create", Actions: []Destination{{Icon: "flag", Label: "A beacon"}}}
	tools := []Destination{{Icon: "pencil", Label: "Edit"}, {Icon: "share-2", Label: "Share"}}
	screen := newScreen(t, func(look Look) {
		look.Menu("ellipsis-vertical", []string{"Tides"})
		look.SplitButton("send", "Send")
		look.FloatingToolbar("tools", tools, 0)
		look.FloatingMenu("create", menu)
	})
	screen.context.SetScreenReader(true)
	screen.frame()
	screen.read("Send")
	screen.read("Edit")
	screen.read("Share")
	screen.read("Create")
	for _, node := range screen.context.ScreenReaderNodes() {
		if node.Label == "chevron-down" || node.Label == "ellipsis-vertical" {
			t.Fatalf("an icon's name is read: %+v", node)
		}
	}
	options := 0
	for _, node := range screen.context.ScreenReaderNodes() {
		if node.Label == "More options" {
			options++
		}
	}
	if options != 2 {
		t.Fatalf("the menu's button and the split button's arrow read %q, found %d", "More options", options)
	}
	screen.context.ScreenReaderAct(screen.read("Create").ID, ui.ReaderPress)
	screen.frame()
	if !menu.Open {
		t.Fatal("a double tap opens the floating menu")
	}
	screen.read("Close")
	screen.context.ScreenReaderAct(screen.read("More options").ID, ui.ReaderPress)
	screen.frame()
	screen.read("Tides")
	screen.read("More options") // the menu's button, while its menu is open
}
