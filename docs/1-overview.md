<p align="center"><img alt="" src="images/logo.svg" width="120"></p>

# Bubble

Bubble is a design system: the colours, sizes and parts every interface of akrck02's is drawn with.
It gives **the same look in a web page and in a Go window**: its CSS dresses HTML, and its Go
module dresses [Go UI](https://github.com/akrck02/go-ui).

![The same screen in the browser (left) and in a Go UI window (right)](images/bubble-paper-light.jpg)

## 🫧 Is it for me?

| You want to… | Bubble gives you |
|---|---|
| style a web app without designing it | two CSS files: the tokens and the components |
| make a Go UI program look like your web app | `goui.Theme`, one line before the first control |
| use your logo's colour without failing contrast | the accent, adjusted until every text reads |
| offer dark mode, a compact view or Material | one attribute on `<html>`, or one option in Go |

You do not need to know design. You need to know a little HTML, or a little Go.

## 🧭 Where to start

| If you… | Read | Time |
|---|---|---|
| want to see it working | [Getting started](2-getting-started.md) | 10 min |
| learn by building | 🏴‍☠️ [The quest](4-the-quest/0-the-quest.md): a herbarium, in 8 levels | 1 h 30 |
| want to know why it behaves as it does | [Core ideas](3-core-ideas.md) | 15 min |
| have one task in mind | [The guides](5-guides/1-the-accent.md) | — |
| look up a token, a class or a function | [The reference](6-reference/1-tokens.md) | — |

## 📚 All the pages, in reading order

**Start**

- [Overview](1-overview.md): this page
- [Getting started](2-getting-started.md): Bubble in a web page, then in a Go window

**Understand**

- [Core ideas](3-core-ideas.md): tokens and parts, the four choices, the accent

**Learn**

- [The quest](4-the-quest/0-the-quest.md): build a herbarium, from one page to a Go program

**Use**

- [Your logo's colour](5-guides/1-the-accent.md)
- [Light, dark and compact](5-guides/2-mode-and-density.md)
- [Glass](5-guides/3-glass.md)
- [Drawing Go UI in Bubble](5-guides/4-go-ui.md)
- [Material's components](5-guides/5-material-components.md) (coming)

**Look up**

- [Tokens](6-reference/1-tokens.md) · [Components](6-reference/2-components.md) · [Go](6-reference/3-go.md)

**The project**

- [Contributing](8-contributing.md) · [How it was built](9-journey/0-journey.md)

## 🎨 Three looks

| Look | What it feels like | How to ask for it |
|---|---|---|
| **Paper** (the default) | a warm page by day, the slate of the logos by night | nothing |
| **Material** (half done) | what Android users know: tinted cards, a rail of icons | `data-style="material"` |
| **Glass** | frosted panes over a soft light or a wallpaper | `data-glass` (paper only) |

Each comes in light and dark, and in a comfortable or a compact density.

## 📖 Words you will meet

| Word | What it means |
|---|---|
| Token | A named value, such as `--accent` or `--card-radius`. Parts read tokens, never raw colours. |
| Accent | The project's own colour, taken from its logo: the main button, links, the focus ring. |
| Contrast | How much text stands out from what is behind it. Bubble asks for 4.5 for text, 3 for edges. |
| Density | How much fits on screen: comfortable for fingers, compact for a mouse. |
| Go UI | akrck02's library for interfaces drawn in Go, in a window or a phone. |

## 🚧 Status

**In development.** Bubble 4 is a new start: Bubble 3, a CSS framework of classes, is gone.
The tokens, the CSS and the Go UI module work and are tested. Coming next:

- Material's full look and its components for Go UI (navigation bars, dialogs, snackbars),
  moving here from Go UI for Android, whose Material is the reference. See [Material's components](5-guides/5-material-components.md).
- A release with a version number. Until then, take the files from the `main` branch.

## 👪 The family

- [Go UI](https://github.com/akrck02/go-ui): interfaces in Go, which Bubble dresses.
- [Go Renderer](https://github.com/akrck02/go-renderer): the window and drawing under Go UI.
- [GTD](https://github.com/akrck02/gtd): the web library whose apps wear Bubble's CSS.
- [Papiro](https://github.com/akrck02/papiro): the wiki that publishes these pages.
