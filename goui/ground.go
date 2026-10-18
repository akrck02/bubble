package goui

import (
	"image"
	"image/color"
	"math"

	"github.com/akrck02/bubble"
)

// glow is one of the soft lights of Bubble's glass ground, as css/bubble.css draws it: an ellipse
// of a colour, fully there at its centre and gone at 70 % of its radius.
type glow struct {
	colour           bubble.Color
	centreX, centreY float64 // shares of the window: 0.1 is 10 %
	radiusX, radiusY float64 // pixels
}

// glowFadesBy is where CSS's "transparent 70%" stop puts the end of each light.
const glowFadesBy = 0.7

// remInPixels turns the CSS sizes of the lights (52rem 44rem) into pixels.
const remInPixels = 16

// Ground paints glass's ground the way the browser does, layer by layer: the background, the
// second light over it, the first light on top. Panels are drawn over this picture; it is smooth,
// so blurring what is behind them would change nothing.
func Ground(theme bubble.Theme, width, height int) *image.RGBA {
	lights := []glow{
		{theme.Glass.Glows[1], 0.92, 0.96, 48 * remInPixels, 44 * remInPixels},
		{theme.Glass.Glows[0], 0.10, 0.06, 52 * remInPixels, 44 * remInPixels},
	}
	picture := image.NewRGBA(image.Rect(0, 0, width, height))
	background := theme.Colors.Background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			red, green, blue := float64(background.Red), float64(background.Green), float64(background.Blue)
			for _, light := range lights {
				across := (float64(x) + 0.5 - light.centreX*float64(width)) / light.radiusX
				down := (float64(y) + 0.5 - light.centreY*float64(height)) / light.radiusY
				strength := 1 - math.Sqrt(across*across+down*down)/glowFadesBy
				if strength <= 0 {
					continue
				}
				red += (float64(light.colour.Red) - red) * strength
				green += (float64(light.colour.Green) - green) * strength
				blue += (float64(light.colour.Blue) - blue) * strength
			}
			picture.SetRGBA(x, y, color.RGBA{uint8(math.Round(red)), uint8(math.Round(green)), uint8(math.Round(blue)), 255})
		}
	}
	return picture
}
