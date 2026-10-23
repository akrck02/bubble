package bubble

import (
	"encoding/json"
	"os"
	"testing"
)

// materialSchemes is testdata/material-schemes.json: the roles Go UI for Android's material package
// made from each seed (at commit de05770), the look Bubble's material style must keep.
type materialSchemes map[string]map[Mode]map[string]string

// channelTolerance is how far a channel may be from Go UI for Android's, which mixed its shades
// before rounding them.
const channelTolerance = 2

func TestMaterialRolesComeFromTheSeedAsGoUIForAndroidMadeThem(test *testing.T) {
	data, err := os.ReadFile("testdata/material-schemes.json")
	if err != nil {
		test.Fatal(err)
	}
	var schemes materialSchemes
	if err := json.Unmarshal(data, &schemes); err != nil {
		test.Fatal(err)
	}
	for seed, modes := range schemes {
		for mode, expected := range modes {
			theme, err := NewTheme(Options{Style: StyleMaterial, Mode: mode, Density: DensityComfortable, Accent: hex(seed)})
			if err != nil {
				test.Fatal(err)
			}
			for role, got := range materialRolesOf(theme.Colors) {
				if want := hex(expected[role]); !closeColors(got, want) {
					test.Errorf("seed %s, %s: %s is %s, Go UI for Android made %s", seed, mode, role, got.Hex(), want.Hex())
				}
			}
		}
	}
}

// materialRolesOf names Bubble's colours by the Material role each one carries.
func materialRolesOf(colors Colors) map[string]Color {
	return map[string]Color{
		"Primary": colors.Accent, "OnPrimary": colors.OnAccent,
		"PrimaryContainer": colors.AccentContainer, "OnPrimaryContainer": colors.OnAccentContainer,
		"SecondaryContainer": colors.AccentTint, "OnSecondaryContainer": colors.OnAccentTint,
		"Error": colors.Danger, "OnError": colors.OnDanger,
		"ErrorContainer": colors.DangerContainer, "OnErrorContainer": colors.OnDangerContainer,
		"Surface": colors.Background, "OnSurface": colors.Ink, "OnSurfaceVariant": colors.InkSoft,
		"SurfaceContainerLow": colors.ContainerLow, "SurfaceContainer": colors.Container,
		"SurfaceContainerHigh": colors.ContainerHigh, "SurfaceContainerHighest": colors.ContainerHighest,
		"Outline": colors.ControlLine, "OutlineVariant": colors.Line,
		"InverseSurface": colors.Inverse, "InverseOnSurface": colors.OnInverse, "InversePrimary": colors.InverseAccent,
	}
}

func closeColors(first, second Color) bool {
	near := func(a, b uint8) bool { return int(a)-int(b) <= channelTolerance && int(b)-int(a) <= channelTolerance }
	return near(first.Red, second.Red) && near(first.Green, second.Green) && near(first.Blue, second.Blue) && near(first.Alpha, second.Alpha)
}

func TestAPhonesPalettesColourTheMaterialStyle(test *testing.T) {
	palettes := PalettesFromColor(hex("#2E7D32"))
	fromPalettes, _ := NewTheme(Options{Style: StyleMaterial, Mode: ModeLight, Density: DensityComfortable, Accent: hex("#C2185B"), Palettes: &palettes})
	fromSeed, _ := NewTheme(Options{Style: StyleMaterial, Mode: ModeLight, Density: DensityComfortable, Accent: hex("#2E7D32")})
	if fromPalettes.Colors.Accent != fromSeed.Colors.Accent {
		test.Errorf("the phone's palettes give the accent %s, the seed %s", fromPalettes.Colors.Accent.Hex(), fromSeed.Colors.Accent.Hex())
	}
}

func TestInkReadsOnEveryStepOfTheLadder(test *testing.T) {
	for _, theme := range everyTheme(test) {
		colors := theme.Colors
		for _, step := range []struct {
			name  string
			color Color
		}{{"the low container", colors.ContainerLow}, {"the container", colors.Container}, {"the high container", colors.ContainerHigh}, {"the highest container", colors.ContainerHighest}} {
			expectContrast(test, theme, "ink on "+step.name, colors.Ink, step.color, textContrast)
			expectContrast(test, theme, "soft ink on "+step.name, colors.InkSoft, step.color, textContrast)
		}
		expectContrast(test, theme, "ink on the accent container", colors.OnAccentContainer, colors.AccentContainer, textContrast)
		expectContrast(test, theme, "ink on the accent tint", colors.OnAccentTint, colors.AccentTint, textContrast)
		expectContrast(test, theme, "text on danger", colors.OnDanger, colors.Danger, textContrast)
		expectContrast(test, theme, "ink on the danger container", colors.OnDangerContainer, colors.DangerContainer, textContrast)
		expectContrast(test, theme, "ink on the inverse surface", colors.OnInverse, colors.Inverse, textContrast)
		expectContrast(test, theme, "the accent on the inverse surface", colors.InverseAccent, colors.Inverse, textContrast)
	}
}
