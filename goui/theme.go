// Package goui draws go-ui (https://github.com/akrck02/go-ui) in Bubble 4: it turns a Bubble theme
// into go-ui's own, so every control follows it, and paints glass's ground. It is a module of its
// own, so Bubble's package stays free of go-ui, and go-ui knows nothing of Bubble.
//
//	theme, err := bubble.NewTheme(bubble.Options{Style: bubble.StylePaper, Mode: bubble.ModeLight,
//		Density: bubble.DensityComfortable, Accent: accent})
//	context.Theme = goui.Theme(theme)
package goui

import (
	"github.com/akrck02/bubble"
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// Bubble's measures that its Go theme leaves to the CSS (../css/bubble.css), kept here so a Go
// screen and a web screen of the same design match.
const (
	lineHeight        = 1.5  // body text, from the tokens
	headingLineHeight = 1.2  // titles and headings (h1, h2)
	buttonPadding     = 0.9  // a button's room on each side: the panels' padding × this
	buttonGap         = 8    // between a button's icon and its label
	itemGap           = 12   // between an item's icon and its label, in the index
	itemPadding       = 12   // left and right inside an item of the index
	checkSize         = 22   // a check box's side
	checkRadius       = 6    // its corners
	checkBorder       = 2    // its edge
	checkMark         = 14   // its tick
	glassFieldEdge    = 0.16 // a glass field's soft edge: this much ink over nothing
)

// Theme turns a Bubble theme into go-ui's. The colours carry the same names, mostly; the sizes come
// from Bubble's density. What is chosen is brighter and heavier, never coloured, and the main
// button stands on its darker edge. On glass, panels and wells are the translucent panes, and the
// window's background is glass's ground (paint it with Ground).
func Theme(theme bubble.Theme) ui.Theme {
	colours, shape, space, kind := theme.Colors, theme.Shape, theme.Space, theme.Type
	result := ui.LightTheme()

	result.Background, result.Surface, result.Sunken = colour(colours.Background), colour(colours.Surface), colour(colours.Sunken)
	result.Ink, result.InkSoft, result.InkQuiet = colour(colours.Ink), colour(colours.InkSoft), colour(colours.InkQuiet)
	result.Line, result.ControlLine = colour(colours.Line), colour(colours.ControlLine)
	result.Primary, result.OnPrimary, result.PrimaryTint = colour(colours.Accent), colour(colours.OnAccent), colour(colours.AccentTint)
	result.Positive, result.Warning, result.Danger = colour(colours.Positive), colour(colours.Warning), colour(colours.Danger)
	result.FocusRing, result.Selection = colour(colours.FocusRing), colour(colours.Selection)
	result.Shadow, result.Scrim = colour(colours.Shadow), colour(colours.Scrim)

	result.TextSize, result.SmallTextSize = float32(kind.BodySize), float32(kind.SmallSize)
	result.HeadingSize, result.IconSize = float32(kind.HeadingSize), float32(kind.IconSize)
	result.LineHeight, result.HeadingLineHeight = lineHeight, headingLineHeight
	result.ControlHeight, result.Gap, result.Padding = float32(space.ControlHeight), float32(space.Gap), float32(space.Padding)
	result.ControlRadius, result.PanelRadius = float32(shape.ControlRadius), float32(shape.PanelRadius)
	if !shape.Elevated {
		result.ShadowBlur, result.ShadowDrop = 0, 0 // paper has no shadows: lines and fills only
	}

	result.ButtonRadius = float32(shape.ButtonRadius)
	result.ButtonPadding = float32(space.Padding) * buttonPadding
	result.ButtonGap, result.ButtonWeight = buttonGap, draw2d.Bold
	result.PrimaryEdge, result.EdgeWidth = colour(colours.AccentEdge), float32(shape.EdgeWidth)
	result.FlatHeight, result.FlatPadding, result.FlatInk = float32(space.ControlHeight), result.ButtonPadding, result.Ink

	result.Chosen, result.ChosenInk, result.ChosenWeight = result.Sunken, result.Ink, draw2d.Bold
	result.ItemPadding, result.ItemGap, result.ItemIconInk = itemPadding, itemGap, result.Ink

	result.CheckSize, result.CheckRadius, result.CheckBorder, result.CheckMark = checkSize, checkRadius, checkBorder, checkMark
	result.CheckLine, result.CheckedBorder = result.ControlLine, result.PrimaryEdge

	if theme.Glass.Enabled {
		glassy(&result, theme)
	}
	return result
}

// glassy makes the panels crystal: nothing on glass is opaque, and what Bubble sinks takes the
// raised pane, as css/bubble.css does.
func glassy(result *ui.Theme, theme bubble.Theme) {
	glass := theme.Glass
	result.Surface = colour(glass.Pane)
	result.Sunken = colour(glass.Raised)
	result.Chosen = result.Sunken
	result.ControlLine = colour(theme.Colors.Ink.WithAlpha(glassFieldEdge))
}

func colour(value bubble.Color) draw2d.Color { return draw2d.Hex(value.Hex()) }
