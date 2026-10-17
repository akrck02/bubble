package bubble

// How the project's accent becomes the accent tokens. A logo's colour is chosen to look good on the
// dark tile, not to carry white text, so the theme moves it, step by step, until every pairing
// reaches its contrast: text 4.5, the focus ring and the edges of controls 3.

const (
	textContrast    = 4.5
	controlContrast = 3.0
	stepToward      = 0.06 // how far each step moves a colour toward black or white
	maximumSteps    = 40   // more than enough to reach black or white from any colour
	sideShade       = 0.85 // the logo's lighter side
	edgeShade       = 0.70 // the logo's darker side
	paperTint       = 0.82 // how much of the background a paper tint keeps
	materialTint    = 0.75 // Material's containers carry more colour
	darkTint        = 0.65
	selectionAlpha  = 0.3
)

var (
	black = hex("#000000")
	white = hex("#FFFFFF")
	// onDarkAccent is the text on a light accent fill in dark mode, a near-black with a hint of slate.
	onDarkAccent = hex("#0E1116")
)

func withAccent(colors Colors, options Options) Colors {
	if options.Mode == ModeLight {
		return lightAccent(colors, options)
	}
	return darkAccent(colors, options)
}

// In light mode a dark accent carries white text. A bright one (a yellow, a gold) keeps the logo's
// own colour and carries dark text instead, as Material does: darkening it would turn it brown.
func lightAccent(colors Colors, options Options) Colors {
	fill, onFill := options.Accent, onDarkAccent
	if Contrast(options.Accent, white) >= textContrast || Contrast(options.Accent, onDarkAccent) < textContrast {
		onFill = white
		fill = towardUntil(options.Accent, black, func(candidate Color) bool {
			return Contrast(candidate, white) >= textContrast
		})
	}
	colors.Accent = fill
	colors.OnAccent = onFill
	colors.AccentEdge = fill.Shade(edgeShade)
	colors.AccentSide = fill.Shade(sideShade)
	colors.AccentTint = readableTint(fill.Mix(tintBase(colors, options.Style), tintAmount(options.Style)), colors)
	colors.AccentInk = towardUntil(fill, black, func(candidate Color) bool {
		return Contrast(candidate, colors.Background) >= textContrast && Contrast(candidate, colors.AccentTint) >= textContrast
	})
	colors.FocusRing = colors.AccentInk
	colors.Selection = fill.WithAlpha(selectionAlpha)
	return colors
}

// In dark mode the accent is a light fill under dark text.
func darkAccent(colors Colors, options Options) Colors {
	fill := towardUntil(options.Accent, white, func(candidate Color) bool {
		return Contrast(candidate, onDarkAccent) >= textContrast && Contrast(candidate, colors.Background) >= controlContrast
	})
	colors.Accent = fill
	colors.OnAccent = onDarkAccent
	colors.AccentEdge = fill.Shade(edgeShade)
	colors.AccentSide = fill.Shade(sideShade)
	colors.AccentTint = readableTint(options.Accent.Mix(colors.Background, darkTint), colors)
	colors.AccentInk = towardUntil(fill, white, func(candidate Color) bool {
		return Contrast(candidate, colors.Background) >= textContrast && Contrast(candidate, colors.AccentTint) >= textContrast
	})
	colors.FocusRing = colors.AccentInk
	colors.Selection = fill.WithAlpha(selectionAlpha)
	return colors
}

func tintBase(colors Colors, style Style) Color {
	if style == StyleMaterial {
		return colors.Surface
	}
	return colors.Background
}

func tintAmount(style Style) float64 {
	if style == StyleMaterial {
		return materialTint
	}
	return paperTint
}

// readableTint moves a tint toward the background until body text still reads on it: a very light
// accent would otherwise make a dark-mode tint too bright for light ink.
func readableTint(tint Color, colors Colors) Color {
	return towardUntil(tint, colors.Background, func(candidate Color) bool {
		return Contrast(colors.Ink, candidate) >= textContrast && Contrast(colors.InkSoft, candidate) >= textContrast
	})
}

// towardUntil moves a colour toward a target in small steps until it is good enough, or until it
// becomes the target, which always is.
func towardUntil(color Color, target Color, goodEnough func(Color) bool) Color {
	candidate := color
	for step := 0; step < maximumSteps && !goodEnough(candidate); step++ {
		candidate = candidate.Mix(target, stepToward)
	}
	if goodEnough(candidate) {
		return candidate
	}
	return target
}
