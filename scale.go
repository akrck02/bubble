package bubble

// Corners, type and space for each style and density. Comfortable sizes are for pages and phones,
// where a finger needs 44 pixels; compact sizes are for desktop tools used with a pointer.

const pill = 999 // a corner so round that the ends become half circles

func shapeOf(style Style, density Density) Shape {
	// The focus ring sits on the control's own edge, with no gap: the edge lights up (after Bubble 3).
	shape := Shape{LineWidth: 1, FocusWidth: 2}
	if style == StylePaper {
		// Paper is cut, not moulded: buttons share the fields' medium corners, like cards on a desk.
		shape.ControlRadius, shape.ButtonRadius, shape.CardRadius, shape.PanelRadius, shape.FloatingRadius = 10, 10, 12, 16, pill
		shape.EdgeWidth = 2
	} else {
		// Material's own corners, as Go UI for Android drew them: rounded fields, pill buttons,
		// 12-pixel cards, 28-pixel dialogs and sheets, a 16-pixel floating button.
		shape.ControlRadius, shape.ButtonRadius, shape.CardRadius, shape.PanelRadius, shape.FloatingRadius = 20, pill, 12, 28, 16
		shape.Elevated = true
	}
	if density == DensityCompact {
		shape.ControlRadius *= 0.6
		if shape.ButtonRadius != pill {
			shape.ButtonRadius *= 0.6
		}
		shape.CardRadius *= 0.6
		shape.PanelRadius *= 0.6
	}
	return shape
}

const (
	systemFamily    = "system-ui, -apple-system, \"Segoe UI\", sans-serif"
	materialFamily  = "Roboto, " + systemFamily // Android's own typeface; goui carries it for Go
	monospaceFamily = "ui-monospace, \"SF Mono\", Menlo, Consolas, monospace"
)

func typeOf(style Style, density Density) Type {
	if style == StyleMaterial {
		return materialTypeOf(density)
	}
	if density == DensityCompact {
		return Type{Family: systemFamily, MonospaceFamily: monospaceFamily, SmallSize: 11, BodySize: 13, HeadingSize: 15, TitleSize: 20, DisplaySize: 28, IconSize: 16, LineHeight: 1.4}
	}
	return Type{Family: systemFamily, MonospaceFamily: monospaceFamily, SmallSize: 13, BodySize: 16, HeadingSize: 20, TitleSize: 30, DisplaySize: 44, IconSize: 20, LineHeight: 1.5}
}

// materialTypeOf is Material 3's type scale: body large 16, label large 14, title large 22,
// headline large 32, display medium 45, with 24-pixel icons.
func materialTypeOf(density Density) Type {
	if density == DensityCompact {
		return Type{Family: materialFamily, MonospaceFamily: monospaceFamily, SmallSize: 12, BodySize: 14, HeadingSize: 16, TitleSize: 22, DisplaySize: 32, IconSize: 20, LineHeight: 1.4}
	}
	return Type{Family: materialFamily, MonospaceFamily: monospaceFamily, SmallSize: 14, BodySize: 16, HeadingSize: 22, TitleSize: 32, DisplaySize: 45, IconSize: 24, LineHeight: 1.5}
}

func spaceOf(style Style, density Density) Space {
	switch {
	case style == StyleMaterial && density == DensityCompact:
		return Space{Gap: 8, Padding: 12, ControlHeight: 36, TouchTarget: 32}
	case style == StyleMaterial:
		return Space{Gap: 12, Padding: 16, ControlHeight: 48, TouchTarget: 48} // Material's 48-pixel targets
	case density == DensityCompact:
		return Space{Gap: 8, Padding: 12, ControlHeight: 28, TouchTarget: 24}
	}
	return Space{Gap: 12, Padding: 20, ControlHeight: 44, TouchTarget: 44}
}
