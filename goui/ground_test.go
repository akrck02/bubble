package goui

import (
	"image/color"
	"testing"

	"github.com/akrck02/bubble"
)

func rgbaOf(value bubble.Color) color.RGBA {
	return color.RGBA{value.Red, value.Green, value.Blue, 255}
}

func TestGlassesGroundIsTheFirstLightAtItsCentreAndTheBackgroundFarFromBoth(t *testing.T) {
	theme := themeFor(t, bubble.StylePaper, bubble.ModeDark, true)
	width, height := 1280, 760
	ground := Ground(theme, width, height)
	centre := ground.RGBAAt(int(0.10*float64(width)), int(0.06*float64(height)))
	if wanted := rgbaOf(theme.Glass.Glows[0]); !nearColour(centre, wanted) {
		t.Fatalf("the first light's centre is its own colour, got %v, wanted %v", centre, wanted)
	}
	// the top right corner is beyond 70 % of both lights' radii
	corner := ground.RGBAAt(width-1, 0)
	if wanted := rgbaOf(theme.Colors.Background); corner != wanted {
		t.Fatalf("far from both lights the ground is the background, got %v, wanted %v", corner, wanted)
	}
}

func TestTheGroundIsTheSizeOfTheWindow(t *testing.T) {
	ground := Ground(themeFor(t, bubble.StylePaper, bubble.ModeLight, true), 300, 200)
	if size := ground.Bounds().Size(); size.X != 300 || size.Y != 200 {
		t.Fatalf("got %v", size)
	}
}

// nearColour allows a step of rounding in each channel: the centre pixel is half a pixel off the centre.
func nearColour(first, second color.RGBA) bool {
	near := func(one, other uint8) bool { return int(one)-int(other) <= 2 && int(other)-int(one) <= 2 }
	return near(first.R, second.R) && near(first.G, second.G) && near(first.B, second.B)
}
