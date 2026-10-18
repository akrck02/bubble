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
		shape.ControlRadius, shape.ButtonRadius, shape.CardRadius, shape.PanelRadius, shape.FloatingRadius = 12, pill, 20, 28, 16
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
	monospaceFamily = "ui-monospace, \"SF Mono\", Menlo, Consolas, monospace"
)

func typeOf(density Density) Type {
	if density == DensityCompact {
		return Type{Family: systemFamily, MonospaceFamily: monospaceFamily, SmallSize: 11, BodySize: 13, HeadingSize: 15, TitleSize: 20, DisplaySize: 28, IconSize: 16, LineHeight: 1.4}
	}
	return Type{Family: systemFamily, MonospaceFamily: monospaceFamily, SmallSize: 13, BodySize: 16, HeadingSize: 20, TitleSize: 30, DisplaySize: 44, IconSize: 20, LineHeight: 1.5}
}

func spaceOf(density Density) Space {
	if density == DensityCompact {
		return Space{Gap: 8, Padding: 12, ControlHeight: 28, TouchTarget: 24}
	}
	return Space{Gap: 12, Padding: 20, ControlHeight: 44, TouchTarget: 44}
}
