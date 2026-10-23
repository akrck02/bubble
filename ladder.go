package bubble

// The ladder of containers and the paired colours Material's parts need, made for paper and glass
// from their own colours, so the same parts draw in every style.

// How far each step of the ladder goes from the surface toward the sunken colour.
const (
	containerStep        = 1.0 / 3
	dangerWash           = 0.85 // how much of the background a danger container keeps
	inverseAccentTowards = 0.45 // how far the accent moves toward the inverse surface's opposite
)

func ladderOf(colors Colors, mode Mode) Colors {
	colors.ContainerLow, colors.ContainerHighest = colors.Surface, colors.Sunken
	colors.Container = colors.Surface.Mix(colors.Sunken, containerStep)
	colors.ContainerHigh = colors.Surface.Mix(colors.Sunken, 2*containerStep)
	colors.OnAccentTint = colors.Ink
	colors.AccentContainer, colors.OnAccentContainer = colors.AccentTint, colors.Ink
	colors.OnDanger = white
	if mode == ModeDark {
		colors.OnDanger = onDarkAccent
	}
	colors.DangerContainer = colors.Danger.Mix(colors.Background, dangerWash)
	colors.OnDangerContainer = colors.Ink
	colors.Inverse, colors.OnInverse = colors.Ink, colors.Background
	toward := white
	if mode == ModeDark {
		toward = black
	}
	colors.InverseAccent = towardUntil(colors.Accent.Mix(toward, inverseAccentTowards), toward, func(candidate Color) bool {
		return Contrast(candidate, colors.Inverse) >= textContrast
	})
	return colors
}
