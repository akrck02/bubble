package material

import (
	"testing"

	"github.com/akrck02/bubble"
	"github.com/akrck02/go-renderer/input"
	ui "github.com/akrck02/go-ui"
)

func TestThePiecesDrawInEveryStyleAndFinish(t *testing.T) {
	for _, options := range []bubble.Options{
		{Style: bubble.StylePaper, Mode: bubble.ModeLight, Density: bubble.DensityComfortable, Accent: materialPurple},
		{Style: bubble.StylePaper, Mode: bubble.ModeDark, Density: bubble.DensityComfortable, Accent: materialPurple, Glass: true},
		{Style: bubble.StyleMaterial, Mode: bubble.ModeDark, Density: bubble.DensityCompact, Accent: materialPurple},
	} {
		theme, err := bubble.NewTheme(options)
		if err != nil {
			t.Fatal(err)
		}
		roles := RolesOf(theme)
		if roles.Surface != colour(theme.Colors.Background) || roles.SurfaceContainerHighest != colour(theme.Colors.ContainerHighest) {
			t.Errorf("%s %s: the roles are not the theme's colours", options.Style, options.Mode)
		}
		state := input.NewState()
		state.WindowWidth, state.WindowHeight, state.FramebufferWidth, state.FramebufferHeight = 400, 800, 400, 800
		context := ui.NewWithMeasurer(measurer{})
		on, chosen := true, 1
		context.Begin(state)
		look := New(context, theme)
		look.Card(func() {
			look.ListItem("Lighthouse", Item{Headline: "North Point"})
			look.Switch("Lit", &on)
			look.SegmentedButtons("Span", []string{"Day", "Week"}, &chosen)
			look.FilterChip("Open", true)
			look.Banner("info", "The light is off")
		})
		context.End()
		if len(context.Drawing().Commands) == 0 {
			t.Errorf("%s %s: nothing was drawn", options.Style, options.Mode)
		}
	}
}
