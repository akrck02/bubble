package material

import (
	"fmt"
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
	"math"
)

// TopBar is what the top app bar shows.
type TopBar struct {
	Title    string
	Subtitle string // under the title, smaller; none when empty
	Back     bool   // an arrow before the title, to go back
	// Close makes that arrow a cross, for leaving a mode (rows chosen) rather than a screen.
	Close bool
	// Actions draws what goes at the bar's end, usually IconButtons; none when nil.
	Actions func()
	// OverScene paints the bar see-through (Roles.BarScrim), over a scene behind the interface.
	OverScene bool
	// Large writes the title big under the arrow and the actions, for a screen's first level.
	Large bool
	// Scrolled is how far the content under the bar has scrolled up (ui.Context.ScrollOffset): a
	// large bar shrinks into a small one as it goes, and any bar takes a container's colour once
	// content is under it.
	Scrolled float32
}

// largeCollapse is how far content scrolls for a large top bar to become a small one.
const largeCollapse = 88

// TopAppBar is the bar at the top, painted up under the status bar, so the screen reaches the top
// edge as Android's own apps do. Inside Body, it keeps clear of a cutout at the sides too. It returns true when the back arrow is touched.
func (look Look) TopAppBar(bar TopBar) (back bool) {
	insets := look.Context.Insets()
	roles := look.Roles
	fill := roles.Surface
	if bar.Scrolled > 0 {
		fill = roles.SurfaceContainer // content goes under it
	}
	if bar.OverScene {
		fill = roles.BarScrim
	}
	if bar.Large {
		return look.largeTopAppBar(bar, fill)
	}
	layout := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(64 + insets.Top),
		Padding: ui.Padding{Left: 4, Top: insets.Top, Right: 4}, Gap: 4, AlignY: ui.Center}
	look.Context.Frame(layout, ui.Look{Fill: fill}, func() {
		if bar.Back {
			back = look.NamedIconButton(bar.backIcon(), look.backWord(bar))
		} else {
			look.Context.Box(ui.Layout{Width: ui.Fixed(12)}, nil)
		}
		look.Context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Gap: 2}, func() {
			look.Text(bar.Title, 22, draw2d.Regular, roles.OnSurface)
			look.Context.Describe(ui.Description{Role: ui.HeadingRole})
			if bar.Subtitle != "" {
				look.Text(bar.Subtitle, 14, draw2d.Regular, roles.OnSurfaceVariant)
			}
		})
		if bar.Actions != nil {
			bar.Actions()
		}
	})
	return back
}

// backWord is what a screen reader calls the bar's arrow, or its cross.
func (look Look) backWord(bar TopBar) string {
	if bar.Close {
		return look.words().Close
	}
	return look.words().Back
}

func (bar TopBar) backIcon() string {
	if bar.Close {
		return "x"
	}
	return "arrow-left"
}

// Destination is a place the navigation bar leads to: an icon and a short label.
type Destination struct {
	Icon, Label string
	// Badge marks the icon: a count ("3"), or a dot with "•"; none when empty.
	Badge string
}

// largeTopAppBar is the large top bar: the arrow and the actions on a row, the title big under them.
func (look Look) largeTopAppBar(bar TopBar, fill draw2d.Color) (back bool) {
	insets := look.Context.Insets()
	roles := look.Roles
	collapsed := clampFloat(bar.Scrolled/largeCollapse, 0, 1)
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Fixed(64 + largeCollapse*(1-collapsed) + insets.Top),
		Padding: ui.Padding{Left: 4, Top: insets.Top, Right: 4, Bottom: 24 * (1 - collapsed)}}
	look.Context.Frame(layout, ui.Look{Fill: fill, Clip: true}, func() {
		look.Context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(64), Gap: 4, AlignY: ui.Center}, func() {
			if bar.Back {
				back = look.NamedIconButton(bar.backIcon(), look.backWord(bar))
			} else {
				look.Context.Box(ui.Layout{Width: ui.Fixed(12)}, nil)
			}
			// the small title comes in as the large one goes; a screen reader reads the one shown more
			look.Context.Box(ui.Layout{Width: ui.Grow()}, func() {
				look.Text(bar.Title, 22, draw2d.Regular, faded(roles.OnSurface, max(collapsed*2-1, 0)))
				look.Context.Describe(ui.Description{Role: ui.HeadingRole, Hidden: collapsed < 0.5})
			})
			if bar.Actions != nil {
				bar.Actions()
			}
		})
		look.Context.Spacer() // the large title keeps to the bar's foot as it shrinks
		look.Context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.PadXY(12, 0), Gap: 4}, func() {
			shown := 1 - collapsed
			look.Text(bar.Title, 32, draw2d.Regular, faded(roles.OnSurface, shown))
			look.Context.Describe(ui.Description{Role: ui.HeadingRole, Hidden: collapsed >= 0.5})
			if bar.Subtitle != "" {
				look.Text(bar.Subtitle, 14, draw2d.Regular, faded(roles.OnSurfaceVariant, shown))
			}
		})
	})
	return back
}

// destinationIcon is a destination's icon with its badge over the icon's top end, as Material
// places it.
func (look Look) destinationIcon(place Destination, ink draw2d.Color) {
	look.Context.IconSized(place.Icon, 24, ink)
	if place.Badge == "" {
		return
	}
	roles := look.Roles
	// a negative gap lays the badge over the icon's corner
	look.Context.Box(ui.Layout{Width: ui.Fixed(0), Height: ui.Fixed(24)}, nil)
	if place.Badge == "•" {
		look.Context.Box(ui.Layout{Height: ui.Fixed(24), Padding: ui.Padding{Left: -8, Top: 2}}, func() {
			look.Context.Frame(ui.Layout{Width: ui.Fixed(6), Height: ui.Fixed(6)}, ui.Look{Fill: roles.Error, Radius: 3}, nil)
		})
		return
	}
	look.Context.Box(ui.Layout{Height: ui.Fixed(24), Padding: ui.Padding{Left: -12}}, func() {
		look.Context.Frame(ui.Layout{Height: ui.Fixed(16), Padding: ui.PadXY(4, 0), AlignY: ui.Center}, ui.Look{Fill: roles.Error, Radius: 8}, func() {
			look.Text(place.Badge, 11, draw2d.Medium, roles.OnError)
		})
	})
}

// NavigationBar is the bar at the bottom, painted down under the gesture bar: three to five
// destinations, the selected one with its icon on a pill. It changes selected and returns true
// when another destination is touched.
func (look Look) NavigationBar(destinations []Destination, selected *int) bool {
	insets := look.Context.Insets()
	layout := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(80 + insets.Bottom),
		Padding: ui.Padding{Left: insets.Left, Top: 12, Right: insets.Right, Bottom: 16 + insets.Bottom}}
	changed := false
	look.Context.Frame(layout, ui.Look{Fill: look.Roles.SurfaceContainer}, func() {
		for index, place := range destinations {
			if look.navigationItem(place, index == *selected) && index != *selected {
				*selected, changed = index, true
			}
		}
	})
	return changed
}

func (look Look) navigationItem(place Destination, selected bool) bool {
	roles := look.Roles
	ink, label, weight := roles.OnSurfaceVariant, roles.OnSurfaceVariant, draw2d.Medium
	if selected {
		ink, label, weight = roles.OnSecondaryContainer, roles.OnSurface, draw2d.Bold
	}
	// the pill grows out from the icon as its destination is chosen
	shown := look.Context.Animate("pill "+place.Label, boolFloat(selected), arrivingMotion)
	pill := ui.Look{Fill: stateLayer(roles.SecondaryContainer, shown*roles.SecondaryContainer[3]), Radius: 16}
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Share(1), Height: ui.Grow(), Gap: 4, AlignX: ui.Center}
	touched := look.Context.Clickable("navigate "+place.Label, layout, ui.Looks{}, func() {
		indicator := ui.Layout{Width: ui.Fixed(32 + 32*shown), Height: ui.Fixed(32), AlignX: ui.Center, AlignY: ui.Center}
		look.Context.Frame(indicator, pill, func() { look.destinationIcon(place, ink) })
		look.Text(place.Label, 12, weight, label)
	})
	// read by its label, then its badge, rather than the badge first as it is drawn
	look.Context.Describe(ui.Description{Role: ui.TabRole, Label: place.Label, Value: place.Badge, Selected: selected})
	return touched
}

// StatusBarScrim paints the strip under the status bar alone, for a screen whose top is a scene
// rather than a bar: the system's icons stay readable over it.
func (look Look) StatusBarScrim() {
	top := look.Context.Insets().Top
	if top <= 0 {
		return
	}
	look.Context.Frame(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(top)}, ui.Look{Fill: look.Roles.BarScrim}, nil)
}

// NavigationRail is the navigation bar's form on a wide screen (a phone on its side, a tablet): the
// destinations down the left edge, painted under the status bar and beside a camera's cutout. It
// changes selected and returns true when another destination is touched. Wide tells when to use it.
func (look Look) NavigationRail(destinations []Destination, selected *int) bool {
	insets := look.Context.Insets()
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Fixed(80 + insets.Left), Height: ui.Grow(),
		Padding: ui.Padding{Left: insets.Left, Top: insets.Top + 12, Bottom: insets.Bottom + 12}, Gap: 12}
	changed := false
	look.Context.Frame(layout, ui.Look{Fill: look.Roles.SurfaceContainer}, func() {
		for index, place := range destinations {
			if look.railItem(place, index == *selected) && index != *selected {
				*selected, changed = index, true
			}
		}
	})
	return changed
}

func (look Look) railItem(place Destination, selected bool) bool {
	roles := look.Roles
	pill, ink, label, weight := ui.Look{Radius: 16}, roles.OnSurfaceVariant, roles.OnSurfaceVariant, draw2d.Medium
	if selected {
		pill.Fill, ink, label, weight = roles.SecondaryContainer, roles.OnSecondaryContainer, roles.OnSurface, draw2d.Bold
	}
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Fixed(56), Gap: 4, AlignX: ui.Center}
	touched := look.Context.Clickable("navigate "+place.Label, layout, ui.Looks{}, func() {
		indicator := ui.Layout{Width: ui.Fixed(56), Height: ui.Fixed(32), AlignX: ui.Center, AlignY: ui.Center}
		look.Context.Frame(indicator, pill, func() { look.destinationIcon(place, ink) })
		look.Text(place.Label, 12, weight, label)
	})
	// read by its label, then its badge, rather than the badge first as it is drawn
	look.Context.Describe(ui.Description{Role: ui.TabRole, Label: place.Label, Value: place.Badge, Selected: selected})
	return touched
}

// Wide tells whether the screen is wide enough for a navigation rail instead of a bar: wider than
// tall, or 600 pixels across (Material's medium width, a small tablet).
func (look Look) Wide() bool {
	width, height := look.Context.Size()
	return width > height || width >= 600
}

// Body holds a screen's content between the bars: it grows to fill them and keeps its content
// clear of what the system draws at the sides (a camera's cutout, a navigation bar on the side).
// Beside a navigation rail (besideRail), the rail keeps the left side clear, and with no
// navigation bar under it, Body keeps clear of the gesture bar itself.
func (look Look) Body(besideRail bool, content func()) {
	insets := look.Context.Insets()
	padding := ui.Padding{Left: insets.Left, Right: insets.Right}
	if besideRail {
		padding.Left, padding.Bottom = 0, insets.Bottom
	}
	look.Context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow(), Padding: padding}, content)
}

// Back answers the system's back (a phone's back gesture) for the screen: true when there was one,
// and the program then goes back (to the first destination, to the list a detail came from).
// Dialogs and sheets answer it first; a back nobody answers takes the program to the background.
func (look Look) Back() bool { return look.Context.Back() }

// Enter shows a screen rising gently into place when it is first shown, as Material's screens
// arrive: give each screen its own key.
func (look Look) Enter(key string, content func()) {
	arrived := look.Context.Appear("enter "+key, arrivingMotion)
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow(), Padding: ui.Padding{Top: (1 - arrived) * 24}}
	look.Context.Box(layout, content)
}

// NavigationBarHeight is a navigation bar's height without the gesture bar under it, for a snackbar
// or a floating button to keep above.
const NavigationBarHeight = 80

// NavigationDrawer slides in from the left over the dimmed screen, with a title and the program's
// destinations, for more places than a navigation bar holds. It changes selected when one is
// touched, and returns true when it should close: a destination chosen, a touch outside, the
// system's back. name tells drawers apart.
func (look Look) NavigationDrawer(name, title string, destinations []Destination, selected *int) (closed bool) {
	roles := look.Roles
	insets := look.Context.Insets()
	width, _ := look.Context.Size()
	drawerWidth := min(360, width-56)
	arrived := look.Context.Appear("drawer "+name, arrivingMotion)
	options := ui.LayerOptions{AlignX: ui.Start, AlignY: ui.Start, Padding: ui.Padding{Left: -(1 - arrived) * drawerWidth},
		Width: ui.Fixed(drawerWidth), Height: ui.Grow(), Hidden: 1 - arrived,
		Leave: ui.Leave{Duration: leavingMotion, X: -drawerWidth}} // back off the side it came from
	dismissed := look.Context.Modal(name, options, func() {
		sheet := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow(),
			Padding: ui.Padding{Left: 12 + insets.Left, Top: 12 + insets.Top, Right: 12, Bottom: 12 + insets.Bottom}, Gap: 4}
		look.Context.Frame(sheet, ui.Look{Fill: roles.SurfaceContainerLow, Radius: 16, Shadow: true}, func() {
			look.Context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.Padding{Left: 16, Top: 16, Bottom: 16}}, func() {
				look.Text(title, 14, draw2d.Medium, roles.OnSurfaceVariant)
			})
			for index, place := range destinations {
				if look.drawerItem(place, index == *selected) {
					*selected, closed = index, true
				}
			}
		})
	})
	return closed || dismissed
}

func (look Look) drawerItem(place Destination, selected bool) bool {
	roles := look.Roles
	normal, ink := ui.Look{Radius: 28}, roles.OnSurfaceVariant
	if selected {
		normal, ink = ui.Look{Fill: roles.SecondaryContainer, Radius: 28}, roles.OnSecondaryContainer
	}
	layout := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(56), Padding: ui.Padding{Left: 16, Right: 24}, Gap: 12, AlignY: ui.Center}
	touched := look.Context.Clickable("drawer "+place.Label, layout, pressedOnly(normal, stateLayer(ink, 0.12)), func() {
		look.Context.IconSized(place.Icon, 24, ink)
		look.Context.Box(ui.Layout{Width: ui.Grow()}, func() { look.Text(place.Label, 14, draw2d.Medium, ink) })
		if place.Badge != "" {
			look.Text(place.Badge, 14, draw2d.Medium, ink)
		}
	})
	look.Context.Describe(ui.Description{Selected: selected})
	return touched
}

// Tabs are the primary tabs under a top bar: labels sharing the width, an underline under the
// chosen one that slides to the next one chosen. It changes selected and returns true when another
// is touched.
func (look Look) Tabs(name string, labels []string, selected *int) bool {
	at := look.Context.Animate(name+" underline", float32(*selected), arrivingMotion)
	return look.tabStrip(name, labels, selected, at)
}

// tabStrip draws the tabs with the underline at a position, in tabs: between two while a page is
// dragged from one to the next.
func (look Look) tabStrip(name string, labels []string, selected *int, at float32) bool {
	roles := look.Roles
	changed := false
	look.Context.Push(name)
	defer look.Context.Pop()
	look.Context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow()}, func() {
		look.Context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(45)}, func() {
			for index, label := range labels {
				chosen := index == *selected
				ink := roles.OnSurfaceVariant
				if chosen {
					ink = roles.Primary
				}
				layout := ui.Layout{Width: ui.Share(1), Height: ui.Grow(), AlignX: ui.Center, AlignY: ui.Center}
				touched := look.Context.Clickable(label, layout, ui.Looks{Pressed: ui.Look{Fill: stateLayer(ink, 0.1)}}, func() {
					look.Text(label, 14, draw2d.Medium, ink)
				})
				look.Context.Describe(ui.Description{Role: ui.TabRole, Selected: chosen})
				if touched && !chosen {
					*selected, changed = index, true
				}
			}
		})
		// the underline, painted where the position says, as wide as a short label
		look.Context.Paint(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(3)}, func(area draw2d.Rectangle, drawing *draw2d.List) {
			if len(labels) == 0 {
				return
			}
			tab := area.Width / float32(len(labels))
			width := min(tab*0.6, 96)
			x := area.X + at*tab + (tab-width)/2
			drawing.Rectangle(draw2d.Rectangle{X: x, Y: area.Y, Width: width, Height: area.Height}, 1.5, roles.Primary, draw2d.Transparent, 0)
		})
		look.Divider()
	})
	return changed
}

// Pager is what tabs with swiping pages remember: the page chosen, and where the pages were shown
// in the last frame (between two while one is dragged).
type Pager struct {
	Selected int
	shown    float32
}

// TabPager is tabs over pages a finger swipes from side to side: the underline follows the pages
// as they move, a page let go past a third of the way (or flicked) turns to the next, and touching
// a tab slides to its page. page draws one page, given its index; lists inside it still scroll up
// and down. name tells pagers apart.
func (look Look) TabPager(name string, labels []string, pager *Pager, page func(index int)) {
	context := look.Context
	look.tabStrip(name, labels, &pager.Selected, pager.shown)
	context.SwipeBox(name+" pages", ui.Layout{Width: ui.Grow(), Height: ui.Grow()}, ui.Horizontal, func(swipe ui.Swipe) {
		width := swipe.Area.Width
		if width <= 0 {
			width, _ = context.Size()
		}
		dragged := -swipe.Offset / width // in pages: a finger moving left shows the next
		last := float32(len(labels) - 1)
		if swipe.Released {
			next := pager.Selected
			switch {
			case dragged > 0.33 || swipe.Velocity < -600:
				next++
			case dragged < -0.33 || swipe.Velocity > 600:
				next--
			}
			pager.Selected = int(clampFloat(float32(next), 0, last))
		}
		target, motion := float32(pager.Selected), 3*arrivingMotion/2
		if swipe.Held {
			// past the first and the last page the pages give only a little
			target, motion = clampFloat(float32(pager.Selected)+dragged, -0.1, last+0.1), 0
		}
		pager.shown = context.Animate(name+" shown", target, motion)
		first := int(math.Floor(float64(pager.shown)))
		strip := ui.Layout{Width: ui.Grow(), Height: ui.Grow(), Padding: ui.Padding{Left: -(pager.shown - float32(first)) * width}}
		context.Box(strip, func() {
			for index := first; index <= first+1; index++ {
				if index < 0 || index > int(last) {
					context.Box(ui.Layout{Width: ui.Fixed(width), Height: ui.Grow()}, nil)
					continue
				}
				context.Push(fmt.Sprint(name, " page ", index))
				context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Fixed(width), Height: ui.Grow()}, func() { page(index) })
				context.Pop()
			}
		})
	})
}

// Navigation is a program's screens, one gone to from another: a list, then one of its rows'
// details, then a form about it. The program keeps one and names its screens as it likes; Screens
// draws the one on top. The zero value starts at "".
type Navigation struct {
	stack   []string
	leaving string // the screen going away, while it slides out
	forward bool   // the last change went to a screen, rather than back from one
	moves   int    // counts the changes, which Screens eases between
}

// Current is the screen on top.
func (navigation *Navigation) Current() string {
	if len(navigation.stack) == 0 {
		return ""
	}
	return navigation.stack[len(navigation.stack)-1]
}

// Depth is how many screens there are, the first one included.
func (navigation *Navigation) Depth() int { return max(len(navigation.stack), 1) }

// Go goes to a screen, over the one shown.
func (navigation *Navigation) Go(screen string) {
	navigation.leaving, navigation.forward = navigation.Current(), true
	if len(navigation.stack) == 0 {
		navigation.stack = []string{""} // the zero value's first screen
	}
	navigation.stack = append(navigation.stack, screen)
	navigation.moves++
}

// Back goes back to the screen before, and tells whether there was one.
func (navigation *Navigation) Back() bool {
	if len(navigation.stack) < 2 {
		return false
	}
	navigation.leaving, navigation.forward = navigation.Current(), false
	navigation.stack = navigation.stack[:len(navigation.stack)-1]
	navigation.moves++
	return true
}

// Reset starts again at one screen, with nothing to go back to: a destination of the navigation
// bar.
func (navigation *Navigation) Reset(screen string) {
	navigation.stack, navigation.leaving = []string{screen}, ""
}

// Screens draws the screen on top of a navigation: one gone to slides in from the end over the one
// it was gone to from, and slides back out when the program goes back. screen draws a screen by
// its name. The program answers the system's back with GoBack, after what goes over the screens.
func (look Look) Screens(name string, navigation *Navigation, screen func(key string)) {
	context := look.Context
	target := float32(navigation.moves)
	at := context.Animate(name+" moves", target, arrivingMotion)
	travelled := clampFloat(1-(target-at), 0, 1)
	if travelled >= 1 || navigation.moves == 0 {
		look.screenBox(name, navigation.Current(), screen, ui.Grow())
		return
	}
	width, _ := context.Size()
	// the two screens side by side, the strip moved as far as the change has gone
	first, second := navigation.leaving, navigation.Current()
	shift := travelled
	if !navigation.forward {
		first, second, shift = navigation.Current(), navigation.leaving, 1-travelled
	}
	strip := ui.Layout{Width: ui.Grow(), Height: ui.Grow(), Padding: ui.Padding{Left: -shift * width}}
	context.Frame(ui.Layout{Width: ui.Grow(), Height: ui.Grow()}, ui.Look{Clip: true}, func() {
		context.Box(strip, func() {
			look.screenBox(name, first, screen, ui.Fixed(width))
			look.screenBox(name, second, screen, ui.Fixed(width))
		})
	})
}

// screenBox draws one screen in a scope of its own, so two screens drawn at once keep their names
// apart.
func (look Look) screenBox(name, key string, screen func(key string), width ui.Sizing) {
	look.Context.Push(name + " screen " + key)
	defer look.Context.Pop()
	look.Context.Box(ui.Layout{Direction: ui.Vertical, Width: width, Height: ui.Grow()}, func() { screen(key) })
}

// GoBack answers the system's back by going back a screen, while there is one to go back to. Ask
// for it after drawing what goes over the screens, so that answers first.
func (look Look) GoBack(navigation *Navigation) bool {
	return navigation.Depth() > 1 && look.Back() && navigation.Back()
}
