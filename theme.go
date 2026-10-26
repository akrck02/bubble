// Package bubble is the Bubble design system: the tokens every Bubble interface is drawn with, in
// any technology. A theme comes from four choices (a style, a mode, a density and the project's
// accent colour) and gives the same names everywhere: the CSS of BubbleUI and go-ui's controls read
// them alike.
package bubble

import "time"

// Style is the visual language: paper is the default, material the alternative that Android users
// expect. Both give the same tokens, so an interface switches style without changing its controls.
type Style string

// The two styles.
const (
	StylePaper    Style = "paper"
	StyleMaterial Style = "material"
)

// Mode is light or dark; interfaces follow the one the system asks for.
type Mode string

// The two modes.
const (
	ModeLight Mode = "light"
	ModeDark  Mode = "dark"
)

// Density is how much fits on screen: comfortable for web pages and phones, compact for desktop
// tools used with a pointer.
type Density string

// The two densities.
const (
	DensityComfortable Density = "comfortable"
	DensityCompact     Density = "compact"
)

// DefaultAccent is the accent of a project that has not chosen one: Bubble's own blue.
var DefaultAccent = hex("#2C5C9E")

// Options are the choices a theme is made from.
type Options struct {
	Style   Style
	Mode    Mode
	Density Density
	Accent  Color // the main colour of the project's logo; the theme adjusts it until it reads well
	Glass   bool  // frosted, translucent panels over the background (paper only)
	// Palettes, when a phone gave them (Go for Android reads its wallpaper's), make the material
	// style's colours instead of the accent, as the system's own applications are coloured.
	Palettes *Palettes
}

// DefaultOptions is paper, light, comfortable, with the default accent and no glass.
func DefaultOptions() Options {
	return Options{Style: StylePaper, Mode: ModeLight, Density: DensityComfortable, Accent: DefaultAccent}
}

// Theme holds every token of one combination of options.
type Theme struct {
	Options Options
	Colors  Colors
	Shape   Shape
	Type    Type
	Space   Space
	Motion  Motion
	Glass   Glass
}

// Colors are the theme's colours, named by what they are for.
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
	AccentEdge  Color // the logo's darker side, ×0.70: the edge of a checked box and of a raised button
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

	// The ladder of containers, from the lowest surface to the highest: where a card, a navigation
	// bar, a field and a switch's track sit. Material makes it from its scheme; paper and glass
	// from their own surface and sunken colours.
	ContainerLow, Container, ContainerHigh, ContainerHighest Color

	OnAccentTint      Color // text and icons on the accent tint
	AccentContainer   Color // a stronger wash of the accent: a chosen chip, a floating button
	OnAccentContainer Color
	OnDanger          Color // text on a danger fill
	DangerContainer   Color // a wash of danger: an error banner
	OnDangerContainer Color
	Inverse           Color // what stands out against the page: a snackbar
	OnInverse         Color
	InverseAccent     Color // the accent on the inverse surface: a snackbar's action
}

// Shape holds corners, edges and elevation, in layout pixels.
type Shape struct {
	ControlRadius  float64 // fields and selects
	ButtonRadius   float64 // buttons and chips
	CardRadius     float64
	PanelRadius    float64 // windows, sheets, dialogs
	FloatingRadius float64 // the floating action button
	LineWidth      float64
	EdgeWidth      float64 // a block edge under the main button; 0 in every style today, which draws it flat
	FocusWidth     float64
	FocusOffset    float64 // the gap between a control and its focus ring; 0, so the edge lights up, kept for go-ui
	Elevated       bool    // whether raised things cast a shadow (material) or sit on a line (paper)
}

// Type holds the typefaces and the type scale, in layout pixels.
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

// Space holds distances and the size of controls, in layout pixels.
type Space struct {
	Gap           float64 // between controls
	Padding       float64 // inside cards and panels
	ControlHeight float64
	TouchTarget   float64 // the smallest area a finger or pointer must be able to hit
}

// Motion holds how long transitions take. Interfaces skip them when the system asks for less motion.
type Motion struct {
	Fast   time.Duration
	Medium time.Duration
	Slow   time.Duration
}

// Glass holds the crystalline panels of the glass variant; Enabled is false outside it.
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

// NewTheme makes the theme of the given options.
func NewTheme(options Options) (Theme, error) {
	if err := checkOptions(options); err != nil {
		return Theme{}, err
	}
	base := baseColors(options.Style, options.Mode)
	var glass Glass
	if options.Glass && options.Style == StylePaper {
		glass = glassOf(options.Mode)
		base = glassColors(base, glass, options.Mode)
	}
	theme := Theme{
		Options: options,
		Colors:  colorsOf(base, options),
		Shape:   shapeOf(options.Style, options.Density),
		Type:    typeOf(options.Style, options.Density),
		Space:   spaceOf(options.Style, options.Density),
		Motion:  Motion{Fast: 250 * time.Millisecond, Medium: 500 * time.Millisecond, Slow: 750 * time.Millisecond},
	}
	if glass.Enabled {
		theme.Glass = glass
		theme.Colors = soberAccent(theme.Colors, glass, options.Mode)
		theme.Colors = ladderOf(theme.Colors, options.Mode)
	}
	return theme, nil
}

// colorsOf adds the accent to a style's neutral colours: Material makes every colour from it, as
// its own tools do; paper moves the accent until it reads, and builds its ladder from its surfaces.
func colorsOf(base Colors, options Options) Colors {
	if options.Style == StyleMaterial {
		palettes := PalettesFromColor(options.Accent)
		if options.Palettes != nil {
			palettes = *options.Palettes
		}
		return materialColors(base, rolesFromPalettes(palettes, options.Mode))
	}
	return ladderOf(withAccent(base, options), options.Mode)
}

func checkOptions(options Options) error {
	if options.Style != StylePaper && options.Style != StyleMaterial {
		return UnknownStyle
	}
	if options.Mode != ModeLight && options.Mode != ModeDark {
		return UnknownMode
	}
	if options.Density != DensityComfortable && options.Density != DensityCompact {
		return UnknownDensity
	}
	if options.Palettes != nil && (options.Style != StyleMaterial || !options.Palettes.opaque()) {
		return InvalidPalettes
	}
	return nil
}
