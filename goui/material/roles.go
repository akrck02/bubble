package material

import (
	"github.com/akrck02/bubble"
	"github.com/akrck02/go-renderer/draw2d"
)

// Roles are a Bubble theme's colours by Material 3's names, the names the pieces were written
// with. In the material style they are Material's own; in paper and glass, the colours doing the
// same job, so every piece draws in every style.
type Roles struct {
	Primary, OnPrimary, PrimaryContainer, OnPrimaryContainer draw2d.Color
	SecondaryContainer, OnSecondaryContainer                 draw2d.Color
	Error, OnError, ErrorContainer, OnErrorContainer         draw2d.Color
	// Surface is the screen's background; the containers are the surfaces raised over it, from the
	// lowest to the highest (a card, a navigation bar, a switch's track).
	Surface, OnSurface, OnSurfaceVariant                                                 draw2d.Color
	SurfaceContainerLow, SurfaceContainer, SurfaceContainerHigh, SurfaceContainerHighest draw2d.Color
	Outline, OutlineVariant                                                              draw2d.Color
	// The inverse roles paint what stands out against the screen: a snackbar.
	InverseSurface, InverseOnSurface, InversePrimary draw2d.Color
	// BarScrim is the surface's colour, see-through: under a bar painted over a scene. It is not
	// the theme's Scrim, the darkening behind a dialog.
	BarScrim draw2d.Color
}

// scrimOpacity is how much of the surface a bar over a scene keeps.
const scrimOpacity = 0.8

// RolesOf names a Bubble theme's colours as Material does.
func RolesOf(theme bubble.Theme) Roles {
	colors := theme.Colors
	return Roles{
		Primary: colour(colors.Accent), OnPrimary: colour(colors.OnAccent),
		PrimaryContainer: colour(colors.AccentContainer), OnPrimaryContainer: colour(colors.OnAccentContainer),
		SecondaryContainer: colour(colors.AccentTint), OnSecondaryContainer: colour(colors.OnAccentTint),
		Error: colour(colors.Danger), OnError: colour(colors.OnDanger),
		ErrorContainer: colour(colors.DangerContainer), OnErrorContainer: colour(colors.OnDangerContainer),
		Surface: colour(colors.Background), OnSurface: colour(colors.Ink), OnSurfaceVariant: colour(colors.InkSoft),
		SurfaceContainerLow: colour(colors.ContainerLow), SurfaceContainer: colour(colors.Container),
		SurfaceContainerHigh: colour(colors.ContainerHigh), SurfaceContainerHighest: colour(colors.ContainerHighest),
		Outline: colour(colors.ControlLine), OutlineVariant: colour(colors.Line),
		InverseSurface: colour(colors.Inverse), InverseOnSurface: colour(colors.OnInverse), InversePrimary: colour(colors.InverseAccent),
		BarScrim: colour(colors.Background.WithAlpha(scrimOpacity)),
	}
}

func colour(value bubble.Color) draw2d.Color { return draw2d.MustHex(value.Hex()) }

// mix is a colour a share of the way from one to another.
func mix(from, to draw2d.Color, share float32) draw2d.Color {
	var mixed draw2d.Color
	for channel := range mixed {
		mixed[channel] = from[channel] + (to[channel]-from[channel])*share
	}
	return mixed
}
