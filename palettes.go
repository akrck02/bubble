package bubble

// The colours that do not depend on the project's accent, for each style and mode.
// Paper is warm: a paper page by day, the slate of the logos' tile by night. Material is the
// neutral grey of Material Design's tonal surfaces.

func baseColors(style Style, mode Mode) Colors {
	colors := neutralColors(style, mode)
	if mode == ModeLight {
		colors.Positive = hex("#1E7D4F")
		colors.Warning = hex("#8A5300")
		colors.Danger = hex("#B3261E")
		colors.Shadow = hex("#1F232929")
		colors.Scrim = hex("#1F232966")
	} else {
		colors.Positive = hex("#4CC38A")
		colors.Warning = hex("#E5A64A")
		colors.Danger = hex("#EF7A6C")
		colors.Shadow = hex("#00000066")
		colors.Scrim = hex("#0000008C")
	}
	colors.Tile = hex("#2F3540")
	return colors
}

func neutralColors(style Style, mode Mode) Colors {
	switch {
	case style == StylePaper && mode == ModeLight:
		return Colors{
			Background: hex("#F7F5F0"), Surface: hex("#FDFCF9"), Sunken: hex("#EFEAE0"),
			Ink: hex("#1F2329"), InkSoft: hex("#555A62"), InkQuiet: hex("#62666E"),
			Line: hex("#DDD8CE"), ControlLine: hex("#8C8578"),
		}
	case style == StylePaper:
		return Colors{
			Background: hex("#1A1D24"), Surface: hex("#22262E"), Sunken: hex("#15181E"),
			Ink: hex("#E9E7E1"), InkSoft: hex("#B0B3BB"), InkQuiet: hex("#8F939C"),
			Line: hex("#2F3540"), ControlLine: hex("#6E7480"),
		}
	case mode == ModeLight:
		return Colors{
			Background: hex("#F4F4F6"), Surface: hex("#FFFFFF"), Sunken: hex("#E9E9EE"),
			Ink: hex("#1B1B1F"), InkSoft: hex("#45474F"), InkQuiet: hex("#5E616A"),
			Line: hex("#D5D6DD"), ControlLine: hex("#74777F"),
		}
	default:
		return Colors{
			Background: hex("#121316"), Surface: hex("#1E1F24"), Sunken: hex("#2A2B31"),
			Ink: hex("#E4E2E6"), InkSoft: hex("#C5C6D0"), InkQuiet: hex("#9A9BA5"),
			Line: hex("#34363C"), ControlLine: hex("#8E9099"),
		}
	}
}

// Glass is not paper made transparent: it is paper's standards reimagined on glass, after the
// first Valhalla. Panels are crystalline fills with no outline, told apart by how bright they
// are; the ground is near black by night and a pale frost by day, lit by two soft, sober lights.
// Text is measured over a panel over the ground and over each light, the worst cases.
const glassBlur = 24

// A wallpaper can be any picture, from black to white, so the veil over it is measured against
// both: it keeps enough of the ground to carry every ink, and still lets the picture glow through.
// A dialog floats over other text, which may be anything from black to white: its fill is dense
// enough to carry every ink over both, behind the scrim, and lets only a blur of the page through.
// By day a dark page behind shows more, so the frost is a little denser.
const (
	glassDialogOpacity = 0.86
	frostDialogOpacity = 0.9
)

// By night the picture is dimmed first, as the first Valhalla's dark wallpapers were.
const (
	glassVeilOpacity = 0.72
	nightDimming     = 0.45
)

func glassOf(mode Mode) Glass {
	if mode == ModeLight {
		return Glass{
			Enabled: true, Blur: glassBlur,
			Pane: white.WithAlpha(0.55), Raised: white.WithAlpha(0.85),
			Dialog: hex("#EFECE6").WithAlpha(frostDialogOpacity),
			Glows:  [2]Color{hex("#FFFFFF"), hex("#CBC6BB")},
			Veil:   hex("#E6E3DC").WithAlpha(glassVeilOpacity), Dimming: 1,
		}
	}
	return Glass{
		Enabled: true, Blur: glassBlur,
		Pane: white.WithAlpha(0.06), Raised: white.WithAlpha(0.12),
		Dialog: hex("#16181D").WithAlpha(glassDialogOpacity),
		Glows:  [2]Color{hex("#24272D"), hex("#181A1F")},
		Veil:   hex("#0B0C0F").WithAlpha(glassVeilOpacity), Dimming: nightDimming,
	}
}

// glassColors gives the neutral colours of glass: its own ground and inks, and the solid surface a
// panel becomes when the system asks for less transparency.
func glassColors(colors Colors, glass Glass, mode Mode) Colors {
	if mode == ModeLight {
		colors.Background, colors.Ink, colors.InkSoft, colors.InkQuiet = hex("#E6E3DC"), hex("#1C1F24"), hex("#4A4F57"), hex("#585D65")
		colors.ControlLine = hex("#77726A")
		colors.Line = colors.Ink.WithAlpha(0.1).Over(colors.Background)
		colors.Sunken = colors.Background.Mix(colors.Ink, 0.06)
		colors.Positive, colors.Warning, colors.Danger = hex("#17663F"), hex("#7A4900"), hex("#A1211A")
	} else {
		colors.Background, colors.Ink, colors.InkSoft, colors.InkQuiet = hex("#0B0C0F"), hex("#EDEEF0"), hex("#BDC1C8"), hex("#ADB2BA")
		colors.ControlLine = hex("#8E949D")
		colors.Line = white.WithAlpha(0.08).Over(colors.Background)
		colors.Sunken = colors.Background.Mix(black, 0.4)
	}
	colors.Surface = glass.Pane.Over(colors.Background)
	return colors
}

// soberAccent keeps the accent for the main action, the focus ring and links, and makes the
// chosen state neutral: on glass, what is chosen is brighter, not coloured.
func soberAccent(colors Colors, glass Glass, mode Mode) Colors {
	colors.AccentTint = glass.Raised.Over(colors.Background)
	target := white
	if mode == ModeLight {
		target = black
	}
	colors.AccentInk = towardUntil(colors.AccentInk, target, func(candidate Color) bool {
		return Contrast(candidate, colors.Background) >= textContrast && Contrast(candidate, colors.AccentTint) >= textContrast
	})
	colors.FocusRing = colors.AccentInk
	return colors
}
