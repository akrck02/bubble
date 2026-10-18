package goui

import (
	"testing"

	"github.com/akrck02/bubble"
	"github.com/akrck02/go-renderer/draw2d"
)

func themeFor(t *testing.T, style bubble.Style, mode bubble.Mode, glass bool) bubble.Theme {
	t.Helper()
	accent, err := bubble.ParseColor("#81D4FA")
	if err != nil {
		t.Fatal(err)
	}
	theme, err := bubble.NewTheme(bubble.Options{Style: style, Mode: mode, Density: bubble.DensityComfortable, Accent: accent, Glass: glass})
	if err != nil {
		t.Fatal(err)
	}
	return theme
}

func TestEveryBubbleColourLandsOnTheGoUITokenOfTheSameRole(t *testing.T) {
	theme := themeFor(t, bubble.StylePaper, bubble.ModeLight, false)
	result, colours := Theme(theme), theme.Colors
	for _, pair := range []struct {
		name   string
		got    draw2d.Color
		wanted bubble.Color
	}{
		{"background", result.Background, colours.Background}, {"surface", result.Surface, colours.Surface},
		{"sunken", result.Sunken, colours.Sunken}, {"ink", result.Ink, colours.Ink}, {"soft ink", result.InkSoft, colours.InkSoft},
		{"quiet ink", result.InkQuiet, colours.InkQuiet}, {"line", result.Line, colours.Line},
		{"control line", result.ControlLine, colours.ControlLine}, {"accent", result.Primary, colours.Accent},
		{"text on the accent", result.OnPrimary, colours.OnAccent}, {"focus ring", result.FocusRing, colours.FocusRing},
		{"danger", result.Danger, colours.Danger}, {"main button's edge", result.PrimaryEdge, colours.AccentEdge},
	} {
		if pair.got != colour(pair.wanted) {
			t.Errorf("%s: got %v, wanted Bubble's %s", pair.name, pair.got, pair.wanted.Hex())
		}
	}
}

func TestWhatIsChosenIsBrighterAndHeavierNeverColoured(t *testing.T) {
	result := Theme(themeFor(t, bubble.StylePaper, bubble.ModeLight, false))
	if result.Chosen != result.Sunken || result.ChosenInk != result.Ink || result.ChosenWeight != draw2d.Bold {
		t.Fatalf("chosen is the sunken fill with bold ink, got %v %v %v", result.Chosen, result.ChosenInk, result.ChosenWeight)
	}
}

func TestSizesComeFromBubblesDensityAndStyle(t *testing.T) {
	paper := themeFor(t, bubble.StylePaper, bubble.ModeLight, false)
	material := themeFor(t, bubble.StyleMaterial, bubble.ModeLight, false)
	paperResult, materialResult := Theme(paper), Theme(material)
	if paperResult.ControlHeight != float32(paper.Space.ControlHeight) || paperResult.TextSize != float32(paper.Type.BodySize) {
		t.Fatal("controls and text take Bubble's sizes")
	}
	if paperResult.ButtonRadius != float32(paper.Shape.ButtonRadius) || materialResult.ButtonRadius != float32(material.Shape.ButtonRadius) {
		t.Fatal("buttons take Bubble's corners: medium on paper, a pill on material")
	}
	if paperResult.EdgeWidth != float32(paper.Shape.EdgeWidth) || paperResult.ShadowBlur != 0 {
		t.Fatal("paper's main button stands on its edge, and paper has no shadows")
	}
	if materialResult.ShadowBlur == 0 {
		t.Fatal("material keeps its shadows")
	}
	if paperResult.LineHeight != lineHeight || paperResult.HeadingLineHeight != headingLineHeight {
		t.Fatal("lines take Bubble's line heights")
	}
}

func TestOnGlassPanelsAndWellsAreTheTranslucentPanes(t *testing.T) {
	theme := themeFor(t, bubble.StylePaper, bubble.ModeDark, true)
	result := Theme(theme)
	if result.Surface != colour(theme.Glass.Pane) || result.Sunken != colour(theme.Glass.Raised) || result.Chosen != result.Sunken {
		t.Fatal("on glass, surfaces are the pane and wells the raised pane")
	}
	if result.Surface[3] >= 1 || result.Sunken[3] >= 1 {
		t.Fatal("nothing on glass is opaque")
	}
	if result.Backdrop != float32(theme.Glass.Blur) {
		t.Fatal("glass blurs what is behind its panes by Bubble's blur")
	}
	if Theme(themeFor(t, bubble.StylePaper, bubble.ModeDark, false)).Backdrop != 0 {
		t.Fatal("paper without glass blurs nothing")
	}
}

func TestABubbleThemeLightsFieldsAndShrinksWhenPressed(t *testing.T) {
	theme := themeFor(t, bubble.StylePaper, bubble.ModeLight, false)
	result := Theme(theme)
	if result.FocusWidth != float32(theme.Shape.FocusWidth) || !result.FieldFocusEdge {
		t.Fatalf("the focus is Bubble's: a %v pixel ring on the edge, fields lighting their own (%v)", result.FocusWidth, result.FieldFocusEdge)
	}
	if result.PressScale != 0.97 {
		t.Fatalf("a held control is drawn at 0.97, as Bubble's CSS presses it: %v", result.PressScale)
	}
}
