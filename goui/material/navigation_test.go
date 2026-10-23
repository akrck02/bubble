package material

import (
	"testing"

	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

func TestTheNavigationBarGoesWhereItIsTouched(t *testing.T) {
	selected, changes := 0, 0
	places := []Destination{{Icon: "home", Label: "Home"}, {Icon: "map", Label: "Map"}, {Icon: "settings", Label: "Settings"}}
	screen := newScreen(t, func(look Look) {
		if look.NavigationBar(places, &selected) {
			changes++
		}
	})
	screen.frame()
	screen.tap("navigate Map")
	screen.tap("navigate Map")
	if selected != 1 || changes != 1 {
		t.Fatalf("selected %d after %d changes; want Map, changed once", selected, changes)
	}
}

func TestBarsReachUnderTheSystemsBars(t *testing.T) {
	var topHeight, bottomHeight float32
	screen := newScreen(t, func(look Look) {
		look.Context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow()}, func() {
			look.TopAppBar(TopBar{Title: "Inbox"})
			look.Context.Spacer()
			selected := 0
			look.NavigationBar([]Destination{{Icon: "home", Label: "Home"}}, &selected)
		})
	})
	screen.context.SetInsets(ui.Insets{Top: 24, Bottom: 48})
	screen.frame()
	if area, found := screen.context.Area("navigate Home"); found {
		bottomHeight = 800 - area.Y
	}
	topHeight = 64 + screen.context.Insets().Top
	if topHeight != 88 || bottomHeight < 48 {
		t.Fatalf("top bar %v tall, navigation reaching %v from the bottom; want 88 and below the gesture bar", topHeight, bottomHeight)
	}
}

func TestTheBackArrowSaysItWasTouched(t *testing.T) {
	wentBack := false
	screen := newScreen(t, func(look Look) {
		if look.TopAppBar(TopBar{Title: "Details", Back: true}) {
			wentBack = true
		}
	})
	screen.frame()
	screen.tap("icon arrow-left")
	if !wentBack {
		t.Fatal("touching the back arrow did not say so")
	}
}

func TestTheCrossLeavesAModeAsTheArrowLeavesAScreen(t *testing.T) {
	left := false
	screen := newScreen(t, func(look Look) {
		if look.TopAppBar(TopBar{Title: "2 chosen", Back: true, Close: true}) {
			left = true
		}
	})
	screen.frame()
	screen.tap("icon x")
	if !left {
		t.Fatal("touching the cross did not say so")
	}
}

func TestOnItsSideTheScreenUsesARailAndKeepsClearOfTheEdges(t *testing.T) {
	selected := 0
	screen := newScreen(t, func(look Look) {
		if !look.Wide() {
			t.Fatal("a screen wider than tall is not wide")
		}
		look.Context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Grow()}, func() {
			look.NavigationRail([]Destination{{Icon: "home", Label: "Home"}, {Icon: "map", Label: "Map"}}, &selected)
			look.Body(true, func() { look.FilledButton("", "Inside") })
		})
	})
	screen.state.WindowWidth, screen.state.WindowHeight = 800, 400
	screen.context.SetInsets(ui.Insets{Top: 24, Bottom: 20, Left: 30, Right: 40})
	screen.frame()
	screen.tap("navigate Map")
	if selected != 1 {
		t.Fatal("the rail did not go to Map")
	}
	rail, _ := screen.context.Area("navigate Home")
	button, _ := screen.context.Area("button Inside")
	if rail.X < 30 || button.X < 30+80 {
		t.Fatalf("rail item at %v, content at %v: want them clear of the 30-pixel cutout and the rail", rail.X, button.X)
	}
}

func TestADrawerChoosesAndCloses(t *testing.T) {
	selected, open, closings := 0, true, 0
	places := []Destination{{Icon: "inbox", Label: "Inbox"}, {Icon: "send", Label: "Sent"}}
	screen := newScreen(t, func(look Look) {
		if open && look.NavigationDrawer("menu", "Mail", places, &selected) {
			open = false
			closings++
		}
	})
	screen.frame()
	screen.tap("menu", "drawer Sent")
	if selected != 1 || open || closings != 1 {
		t.Fatalf("selected %d, open %v: want Sent chosen and the drawer closed", selected, open)
	}
}

func TestTheBackGoesToTheScreenWhenNoDialogIsOpen(t *testing.T) {
	dialog, wentBack := true, false
	screen := newScreen(t, func(look Look) {
		if dialog {
			dialog = look.Dialog("ask", Dialog{Title: "Sure?", Confirm: "Yes"}) == Open
		}
		if look.Back() {
			wentBack = true
		}
	})
	screen.frame()
	screen.state.Feed().PressBack()
	screen.frame()
	if dialog || wentBack {
		t.Fatalf("dialog open %v, screen went back %v: the dialog closes first", dialog, wentBack)
	}
	screen.state.Feed().PressBack()
	screen.frame()
	if !wentBack {
		t.Fatal("the screen did not go back")
	}
}

func TestALargeBarShrinksIntoASmallOneAsTheContentScrolls(t *testing.T) {
	scrolled := float32(0)
	screen := newScreen(t, func(look Look) {
		look.TopAppBar(TopBar{Title: "Lighthouses", Large: true, Scrolled: scrolled})
		look.Context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Grow()}, nil)
	})
	heightOfBar := func() float32 {
		screen.frame()
		for _, command := range screen.context.Drawing().Commands {
			if command.Kind == draw2d.RectangleKind && command.Area.Y == 0 && command.Area.Width == 400 {
				return command.Area.Height
			}
		}
		t.Fatal("no bar across the top")
		return 0
	}
	if height := heightOfBar(); height != 152 {
		t.Fatalf("large bar %v tall at rest; want 152", height)
	}
	scrolled = 44
	if height := heightOfBar(); height != 108 {
		t.Fatalf("large bar %v tall half way; want 108", height)
	}
	scrolled = 400
	if height := heightOfBar(); height != 64 {
		t.Fatalf("large bar %v tall with the content well under it; want a small bar's 64", height)
	}
}

func TestALargeBarAtRestShowsNoSmallTitle(t *testing.T) {
	screen := newScreen(t, func(look Look) { look.TopAppBar(TopBar{Title: "Lighthouses", Large: true}) })
	screen.frame()
	for _, command := range screen.context.Drawing().Commands {
		if command.Kind == draw2d.TextKind && command.Size == 22 && command.Fill[3] > 0.01 {
			t.Fatalf("the small title shows with opacity %v on a large bar at rest", command.Fill[3])
		}
	}
}

func TestScreensAreGoneToAndBackFrom(t *testing.T) {
	var navigation Navigation
	navigation.Reset("list")
	drawn := map[string]bool{}
	screen := newScreen(t, func(look Look) {
		clear(drawn)
		look.Screens("lights", &navigation, func(key string) {
			drawn[key] = true
			look.Text("screen "+key, 16, draw2d.Regular, look.Roles.OnSurface)
		})
		look.GoBack(&navigation)
	})
	screen.clocked()
	screen.frame()
	navigation.Go("seal rock")
	screen.frame()
	if !drawn["list"] || !drawn["seal rock"] {
		t.Fatalf("drawn %v while going to a screen; want both, the new one sliding in", drawn)
	}
	for range 30 {
		screen.frame()
	}
	if drawn["list"] || !drawn["seal rock"] || navigation.Depth() != 2 {
		t.Fatalf("drawn %v at depth %d once there; want the new screen alone", drawn, navigation.Depth())
	}
	screen.state.Feed().PressBack()
	screen.frame()
	if navigation.Current() != "list" || !drawn["seal rock"] {
		t.Fatalf("on %q, drawn %v after the back; want the list, the details sliding out", navigation.Current(), drawn)
	}
	for range 30 {
		screen.frame()
	}
	screen.state.Feed().PressBack()
	screen.frame()
	if navigation.Current() != "list" || navigation.Back() {
		t.Fatal("the first screen was gone back from")
	}
}

func TestANavigationStartsAtItsZeroValuesScreen(t *testing.T) {
	var navigation Navigation
	navigation.Go("details")
	if navigation.Depth() != 2 || !navigation.Back() || navigation.Current() != "" {
		t.Fatal("going back from the first screen gone to does not return to the zero value's")
	}
}
