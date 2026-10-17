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

func glassOf(colors Colors, mode Mode) Glass {
	if mode == ModeLight {
		return Glass{Enabled: true, Surface: colors.Surface.WithAlpha(0.55), Line: colors.Ink.WithAlpha(0.12), Blur: 24}
	}
	return Glass{Enabled: true, Surface: colors.Surface.WithAlpha(0.6), Line: hex("#FFFFFF").WithAlpha(0.1), Blur: 24}
}
