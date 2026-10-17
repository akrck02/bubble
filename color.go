package bubble

import (
	"fmt"
	"math"
	"strconv"
)

// Color is a colour in sRGB with its opacity, each channel from 0 to 255.
type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}

// ParseColor reads a colour written as #RRGGBB or #RRGGBBAA.
func ParseColor(text string) (Color, error) {
	if len(text) != 7 && len(text) != 9 || text[0] != '#' {
		return Color{}, InvalidColour
	}
	value, err := strconv.ParseUint(text[1:], 16, 32)
	if err != nil {
		return Color{}, InvalidColour
	}
	if len(text) == 7 {
		value = value<<8 | 0xFF
	}
	return Color{Red: uint8(value >> 24), Green: uint8(value >> 16), Blue: uint8(value >> 8), Alpha: uint8(value)}, nil
}

// hex reads a colour the package itself writes; a typo in one of its own values is a bug that the
// tests catch on the first run, so it stops there instead of returning an error every caller ignores.
func hex(text string) Color {
	color, err := ParseColor(text)
	if err != nil {
		panic(fmt.Sprintf("bubble: %s is not a colour: %v", text, err))
	}
	return color
}

// Hex writes the colour as #RRGGBB, or #RRGGBBAA when it is not opaque.
func (color Color) Hex() string {
	if color.Alpha == 0xFF {
		return fmt.Sprintf("#%02X%02X%02X", color.Red, color.Green, color.Blue)
	}
	return fmt.Sprintf("#%02X%02X%02X%02X", color.Red, color.Green, color.Blue, color.Alpha)
}

// Shade multiplies each channel by amount: 0.85 and 0.70 give the left and right sides of a voxel
// block, and the edge under a primary button.
func (color Color) Shade(amount float64) Color {
	scale := func(channel uint8) uint8 {
		return uint8(math.Round(math.Min(255, float64(channel)*amount)))
	}
	return Color{Red: scale(color.Red), Green: scale(color.Green), Blue: scale(color.Blue), Alpha: color.Alpha}
}

// Mix moves the colour toward another by amount, from 0 (itself) to 1 (the other one).
func (color Color) Mix(toward Color, amount float64) Color {
	blend := func(from uint8, to uint8) uint8 {
		return uint8(math.Round(float64(from) + (float64(to)-float64(from))*amount))
	}
	return Color{
		Red:   blend(color.Red, toward.Red),
		Green: blend(color.Green, toward.Green),
		Blue:  blend(color.Blue, toward.Blue),
		Alpha: blend(color.Alpha, toward.Alpha),
	}
}

// WithAlpha gives the same colour with another opacity, from 0 to 1.
func (color Color) WithAlpha(opacity float64) Color {
	color.Alpha = uint8(math.Round(opacity * 255))
	return color
}

// Luminance is the colour's relative luminance, as WCAG defines it: 0 for black, 1 for white.
func (color Color) Luminance() float64 {
	linear := func(channel uint8) float64 {
		value := float64(channel) / 255
		if value <= 0.04045 {
			return value / 12.92
		}
		return math.Pow((value+0.055)/1.055, 2.4)
	}
	return 0.2126*linear(color.Red) + 0.7152*linear(color.Green) + 0.0722*linear(color.Blue)
}

// Over paints the colour, with its opacity, over an opaque background and gives what the eye sees.
func (color Color) Over(background Color) Color {
	opaque := color
	opaque.Alpha = 0xFF
	return background.Mix(opaque, float64(color.Alpha)/255)
}

// Contrast is the WCAG contrast ratio between two colours, from 1 to 21. Text needs 4.5 (3 when
// it is large), the edges of controls and the focus ring need 3.
func Contrast(first Color, second Color) float64 {
	lighter, darker := first.Luminance(), second.Luminance()
	if darker > lighter {
		lighter, darker = darker, lighter
	}
	return (lighter + 0.05) / (darker + 0.05)
}
