# 🎨 Your logo's colour

This guide is for someone giving Bubble their project's colour. At the end you will have tokens
for it, on the web and in Go, and know how to see what Bubble made of it.

## On the web

```sh
go run github.com/akrck02/bubble/cmd/bubble-css -accent "#F6BD60" > styles/tokens.css
```

Inside a clone of Bubble, `go run ./cmd/bubble-css` does the same. Load the file before
`bubble.css`. Make it again whenever the colour changes; never edit it by hand.

Without `-accent`, you get Bubble's own blue, `#2C5C9E`.

## In Go

```go
accent, err := bubble.ParseColor("#F6BD60")
if err != nil {
	return err // bubble.invalid-colour: write it as #RRGGBB or #RRGGBBAA
}
options := bubble.DefaultOptions()
options.Accent = accent
theme, err := bubble.NewTheme(options)
```

## See what Bubble made

Bubble may move your colour. The tokens say what it chose:

| Token (CSS) | Field (Go) | What it is |
|---|---|---|
| `--accent` | `Colors.Accent` | the main button's fill |
| `--on-accent` | `Colors.OnAccent` | text and icons on it: white, or a near black |
| `--accent-edge` | `Colors.AccentEdge` | the block edge under the main button (×0.70) |
| `--accent-side` | `Colors.AccentSide` | the pressed button (×0.85) |
| `--accent-ink` | `Colors.AccentInk` | links and accent text on the page |
| `--accent-tint` | `Colors.AccentTint` | a wash for what is chosen, and Material's cards |
| `--focus-ring` | `Colors.FocusRing` | the ring that follows the keyboard |

`grep -- "--accent:" styles/tokens.css` shows the fill for each theme.

## What to expect from a colour

| Your colour | Light mode | Dark mode |
|---|---|---|
| dark (`#2E5E3A`) | kept, white text | lightened, dark text |
| middle (`#4C9A5B`) | kept, dark text | kept, dark text |
| bright (`#F6BD60`, a gold) | kept, dark text | kept, dark text |
| very dark (`#1B2A4A`) | kept, white text | lightened until it stands out from the night |
| very pale (`#F7EBD0`) | kept, dark text | kept, dark text |

The [core ideas](../3-core-ideas.md#the-accent) explain the steps behind this table.

## Variations

- **A colour per page or per user.** Make one `tokens.css` per colour, or call `bubble.NewTheme`
  again in Go; a theme is a plain value.
- **Check a pair yourself.** `bubble.Contrast(a, b)` gives the WCAG contrast of two colours.
