package material

import (
	"fmt"
	"time"

	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// Buttons are 40 pixels tall pills with a 48-pixel target around them, as Material's are; the
// icon, when given, goes before the label. An empty icon draws none.

// FilledButton is the screen's main action, in the primary colour.
func (look Look) FilledButton(icon, label string) bool {
	scheme := look.Roles
	return look.button(icon, label, ui.Look{Fill: scheme.Primary, Radius: 20}, scheme.PrimaryContainer, scheme.OnPrimary, 24)
}

// TonalButton is a main action that should not stand out as much, on the secondary container.
func (look Look) TonalButton(icon, label string) bool {
	scheme := look.Roles
	normal := ui.Look{Fill: scheme.SecondaryContainer, Radius: 20}
	return look.button(icon, label, normal, stateLayer(scheme.OnSecondaryContainer, 0.16), scheme.OnSecondaryContainer, 24)
}

// OutlinedButton is a second action beside the main one.
func (look Look) OutlinedButton(icon, label string) bool {
	scheme := look.Roles
	normal := ui.Look{Border: scheme.Outline, BorderWidth: 1, Radius: 20}
	return look.button(icon, label, normal, stateLayer(scheme.Primary, 0.12), scheme.Primary, 24)
}

// TextButton is the quietest action: a dialog's "Cancel", a card's "More".
func (look Look) TextButton(icon, label string) bool {
	scheme := look.Roles
	return look.button(icon, label, ui.Look{Radius: 20}, stateLayer(scheme.Primary, 0.12), scheme.Primary, 12)
}

func (look Look) button(icon, label string, normal ui.Look, pressedFill, ink draw2d.Color, padding float32) bool {
	layout := ui.Layout{Height: ui.Fixed(40), Padding: ui.PadXY(padding, 0), Gap: 8, AlignY: ui.Center}
	return look.Context.Clickable("button "+label, layout, pressedOnly(normal, pressedFill), func() {
		if icon != "" {
			look.Context.IconSized(icon, 18, ink)
		}
		look.Text(label, 14, draw2d.Medium, ink)
	})
}

// IconButton is a round 48-pixel target with a 24-pixel icon, in the bars and beside list items.
func (look Look) IconButton(icon string) bool {
	ink := look.Roles.OnSurfaceVariant
	layout := ui.Layout{Width: ui.Fixed(48), Height: ui.Fixed(48), AlignX: ui.Center, AlignY: ui.Center}
	looks := ui.Looks{Pressed: ui.Look{Fill: stateLayer(ink, 0.12), Radius: 24}}
	return look.Context.Clickable("icon "+icon, layout, looks, func() { look.Context.IconSized(icon, 24, ink) })
}

// FloatingButton is the floating action button: 56 pixels, rounded square corners, raised over
// what scrolls beneath it. Place it at the end of the screen, above the navigation bar.
func (look Look) FloatingButton(icon string) bool {
	scheme := look.Roles
	normal := ui.Look{Fill: scheme.PrimaryContainer, Radius: 16, Shadow: true}
	layout := ui.Layout{Width: ui.Fixed(56), Height: ui.Fixed(56), AlignX: ui.Center, AlignY: ui.Center}
	return look.Context.Clickable("floating "+icon, layout, pressedOnly(normal, scheme.SecondaryContainer), func() {
		look.Context.IconSized(icon, 24, scheme.OnPrimaryContainer)
	})
}

// RoundButton is a round floating button of any size, over a map or a scene (zoom in, zoom out,
// change the view). primary fills it with the main colour; otherwise it is a raised surface.
func (look Look) RoundButton(icon string, size float32, primary bool) bool {
	scheme := look.Roles
	normal := ui.Look{Fill: scheme.SurfaceContainerLow, Border: scheme.OutlineVariant, BorderWidth: 1, Radius: size / 2, Shadow: true}
	ink := scheme.OnSurface
	if primary {
		normal, ink = ui.Look{Fill: scheme.Primary, Radius: size / 2, Shadow: true}, scheme.OnPrimary
	}
	layout := ui.Layout{Width: ui.Fixed(size), Height: ui.Fixed(size), AlignX: ui.Center, AlignY: ui.Center}
	return look.Context.Clickable("round "+icon, layout, pressedOnly(normal, scheme.SecondaryContainer), func() {
		look.Context.IconSized(icon, size*0.45, ink)
	})
}

// ExtendedFloatingButton is the floating action button with its action written beside the icon:
// clearer where the icon alone would not say it.
func (look Look) ExtendedFloatingButton(icon, label string) bool {
	scheme := look.Roles
	normal := ui.Look{Fill: scheme.PrimaryContainer, Radius: 16, Shadow: true}
	layout := ui.Layout{Height: ui.Fixed(56), Padding: ui.Padding{Left: 16, Right: 20}, Gap: 12, AlignY: ui.Center}
	return look.Context.Clickable("floating "+label, layout, pressedOnly(normal, scheme.SecondaryContainer), func() {
		look.Context.IconSized(icon, 24, scheme.OnPrimaryContainer)
		look.Text(label, 14, draw2d.Medium, scheme.OnPrimaryContainer)
	})
}

// Pieces Material 3 Expressive added for actions: a floating button that opens into a few actions,
// a button split from its menu, buttons joined in a group, and a toolbar floating over the content.

// FloatingMenu is a floating action button that opens into a few related actions stacked above
// it, its icon turning into a cross. A program keeps one, and draws it with Look.FloatingMenu.
type FloatingMenu struct {
	Icon    string        // the button's icon while closed
	Label   string        // what a screen reader reads for the button while closed ("Create")
	Actions []Destination // what it opens into: an icon and a label each
	// Above is how far over the bottom of the window it sits: a navigation bar's height
	// (NavigationBarHeight) when there is one.
	Above float32
	Open  bool
}

// FloatingMenu draws the button at the bottom end of the screen, over everything, and its actions
// while open. It returns the index of the action touched in this frame, or -1; touching one, the
// cross or the system's back closes it.
func (look Look) FloatingMenu(name string, menu *FloatingMenu) int {
	scheme := look.Roles
	insets := look.Context.Insets()
	if menu.Open && look.Context.Back() {
		menu.Open = false
	}
	chosen := -1
	padding := ui.Padding{Right: 16 + insets.Right, Bottom: 16 + insets.Bottom + menu.Above}
	look.Context.Overlay(name, ui.LayerOptions{AlignX: ui.End, AlignY: ui.End, Padding: padding}, func() {
		look.Context.Box(ui.Layout{Direction: ui.Vertical, Gap: 8, AlignX: ui.End}, func() {
			// closed, the actions sink and fade before they go
			shown := look.Context.Animate(name+" open", boolFloat(menu.Open), leavingMotion)
			for index, action := range menu.Actions {
				switch {
				case menu.Open && look.floatingAction(name, index, action):
					chosen, menu.Open = index, false
				case !menu.Open && shown > 0:
					look.leavingAction(action, shown)
				}
			}
			// open, the button is a smaller round cross in the main colour
			icon, size, normal := menu.Icon, float32(56), ui.Look{Fill: scheme.PrimaryContainer, Radius: 16, Shadow: true}
			ink := scheme.OnPrimaryContainer
			if menu.Open {
				icon, size, normal, ink = "x", 56, ui.Look{Fill: scheme.Primary, Radius: 28, Shadow: true}, scheme.OnPrimary
			}
			layout := ui.Layout{Width: ui.Fixed(size), Height: ui.Fixed(size), AlignX: ui.Center, AlignY: ui.Center}
			if look.Context.Clickable("toggle", layout, pressedOnly(normal, scheme.SecondaryContainer), func() {
				look.Context.IconSized(icon, 24, ink)
			}) {
				menu.Open = !menu.Open
			}
			look.Context.Describe(ui.Description{Label: look.floatingMenuName(menu)})
		})
	})
	return chosen
}

// floatingMenuName is what a screen reader reads for a floating menu's button: its label while
// closed (its icon's name if it has none), and "Close" while open, as it shows a cross.
func (look Look) floatingMenuName(menu *FloatingMenu) string {
	if menu.Open {
		return look.words().Close
	}
	return menu.Label
}

// floatingAction is one of an open floating menu's actions, rising in after the ones before it. It
// returns true when touched.
func (look Look) floatingAction(name string, index int, action Destination) bool {
	scheme := look.Roles
	arrived := look.Context.Appear(fmt.Sprint(name, " action ", index), time.Duration(150+50*index)*time.Millisecond)
	layout := ui.Layout{Height: ui.Fixed(56), Padding: ui.Padding{Left: 16, Right: 20, Top: (1 - arrived) * 16}, Gap: 12, AlignY: ui.Center}
	normal := ui.Look{Fill: scheme.PrimaryContainer, Radius: 28, Shadow: true}
	return look.Context.Clickable("action "+action.Label, layout, pressedOnly(normal, scheme.SecondaryContainer), func() {
		look.Context.IconSized(action.Icon, 24, scheme.OnPrimaryContainer)
		look.Text(action.Label, 16, draw2d.Medium, scheme.OnPrimaryContainer)
	})
}

// leavingAction is an action of a menu just closed, on its way out: shown is how much of it is
// left, from 1 to 0. It cannot be touched, and casts no shadow, which would not fade with it.
func (look Look) leavingAction(action Destination, shown float32) {
	scheme := look.Roles
	layout := ui.Layout{Height: ui.Fixed(56), Padding: ui.Padding{Left: 16, Right: 20, Top: (1 - shown) * 16}, Gap: 12, AlignY: ui.Center}
	ink := faded(scheme.OnPrimaryContainer, shown)
	look.Context.Frame(layout, ui.Look{Fill: faded(scheme.PrimaryContainer, shown), Radius: 28}, func() {
		look.Context.IconSized(action.Icon, 24, ink)
		look.Text(action.Label, 16, draw2d.Medium, ink)
	})
}

// SplitButton is a main action joined to a small button that opens its other choices (send, and
// schedule): it tells which half was touched.
func (look Look) SplitButton(icon, label string) (action, more bool) {
	scheme := look.Roles
	look.Context.Push("split " + label)
	defer look.Context.Pop()
	look.Context.Box(ui.Layout{Height: ui.Fixed(40), Gap: 2}, func() {
		leading := ui.Look{Fill: scheme.Primary, Radius: 20}
		layout := ui.Layout{Height: ui.Grow(), Padding: ui.Padding{Left: 16, Right: 12}, Gap: 8, AlignY: ui.Center}
		action = look.Context.Clickable("action", layout, pressedOnly(leading, scheme.PrimaryContainer), func() {
			if icon != "" {
				look.Context.IconSized(icon, 18, scheme.OnPrimary)
			}
			look.Text(label, 14, draw2d.Medium, scheme.OnPrimary)
		})
		trailing := ui.Look{Fill: scheme.Primary, Radius: 20}
		layout = ui.Layout{Width: ui.Fixed(40), Height: ui.Grow(), AlignX: ui.Center, AlignY: ui.Center}
		more = look.Context.Clickable("more", layout, pressedOnly(trailing, scheme.PrimaryContainer), func() {
			look.Context.IconSized("chevron-down", 18, scheme.OnPrimary)
		})
		look.Context.Describe(ui.Description{Label: look.words().MoreOptions})
	})
	return action, more
}

// ButtonGroup is a few choices joined side by side, the chosen one filled and fully rounded, the
// others squarer: Material 3 Expressive's connected group (a view, a size, an alignment). It
// changes selected and returns true when another is touched.
func (look Look) ButtonGroup(name string, labels []string, selected *int) bool {
	scheme := look.Roles
	changed := false
	look.Context.Push(name)
	defer look.Context.Pop()
	look.Context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(40), Gap: 2}, func() {
		for index, label := range labels {
			chosen := index == *selected
			// the chosen one rounds out fully; its corners ease as the choice moves
			roundness := look.Context.Animate("round "+label, boolFloat(chosen), switchMotion)
			normal := ui.Look{Fill: scheme.SecondaryContainer, Radius: 8 + 12*roundness}
			ink := scheme.OnSecondaryContainer
			if chosen {
				normal.Fill, ink = scheme.Primary, scheme.OnPrimary
			}
			layout := ui.Layout{Width: ui.Share(1), Height: ui.Grow(), AlignX: ui.Center, AlignY: ui.Center}
			touched := look.Context.Clickable(label, layout, pressedOnly(normal, stateLayer(ink, 0.2)), func() {
				look.Text(label, 14, draw2d.Medium, ink)
			})
			look.Context.Describe(ui.Description{Selected: chosen})
			if touched && !chosen {
				*selected, changed = index, true
			}
		}
	})
	return changed
}

// FloatingToolbar is a row of actions in a pill floating over the content at the bottom of the
// screen, for what is done to the content on show (edit, share, delete): each tool shows its icon,
// and a screen reader reads its label. above is how far over the bottom it sits. It returns the
// index of the tool touched in this frame, or -1.
func (look Look) FloatingToolbar(name string, tools []Destination, above float32) int {
	scheme := look.Roles
	insets := look.Context.Insets()
	touched := -1
	padding := ui.Padding{Bottom: 16 + insets.Bottom + above}
	look.Context.Overlay(name, ui.LayerOptions{AlignX: ui.Center, AlignY: ui.End, Padding: padding}, func() {
		bar := ui.Layout{Height: ui.Fixed(64), Padding: ui.PadXY(8, 0), Gap: 4, AlignY: ui.Center}
		look.Context.Frame(bar, ui.Look{Fill: scheme.PrimaryContainer, Radius: 32, Shadow: true}, func() {
			for index, tool := range tools {
				layout := ui.Layout{Width: ui.Fixed(48), Height: ui.Fixed(48), AlignX: ui.Center, AlignY: ui.Center}
				if look.Context.Clickable("tool "+tool.Icon, layout, ui.Looks{Pressed: ui.Look{Fill: stateLayer(scheme.OnPrimaryContainer, 0.12), Radius: 24}}, func() {
					look.Context.IconSized(tool.Icon, 24, scheme.OnPrimaryContainer)
				}) {
					touched = index
				}
				look.Context.Describe(ui.Description{Label: tool.Label})
			}
		})
	})
	return touched
}
