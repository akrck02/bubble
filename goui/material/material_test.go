package material

import (
	"testing"

	"github.com/akrck02/bubble"
	"unicode/utf8"

	"github.com/akrck02/go-renderer/draw2d"
	"github.com/akrck02/go-renderer/input"
	ui "github.com/akrck02/go-ui"
)

func fixed(size float32) ui.Sizing { return ui.Fixed(size) }

// measurer stands in for the fonts: no display needed.
type measurer struct{}

func (measurer) MeasureText(text string, size float32, _ draw2d.Weight, _ float32) draw2d.TextMetrics {
	return draw2d.TextMetrics{Width: float32(utf8.RuneCountInString(text)) * size / 2, Ascent: size * 0.8, Descent: size * 0.2}
}

func (measurer) HasIcon(string) bool { return true }

// screen drives a phone-sized interface frame by frame.
type screen struct {
	t       *testing.T
	context *ui.Context
	state   *input.State
	build   func(look Look)
}

func newScreen(t *testing.T, build func(look Look)) *screen {
	state := input.NewState()
	state.WindowWidth, state.WindowHeight = 400, 800
	state.FramebufferWidth, state.FramebufferHeight = 400, 800
	return &screen{t: t, context: ui.NewWithMeasurer(measurer{}), state: state, build: build}
}

func (screen *screen) frame() {
	screen.context.Begin(screen.state)
	screen.build(New(screen.context, testTheme()))
	screen.context.End()
	screen.state.Feed().EndFrame()
}

// tap presses and lifts the pointer in the middle of a piece laid out in the last frame.
func (screen *screen) tap(path ...string) {
	screen.t.Helper()
	area, found := screen.context.Area(path...)
	if !found {
		screen.t.Fatalf("nothing called %v was laid out", path)
	}
	screen.state.Feed().MoveCursor(float64(area.X+area.Width/2), float64(area.Y+area.Height/2))
	screen.state.Feed().SetButton(input.MouseLeft, true)
	screen.frame()
	screen.state.Feed().SetButton(input.MouseLeft, false)
	screen.frame()
}

func distance(first, second draw2d.Color) float32 {
	return abs(first[0]-second[0]) + abs(first[1]-second[1]) + abs(first[2]-second[2])
}

func abs(value float32) float32 { return max(value, -value) }

// testTheme is the theme the pieces' tests draw in: Bubble's material style, light, in Material's
// own purple, the seed of its baseline scheme.
func testTheme() bubble.Theme {
	theme, err := bubble.NewTheme(bubble.Options{Style: bubble.StyleMaterial, Mode: bubble.ModeLight, Density: bubble.DensityComfortable, Accent: materialPurple})
	if err != nil {
		panic(err) // fixed options: an error is a bug in Bubble
	}
	return theme
}

// materialPurple is the seed of Material's baseline scheme.
var materialPurple, _ = bubble.ParseColor("#6750A4")
