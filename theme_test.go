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
					for _, glass := range []bool{false, true} {
						theme, err := NewTheme(Options{Style: style, Mode: mode, Density: density, Accent: hex(accent), Glass: glass})
						if err != nil {
							test.Fatalf("NewTheme(%s %s %s %s): %v", style, mode, density, accent, err)
						}
						themes = append(themes, theme)
					}
				}
			}
		}
	}
	return themes
}

func name(theme Theme) string {
	glass := ""
	if theme.Glass.Enabled {
		glass = " glass"
	}
	return fmt.Sprintf("%s%s %s %s accent %s", theme.Options.Style, glass, theme.Options.Mode, theme.Options.Density, theme.Options.Accent.Hex())
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
		expectContrast(test, theme, "soft ink on the accent tint", colors.InkSoft, colors.AccentTint, textContrast)
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
		expectContrast(test, theme, "danger on a surface", colors.Danger, colors.Surface, textContrast)
		expectContrast(test, theme, "danger on the sunken ground", colors.Danger, colors.Sunken, textContrast)
	}
}

func TestTextOnGlassReadsOverTheGroundEachLightAndAnyWallpaper(test *testing.T) {
	for _, theme := range everyTheme(test) {
		glass, colors := theme.Glass, theme.Colors
		if !glass.Enabled {
			continue
		}
		veiled := []Color{glass.Veil.Over(black), glass.Veil.Over(white.Shade(glass.Dimming))}
		for _, behind := range append([]Color{colors.Background, glass.Glows[0], glass.Glows[1]}, veiled...) {
			for _, fill := range []struct {
				name  string
				color Color
			}{{"a pane", glass.Pane}, {"a raised pane", glass.Raised}} {
				seen := fill.color.Over(behind)
				where := fill.name + " over " + behind.Hex()
				expectContrast(test, theme, "ink on "+where, colors.Ink, seen, textContrast)
				expectContrast(test, theme, "soft ink on "+where, colors.InkSoft, seen, textContrast)
				expectContrast(test, theme, "quiet ink on "+where, colors.InkQuiet, seen, textContrast)
				expectContrast(test, theme, "the edge of a field on "+where, colors.ControlLine, seen, controlContrast)
				expectContrast(test, theme, "the focus ring on "+where, colors.FocusRing, seen, controlContrast)
			}
		}
	}
}

func TestTextOnAGlassDialogReadsOverAnyPageBehindIt(test *testing.T) {
	for _, theme := range everyTheme(test) {
		glass, colors := theme.Glass, theme.Colors
		if !glass.Enabled {
			continue
		}
		for _, behind := range []Color{black, white} {
			page := colors.Scrim.Over(behind)
			for _, fill := range []struct {
				name  string
				color Color
			}{{"a dialog", glass.Dialog}, {"a raised field in a dialog", glass.Raised}} {
				seen := fill.color.Over(glass.Dialog.Over(page))
				if fill.name == "a dialog" {
					seen = glass.Dialog.Over(page)
				}
				where := fill.name + " over " + behind.Hex()
				expectContrast(test, theme, "ink on "+where, colors.Ink, seen, textContrast)
				expectContrast(test, theme, "soft ink on "+where, colors.InkSoft, seen, textContrast)
				expectContrast(test, theme, "quiet ink on "+where, colors.InkQuiet, seen, textContrast)
				expectContrast(test, theme, "the focus ring on "+where, colors.FocusRing, seen, controlContrast)
			}
		}
	}
}

func TestOnGlassWhatIsChosenIsBrighterNotColoured(test *testing.T) {
	for _, theme := range everyTheme(test) {
		if theme.Glass.Enabled && theme.Colors.AccentTint != theme.Glass.Raised.Over(theme.Colors.Background) {
			test.Errorf("%s: the chosen fill is %s, not the raised pane", name(theme), theme.Colors.AccentTint.Hex())
		}
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

func TestABrightAccentKeepsTheLogosColourAndCarriesDarkText(test *testing.T) {
	gold := hex("#F6BD60")
	theme, err := NewTheme(Options{Style: StylePaper, Mode: ModeLight, Density: DensityComfortable, Accent: gold})
	if err != nil {
		test.Fatal(err)
	}
	if theme.Colors.Accent != gold {
		test.Errorf("accent %s, want the logo's own %s", theme.Colors.Accent.Hex(), gold.Hex())
	}
	if theme.Colors.OnAccent == white {
		test.Error("white text on a gold accent; it should be dark")
	}
}

func TestPaperButtonsHaveMediumCornersAndMaterialButtonsArePills(test *testing.T) {
	for _, theme := range everyTheme(test) {
		isPill := theme.Shape.ButtonRadius == pill
		if theme.Options.Style == StylePaper && isPill {
			test.Errorf("%s: paper buttons are pills", name(theme))
		}
		if theme.Options.Style == StyleMaterial && !isPill {
			test.Errorf("%s: material buttons are not pills", name(theme))
		}
	}
}

func TestTheFocusRingTouchesTheControl(test *testing.T) {
	for _, theme := range everyTheme(test) {
		if theme.Shape.FocusOffset != 0 {
			test.Errorf("%s: the focus ring floats %v pixels away from the control", name(theme), theme.Shape.FocusOffset)
		}
	}
}
