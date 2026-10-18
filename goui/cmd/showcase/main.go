// Command showcase draws one screen of Bubble 4 with go-ui: the theme comes from the bubble package,
// package goui turns it into go-ui's, and the screen is made of go-ui's own controls. web.html is
// the same screen in a browser, with Bubble's CSS; compare.sh puts the two side by side
// (docs/5-guides/1-go-ui.md).
//
//	go run ./cmd/showcase -capture showcase.png
//	go run ./cmd/showcase -dark -capture showcase-dark.png
//	go run ./cmd/showcase -style material -capture showcase-material.png
//	go run ./cmd/showcase -glass -capture showcase-glass.png
package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/akrck02/bubble"
	"github.com/akrck02/bubble/goui"
	"github.com/akrck02/go-renderer/draw2d"
	"github.com/akrck02/go-renderer/opengl/overlay"
	ui "github.com/akrck02/go-ui"
	"github.com/akrck02/go-ui/window"
)

// logoPicture is Bubble's logo (github.com/akrck02/bubble, docs/images/logo.svg), drawn to a PNG.
//
//go:embed logo.png
var logoPicture []byte

// Measures of this screen that are not controls of their own, as in web.html.
const (
	sidebarWidth = 240 // 15rem
	railWidth    = 96  // 6rem: material's index
	brandLogo    = 28
	titleSize    = 30
	swatchHeight = 36
	letterRoom   = 0.14 // the name's letter spacing, in ems
)

// screen is the showcase's state: the Bubble theme it draws with, and what the controls change.
type screen struct {
	design bubble.Theme
	theme  ui.Theme
	logo   *draw2d.Image
	ground *draw2d.Image // glass's ground; nil on opaque finishes
	title  string
	done   bool
}

func main() {
	options := window.Options{Title: "Bubble 4 in Go", MinWidth: 1280, MinHeight: 760}
	options.RegisterFlags(flag.CommandLine)
	style := flag.String("style", "paper", "paper or material")
	accent := flag.String("accent", "#81D4FA", "the project's accent, from its logo")
	glass := flag.Bool("glass", false, "the glass finish (paper only, as in Bubble)")
	stripes := flag.Bool("stripes", false, "on glass, a ground of sharp stripes, to see the blur behind the panes")
	font := flag.String("font", "", "a TrueType file to write in, its bold beside it as NAME-Bold.ttf (to match a browser)")
	flag.Parse()

	if *font != "" {
		if err := useFont(*font); err != nil {
			fail(err)
		}
	}
	showcase, err := newScreen(*style, *accent, *glass, options)
	if err != nil {
		fail(err)
	}
	if *stripes && showcase.ground != nil {
		showcase.ground = draw2d.NewImage(stripedGround(options.Width, options.Height))
	}
	if err := window.Run(options, showcase.draw); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "showcase:", err)
	os.Exit(1)
}

// useFont writes text in a face of the system's, such as the one a browser picks for system-ui,
// so a comparison with web.html is never just about the typeface.
func useFont(path string) error {
	regular, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	bold, err := os.ReadFile(path[:len(path)-len(".ttf")] + "-Bold.ttf")
	if err != nil {
		return err
	}
	overlay.UseFonts(overlay.Fonts{Regular: regular, Medium: bold, Bold: bold})
	return nil
}

func newScreen(style, accent string, glass bool, options window.Options) (*screen, error) {
	colour, err := bubble.ParseColor(accent)
	if err != nil {
		return nil, err
	}
	mode := bubble.ModeLight
	if options.Dark {
		mode = bubble.ModeDark
	}
	theme, err := bubble.NewTheme(bubble.Options{Style: bubble.Style(style), Mode: mode, Density: bubble.DensityComfortable, Accent: colour, Glass: glass})
	if err != nil {
		return nil, err
	}
	logo, err := draw2d.DecodeImage(bytes.NewReader(logoPicture))
	if err != nil {
		return nil, err
	}
	showcase := &screen{design: theme, theme: goui.Theme(theme), logo: logo, title: "Write the release notes", done: true}
	if theme.Glass.Enabled {
		showcase.ground = draw2d.NewImage(goui.Ground(theme, options.Width, options.Height))
	}
	return showcase, nil
}

func (showcase *screen) material() bool { return showcase.design.Options.Style == bubble.StyleMaterial }

func (showcase *screen) draw(context *ui.Context) {
	context.Theme = showcase.theme
	context.KeyboardNavigation = ui.Always // the showcase is only an interface: Tab walks every control
	if showcase.ground == nil {
		showcase.frame(context)
		return
	}
	// On glass the ground is a picture, and the whole screen a layer of panes over it.
	width, height := context.Size()
	context.Image(showcase.ground, ui.Layout{Width: ui.Grow(), Height: ui.Grow()}, 0)
	context.Overlay("screen", ui.LayerOptions{Width: ui.Fixed(width), Height: ui.Fixed(height)}, func() { showcase.frame(context) })
}

// frame is Bubble's: the index on the left, the whole height, and the page beside it.
func (showcase *screen) frame(context *ui.Context) {
	context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Grow()}, func() {
		switch {
		case showcase.material():
			showcase.rail(context)
		case showcase.ground != nil:
			// on glass the index is a pane of its own, with no line beside it
			context.Frame(ui.Layout{Width: ui.Fixed(sidebarWidth), Height: ui.Grow()}, ui.Look{Fill: context.Theme.Surface}, func() {
				showcase.sidebar(context, sidebarWidth)
			})
		default:
			showcase.sidebar(context, sidebarWidth-1) // the line beside it takes the last pixel
			context.Separator()
		}
		showcase.page(context)
	})
}

var places = []struct{ icon, name string }{
	{"palette", "Tokens"}, {"layout-grid", "Components"}, {"layers", "Finishes"}, {"shield-check", "Accessibility"},
}

const chosenPlace = "Components"

func (showcase *screen) sidebar(context *ui.Context, width float32) {
	theme := context.Theme
	context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Fixed(width), Height: ui.Grow(), Padding: ui.PadXY(12, 14), Gap: theme.Gap}, func() {
		context.Box(ui.Layout{Height: ui.Fixed(theme.ControlHeight), Padding: ui.PadXY(8, 0), Gap: 10, AlignY: ui.Center}, func() {
			context.Image(showcase.logo, ui.Layout{Width: ui.Fixed(brandLogo), Height: ui.Fixed(brandLogo)}, 0)
			spaced(context, "BUBBLE", theme.SmallTextSize)
		})
		context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Gap: 2}, func() {
			for _, place := range places {
				context.Selectable(place.icon, place.name, place.name == chosenPlace)
			}
		})
	})
}

// spaced writes the name in capitals with Bubble's letter spacing, letter by letter.
func spaced(context *ui.Context, text string, size float32) {
	context.Box(ui.Layout{Gap: size * letterRoom, AlignY: ui.Center}, func() {
		for _, letter := range text {
			context.Text(string(letter), ui.TextStyle{Size: size, Weight: draw2d.Bold})
		}
	})
}

// rail is material's index: the logo alone on top, then each place as an icon in a pill with its
// name under it; the chosen pill is tinted, as in Material's navigation rail.
func (showcase *screen) rail(context *ui.Context) {
	theme := context.Theme
	context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Fixed(railWidth), Height: ui.Grow(), Padding: ui.PadXY(12, 14), Gap: theme.Gap}, func() {
		context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(theme.ControlHeight), AlignX: ui.Center, AlignY: ui.Center}, func() {
			context.Image(showcase.logo, ui.Layout{Width: ui.Fixed(brandLogo), Height: ui.Fixed(brandLogo)}, 0)
		})
		for _, place := range places {
			chosen := place.name == chosenPlace
			pill, weight := ui.Look{Radius: theme.ButtonRadius}, draw2d.Regular
			if chosen {
				pill.Fill, weight = theme.PrimaryTint, draw2d.Bold
			}
			context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.PadXY(0, 4), Gap: 4, AlignX: ui.Center}, func() {
				context.Frame(ui.Layout{Padding: ui.PadXY(16, 4)}, pill, func() { context.IconSized(place.icon, theme.IconSize, theme.Ink) })
				context.Text(place.name, ui.TextStyle{Size: theme.SmallTextSize, Weight: weight})
			})
		}
	})
}

func (showcase *screen) page(context *ui.Context) {
	theme := context.Theme
	context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow(), Padding: ui.PadXY(theme.Padding, theme.Padding*1.6), Gap: theme.Gap}, func() {
		showcase.pageHeader(context)
		showcase.swatches(context)
		context.Box(ui.Layout{Width: ui.Grow(), Columns: 2, Gap: theme.Gap, AlignY: ui.Start}, func() {
			showcase.controls(context)
			showcase.proofs(context)
		})
	})
}

func (showcase *screen) pageHeader(context *ui.Context) {
	theme := context.Theme
	context.Box(ui.Layout{Width: ui.Grow(), Gap: 16, AlignY: ui.End}, func() {
		context.Box(ui.Layout{Direction: ui.Vertical, Gap: 4}, func() {
			context.Text("Components", ui.TextStyle{Size: titleSize, Weight: draw2d.Bold})
			context.Text("Bubble 4, drawn the same in the browser and in Go", ui.TextStyle{Color: theme.InkSoft})
		})
		context.Spacer()
		context.Box(ui.Layout{Gap: 8, AlignY: ui.Center}, func() {
			showcase.segmented(context, "Paper", "Material")
			context.PrimaryIconButton("plus", "New component")
		})
	})
}

// segmented is Bubble's segmented control: a few buttons in a thin frame, the chosen one brighter.
func (showcase *screen) segmented(context *ui.Context, chosen string, others ...string) {
	theme := context.Theme
	frame := ui.Look{Border: theme.Line, BorderWidth: 1, Radius: theme.ButtonRadius}
	if showcase.ground != nil {
		frame = ui.Look{Fill: theme.Surface, Radius: theme.ButtonRadius}
	}
	context.Frame(ui.Layout{Padding: ui.Pad(2), Gap: 2}, frame, func() {
		for _, label := range append([]string{chosen}, others...) {
			look := ui.Look{Radius: theme.ButtonRadius}
			if label == chosen {
				look.Fill = theme.Sunken
			}
			layout := ui.Layout{Height: ui.Fixed(theme.ControlHeight), Padding: ui.PadXY(theme.ButtonPadding, 0), AlignY: ui.Center}
			context.Clickable(label, layout, ui.Looks{Normal: look, Hovered: look, Pressed: look}, func() {
				context.Text(label, ui.TextStyle{Weight: theme.ButtonWeight})
			})
		}
	})
}

// card is Bubble's .card: on paper a surface with a fine line, on material a tinted surface that
// floats, on glass a pane.
func (showcase *screen) card(context *ui.Context, layout ui.Layout, content func()) {
	context.Frame(layout, showcase.cardLook(context, true), content)
}

// materialTint is how much accent tint Bubble's CSS mixes into a material card.
const materialTint = 0.45

func (showcase *screen) cardLook(context *ui.Context, tinted bool) ui.Look {
	theme, colours := context.Theme, showcase.design.Colors
	radius := float32(showcase.design.Shape.CardRadius)
	switch {
	case showcase.ground != nil:
		return ui.Look{Fill: theme.Surface, Radius: radius}
	case showcase.material() && tinted:
		mixed := colours.Surface.Mix(colours.AccentTint, materialTint)
		return ui.Look{Fill: draw2d.Hex(mixed.Hex()), Radius: radius, Shadow: true}
	default:
		return ui.Look{Fill: theme.Surface, Border: theme.Line, BorderWidth: 1, Radius: radius}
	}
}

func (showcase *screen) swatches(context *ui.Context) {
	theme := context.Theme
	colours := showcase.design.Colors
	context.Box(ui.Layout{Width: ui.Grow(), Columns: 7, Gap: theme.Gap}, func() {
		for _, swatch := range []struct {
			name  string
			color draw2d.Color
		}{
			{"Accent", theme.Primary}, {"Ink", theme.Ink}, {"Soft ink", theme.InkSoft}, {"Surface", draw2d.Hex(colours.Surface.Hex())},
			{"Sunken", theme.Sunken}, {"Positive", theme.Positive}, {"Danger", theme.Danger},
		} {
			context.Frame(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.Pad(12), Gap: 8}, showcase.cardLook(context, false), func() {
				chip := ui.Look{Fill: swatch.color, Border: theme.Line, BorderWidth: 1, Radius: theme.ControlRadius}
				context.Frame(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(swatchHeight)}, chip, nil)
				context.Text(swatch.name, ui.TextStyle{Weight: draw2d.Bold})
			})
		}
	})
}

func (showcase *screen) controls(context *ui.Context) {
	theme := context.Theme
	showcase.card(context, ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.Pad(theme.Padding), Gap: theme.Gap}, func() {
		context.Heading("Controls")
		context.Box(ui.Layout{Gap: 8, AlignY: ui.Center}, func() {
			context.PrimaryButton("Save")
			context.Button("Cancel")
			context.FlatButton("More")
			context.FlatIconButton("pencil", "")
		})
		context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Gap: 6}, func() {
			context.Text("What is there to do?", ui.TextStyle{Size: theme.SmallTextSize, Weight: draw2d.Bold})
			context.TextFieldWith("title", &showcase.title, ui.TextFieldOptions{Width: ui.Grow()})
		})
		context.Checkbox("Done", &showcase.done)
	})
}

func (showcase *screen) proofs(context *ui.Context) {
	theme := context.Theme
	showcase.card(context, ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.Pad(theme.Padding), Gap: theme.Gap}, func() {
		context.Heading("What the tests prove")
		context.Box(ui.Layout{Direction: ui.Vertical, Gap: 8}, func() {
			for _, proof := range []string{
				"Text reads on every ground of every theme",
				"The accent reads whatever colour the logo is",
				"Controls and the focus ring stand out",
				"Glass text reads over any wallpaper",
				"Comfortable controls can be hit with a finger",
			} {
				context.Box(ui.Layout{Gap: 10, AlignY: ui.Center}, func() {
					context.IconSized("check", theme.IconSize, theme.Positive)
					context.Label(proof)
				})
			}
		})
	})
}

// stripeWidth is how wide each stripe of the test ground is, in pixels.
const stripeWidth = 24

// stripedGround is a ground of sharp light and dark stripes: under a frosted pane they turn into a
// soft grey, which shows the blur at work, where glass's own smooth ground would hide it.
func stripedGround(width, height int) *image.RGBA {
	picture := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			shade := uint8(40)
			if (x/stripeWidth+y/stripeWidth)%2 == 0 {
				shade = 200
			}
			picture.SetRGBA(x, y, color.RGBA{shade, shade, shade, 255})
		}
	}
	return picture
}
