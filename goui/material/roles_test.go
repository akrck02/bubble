package material

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/akrck02/bubble"
)

func TestTheBarScrimIsTheOneGoUIForAndroidMade(t *testing.T) {
	data, err := os.ReadFile("../../testdata/material-schemes.json") // Bubble's fixture of Go UI for Android's colours
	if err != nil {
		t.Fatal(err)
	}
	var schemes map[string]map[bubble.Mode]map[string]string
	if err := json.Unmarshal(data, &schemes); err != nil {
		t.Fatal(err)
	}
	for seed, modes := range schemes {
		accent, err := bubble.ParseColor(seed)
		if err != nil {
			t.Fatal(err)
		}
		for mode, roles := range modes {
			theme, err := bubble.NewTheme(bubble.Options{Style: bubble.StyleMaterial, Mode: mode, Density: bubble.DensityComfortable, Accent: accent})
			if err != nil {
				t.Fatal(err)
			}
			if got, want := RolesOf(theme).BarScrim, colour(mustParse(t, roles["Scrim"])); distance(got, want) > 0.02 {
				t.Errorf("seed %s, %s: the bar's scrim is %v, Go UI for Android made %v", seed, mode, got, want)
			}
		}
	}
}

func mustParse(t *testing.T, text string) bubble.Color {
	t.Helper()
	color, err := bubble.ParseColor(text)
	if err != nil {
		t.Fatal(err)
	}
	return color
}
