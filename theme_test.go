package bubble

import (
	"errors"
	"fmt"
	"testing"
)

// Accents that are hard to read on: a pale paper, a bright yellow, a light teal, a saturated red,
// and the default blue.
var accentsToCheck = []string{"#2C5C9E", "#F7EBD0", "#FFCA28", "#2BB3A3", "#EF5350", "#3DDC84", "#000000", "#FFFFFF"}

func everyTheme(test *testing.T) []Theme {
	var themes []Theme
	for _, style := range []Style{StylePaper, StyleMaterial} {
		for _, mode := range []Mode{ModeLight, ModeDark} {
			for _, density := range []Density{DensityComfortable, DensityCompact} {
				for _, accent := range accentsToCheck {
					theme, err := NewTheme(Options{Style: style, Mode: mode, Density: density, Accent: hex(accent), Glass: true})
					if err != nil {
						test.Fatalf("NewTheme(%s %s %s %s): %v", style, mode, density, accent, err)
					}
					themes = append(themes, theme)
				}
			}
		}
	}
	return themes
}

func name(theme Theme) string {
	return fmt.Sprintf("%s %s %s accent %s", theme.Options.Style, theme.Options.Mode, theme.Options.Density, theme.Options.Accent.Hex())
}

func expectContrast(test *testing.T, theme Theme, what string, first Color, second Color, minimum float64) {
	test.Helper()
	if ratio := Contrast(first, second); ratio < minimum {
		test.Errorf("%s: %s has a contrast of %.2f, it needs %.1f", name(theme), what, ratio, minimum)
	}
}

func TestTextReadsOnEveryBackgroundOfEveryTheme(test *testing.T) {
	for _, theme := range everyTheme(test) {
		colors := theme.Colors
		for _, ground := range []struct {
			name  string
			color Color
		}{{"background", colors.Background}, {"surface", colors.Surface}, {"sunken", colors.Sunken}} {
			expectContrast(test, theme, "ink on "+ground.name, colors.Ink, ground.color, textContrast)
			expectContrast(test, theme, "soft ink on "+ground.name, colors.InkSoft, ground.color, textContrast)
			expectContrast(test, theme, "quiet ink on "+ground.name, colors.InkQuiet, ground.color, textContrast)
		}
		expectContrast(test, theme, "ink on the accent tint", colors.Ink, colors.AccentTint, textContrast)
	}
}

func TestTheAccentReadsWhateverColourTheLogoIs(test *testing.T) {
	for _, theme := range everyTheme(test) {
		colors := theme.Colors
		expectContrast(test, theme, "text on the accent", colors.OnAccent, colors.Accent, textContrast)
		expectContrast(test, theme, "accent ink on the background", colors.AccentInk, colors.Background, textContrast)
		expectContrast(test, theme, "accent ink on the accent tint", colors.AccentInk, colors.AccentTint, textContrast)
	}
}

func TestControlsAndTheFocusRingStandOutFromTheBackground(test *testing.T) {
	for _, theme := range everyTheme(test) {
		colors := theme.Colors
		expectContrast(test, theme, "control line", colors.ControlLine, colors.Background, controlContrast)
		expectContrast(test, theme, "control line on a surface", colors.ControlLine, colors.Surface, controlContrast)
		expectContrast(test, theme, "focus ring", colors.FocusRing, colors.Background, controlContrast)
	}
}

func TestStatusColoursReadAsText(test *testing.T) {
	for _, theme := range everyTheme(test) {
		colors := theme.Colors
		expectContrast(test, theme, "positive", colors.Positive, colors.Background, textContrast)
		expectContrast(test, theme, "warning", colors.Warning, colors.Background, textContrast)
		expectContrast(test, theme, "danger", colors.Danger, colors.Background, textContrast)
	}
}

func TestTextOnGlassReadsOverTheBackground(test *testing.T) {
	for _, theme := range everyTheme(test) {
		if !theme.Glass.Enabled {
			continue
		}
		seen := theme.Glass.Surface.Over(theme.Colors.Background)
		expectContrast(test, theme, "ink on glass", theme.Colors.Ink, seen, textContrast)
		expectContrast(test, theme, "soft ink on glass", theme.Colors.InkSoft, seen, textContrast)
	}
}

func TestGlassIsOnlyForPaper(test *testing.T) {
	theme, err := NewTheme(Options{Style: StyleMaterial, Mode: ModeLight, Density: DensityComfortable, Accent: DefaultAccent, Glass: true})
	if err != nil {
		test.Fatal(err)
	}
	if theme.Glass.Enabled {
		test.Error("material with glass: glass is enabled, it should only exist in paper")
	}
}

func TestTheEdgeIsTheLogosDarkerSide(test *testing.T) {
	theme, err := NewTheme(DefaultOptions())
	if err != nil {
		test.Fatal(err)
	}
	if theme.Colors.AccentEdge != theme.Colors.Accent.Shade(0.70) {
		test.Errorf("edge %s, want the accent %s ×0.70", theme.Colors.AccentEdge.Hex(), theme.Colors.Accent.Hex())
	}
	if theme.Shape.EdgeWidth == 0 {
		test.Error("paper has no block edge")
	}
}

func TestComfortableControlsCanBeHitWithAFinger(test *testing.T) {
	theme, err := NewTheme(DefaultOptions())
	if err != nil {
		test.Fatal(err)
	}
	if theme.Space.TouchTarget < 44 || theme.Space.ControlHeight < 44 {
		test.Errorf("touch target %v and control height %v, both need at least 44", theme.Space.TouchTarget, theme.Space.ControlHeight)
	}
}

func TestEachWrongOptionHasItsOwnCode(test *testing.T) {
	cases := []struct {
		options Options
		want    *Code
	}{
		{Options{Style: "glass", Mode: ModeLight, Density: DensityComfortable}, UnknownStyle},
		{Options{Style: StylePaper, Mode: "dim", Density: DensityComfortable}, UnknownMode},
		{Options{Style: StylePaper, Mode: ModeLight, Density: "cosy"}, UnknownDensity},
	}
	for _, testCase := range cases {
		if _, err := NewTheme(testCase.options); !errors.Is(err, testCase.want) {
			test.Errorf("NewTheme(%+v) = %v, want %s", testCase.options, err, testCase.want.Name())
		}
	}
}

func TestAColourIsReadOnlyWhenItIsWrittenAsHex(test *testing.T) {
	for _, text := range []string{"#2C5C9E", "#2C5C9E80"} {
		if _, err := ParseColor(text); err != nil {
			test.Errorf("ParseColor(%q): %v", text, err)
		}
	}
	for _, text := range []string{"", "2C5C9E", "#2C5", "#GGGGGG", "#2C5C9E8"} {
		if _, err := ParseColor(text); !errors.Is(err, InvalidColour) {
			test.Errorf("ParseColor(%q) = %v, want %s", text, err, InvalidColour.Name())
		}
	}
}

func TestBlackAndWhiteHaveTheHighestContrast(test *testing.T) {
	if ratio := Contrast(black, white); ratio < 20.9 || ratio > 21.1 {
		test.Errorf("black on white: %.2f, want 21", ratio)
	}
}
