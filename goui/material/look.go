package material

import (
	"strings"
	"time"

	"github.com/akrck02/bubble"
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// safePadding is a margin from the window's edges plus what the system draws over them.
func (look Look) safePadding(margin float32) ui.Padding {
	insets := look.Context.Insets()
	return ui.Padding{Left: margin + insets.Left, Top: margin + insets.Top, Right: margin + insets.Right, Bottom: margin + insets.Bottom}
}

// Look draws the pieces in one frame, with one scheme. It is small: make one in every frame.
//
// Pieces that can be touched return true in the frame they are touched. Each takes the name go-ui
// tells it apart by: its label or icon, or a name given. Two pieces with the same name in one scope
// are one to go-ui: give repeated ones a scope each (ui.Context.Push).
type Look struct {
	Context *ui.Context
	Theme   bubble.Theme // the Bubble theme the pieces are drawn in
	Roles   Roles        // its colours by Material's names, as the pieces read them
	Words   Words        // the pieces' own words, in the person's language; the zero value is EnglishWords
}

// Words are what the pieces say themselves: what a screen reader (TalkBack) calls the icons that
// have no label beside them (the arrow that goes back, the cross that clears a field), and the bar
// over a text field's selection.
type Words struct {
	Back, Close, Clear          string
	ShowPassword, HidePassword  string
	PreviousMonth, NextMonth    string
	MoreOptions                 string // a menu's button, a split button's arrow
	Cut, Copy, Paste, SelectAll string // the bar over a text field's selection
}

// EnglishWords are the words a screen reader reads by default.
var EnglishWords = Words{Back: "Back", Close: "Close", Clear: "Clear", ShowPassword: "Show password",
	HidePassword: "Hide password", PreviousMonth: "Previous month", NextMonth: "Next month",
	MoreOptions: "More options", Cut: "Cut", Copy: "Copy", Paste: "Paste", SelectAll: "Select all"}

// words are the Look's words, or English ones.
func (look Look) words() Words {
	if look.Words == (Words{}) {
		return EnglishWords
	}
	return look.Words
}

// NamedIconButton is an IconButton a screen reader (TalkBack) reads by a name, what it does ("Share"),
// rather than by its icon's.
func (look Look) NamedIconButton(icon, name string) bool {
	touched := look.IconButton(icon)
	look.Context.Describe(ui.Description{Label: name})
	return touched
}

// New makes the frame's Look.
func New(context *ui.Context, theme bubble.Theme) Look {
	return Look{Context: context, Theme: theme, Roles: RolesOf(theme)}
}

// Text writes a line in a size, a weight and a colour.
func (look Look) Text(text string, size float32, weight draw2d.Weight, colour draw2d.Color) {
	look.Context.Text(text, ui.TextStyle{Size: size, Weight: weight, Color: colour})
}

// stateLayer is the layer Material puts over something touched: its content's colour, faintly.
func stateLayer(colour draw2d.Color, strength float32) draw2d.Color {
	return draw2d.Color{colour[0], colour[1], colour[2], strength}
}

// pressedOnly is a piece's looks when only a touch changes it: a finger has no hover.
func pressedOnly(normal ui.Look, pressedFill draw2d.Color) ui.Looks {
	pressed := normal
	pressed.Fill = pressedFill
	return ui.Looks{Normal: normal, Hovered: normal, Pressed: pressed}
}

// Paragraph writes text that breaks into lines at the width it has, word by word: a card's body, a
// description. Text writes one line, however long.
func (look Look) Paragraph(text string, size float32, weight draw2d.Weight, colour draw2d.Color) {
	style := ui.TextStyle{Size: size, Weight: weight, Color: colour}
	// a wrapping box puts its children on lines with one gap for both ways: a space's width
	look.Context.Box(ui.Layout{Width: ui.Grow(), Gap: size * 0.28, Wrap: true}, func() {
		for _, word := range strings.Fields(text) {
			look.Context.Text(word, style)
		}
	})
}

// Material's short motions: a control's own change, something arriving on the screen, and
// something leaving it, which goes quicker than it came.
const (
	switchMotion   = 150 * time.Millisecond
	arrivingMotion = 250 * time.Millisecond
	leavingMotion  = 200 * time.Millisecond
)

// faded is a colour with only part of its opacity left, for something on its way out. It never
// reaches nothing: go-ui reads a colour with no opacity as unset and writes in the theme's ink,
// while 1/512 shows as nothing on an 8-bit screen.
func faded(colour draw2d.Color, left float32) draw2d.Color {
	colour[3] = max(colour[3]*left, 1.0/512)
	return colour
}

// shownName is the part of a name that is shown and read: what comes before "##", which only
// tells two pieces with the same label apart.
func shownName(name string) string {
	shown, _, _ := strings.Cut(name, "##")
	return shown
}

func boolFloat(value bool) float32 {
	if value {
		return 1
	}
	return 0
}
