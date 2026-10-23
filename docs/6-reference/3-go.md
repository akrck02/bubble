# 🐹 Go

The Go API of Bubble's two modules, as `go doc -all` prints it. This page is made by `docs/6-reference/generate.py`; do not edit it by hand.

| You want to… | Use |
|---|---|
| make a theme | `bubble.NewTheme(bubble.Options{…})`, or start from `bubble.DefaultOptions()` |
| read a colour | `bubble.ParseColor("#4C9A5B")` |
| write the CSS tokens | `bubble.CSS(writer, accent)`, or the `bubble-css` command |
| check two colours | `bubble.Contrast(a, b)` |
| dress Go UI | `goui.Theme(theme)`, and `goui.Ground` for glass |
| write in Material's typeface | `goui.FontsFor(style)`, Roboto for Material |
| draw Material's parts | `material.New(context, theme)` in `goui/material` |
| colour Material from a phone's wallpaper | `Options.Palettes`, or `bubble.PalettesFromColor(seed)` |

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

var InvalidPalettes = define("bubble.invalid-palettes", "palettes colour only the material style, and every shade must be opaque")
    InvalidPalettes means a phone's palettes were given outside the material
    style, which is the only one they colour, or have a shade that is not an
    opaque colour.

var PaletteShades = [13]int{0, 10, 50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000}
    PaletteShades are the shades of each palette in Android's names: 0 is white,
    500 the middle, 1000 black. Shade s has Material's tone 100 - s/10,
    a lightness on CIELAB's scale (L*).

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
	// Palettes, when a phone gave them (Go for Android reads its wallpaper's), make the material
	// style's colours instead of the accent, as the system's own applications are coloured.
	Palettes *Palettes
}
    Options are the choices a theme is made from.

func DefaultOptions() Options
    DefaultOptions is paper, light, comfortable, with the default accent and no
    glass.

type Palettes struct {
	Accent1, Accent2, Accent3, Neutral1, Neutral2 [13]Color
}
    Palettes are the five tonal palettes Material's colours come from,
    each at PaletteShades: the accent's own, a quieter one for containers, a
    third of a turned hue, and two near-greys of the accent's hue for surfaces.
    Go for Android reads them from the phone's wallpaper.

func PalettesFromColor(seed Color) Palettes
    PalettesFromColor makes the five palettes Android makes from a wallpaper,
    from one colour.

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

func FontsFor(style bubble.Style) overlay.Fonts
    FontsFor gives the typeface a theme is written in: Roboto for the material
    style, go-renderer's own otherwise.

func Ground(theme bubble.Theme, width, height int) *image.RGBA
    Ground paints glass's ground the way the browser does, layer by layer:
    the background, the second light over it, the first light on top. Panels are
    drawn over this picture; it is smooth, so blurring what is behind them would
    change nothing.

func Roboto() overlay.Fonts
    Roboto is Android's own typeface (SIL Open Font License, in fonts/),
    the material style's: give it to the window (window.Options.Fonts) and to
    the context (ui.NewInFonts), which measure and draw text in it.

func Theme(theme bubble.Theme) ui.Theme
    Theme turns a Bubble theme into go-ui's. The colours carry the same names,
    mostly; the sizes come from Bubble's density. What is chosen is brighter
    and heavier, never coloured, and the main button stands on its darker edge.
    On glass, panels and wells are the translucent panes, and the window's
    background is glass's ground (paint it with Ground).
```

## Package goui/material

`github.com/akrck02/bubble/goui/material`: Material 3's parts for Go UI.

```
package material // import "github.com/akrck02/bubble/goui/material"

Package material draws Material 3's pieces with go-ui, in a Bubble theme (with
Material 3 Expressive's newer pieces). They came from Go UI for Android.
The pieces are made of go-ui boxes, texts and clickables, as any program can
make controls of its own, so they draw the same on a phone, on a desktop and in
a browser. They read the theme through Roles, Material's names for its colours,
so they draw in every Bubble style: Material's own, paper and glass. A program
makes a Look for each frame and calls its pieces:

    theme, _ := bubble.NewTheme(bubble.Options{Style: bubble.StyleMaterial, Mode: bubble.ModeLight,
    	Density: bubble.DensityComfortable, Accent: accent})
    context.Theme = goui.Theme(theme)
    look := material.New(context, theme)
    look.TopAppBar(material.TopBar{Title: "Inbox"})
    if look.FilledButton("plus", "New message") { … }

On a phone, Go for Android reads the wallpaper's palettes, which go in
bubble.Options.Palettes. The files follow Material's catalogue:

    look.go         the Look, text and paragraphs
    roles.go        a Bubble theme's colours by Material's names
    navigation.go   top bars, the navigation bar, rail and drawer, tabs and swiping pages, the back
    actions.go      buttons, floating buttons and their menu, split buttons, button groups, toolbars
    selection.go    switches, check boxes, radio buttons, segmented buttons, sliders
    chips.go        filter, assist and input chips, and tags
    inputs.go       text fields and the search bar
    containment.go  cards, list items, avatars, banners, section titles, dividers
    overlays.go     dialogs, bottom sheets, snackbars, menus
    pickers.go      date and time pickers
    progress.go     loading indicators, circular, linear and wavy progress, skeletons
    gestures.go     pull to refresh, rows swiped away, carousels
    stick.go        an on-screen stick

CONSTANTS

const Indeterminate = -1
    Progress is how an indicator of known length shows its work: a fraction from
    0 to 1, or Indeterminate for work of unknown length.

const NavigationBarHeight = 80
    NavigationBarHeight is a navigation bar's height without the gesture bar
    under it, for a snackbar or a floating button to keep above.


VARIABLES

var EnglishDates = DateNames{Title: "Select date",
	Months:   [12]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
	Weekdays: [7]string{"M", "T", "W", "T", "F", "S", "S"}, Confirm: "OK", Dismiss: "Cancel"}
    EnglishDates are the words a date picker shows by default.

var EnglishTimes = TimeNames{Title: "Select time", Confirm: "OK", Dismiss: "Cancel"}
    EnglishTimes are the words a time picker shows by default.

var EnglishWords = Words{Back: "Back", Close: "Close", Clear: "Clear", ShowPassword: "Show password",
	HidePassword: "Hide password", PreviousMonth: "Previous month", NextMonth: "Next month",
	MoreOptions: "More options", Cut: "Cut", Copy: "Copy", Paste: "Paste", SelectAll: "Select all"}
    EnglishWords are the words a screen reader reads by default.


TYPES

type Answer int
    Answer is what the person did with a dialog in a frame.

const (
	Open      Answer = iota // nothing yet: keep showing it
	Confirmed               // the main action
	Dismissed               // the other action, a touch outside, Escape
)
type Carousel struct {
	// Has unexported fields.
}
    Carousel is what a carousel remembers: where it rests, as the index of the
    first item shown.

type DateNames struct {
	Title    string     // over the date chosen
	Months   [12]string // January first
	Weekdays [7]string  // one or two letters each, Monday first
	Confirm  string     // the action that keeps the date
	Dismiss  string     // the one that leaves it
}
    DateNames are the words a date picker shows, for a program in another
    language than English.

type DatePicker struct {
	Names DateNames // the zero value shows EnglishDates

	// Has unexported fields.
}
    DatePicker is what a date picker remembers while shown: the month on show
    and the day chosen.

type Destination struct {
	Icon, Label string
	// Badge marks the icon: a count ("3"), or a dot with "•"; none when empty.
	Badge string
}
    Destination is a place the navigation bar leads to: an icon and a short
    label.

type Dialog struct {
	Icon    string // over the title, centred; none when empty
	Title   string
	Text    string // under the title, broken into lines
	Content func() // more under the text (fields, choices); none when nil
	Confirm string // the main action, at the end
	Dismiss string // the other action, before it; none when empty
}
    Dialog is what a dialog says and offers.

type Dropdown struct {
	Selected int
	// Has unexported fields.
}
    Dropdown is what an exposed dropdown keeps: the choice, and whether its
    choices show.

type Field struct {
	Label       string // above the box
	Placeholder string // in the box, quiet, while it is empty
	Supporting  string // under the box: a hint
	Error       string // under the box instead of the hint, in red, with the box red too
	MaxLength   int    // the most characters it takes; 0 = no limit
	// Counter shows how many characters there are of MaxLength, under the box at its end.
	Counter bool
	// Filled draws Material's filled field, a tinted box with a line under it, rather than the
	// outlined one: for a form on a plain screen, where boxes drawn by their border look empty.
	Filled bool
	// Prefix and Suffix are written in the box before and after what is typed: a currency, a unit.
	Prefix, Suffix string
	// Secret shows a dot for each character (a password). With Reveal, an eye at the box's end
	// shows the characters while *Reveal is true, and touching it turns *Reveal over.
	Secret bool
	Reveal *bool
}
    Field is what a text field says around what is typed.

type FloatingMenu struct {
	Icon    string        // the button's icon while closed
	Label   string        // what a screen reader reads for the button while closed ("Create")
	Actions []Destination // what it opens into: an icon and a label each
	// Above is how far over the bottom of the window it sits: a navigation bar's height
	// (NavigationBarHeight) when there is one.
	Above float32
	Open  bool
}
    FloatingMenu is a floating action button that opens into a few related
    actions stacked above it, its icon turning into a cross. A program keeps
    one, and draws it with Look.FloatingMenu.

type Item struct {
	Icon string       // in a round avatar at the start; none when empty
	Mark draw2d.Color // the avatar's colour; the primary container when zero
	// Picture is a photo in the round avatar instead of the icon (a person, a place).
	Picture    *draw2d.Image
	Headline   string
	Supporting string // the line under the headline; none when empty
	// Trailing draws what goes at the end (a Switch, an IconButton, a Tag); none when nil.
	Trailing func()
	// Selected marks the row chosen, in a list whose rows are chosen with a long press: it is tinted,
	// and a tick takes the avatar's place.
	Selected bool
}
    Item is what a list item shows.

type Look struct {
	Context *ui.Context
	Theme   bubble.Theme // the Bubble theme the pieces are drawn in
	Roles   Roles        // its colours by Material's names, as the pieces read them
	Words   Words        // the pieces' own words, in the person's language; the zero value is EnglishWords
}
    Look draws the pieces in one frame, in one Bubble theme. It is small:
    make one in every frame.

    Pieces that can be touched return true in the frame they are touched. Each
    takes the name go-ui tells it apart by: its label or icon, or a name given.
    Two pieces with the same name in one scope are one to go-ui: give repeated
    ones a scope each (ui.Context.Push).

func New(context *ui.Context, theme bubble.Theme) Look
    New makes the frame's Look.

func (look Look) AssistChip(icon, label string) bool
    AssistChip suggests an action related to what is on screen (add to the
    calendar, share), with an icon before its label. It returns true when
    touched.

func (look Look) Avatar(name, icon string, tint draw2d.Color, size float32)
    Avatar is a person or a thing shown in a circle: their initials (up to two
    letters of name) on a colour, or an icon when name is empty. size is its
    side, 40 in lists.

func (look Look) Back() bool
    Back answers the system's back (a phone's back gesture) for the screen:
    true when there was one, and the program then goes back (to the first
    destination, to the list a detail came from). Dialogs and sheets answer it
    first; a back nobody answers takes the program to the background.

func (look Look) Banner(icon, text string, actions ...string) int
    Banner is a message across the top of the content that stays until acted on
    (no connection, an update to install), with an icon and one or two actions
    under its text. It returns the index of the action touched in this frame,
    or -1.

func (look Look) Body(besideRail bool, content func())
    Body holds a screen's content between the bars: it grows to fill them
    and keeps its content clear of what the system draws at the sides (a
    camera's cutout, a navigation bar on the side). Beside a navigation rail
    (besideRail), the rail keeps the left side clear, and with no navigation bar
    under it, Body keeps clear of the gesture bar itself.

func (look Look) BottomSheet(name string, content func()) (dismissed bool)
    BottomSheet shows content in a sheet rising from the bottom of the screen,
    which is dimmed meanwhile: choices for what was touched, a short form.
    It returns true when the person touches outside it (or presses Escape);
    the program then stops showing it. name tells sheets apart.

func (look Look) ButtonGroup(name string, labels []string, selected *int) bool
    ButtonGroup is a few choices joined side by side, the chosen one filled and
    fully rounded, the others squarer: Material 3 Expressive's connected group
    (a view, a size, an alignment). It changes selected and returns true when
    another is touched.

func (look Look) Card(content func())
    Card is a filled container for one subject: a few lines, a picture,
    its actions.

func (look Look) Carousel(name string, carousel *Carousel, count int, itemWidth float32, item func(index int))
    Carousel is a row of items wider than the screen, slid sideways with a
    finger: let go, it glides on with the finger's speed and settles at an
    item's edge. Each item is itemWidth wide, with 8 pixels between them;
    item draws one, given its index. name tells carousels apart.

func (look Look) CarouselItem(height float32, tint draw2d.Color, icon, title string)
    CarouselItem is a carousel's usual item: a rounded tile in a colour,
    with a title over its foot.

func (look Look) CarouselPicture(picture *draw2d.Image, height float32, title string)
    CarouselPicture is a carousel item made of a photo, with a title over its
    foot on a dark veil that keeps the title readable on any photo; an empty
    title leaves the photo alone (a narrow item).

func (look Look) Checkbox(label string, on *bool) bool
    Checkbox is a square that is ticked or not, with a label after it; touching
    either flips it. It returns true when touched.

func (look Look) CircularProgress(fraction float32, size float32, wavy bool)
    CircularProgress is a ring that fills as work gets done (fraction from 0 to
    1), or, with Indeterminate, an arc chasing round. wavy draws the filled part
    as a wave, as Material 3 Expressive does. size is its side, 48 by default in
    Material.

func (look Look) ContainedLoadingIndicator(size float32)
    ContainedLoadingIndicator is the indicator on a round container, to stand
    out over content (at the top of a list being refreshed).

func (look Look) DatePicker(name string, picker *DatePicker, date *time.Time) Answer
    DatePicker asks for a day in a dialog: a month at a time, with arrows to the
    months around it. It starts at the date given; on Confirmed it writes the
    day chosen there. The program shows it with an if while it wants an answer,
    as Dialog.

func (look Look) Dialog(name string, dialog Dialog) Answer
    Dialog asks something over the screen, which is dimmed and does not react
    meanwhile: a card in the middle with a title, a text and one or two actions.
    The program shows it with an if while it wants an answer, and stops when the
    answer is not Open. name tells dialogs apart.

func (look Look) Divider()
    Divider is a thin line across, between groups.

func (look Look) Dropdown(name, label string, choices []string, dropdown *Dropdown) bool
    Dropdown is Material's exposed dropdown menu: a box like a text field's
    showing the choice, which opens the choices under it when touched.
    It returns true when another one is chosen. name tells dropdowns apart.

func (look Look) ElevatedCard(content func())
    ElevatedCard is a card raised over what is beneath it: over a map or a
    scene.

func (look Look) Enter(key string, content func())
    Enter shows a screen rising gently into place when it is first shown,
    as Material's screens arrive: give each screen its own key.

func (look Look) ExpandableItem(name string, item Item, open *bool, more func()) bool
    ExpandableItem is a list row that opens to show more under it (an answer
    under its question, the details of an order), its chevron turning. Touching
    the row opens or closes it; it returns true when touched. name tells rows
    apart.

func (look Look) ExtendedFloatingButton(icon, label string) bool
    ExtendedFloatingButton is the floating action button with its action written
    beside the icon: clearer where the icon alone would not say it.

func (look Look) FilledButton(icon, label string) bool
    FilledButton is the screen's main action, in the primary colour.

func (look Look) FilterChip(label string, selected bool) bool
    FilterChip is a choice that is on or off, among others in a row: outlined,
    or filled with a tick when selected. It returns true when touched; the
    program flips its own value.

func (look Look) FloatingButton(icon string) bool
    FloatingButton is the floating action button: 56 pixels, rounded square
    corners, raised over what scrolls beneath it. Place it at the end of the
    screen, above the navigation bar.

func (look Look) FloatingMenu(name string, menu *FloatingMenu) int
    FloatingMenu draws the button at the bottom end of the screen, over
    everything, and its actions while open. It returns the index of the action
    touched in this frame, or -1; touching one, the cross or the system's back
    closes it.

func (look Look) FloatingToolbar(name string, tools []Destination, above float32) int
    FloatingToolbar is a row of actions in a pill floating over the content at
    the bottom of the screen, for what is done to the content on show (edit,
    share, delete): each tool shows its icon, and a screen reader reads its
    label. above is how far over the bottom it sits. It returns the index of the
    tool touched in this frame, or -1.

func (look Look) GoBack(navigation *Navigation) bool
    GoBack answers the system's back by going back a screen, while there is
    one to go back to. Ask for it after drawing what goes over the screens,
    so that answers first.

func (look Look) Heading(text string)
    Heading is a section's small title in the primary colour, as SectionTitle,
    for content that keeps its own margins (a page of cards and buttons) rather
    than a list across the screen.

func (look Look) IconButton(icon string) bool
    IconButton is a round 48-pixel target with a 24-pixel icon, in the bars and
    beside list items.

func (look Look) InputChip(label string) (touched, removed bool)
    InputChip is something the person entered or chose (a recipient, a tag),
    with a cross to take it out. It tells whether the chip was touched and
    whether its cross was.

func (look Look) LinearProgress(fraction float32)
    LinearProgress is a bar across that fills as work gets done: fraction from 0
    to 1. A negative fraction is work of unknown length: a piece of the bar runs
    along it, and it asks for frames while shown.

func (look Look) ListItem(name string, item Item) bool
    ListItem is a row of a list, 72 pixels tall with a supporting line and 56
    without. It returns true when the row is touched. name tells it apart from
    the other rows.

func (look Look) LoadingIndicator(size float32)
    LoadingIndicator is Material's indicator for a short wait of unknown length:
    a filled shape in the primary colour, turning and morphing from one form to
    the next. size is its side, 48 by default in Material.

func (look Look) LongPressed() bool
    LongPressed tells whether the piece drawn just before it was long-pressed
    in this frame: a finger held still on it for half a second, or the mouse's
    right button. The phone gives its buzz then, as Android's own applications
    do. A list uses it to start choosing rows (Item.Selected); a touch that
    became a long press is no touch when it lifts.

func (look Look) Menu(icon string, choices []string) int
    Menu is a button with an icon (three dots, usually) that opens a list of
    choices under it; a screen reader reads the button as "More options".
    It returns the index of the one chosen in this frame, or -1.

func (look Look) MultiBrowseCarousel(name string, carousel *Carousel, count int, largeWidth, height float32, item func(index int, width float32))
    MultiBrowseCarousel is Material's carousel for browsing many items at once:
    large items, then a medium one and a small one at the edge where the row
    runs off. Items shrink as they reach either edge and grow as they come in,
    so the row always hints at more. largeWidth is the width items have in the
    middle; item draws one at the width it has in this frame, and is clipped to
    it.

func (look Look) NamedIconButton(icon, name string) bool
    NamedIconButton is an IconButton a screen reader (TalkBack) reads by a name,
    what it does ("Share"), rather than by its icon's.

func (look Look) NavigationBar(destinations []Destination, selected *int) bool
    NavigationBar is the bar at the bottom, painted down under the gesture bar:
    three to five destinations, the selected one with its icon on a pill.
    It changes selected and returns true when another destination is touched.

func (look Look) NavigationDrawer(name, title string, destinations []Destination, selected *int) (closed bool)
    NavigationDrawer slides in from the left over the dimmed screen, with a
    title and the program's destinations, for more places than a navigation
    bar holds. It changes selected when one is touched, and returns true when
    it should close: a destination chosen, a touch outside, the system's back.
    name tells drawers apart.

func (look Look) NavigationRail(destinations []Destination, selected *int) bool
    NavigationRail is the navigation bar's form on a wide screen (a phone on
    its side, a tablet): the destinations down the left edge, painted under the
    status bar and beside a camera's cutout. It changes selected and returns
    true when another destination is touched. Wide tells when to use it.

func (look Look) OutlinedButton(icon, label string) bool
    OutlinedButton is a second action beside the main one.

func (look Look) OutlinedCard(content func())
    OutlinedCard is a card drawn by its border, among others on the surface.

func (look Look) Paragraph(text string, size float32, weight draw2d.Weight, colour draw2d.Color)
    Paragraph writes text that breaks into lines at the width it has, word by
    word: a card's body, a description. Text writes one line, however long.

func (look Look) PictureAvatar(picture *draw2d.Image, size float32)
    PictureAvatar is a person or a place shown by a photo in a circle. size is
    its side, 40 in lists.

func (look Look) PullToRefresh(name string, refreshing bool, layout ui.Layout, content func()) (refresh bool)
    PullToRefresh is a scrolling list that refreshes when pulled down past its
    top and let go: a loading indicator follows the pull, and stays over the
    list while refreshing is true. It returns true in the frame the person asks
    for a refresh; the program then refreshes and sets refreshing until it is
    done. name tells it apart; content adds the list's rows.

func (look Look) RadioButton(label string, chosen bool) bool
    RadioButton is one of a few choices of which one is chosen, with a label
    after it. It returns true when touched; the program then chooses it.

func (look Look) RoundButton(icon string, size float32, primary bool) bool
    RoundButton is a round floating button of any size, over a map or a scene
    (zoom in, zoom out, change the view). primary fills it with the main colour;
    otherwise it is a raised surface.

func (look Look) Screens(name string, navigation *Navigation, screen func(key string))
    Screens draws the screen on top of a navigation: one gone to slides in
    from the end over the one it was gone to from, and slides back out when the
    program goes back. screen draws a screen by its name. The program answers
    the system's back with GoBack, after what goes over the screens.

func (look Look) SearchBar(name string, text *string, placeholder string) bool
    SearchBar is a pill to search in, at the top of a screen. It returns true
    when the text changed.

func (look Look) SearchView(name, placeholder string, search *Search, content func()) bool
    SearchView is Material's search: a pill at the top of a screen which,
    touched, opens over the whole screen ready to type in, with a way back and
    a cross to empty it; content shows under it what the program finds (the
    last searches, suggestions, results), and closes the search by setting
    search.Active to false when one is chosen. The system's back closes it too.
    It returns true when the query changed. name tells searches apart.

func (look Look) SectionTitle(text string)
    SectionTitle is a small heading over a group of rows, in the primary colour.

func (look Look) SegmentedButtons(name string, options []string, selected *int) bool
    SegmentedButtons are two to five choices side by side sharing the row's
    width, one of them chosen (a view, a speed, a period). It changes selected
    and returns true when another is touched. name tells the group apart from
    another one on the same screen.

func (look Look) SideSheet(name, title string, content func()) (closed bool)
    SideSheet slides in from the end of the screen over the dimmed screen,
    headed by a title and a cross: filters for a list, the details of what is
    chosen, room for a form on a wide screen. It returns true when it should
    close: the cross touched, a touch outside, the system's back. name tells
    side sheets apart.

func (look Look) Skeleton(width, height ui.Sizing, radius float32)
    Skeleton stands in for content still coming (a picture, a line of text):
    a quiet block with a light passing across it. A row of them in the content's
    shape tells what will be there.

func (look Look) Slider(name string, value *float32, minimum, maximum, step float32) bool
    Slider chooses a value between two ends by dragging, in steps (0 for
    any value): go-ui's slider in the theme's colours and a finger's size.
    It returns true when the value changed.

func (look Look) Snackbar(snackbar *Snackbar) (actionTouched bool)
    Snackbar draws the note while it is up, over everything, and asks for a
    frame when it is due to go. It returns true when its action is touched,
    and the note goes then.

func (look Look) SplitButton(icon, label string) (action, more bool)
    SplitButton is a main action joined to a small button that opens its other
    choices (send, and schedule): it tells which half was touched.

func (look Look) StandardSheet(name string, sheet *Sheet, peek float32, content func())
    StandardSheet is a bottom sheet that stays over the screen without dimming
    it, which goes on working above it (a map under a list of its places):
    it rests showing peek pixels, or expanded to most of the screen, and a
    finger drags it by its handle between the two; let go, it settles where it
    was heading. Its content scrolls inside it. name tells sheets apart.

func (look Look) StatusBarScrim()
    StatusBarScrim paints the strip under the status bar alone, for a screen
    whose top is a scene rather than a bar: the system's icons stay readable
    over it.

func (look Look) Stick(name string, stick *Stick, size float32)
    Stick draws the stick, its knob where the finger points it, and follows the
    fingers of this frame (ui.Context.Touches). size is its side, 128 or so;
    name tells it apart. It may sit anywhere: in a Push scope, a scrolling list
    or swiping pages, whose finger it takes while it steers. It asks for frames
    while held, so what it drives keeps moving under a finger that stays still.

func (look Look) SwipeToDismiss(name string, content func()) (dismissed bool)
    SwipeToDismiss is a row a finger slides sideways to take it away,
    a bin showing in the gap it leaves: pushed past two fifths of its width,
    or flicked, it goes, and SwipeToDismiss returns true for the program to
    remove it; let go short of that, it slides back. name tells rows apart.

func (look Look) Switch(name string, on *bool) bool
    Switch is Material's switch: a 52 × 32 track whose handle grows and moves to
    the end when on. It flips the value and returns true when touched.

func (look Look) TabPager(name string, labels []string, pager *Pager, page func(index int))
    TabPager is tabs over pages a finger swipes from side to side: the underline
    follows the pages as they move, a page let go past a third of the way (or
    flicked) turns to the next, and touching a tab slides to its page. page
    draws one page, given its index; lists inside it still scroll up and down.
    name tells pagers apart.

func (look Look) Tabs(name string, labels []string, selected *int) bool
    Tabs are the primary tabs under a top bar: labels sharing the width,
    an underline under the chosen one that slides to the next one chosen.
    It changes selected and returns true when another is touched.

func (look Look) Tag(text string, ink, tint draw2d.Color)
    Tag is a label that cannot be touched, in a pill of its own colours:
    a state ("3 running"), a count, a kind. ink is the text's colour and tint
    the pill's.

func (look Look) Text(text string, size float32, weight draw2d.Weight, colour draw2d.Color)
    Text writes a line in a size, a weight and a colour.

func (look Look) TextButton(icon, label string) bool
    TextButton is the quietest action: a dialog's "Cancel", a card's "More".

func (look Look) TextField(name string, text *string, field Field) bool
    TextField is Material's text field: a label, the box typed in (outlined,
    or filled), and a hint or an error and a count of characters under it.
    On a phone the on-screen keyboard comes up while it is typed in. It returns
    true when the text changed. name tells fields apart.

func (look Look) TimePicker(name string, picker *TimePicker, date *time.Time) Answer
    TimePicker asks for a time of day in a dialog, on Material's dial: the
    hour first, on a ring of 1 to 12 with 13 to 00 inside it, then the minute,
    by touching or dragging round the dial. The hour and the minute at the top
    choose which the dial shows. It starts at the time given; on Confirmed it
    writes the time chosen there, keeping the day.

func (look Look) TonalButton(icon, label string) bool
    TonalButton is a main action that should not stand out as much, on the
    secondary container.

func (look Look) TopAppBar(bar TopBar) (back bool)
    TopAppBar is the bar at the top, painted up under the status bar,
    so the screen reaches the top edge as Android's own apps do. Inside Body,
    it keeps clear of a cutout at the sides too. It returns true when the back
    arrow is touched.

func (look Look) VideoPlayer(name string, video *window.Video, height float32)
    VideoPlayer shows a film (go-ui's window.Video) in a rounded box: whole,
    over black where their shapes differ, with a round play button over it
    while it is paused and a thin line at its foot of how far it has played.
    Touching it plays or pauses the film; it brings the film up to date itself
    (Video.Update). name tells players apart.

func (look Look) WavyProgress(fraction float32)
    WavyProgress is a line that fills as work gets done (fraction from 0 to 1),
    the filled part a travelling wave and the rest a flat track ending in a
    stop dot, as Material 3 Expressive draws it; with Indeterminate, a wave runs
    along it.

func (look Look) Wide() bool
    Wide tells whether the screen is wide enough for a navigation rail instead
    of a bar: wider than tall, or 600 pixels across (Material's medium width,
    a small tablet).

type Navigation struct {
	// Has unexported fields.
}
    Navigation is a program's screens, one gone to from another: a list,
    then one of its rows' details, then a form about it. The program keeps one
    and names its screens as it likes; Screens draws the one on top. The zero
    value starts at "".

func (navigation *Navigation) Back() bool
    Back goes back to the screen before, and tells whether there was one.

func (navigation *Navigation) Current() string
    Current is the screen on top.

func (navigation *Navigation) Depth() int
    Depth is how many screens there are, the first one included.

func (navigation *Navigation) Go(screen string)
    Go goes to a screen, over the one shown.

func (navigation *Navigation) Reset(screen string)
    Reset starts again at one screen, with nothing to go back to: a destination
    of the navigation bar.

type Pager struct {
	Selected int
	// Has unexported fields.
}
    Pager is what tabs with swiping pages remember: the page chosen, and where
    the pages were shown in the last frame (between two while one is dragged).

type Roles struct {
	Primary, OnPrimary, PrimaryContainer, OnPrimaryContainer draw2d.Color
	SecondaryContainer, OnSecondaryContainer                 draw2d.Color
	Error, OnError, ErrorContainer, OnErrorContainer         draw2d.Color
	// Surface is the screen's background; the containers are the surfaces raised over it, from the
	// lowest to the highest (a card, a navigation bar, a switch's track).
	Surface, OnSurface, OnSurfaceVariant                                                 draw2d.Color
	SurfaceContainerLow, SurfaceContainer, SurfaceContainerHigh, SurfaceContainerHighest draw2d.Color
	Outline, OutlineVariant                                                              draw2d.Color
	// The inverse roles paint what stands out against the screen: a snackbar.
	InverseSurface, InverseOnSurface, InversePrimary draw2d.Color
	// BarScrim is the surface's colour, see-through: under a bar painted over a scene. It is not
	// the theme's Scrim, the darkening behind a dialog.
	BarScrim draw2d.Color
}
    Roles are a Bubble theme's colours by Material 3's names, the names the
    pieces were written with. In the material style they are Material's own;
    in paper and glass, the colours doing the same job, so every piece draws in
    every style.

func RolesOf(theme bubble.Theme) Roles
    RolesOf names a Bubble theme's colours as Material does.

type Search struct {
	Query  string
	Active bool
}
    Search is what a search view keeps: what is typed, and whether it is open
    over the screen.

type Sheet struct {
	Expanded bool
	// Above is how far over the bottom of the window it sits: a navigation bar's height
	// (NavigationBarHeight) when there is one.
	Above float32
	// Has unexported fields.
}
    Sheet is what a standard bottom sheet keeps: whether it is expanded, and how
    far a finger has dragged it from where it rests.

type Snackbar struct {
	Text   string
	Action string // a button at its end; none when empty
	// Above is how far over the bottom of the window it sits: a navigation bar's height
	// (NavigationBarHeight) when there is one.
	Above float32
	// Has unexported fields.
}
    Snackbar is a short note at the bottom of the screen for a few seconds, with
    an action that undoes or follows what was just done. A program keeps one,
    calls Show when something happens and draws it with Look.Snackbar in every
    frame.

func (snackbar *Snackbar) Hide()
    Hide takes the note away before its time.

func (snackbar *Snackbar) Show(text, action string, duration time.Duration)
    Show puts a note up for a while (Material's are 4 to 10 seconds).

func (snackbar *Snackbar) Shown() bool
    Shown tells whether the note is up.

type Stick struct {
	Forward  float32 // from -1 (back) to 1 (ahead): the finger above the middle is ahead
	Sideways float32 // from -1 (left) to 1 (right)
	Held     bool
	Finger   int64 // the finger holding it, while Held

	// Has unexported fields.
}
    Stick is an on-screen stick, for walking or driving something with a thumb:
    the finger that comes down on it holds it, and where that finger is from
    its middle says where to go. Other fingers stay free (to look around,
    to press buttons). A program keeps one Stick and reads it in every step.

func (stick *Stick) Follow(touches []input.Touch)
    Follow gives the stick to the finger that came down on it, and points it at
    that finger; a finger lifted lets it go back to the middle. Look.Stick calls
    it; a program drawing a stick of its own calls it with the touches and where
    its stick is (Place).

func (stick *Stick) Place(area draw2d.Rectangle, size float32)
    Place says where a stick of a program's own drawing is, and its side,
    for Follow.

func (stick *Stick) Touching(touch input.Touch) bool
    Touching tells whether a finger is the one holding the stick: a program that
    turns a view with the other fingers leaves this one out.

type TimeNames struct {
	Title   string
	Confirm string
	Dismiss string
}
    TimeNames are the words a time picker shows, for a program in another
    language than English.

type TimePicker struct {
	Names TimeNames // the zero value shows EnglishTimes

	// Has unexported fields.
}
    TimePicker is what a time picker remembers while shown: the time chosen,
    not yet kept, and whether the dial chooses the hour or the minute.

type TopBar struct {
	Title    string
	Subtitle string // under the title, smaller; none when empty
	Back     bool   // an arrow before the title, to go back
	// Close makes that arrow a cross, for leaving a mode (rows chosen) rather than a screen.
	Close bool
	// Actions draws what goes at the bar's end, usually IconButtons; none when nil.
	Actions func()
	// OverScene paints the bar see-through (Roles.BarScrim), over a scene behind the interface.
	OverScene bool
	// Large writes the title big under the arrow and the actions, for a screen's first level.
	Large bool
	// Scrolled is how far the content under the bar has scrolled up (ui.Context.ScrollOffset): a
	// large bar shrinks into a small one as it goes, and any bar takes a container's colour once
	// content is under it.
	Scrolled float32
}
    TopBar is what the top app bar shows.

type Words struct {
	Back, Close, Clear          string
	ShowPassword, HidePassword  string
	PreviousMonth, NextMonth    string
	MoreOptions                 string // a menu's button, a split button's arrow
	Cut, Copy, Paste, SelectAll string // the bar over a text field's selection
}
    Words are what the pieces say themselves: what a screen reader (TalkBack)
    calls the icons that have no label beside them (the arrow that goes back,
    the cross that clears a field), and the bar over a text field's selection.
```

## Command bubble-css

```
Command bubble-css writes the Bubble 4 tokens as CSS for a project's accent:

    go run github.com/akrck02/bubble/cmd/bubble-css -accent "#F6BD60" > styles/bubble-tokens.css

The accent is the main colour of the project's logo; Bubble adjusts it until it
reads well.
```
