// Command herbarium is the last level of Bubble's quest (docs/4-the-quest): the herbarium of
// examples/herbarium, drawn with go-ui's own controls in a Bubble theme.
//
//	go run ./examples/herbarium                    paper, by day
//	go run ./examples/herbarium -dark              at night
//	go run ./examples/herbarium -style material    in Material
//	go run ./examples/herbarium -glass             glass
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/akrck02/bubble"
	"github.com/akrck02/bubble/goui"
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
	"github.com/akrck02/go-ui/window"
)

// The herbarium's green, the same accent its web pages were generated for.
const accent = "#4C9A5B"

// The widths of the index and of the plant's sheet, as in the web pages (16rem and 26rem).
const (
	sidebarWidth = 256
	sheetWidth   = 416
	titleSize    = 30 // Bubble's --text-title
)

type plant struct {
	name, latin, bed string
	season           int
}

var seasons = []string{"Spring", "Summer", "Autumn", "Winter", "All year"}

var beds = []string{"Herb bed", "South wall", "Vegetable patch", "Rockery", "Under the oak", "Shady corner"}

// herbarium is the program's state: the plants, which one is open, and the search.
type herbarium struct {
	design bubble.Theme // Bubble's measures that go-ui's theme has no name for
	theme  ui.Theme
	ground *draw2d.Image // glass's ground; nil on opaque finishes
	plants []plant
	chosen int
	search string
}

func main() {
	options := window.Options{Title: "Herbarium", MinWidth: 1280, MinHeight: 760}
	options.RegisterFlags(flag.CommandLine)
	style := flag.String("style", "paper", "paper or material")
	glass := flag.Bool("glass", false, "the glass finish (paper only)")
	flag.Parse()
	options.Fonts = goui.FontsFor(bubble.Style(*style)) // Roboto for Material, as Android writes

	garden, err := newHerbarium(bubble.Style(*style), *glass, options)
	if err != nil {
		fmt.Fprintln(os.Stderr, "herbarium:", err)
		os.Exit(1)
	}
	if err := window.Run(options, garden.draw); err != nil {
		fmt.Fprintln(os.Stderr, "herbarium:", err)
		os.Exit(1)
	}
}

// newHerbarium makes Bubble's theme from the four choices and turns it into go-ui's.
func newHerbarium(style bubble.Style, glass bool, options window.Options) (*herbarium, error) {
	colour, err := bubble.ParseColor(accent)
	if err != nil {
		return nil, err
	}
	mode := bubble.ModeLight
	if options.Dark {
		mode = bubble.ModeDark
	}
	theme, err := bubble.NewTheme(bubble.Options{Style: style, Mode: mode, Density: bubble.DensityComfortable, Accent: colour, Glass: glass})
	if err != nil {
		return nil, err
	}
	garden := &herbarium{design: theme, theme: goui.Theme(theme), plants: []plant{
		{"Rosemary", "Salvia rosmarinus", "Herb bed", 0},
		{"Lavender", "Lavandula angustifolia", "South wall", 1},
		{"Sunflower", "Helianthus annuus", "Vegetable patch", 1},
		{"Heather", "Calluna vulgaris", "Rockery", 2},
		{"Snowdrop", "Galanthus nivalis", "Under the oak", 3},
		{"Hart's-tongue fern", "Asplenium scolopendrium", "Shady corner", 4},
	}}
	if theme.Glass.Enabled {
		garden.ground = draw2d.NewImage(goui.Ground(theme, options.Width, options.Height))
	}
	return garden, nil
}

func (garden *herbarium) draw(context *ui.Context) {
	context.Theme = garden.theme
	if garden.ground == nil {
		garden.frame(context)
		return
	}
	// On glass the ground is a picture, and the whole screen a layer of panes over it.
	width, height := context.Size()
	context.Image(garden.ground, ui.Layout{Width: ui.Grow(), Height: ui.Grow()}, 0)
	context.Overlay("screen", ui.LayerOptions{Width: ui.Fixed(width), Height: ui.Fixed(height)}, func() { garden.frame(context) })
}

// frame is Bubble's .app: the index on the left, the plants, and the open plant's sheet.
func (garden *herbarium) frame(context *ui.Context) {
	theme := context.Theme
	context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Grow(), Padding: ui.Pad(theme.Gap), Gap: theme.Gap}, func() {
		garden.pane(context, ui.Layout{Direction: ui.Vertical, Width: ui.Fixed(sidebarWidth), Height: ui.Grow(), Padding: ui.PadXY(12, 14), Gap: theme.Gap}, func() {
			context.Box(ui.Layout{Height: ui.Fixed(theme.ControlHeight), Padding: ui.PadXY(8, 0), Gap: 10, AlignY: ui.Center}, func() {
				context.IconSized("sprout", 28, theme.Ink)
				context.Text("Herbarium", ui.TextStyle{Weight: draw2d.Bold})
			})
			context.Selectable("sprout", "Plants", true)
			context.Selectable("calendar", "Seasons", false)
			context.Selectable("map-pin", "Beds", false)
		})
		context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.PadXY(theme.Padding, theme.Padding), Gap: theme.Gap}, func() {
			context.Text("Plants", ui.TextStyle{Size: titleSize, Weight: draw2d.Bold})
			context.TextFieldWith("search", &garden.search, ui.TextFieldOptions{Width: ui.Grow(), Placeholder: "Search the herbarium"})
			garden.list(context)
		})
		garden.sheet(context)
	})
}

// list is Bubble's .list: one row per plant that matches the search; a click opens its sheet.
func (garden *herbarium) list(context *ui.Context) {
	theme := context.Theme
	for index, each := range garden.plants {
		if !matches(each.name, garden.search) {
			continue
		}
		look := garden.surface(float32(garden.design.Shape.CardRadius))
		chosen := look
		chosen.Fill = theme.Sunken
		looks := ui.Looks{Normal: look, Hovered: chosen, Pressed: chosen}
		if index == garden.chosen {
			looks.Normal = chosen
		}
		layout := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(float32(garden.design.Space.TouchTarget) + 12), Padding: ui.PadXY(theme.Padding, 0), Gap: theme.Gap, AlignY: ui.Center}
		if context.Clickable(each.name, layout, looks, func() {
			context.Text(each.name, ui.TextStyle{Weight: draw2d.Bold})
			context.Spacer()
			context.Text(each.bed, ui.TextStyle{Size: theme.SmallTextSize, Color: theme.InkSoft})
			context.Chip("", seasons[each.season], ui.ChipLook{})
		}) {
			garden.chosen = index
		}
	}
}

// sheet is the open plant's pane: its fields, when it flowers, and what can be done with it.
func (garden *herbarium) sheet(context *ui.Context) {
	theme := context.Theme
	open := &garden.plants[garden.chosen]
	garden.pane(context, ui.Layout{Direction: ui.Vertical, Width: ui.Fixed(sheetWidth), Padding: ui.Pad(theme.Padding), Gap: theme.Gap}, func() {
		context.Heading(open.name)
		field(context, "Name", func() { context.TextFieldWith("name", &open.name, ui.TextFieldOptions{Width: ui.Grow()}) })
		field(context, "Latin name", func() { context.TextFieldWith("latin", &open.latin, ui.TextFieldOptions{Width: ui.Grow()}) })
		bed := indexOf(beds, open.bed)
		field(context, "Where it grows", func() {
			if context.Select("bed", &bed, beds) {
				open.bed = beds[bed]
			}
		})
		field(context, "Flowers in", func() { context.Tabs("season", &open.season, seasons...) })
		context.Box(ui.Layout{Width: ui.Grow(), Gap: 8, AlignY: ui.Center}, func() {
			context.FlatIconButton("trash-2", "Remove")
			context.Spacer()
			context.Button("Cancel")
			context.PrimaryButton("Save")
		})
	})
}

// pane is Bubble's .pane: a surface with a fine line on paper, a pane of glass on glass, and on
// Material a tinted surface that floats.
func (garden *herbarium) pane(context *ui.Context, layout ui.Layout, content func()) {
	theme := context.Theme
	context.Frame(layout, garden.surface(theme.PanelRadius), content)
}

// surface is how Bubble's CSS draws a card or a pane in the theme's style, with the given corners.
func (garden *herbarium) surface(radius float32) ui.Look {
	theme, colours := garden.theme, garden.design.Colors
	switch {
	case garden.ground != nil:
		return ui.Look{Fill: theme.Surface, Radius: radius}
	case garden.design.Options.Style == bubble.StyleMaterial:
		// Material tells surfaces apart by tone: the low container, with no line
		return ui.Look{Fill: draw2d.MustHex(colours.ContainerLow.Hex()), Radius: radius}
	default:
		return ui.Look{Fill: theme.Surface, Border: theme.Line, BorderWidth: 1, Radius: radius}
	}
}

// field writes a field's name over it, as Bubble's .field does.
func field(context *ui.Context, name string, control func()) {
	context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Gap: 6}, func() {
		context.Text(name, ui.TextStyle{Size: context.Theme.SmallTextSize, Weight: draw2d.Bold})
		control()
	})
}

// matches says whether a plant's name holds the search, whatever the case.
func matches(name, search string) bool {
	return strings.Contains(strings.ToLower(name), strings.ToLower(strings.TrimSpace(search)))
}

// indexOf is where value is in values, or 0 when it is not there.
func indexOf(values []string, value string) int {
	for index, each := range values {
		if each == value {
			return index
		}
	}
	return 0
}
