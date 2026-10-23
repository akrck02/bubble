package material

import (
	"strings"
	"testing"

	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

func TestAParagraphBreaksIntoLines(t *testing.T) {
	screen := newScreen(t, func(look Look) {
		look.Context.Box(ui.Layout{Width: ui.Fixed(200)}, func() {
			look.Context.Push("text")
			look.Paragraph("one two three four five six seven eight nine ten eleven twelve", 14, draw2d.Regular, look.Roles.OnSurface)
			look.Context.Pop()
		})
	})
	screen.frame()
	lines := map[float32]bool{}
	for _, command := range screen.context.Drawing().Commands {
		if command.Kind != draw2d.TextKind {
			continue
		}
		lines[command.Baseline] = true
		if command.Area.X+command.Area.Width > 200.5 {
			t.Fatalf("%q ends at %v, past the 200-pixel box", command.Text, command.Area.X+command.Area.Width)
		}
	}
	if len(lines) < 2 {
		t.Fatalf("%d lines; want the words broken into several in 200 pixels", len(lines))
	}
}

func TestATextFieldTakesWhatIsTyped(t *testing.T) {
	name := ""
	screen := newScreen(t, func(look Look) {
		look.TextField("name", &name, Field{Label: "Name", Supporting: "As on the chart"})
	})
	screen.frame()
	screen.tap("name")
	if !screen.context.Typing() {
		t.Fatal("a touched field is not being typed in: the keyboard would not come up")
	}
	screen.state.Feed().AddCharacter('Ñ')
	screen.state.Feed().AddCharacter('o')
	screen.frame()
	if name != "Ño" {
		t.Fatalf("name %q; want what was typed", name)
	}
}

func TestAFieldWritesItsPrefixSuffixAndCount(t *testing.T) {
	reach := "18"
	screen := newScreen(t, func(look Look) {
		look.TextField("reach", &reach, Field{Label: "Reach", Prefix: "≈", Suffix: "miles", MaxLength: 3, Counter: true, Filled: true})
	})
	screen.frame()
	for _, text := range []string{"≈", "miles", "2/3", "18"} {
		if !screen.drawn(text) {
			t.Errorf("the field does not write %q", text)
		}
	}
}

func TestAPasswordShowsWhenRevealed(t *testing.T) {
	password, revealed := "lamp", false
	screen := newScreen(t, func(look Look) {
		look.TextField("password", &password, Field{Label: "Password", Secret: true, Reveal: &revealed})
	})
	screen.frame()
	if screen.drawn("lamp") || !screen.drawn("••••") {
		t.Fatal("a password shows its characters before it is revealed")
	}
	screen.tap("icon eye")
	screen.frame()
	if !revealed || !screen.drawn("lamp") {
		t.Fatal("touching the eye does not show the password")
	}
}

func TestADropdownOpensItsChoicesAndTakesOne(t *testing.T) {
	kinds := []string{"Fixed", "Flashing", "Occulting"}
	var dropdown Dropdown
	changes := 0
	screen := newScreen(t, func(look Look) {
		if look.Dropdown("kind", "Kind of light", kinds, &dropdown) {
			changes++
		}
	})
	screen.frame()
	if !screen.drawn("Fixed") || screen.drawn("Occulting") {
		t.Fatal("a closed dropdown should show its choice only")
	}
	screen.tap("kind")
	screen.frame()
	if !screen.drawn("Occulting") {
		t.Fatal("touching the dropdown does not show its choices")
	}
	screen.tap("kind choices", "choice Occulting")
	screen.frame()
	if dropdown.Selected != 2 || changes != 1 || screen.drawn("Flashing") {
		t.Fatalf("selected %d after %d changes; want Occulting chosen once and the choices closed", dropdown.Selected, changes)
	}
}

func TestASearchOpensReadyToTypeAndCloses(t *testing.T) {
	var search Search
	chosen := ""
	screen := newScreen(t, func(look Look) {
		look.SearchView("lights", "Search the lighthouses", &search, func() {
			for _, name := range []string{"North Cape", "Seal Rock"} {
				if strings.Contains(strings.ToLower(name), strings.ToLower(search.Query)) && look.ListItem(name, Item{Headline: name}) {
					chosen, search.Active = name, false
				}
			}
		})
	})
	screen.clocked()
	screen.frame()
	screen.tap("lights")
	for range 3 {
		screen.frame()
	}
	if !search.Active || !screen.context.Typing() {
		t.Fatalf("active %v, typing %v after touching the pill; want the search open with the keyboard up", search.Active, screen.context.Typing())
	}
	screen.state.Feed().AddCharacter('s')
	screen.state.Feed().AddCharacter('e')
	screen.frame()
	screen.frame()
	if search.Query != "se" || screen.drawn("North Cape") || !screen.drawn("Seal Rock") {
		t.Fatalf("query %q; want what was typed, and only what it finds", search.Query)
	}
	screen.tap("lights open", "Seal Rock")
	if chosen != "Seal Rock" || search.Active {
		t.Fatalf("chosen %q, active %v; want Seal Rock chosen and the search closed", chosen, search.Active)
	}
}

func TestTheSelectionBarSpeaksTheLooksWords(t *testing.T) {
	name := "North light"
	words := EnglishWords
	words.Copy = "Copiar"
	screen := newScreen(t, func(look Look) {
		look.Words = words
		look.TextField("name", &name, Field{Label: "Name"})
	})
	screen.frame()
	if got := screen.context.EditingWords().Copy; got != "Copiar" {
		t.Fatalf("the field's bar has the Look's words: %q", got)
	}
}
