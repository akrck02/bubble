package material

import (
	"testing"

	"github.com/akrck02/go-renderer/draw2d"
)

func TestLoadingIndicatorsPaintAndKeepMoving(t *testing.T) {
	for name, draw := range map[string]func(look Look){
		"loading":   func(look Look) { look.LoadingIndicator(48) },
		"contained": func(look Look) { look.ContainedLoadingIndicator(48) },
		"circular":  func(look Look) { look.CircularProgress(Indeterminate, 48, true) },
		"wavy":      func(look Look) { look.WavyProgress(0.4) },
		"skeleton":  func(look Look) { look.Skeleton(fixed(200), fixed(16), 8) },
	} {
		screen := newScreen(t, draw)
		screen.frame()
		painted := 0
		for _, command := range screen.context.Drawing().Commands {
			if command.Kind == draw2d.RectangleKind {
				painted++
			}
		}
		// a skeleton is one block, and its light only while it crosses; shapes are many rectangles
		if wanted := map[bool]int{true: 1, false: 5}[name == "skeleton"]; painted < wanted {
			t.Fatalf("%s painted %d rectangles; want a shape made of many", name, painted)
		}
		if _, wanted := screen.state.Feed().TakeFrameRequest(); !wanted {
			t.Fatalf("%s does not ask for its next frame", name)
		}
	}
}

func TestAKnownProgressStandsStill(t *testing.T) {
	screen := newScreen(t, func(look Look) { look.CircularProgress(0.5, 48, false) })
	screen.frame()
	if _, wanted := screen.state.Feed().TakeFrameRequest(); wanted {
		t.Fatal("a still ring asks for frames: a screen showing it never sleeps")
	}
}
