package material

import (
	"math"

	"github.com/akrck02/go-renderer/draw2d"
	"github.com/akrck02/go-renderer/input"
	ui "github.com/akrck02/go-ui"
)

// Stick is an on-screen stick, for walking or driving something with a thumb: the finger that
// comes down on it holds it, and where that finger is from its middle says where to go. Other
// fingers stay free (to look around, to press buttons). A program keeps one Stick and reads it in
// every step.
type Stick struct {
	Forward  float32 // from -1 (back) to 1 (ahead): the finger above the middle is ahead
	Sideways float32 // from -1 (left) to 1 (right)
	Held     bool
	Finger   int64 // the finger holding it, while Held

	area draw2d.Rectangle // where it was laid out, in the last frame
	size float32
}

// Stick draws the stick, its knob where the finger points it, and follows the fingers of this frame
// (ui.Context.Touches). size is its side, 128 or so; name tells it apart. It may sit anywhere: in a
// Push scope, a scrolling list or swiping pages, whose finger it takes while it steers. It asks for
// frames while held, so what it drives keeps moving under a finger that stays still.
func (look Look) Stick(name string, stick *Stick, size float32) {
	stick.Follow(look.Context.Touches())
	layout := ui.Layout{Width: ui.Fixed(size), Height: ui.Fixed(size)}
	// Two swipe boxes, one for each way, hold the finger that steers, so the pages or the list
	// around the stick do not move with it; the outer one says where the stick was laid out.
	look.Context.SwipeBox(name+" sideways", layout, ui.Horizontal, func(sideways ui.Swipe) {
		if sideways.Area.Width > 0 {
			stick.area, stick.size = sideways.Area, size
		}
		look.Context.SwipeBox(name+" up and down", layout, ui.Vertical, func(ui.Swipe) { look.stickRing(name, stick, size) })
	})
	if stick.Held {
		look.Context.RequestFrame(0)
	}
}

// stickRing draws the ring and the knob, pushed from the middle as far as the stick is pointed.
func (look Look) stickRing(name string, stick *Stick, size float32) {
	scheme := look.Roles
	knob := size * 0.4
	reach := (size - knob) / 2
	padding := ui.Padding{Left: reach + stick.Sideways*reach, Top: reach - stick.Forward*reach}
	layout := ui.Layout{Width: ui.Fixed(size), Height: ui.Fixed(size), Padding: padding}
	ring := ui.Look{Fill: scheme.Scrim, Border: scheme.OutlineVariant, BorderWidth: 2, Radius: size / 2}
	// a clickable, so the interface keeps its finger from what is behind it
	look.Context.Clickable(name, layout, ui.Looks{Normal: ring, Hovered: ring, Pressed: ring}, func() {
		knobLayout := ui.Layout{Width: ui.Fixed(knob), Height: ui.Fixed(knob), AlignX: ui.Center, AlignY: ui.Center}
		look.Context.Frame(knobLayout, ui.Look{Fill: scheme.Primary, Radius: knob / 2, Shadow: true}, nil)
	})
}

// Follow gives the stick to the finger that came down on it, and points it at that finger; a
// finger lifted lets it go back to the middle. Look.Stick calls it; a program drawing a stick of
// its own calls it with the touches and where its stick is (Place).
func (stick *Stick) Follow(touches []input.Touch) {
	if !stick.Held {
		for _, touch := range touches {
			if stick.area.Contains(float32(touch.X), float32(touch.Y)) {
				stick.Finger, stick.Held = touch.Finger, true
				break
			}
		}
	}
	touch, found := fingerIn(touches, stick.Finger)
	if !stick.Held || !found {
		stick.Held, stick.Forward, stick.Sideways = false, 0, 0
		return
	}
	reach := float64(stick.size) * 0.3 // from the middle to where the knob stops
	if reach <= 0 {
		return
	}
	middleX := float64(stick.area.X + stick.area.Width/2)
	middleY := float64(stick.area.Y + stick.area.Height/2)
	right, down := (touch.X-middleX)/reach, (touch.Y-middleY)/reach
	if length := math.Hypot(right, down); length > 1 {
		right, down = right/length, down/length
	}
	stick.Sideways, stick.Forward = float32(right), float32(-down)
}

// Place says where a stick of a program's own drawing is, and its side, for Follow.
func (stick *Stick) Place(area draw2d.Rectangle, size float32) {
	stick.area, stick.size = area, size
}

// Touching tells whether a finger is the one holding the stick: a program that turns a view with
// the other fingers leaves this one out.
func (stick *Stick) Touching(touch input.Touch) bool {
	return stick.Held && touch.Finger == stick.Finger
}

func fingerIn(touches []input.Touch, finger int64) (input.Touch, bool) {
	for _, touch := range touches {
		if touch.Finger == finger {
			return touch, true
		}
	}
	return input.Touch{}, false
}
