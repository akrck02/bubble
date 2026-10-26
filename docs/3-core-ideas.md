# 💡 Core ideas

This page is for someone who wants to know why Bubble behaves as it does. At the end you will
know the four ideas everything else follows from.

## Tokens and parts

Think of a house painted from one paint chart. The chart names the colours: *wall*, *door*,
*trim*. Every room is painted from the chart, never from a pot mixed on the spot. Change the
chart, and every room changes with it.

Bubble works the same way:

- **Tokens** are the chart: named values such as `--accent`, `--surface` or `--card-radius`.
- **Parts** (buttons, cards, fields) are the rooms. They read tokens, and never a raw colour.

```css
.card {
  border: var(--line-width) solid var(--line);
  border-radius: var(--card-radius);
  background: var(--surface);
}
```

Tokens are named by **what they are for**, not by what they look like: `--ink-soft` is secondary
text, whatever its grey. That is why a dark mode is only a second set of values. No part changes.

## Four choices make a theme

A theme is every token for one combination of four choices:

| Choice | Values | Who makes it |
|---|---|---|
| Style | paper, material | the app (material is what Android users expect) |
| Mode | light, dark | the person, or their system |
| Density | comfortable, compact | the app: compact suits desktop tools used with a mouse |
| Accent | any colour | the project: the main colour of its logo |

Paper is a warm page by day and a slate one by night; material is the look of Android apps. Glass
is an extra option of paper only: frosted, translucent panes.

In Go, the choices are the fields of `bubble.Options`, and `bubble.NewTheme` gives the tokens.
On the web, `tokens.css` holds every combination at once. The page picks one with attributes on
`<html>`: `data-mode="dark"`, `data-style="material"`, `data-density="compact"` and `data-glass`.
Without `data-mode`, it follows the system. In Go, glass is `Options.Glass`.

## The accent

akrck02's projects share a family of logos: blocks drawn in 3D (voxels) standing on a dark tile.
A logo's colour is chosen to look good on that tile. It was not chosen to carry text. A pale blue makes a lovely logo and an unreadable button with white letters.

Readability is measured as a **contrast ratio**, from 1 (the same colour) to 21 (black on white).
The web's accessibility guidelines, the WCAG, set the levels. Bubble never uses the accent without
checking them. It **moves it, step by step**, until every pair it is part of reads well:

| Pair | Contrast it needs |
|---|---|
| text on the accent, accent-coloured text on the page | 4.5, the WCAG level for text |
| the focus ring and the edges of controls | 3, the level for shapes |

Each step moves the colour 6 % of the way toward black (in light mode, under white text) or toward
white (in dark mode, under dark text). A bright accent, such as a yellow, keeps its colour and
gets dark text instead: darkening a yellow would turn it brown.

The darker sides borrow from the logos too. In every voxel logo, the right side of a block is its
colour ×0.70 and the left side ×0.85. A checked box wears the ×0.70 edge, and the main button
takes the ×0.85 side under the pointer. The main button itself is flat: its colour marks it.

The tests check every pairing, in every theme, for a range of accents. That is the promise: pick
any colour, and the text still reads.

## One source, two ends

The tokens are written once, in Go. From there they reach two places:

```
          ┌──> cmd/bubble-css ──> tokens.css ──> bubble.css (the web)
bubble ───┤
          └──> goui.Theme ──────> Go UI's theme (Go windows)
```

The CSS parts and Go UI's controls are drawn separately, one by a browser and one by Go Renderer.
`goui/cmd/showcase/compare.sh` draws the same screen in a browser and in Go UI, to keep them
honest ([contributing](8-contributing.md#comparing-the-browser-and-go) shows how to run it). The two pictures differ in
4 to 6 % of their pixels, mostly along the edges of letters.

### Why the tokens are code

A list of values in a file (JSON, for example) was the obvious alternative. It cannot do what the
accent needs: the values are not written down, they are **worked out** from the logo's colour and
checked for contrast. Code can do that, and its tests can prove it.

### Why `goui` is a module of its own

Bubble's package knows nothing of Go UI, and Go UI knows nothing of Bubble. Making the CSS needs
no windowing library, and Go UI can be dressed by other themes. The small `goui` module, with its
own `go.mod`, is the only piece that knows both.

## Where to go from here

- See it all at work: 🏴‍☠️ [the quest](4-the-quest/0-the-quest.md).
- Every token and its value: [the reference](6-reference/1-tokens.md).
