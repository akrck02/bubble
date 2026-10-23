# 🐹 Go

The Go API of Bubble's two modules, as `go doc -all` prints it. This page is made by `docs/6-reference/generate.py`; do not edit it by hand.

| You want to… | Use |
|---|---|
| make a theme | `bubble.NewTheme(bubble.Options{…})`, or start from `bubble.DefaultOptions()` |
| read a colour | `bubble.ParseColor("#4C9A5B")` |
| write the CSS tokens | `bubble.CSS(writer, accent)`, or the `bubble-css` command |
| check two colours | `bubble.Contrast(a, b)` |
| dress Go UI | `goui.Theme(theme)`, and `goui.Ground` for glass |

## Package bubble

`github.com/akrck02/bubble`

```
package bubble // import "github.com/akrck02/bubble"

Package bubble is the Bubble design system: the tokens every Bubble interface is
drawn with, in any technology. A theme comes from four choices (a style, a mode,
a density and the project's accent colour) and gives the same names everywhere:
the CSS of BubbleUI and go-ui's controls read them alike.

VARIABLES

var DefaultAccent = hex("#2C5C9E")
    DefaultAccent is the accent of a project that has not chosen one: Bubble's
    own blue.

var InvalidColour = define("bubble.invalid-colour", "a colour must be written as #RRGGBB or #RRGGBBAA")
    InvalidColour means a colour is not written as #RRGGBB or #RRGGBBAA.

var UnknownDensity = define("bubble.unknown-density", "the density must be comfortable or compact")
    UnknownDensity means the density is neither comfortable nor compact.

var UnknownMode = define("bubble.unknown-mode", "the mode must be light or dark")
    UnknownMode means the mode is neither light nor dark.

var UnknownStyle = define("bubble.unknown-style", "the style must be paper or material")
    UnknownStyle means the style is neither paper nor material.


FUNCTIONS

func CSS(writer io.Writer, accent Color) error
    CSS writes the tokens of every style, mode and density as CSS custom
    properties, for one accent.

    The page chooses with attributes on <html>: data-style="material",
    data-mode="light" or "dark" (without it, the system's choice),
    data-density="compact" and data-glass. Paper, light and comfortable need no
    attribute.

func Contrast(first Color, second Color) float64
    Contrast is the WCAG contrast ratio between two colours, from 1 to 21.
    Text needs 4.5 (3 when it is large), the edges of controls and the focus
    ring need 3.


TYPES

type Code struct {
	// Has unexported fields.
}
    Code is an error with a stable name that programs check and people can
    search for. A code is made once, so returning or comparing one costs no
    allocation.

func (code *Code) Error() string
    Error gives the code's name and what it means.

func (code *Code) Name() string
    Name is the code's stable name, such as "bubble.invalid-colour".

type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}
    Color is a colour in sRGB with its opacity, each channel from 0 to 255.

func ParseColor(text string) (Color, error)
    ParseColor reads a colour written as #RRGGBB or #RRGGBBAA.

func (color Color) Hex() string
    Hex writes the colour as #RRGGBB, or #RRGGBBAA when it is not opaque.

func (color Color) Luminance() float64
    Luminance is the colour's relative luminance, as WCAG defines it: 0 for
    black, 1 for white.

func (color Color) Mix(toward Color, amount float64) Color
    Mix moves the colour toward another by amount, from 0 (itself) to 1 (the
    other one).

func (color Color) Over(background Color) Color
    Over paints the colour, with its opacity, over an opaque background and
    gives what the eye sees.

func (color Color) Shade(amount float64) Color
    Shade multiplies each channel by amount: 0.85 and 0.70 give the left and
    right sides of a voxel block, and the edge under a primary button.

func (color Color) WithAlpha(opacity float64) Color
    WithAlpha gives the same colour with another opacity, from 0 to 1.

type Colors struct {
	Background  Color // the page or window
	Surface     Color // fields, raised cards, menus
	Sunken      Color // code blocks, tracks, wells
	Ink         Color // body text
	InkSoft     Color // secondary text and icons
	InkQuiet    Color // hints and captions
	Line        Color // separators and the borders of cards
	ControlLine Color // borders of fields and secondary buttons: strong enough to see the control
	Accent      Color // the main action's fill
	OnAccent    Color // text and icons on the accent
	AccentEdge  Color // the block edge under a primary button: the logo's darker side, ×0.70
	AccentSide  Color // the logo's lighter side, ×0.85, for pressed states
	AccentInk   Color // accent-coloured text and icons on the background
	AccentTint  Color // a wash of the accent: the selected item, a checkpoint, a Material container
	Positive    Color
	Warning     Color
	Danger      Color
	FocusRing   Color // the ring around whatever the keyboard is on
	Selection   Color // behind selected text
	Shadow      Color
	Scrim       Color // over everything behind a sheet or a dialog
	Tile        Color // the dark tile every voxel logo stands on
}
    Colors are the theme's colours, named by what they are for.

type Density string
    Density is how much fits on screen: comfortable for web pages and phones,
    compact for desktop tools used with a pointer.

const (
	DensityComfortable Density = "comfortable"
	DensityCompact     Density = "compact"
)
    The two densities.

type Glass struct {
	Enabled bool
	Pane    Color    // the translucent fill of panels, cards and fields, with no outline
	Raised  Color    // a brighter fill, for what is hovered or chosen
	Dialog  Color    // the denser fill of a dialog, which floats over other text
	Glows   [2]Color // the two soft lights behind the panels, drawn blurred on the ground
	Veil    Color    // laid over a wallpaper, so any picture behind keeps the text readable
	Dimming float64  // how bright a wallpaper is kept before the veil, from 0 to 1 (CSS brightness)
	Blur    float64  // how much of what is behind shows through, in layout pixels
}
    Glass holds the crystalline panels of the glass variant; Enabled is false
    outside it.

type Mode string
    Mode is light or dark; interfaces follow the one the system asks for.

const (
	ModeLight Mode = "light"
	ModeDark  Mode = "dark"
)
    The two modes.

type Motion struct {
	Fast   time.Duration
	Medium time.Duration
	Slow   time.Duration
}
    Motion holds how long transitions take. Interfaces skip them when the system
    asks for less motion.

type Options struct {
	Style   Style
	Mode    Mode
	Density Density
	Accent  Color // the main colour of the project's logo; the theme adjusts it until it reads well
	Glass   bool  // frosted, translucent panels over the background (paper only)
}
    Options are the choices a theme is made from.

func DefaultOptions() Options
    DefaultOptions is paper, light, comfortable, with the default accent and no
    glass.

type Shape struct {
	ControlRadius  float64 // fields and selects
	ButtonRadius   float64 // buttons and chips
	CardRadius     float64
	PanelRadius    float64 // windows, sheets, dialogs
	FloatingRadius float64 // the floating action button
	LineWidth      float64
	EdgeWidth      float64 // the block edge under primary buttons; 0 when the style has none
	FocusWidth     float64
	FocusOffset    float64 // the gap between a control and its focus ring; 0, so the edge lights up, kept for go-ui
	Elevated       bool    // whether raised things cast a shadow (material) or sit on a line (paper)
}
    Shape holds corners, edges and elevation, in layout pixels.

type Space struct {
	Gap           float64 // between controls
	Padding       float64 // inside cards and panels
	ControlHeight float64
	TouchTarget   float64 // the smallest area a finger or pointer must be able to hit
}
    Space holds distances and the size of controls, in layout pixels.

type Style string
    Style is the visual language: paper is the default, material the alternative
    that Android users expect. Both give the same tokens, so an interface
    switches style without changing its controls.

const (
	StylePaper    Style = "paper"
	StyleMaterial Style = "material"
)
    The two styles.

type Theme struct {
	Options Options
	Colors  Colors
	Shape   Shape
	Type    Type
	Space   Space
	Motion  Motion
	Glass   Glass
}
    Theme holds every token of one combination of options.

func NewTheme(options Options) (Theme, error)
    NewTheme makes the theme of the given options.

type Type struct {
	Family          string // the system's own interface typeface
	MonospaceFamily string
	SmallSize       float64
	BodySize        float64
	HeadingSize     float64
	TitleSize       float64
	DisplaySize     float64
	IconSize        float64
	LineHeight      float64 // a multiple of the size
}
    Type holds the typefaces and the type scale, in layout pixels.
```

## Package goui

`github.com/akrck02/bubble/goui`, a module of its own (`goui/go.mod`).

```
package goui // import "github.com/akrck02/bubble/goui"

Package goui draws go-ui (https://github.com/akrck02/go-ui) in Bubble 4:
it turns a Bubble theme into go-ui's own, so every control follows it, and
paints glass's ground. It is a module of its own, so Bubble's package stays free
of go-ui, and go-ui knows nothing of Bubble.

    theme, err := bubble.NewTheme(bubble.Options{Style: bubble.StylePaper, Mode: bubble.ModeLight,
    	Density: bubble.DensityComfortable, Accent: accent})
    context.Theme = goui.Theme(theme)

FUNCTIONS

func Ground(theme bubble.Theme, width, height int) *image.RGBA
    Ground paints glass's ground the way the browser does, layer by layer:
    the background, the second light over it, the first light on top. Panels are
    drawn over this picture; it is smooth, so blurring what is behind them would
    change nothing.

func Theme(theme bubble.Theme) ui.Theme
    Theme turns a Bubble theme into go-ui's. The colours carry the same names,
    mostly; the sizes come from Bubble's density. What is chosen is brighter
    and heavier, never coloured, and the main button stands on its darker edge.
    On glass, panels and wells are the translucent panes, and the window's
    background is glass's ground (paint it with Ground).
```

## Command bubble-css

```
Command bubble-css writes the Bubble 4 tokens as CSS for a project's accent:

    go run github.com/akrck02/bubble/cmd/bubble-css -accent "#F6BD60" > styles/bubble-tokens.css

The accent is the main colour of the project's logo; Bubble adjusts it until it
reads well.
```
