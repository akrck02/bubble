package bubble

import "math"

// Material's colours, as Go UI for Android made them and Android's own dynamic colours do: one
// seed colour becomes five tonal palettes (its hue at every lightness, in OKLCH), and the roles
// of Material 3 come from fixed shades of them. A phone gives the palettes it made from its
// wallpaper instead of a seed; either way the roles come out the same.

// PaletteShades are the shades of each palette in Android's names: 0 is white, 500 the middle,
// 1000 black. Shade s has Material's tone 100 - s/10, a lightness on CIELAB's scale (L*).
var PaletteShades = [13]int{0, 10, 50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000}

// Palettes are the five tonal palettes Material's colours come from, each at PaletteShades: the
// accent's own, a quieter one for containers, a third of a turned hue, and two near-greys of the
// accent's hue for surfaces. Go for Android reads them from the phone's wallpaper.
type Palettes struct {
	Accent1, Accent2, Accent3, Neutral1, Neutral2 [13]Color
}

// opaque tells whether every shade of every palette is an opaque colour, as a phone's are.
func (palettes Palettes) opaque() bool {
	for _, palette := range [][13]Color{palettes.Accent1, palettes.Accent2, palettes.Accent3, palettes.Neutral1, palettes.Neutral2} {
		for _, shade := range palette {
			if shade.Alpha != 0xFF {
				return false
			}
		}
	}
	return true
}

// The indexes of the shades in a palette, by their Android names.
const (
	shade0 = iota
	shade10
	shade50
	shade100
	shade200
	shade300
	shade400
	shade500
	shade600
	shade700
	shade800
	shade900
	shade1000
)

// The chromas of Material's "tonal spot" palettes, in OKLCH's units, and the third palette's turn.
const (
	accentChroma    = 0.115
	quietChroma     = 0.05
	thirdChroma     = 0.077
	neutralChroma   = 0.016
	variantChroma   = 0.026
	thirdHueTurn    = math.Pi / 3
	materialScrim   = 0.32 // Material's scrim behind a dialog or a sheet: black, in both modes
	toneSearchSteps = 30
	gamutSteps      = 20
)

// PalettesFromColor makes the five palettes Android makes from a wallpaper, from one colour.
func PalettesFromColor(seed Color) Palettes {
	_, _, hue := oklch(seed)
	return Palettes{
		Accent1:  tonalPalette(hue, accentChroma),
		Accent2:  tonalPalette(hue, quietChroma),
		Accent3:  tonalPalette(hue+thirdHueTurn, thirdChroma),
		Neutral1: tonalPalette(hue, neutralChroma),
		Neutral2: tonalPalette(hue, variantChroma),
	}
}

// materialRoles are Material 3's colour roles, before they become Bubble's colours.
type materialRoles struct {
	Primary, OnPrimary, PrimaryContainer, OnPrimaryContainer                             Color
	SecondaryContainer, OnSecondaryContainer                                             Color
	Error, OnError, ErrorContainer, OnErrorContainer                                     Color
	Surface, OnSurface, OnSurfaceVariant                                                 Color
	SurfaceContainerLow, SurfaceContainer, SurfaceContainerHigh, SurfaceContainerHighest Color
	Outline, OutlineVariant                                                              Color
	InverseSurface, InverseOnSurface, InversePrimary                                     Color
}

// rolesFromPalettes gives the roles from the palettes, at the shades Android's dynamic colours
// take. The error colours stay Material's own, which no wallpaper changes.
func rolesFromPalettes(palettes Palettes, mode Mode) materialRoles {
	accent, second, neutral, variant := palettes.Accent1, palettes.Accent2, palettes.Neutral1, palettes.Neutral2
	if mode == ModeDark {
		return materialRoles{
			Primary: accent[shade200], OnPrimary: accent[shade800],
			PrimaryContainer: accent[shade700], OnPrimaryContainer: accent[shade100],
			SecondaryContainer: second[shade700], OnSecondaryContainer: second[shade100],
			Error: hex("#F2B8B5"), OnError: hex("#601410"), ErrorContainer: hex("#8C1D18"), OnErrorContainer: hex("#F9DEDC"),
			Surface:   mixFloat(neutral[shade900], neutral[shade1000], 0.4), // tone 6
			OnSurface: neutral[shade100], OnSurfaceVariant: variant[shade200],
			SurfaceContainerLow:     neutral[shade900],                                   // tone 10
			SurfaceContainer:        mixFloat(neutral[shade900], neutral[shade800], 0.2), // tone 12
			SurfaceContainerHigh:    mixFloat(neutral[shade900], neutral[shade800], 0.7), // tone 17
			SurfaceContainerHighest: mixFloat(neutral[shade800], neutral[shade700], 0.2),
			Outline:                 variant[shade400], OutlineVariant: variant[shade700],
			InverseSurface: neutral[shade100], InverseOnSurface: neutral[shade800], InversePrimary: accent[shade600],
		}
	}
	return materialRoles{
		Primary: accent[shade600], OnPrimary: accent[shade0],
		PrimaryContainer: accent[shade100], OnPrimaryContainer: accent[shade900],
		SecondaryContainer: second[shade100], OnSecondaryContainer: second[shade900],
		Error: hex("#B3261E"), OnError: hex("#FFFFFF"), ErrorContainer: hex("#F9DEDC"), OnErrorContainer: hex("#410E0B"),
		Surface:   mixFloat(neutral[shade10], neutral[shade50], 0.25), // tone 98
		OnSurface: neutral[shade900], OnSurfaceVariant: variant[shade700],
		SurfaceContainerLow:     mixFloat(neutral[shade10], neutral[shade50], 0.75), // tone 96
		SurfaceContainer:        mixFloat(neutral[shade50], neutral[shade100], 0.2), // tone 94
		SurfaceContainerHigh:    mixFloat(neutral[shade50], neutral[shade100], 0.6), // tone 92
		SurfaceContainerHighest: neutral[shade100],                                  // tone 90
		Outline:                 variant[shade500], OutlineVariant: variant[shade200],
		InverseSurface: neutral[shade800], InverseOnSurface: neutral[shade50], InversePrimary: accent[shade200],
	}
}

// materialColors are Bubble's colours in the material style: each Material role on the Bubble
// colour that does its job, and the ladder of containers as Material makes it.
func materialColors(base Colors, roles materialRoles) Colors {
	colors := base
	colors.Background, colors.Surface, colors.Sunken = roles.Surface, roles.SurfaceContainerLow, roles.SurfaceContainerHighest
	colors.Ink, colors.InkSoft, colors.InkQuiet = roles.OnSurface, roles.OnSurfaceVariant, roles.OnSurfaceVariant
	colors.Line, colors.ControlLine = roles.OutlineVariant, roles.Outline
	colors.Accent, colors.OnAccent, colors.AccentInk = roles.Primary, roles.OnPrimary, roles.Primary
	colors.AccentEdge, colors.AccentSide = roles.Primary, roles.Primary.Mix(roles.OnPrimary, 0.12) // no edge; a pressed state
	colors.AccentTint, colors.OnAccentTint = roles.SecondaryContainer, roles.OnSecondaryContainer
	colors.AccentContainer, colors.OnAccentContainer = roles.PrimaryContainer, roles.OnPrimaryContainer
	colors.Danger, colors.OnDanger = roles.Error, roles.OnError
	colors.DangerContainer, colors.OnDangerContainer = roles.ErrorContainer, roles.OnErrorContainer
	colors.ContainerLow, colors.Container = roles.SurfaceContainerLow, roles.SurfaceContainer
	colors.ContainerHigh, colors.ContainerHighest = roles.SurfaceContainerHigh, roles.SurfaceContainerHighest
	colors.Inverse, colors.OnInverse, colors.InverseAccent = roles.InverseSurface, roles.InverseOnSurface, roles.InversePrimary
	colors.FocusRing, colors.Selection = roles.Primary, roles.Primary.WithAlpha(selectionAlpha)
	colors.Scrim = black.WithAlpha(materialScrim)
	return colors
}

// tonalPalette is a hue at a chroma, at each of PaletteShades.
func tonalPalette(hue, chroma float64) [13]Color {
	var palette [13]Color
	for index, shade := range PaletteShades {
		palette[index] = colorAtTone(hue, chroma, 100-float64(shade)/10)
	}
	return palette
}

// colorAtTone finds the colour of a hue and chroma whose L* is a tone, its chroma lowered as far
// as needed to be a colour a screen shows.
func colorAtTone(hue, chroma, tone float64) Color {
	if tone >= 100 {
		return white
	}
	if tone <= 0 {
		return black
	}
	low, high := 0.0, 1.0
	var red, green, blue float64
	for range toneSearchSteps {
		lightness := (low + high) / 2
		red, green, blue = inGamut(lightness, chroma, hue)
		if cielabLightness(red, green, blue) < tone {
			low = lightness
		} else {
			high = lightness
		}
	}
	return fromLinear(red, green, blue)
}

// inGamut is the linear colour of an OKLCH one, its chroma lowered until it fits sRGB.
func inGamut(lightness, chroma, hue float64) (red, green, blue float64) {
	red, green, blue = fromOKLab(lightness, chroma*math.Cos(hue), chroma*math.Sin(hue))
	if fits(red, green, blue) {
		return red, green, blue
	}
	low, high := 0.0, chroma
	for range gamutSteps {
		middle := (low + high) / 2
		if candidateRed, candidateGreen, candidateBlue := fromOKLab(lightness, middle*math.Cos(hue), middle*math.Sin(hue)); fits(candidateRed, candidateGreen, candidateBlue) {
			low = middle
		} else {
			high = middle
		}
	}
	return fromOKLab(lightness, low*math.Cos(hue), low*math.Sin(hue))
}

func fits(red, green, blue float64) bool {
	const margin = 1e-4
	return red >= -margin && red <= 1+margin && green >= -margin && green <= 1+margin && blue >= -margin && blue <= 1+margin
}

// oklch is a colour's OKLCH lightness, chroma and hue (in radians).
func oklch(color Color) (lightness, chroma, hue float64) {
	lightness, a, b := toOKLab(decodeGamma(color.Red), decodeGamma(color.Green), decodeGamma(color.Blue))
	return lightness, math.Hypot(a, b), math.Atan2(b, a)
}

// toOKLab and fromOKLab turn linear sRGB into OKLab and back (Björn Ottosson's matrices,
// https://bottosson.github.io/posts/oklab/, public domain).
func toOKLab(red, green, blue float64) (lightness, a, b float64) {
	long := math.Cbrt(0.4122214708*red + 0.5363325363*green + 0.0514459929*blue)
	medium := math.Cbrt(0.2119034982*red + 0.6806995451*green + 0.1073969566*blue)
	short := math.Cbrt(0.0883024619*red + 0.2817188376*green + 0.6299787005*blue)
	return 0.2104542553*long + 0.7936177850*medium - 0.0040720468*short,
		1.9779984951*long - 2.4285922050*medium + 0.4505937099*short,
		0.0259040371*long + 0.7827717662*medium - 0.8086757660*short
}

func fromOKLab(lightness, a, b float64) (red, green, blue float64) {
	long := cube(lightness + 0.3963377774*a + 0.2158037573*b)
	medium := cube(lightness - 0.1055613458*a - 0.0638541728*b)
	short := cube(lightness - 0.0894841775*a - 1.2914855480*b)
	return 4.0767416621*long - 3.3077115913*medium + 0.2309699292*short,
		-1.2684380046*long + 2.6097574011*medium - 0.3413193965*short,
		-0.0041960863*long - 0.7034186147*medium + 1.7076147010*short
}

func cube(value float64) float64 { return value * value * value }

// cielabLightness is a linear colour's L*: Material's tone.
func cielabLightness(red, green, blue float64) float64 {
	luminance := 0.2126*red + 0.7152*green + 0.0722*blue
	if luminance <= 216.0/24389 {
		return luminance * 24389.0 / 27
	}
	return 116*math.Cbrt(luminance) - 16
}

func decodeGamma(channel uint8) float64 {
	value := float64(channel) / 255
	if value <= 0.04045 {
		return value / 12.92
	}
	return math.Pow((value+0.055)/1.055, 2.4)
}

func encodeGamma(channel float64) float64 {
	channel = min(max(channel, 0), 1)
	if channel <= 0.0031308 {
		return channel * 12.92
	}
	return 1.055*math.Pow(channel, 1/2.4) - 0.055
}

// fromLinear is the sRGB colour of a linear one.
func fromLinear(red, green, blue float64) Color {
	channel := func(value float64) uint8 { return uint8(math.Round(encodeGamma(value) * 255)) }
	return Color{Red: channel(red), Green: channel(green), Blue: channel(blue), Alpha: 0xFF}
}

// mixFloat is a colour a share of the way from one to another, as Material mixes two shades.
func mixFloat(from, to Color, share float64) Color {
	mixChannel := func(first, second uint8) uint8 {
		return uint8(math.Round(float64(first) + (float64(second)-float64(first))*share))
	}
	return Color{Red: mixChannel(from.Red, to.Red), Green: mixChannel(from.Green, to.Green), Blue: mixChannel(from.Blue, to.Blue), Alpha: 0xFF}
}
