package goui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/akrck02/bubble"
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
	"github.com/akrck02/go-ui/uitest"
)

func TestATonalButtonIsTheAccentsTintInItsInk(t *testing.T) {
	for _, mode := range []bubble.Mode{bubble.ModeLight, bubble.ModeDark} {
		theme := themeFor(t, bubble.StylePaper, mode, false)
		look := TonalLook(theme)
		if look.Fill != colour(theme.Colors.AccentTint) || look.Ink != colour(theme.Colors.AccentInk) {
			t.Errorf("%s: fill %v and ink %v, want the accent's tint and ink", mode, look.Fill, look.Ink)
		}
		if contrast := bubble.Contrast(theme.Colors.AccentInk, theme.Colors.AccentTint); contrast < 4.5 {
			t.Errorf("%s: its label reads at %.2f:1, needs 4.5", mode, contrast)
		}
	}
}

func TestATonalButtonCanBeHitWithAFinger(t *testing.T) {
	theme := themeFor(t, bubble.StylePaper, bubble.ModeLight, false)
	clicks := 0
	window, err := uitest.New(func(context *ui.Context) {
		if TonalButton(context, theme, "Sell one") {
			clicks++
		}
		TonalIconButton(context, theme, "plus", "##Add a tag")
	}, uitest.Size(400, 300), uitest.Theme(Theme(theme)))
	if err != nil {
		t.Fatal(err)
	}
	window.Frame()
	for _, label := range []string{"Sell one", "Add a tag"} {
		button, err := window.Element(label)
		if err != nil {
			t.Fatal(err)
		}
		if button.Rect.Width < 44 || button.Rect.Height < 44 {
			t.Errorf("%s is %v × %v, needs 44 × 44", label, button.Rect.Width, button.Rect.Height)
		}
	}
	if err := window.Click("Sell one"); err != nil || clicks != 1 {
		t.Errorf("a click sold %d (%v), want 1", clicks, err)
	}
}

func TestTonalButtonsRepeatedInAListAreToldApartAndShowTheirWordsOnly(t *testing.T) {
	theme := themeFor(t, bubble.StylePaper, bubble.ModeLight, false)
	sold := []int{0, 0}
	window, err := uitest.New(func(context *ui.Context) {
		for row := range sold {
			if TonalButton(context, theme, fmt.Sprintf("Sell one##%d", row)) {
				sold[row]++
			}
		}
	}, uitest.Size(400, 300), uitest.Theme(Theme(theme)))
	if err != nil {
		t.Fatal(err)
	}
	window.Frame()
	for _, text := range window.Texts() {
		if strings.Contains(text, "##") {
			t.Errorf("a button shows %q", text)
		}
	}
	if _, err := window.Element("Sell one"); uitest.CodeOf(err) != uitest.ElementAmbiguous {
		t.Fatalf("want two buttons read as «Sell one», got %v", err)
	}
	var buttons []draw2d.Rectangle
	for _, node := range window.Context().ScreenReaderNodes() {
		if node.Pressable {
			buttons = append(buttons, node.Area)
		}
	}
	if len(buttons) != 2 {
		t.Fatalf("%d buttons, want 2", len(buttons))
	}
	window.ClickAt(buttons[1].X+buttons[1].Width/2, buttons[1].Y+buttons[1].Height/2)
	if sold[0] != 0 || sold[1] != 1 {
		t.Errorf("a click on the second row sold %v, want only the second", sold)
	}
}
