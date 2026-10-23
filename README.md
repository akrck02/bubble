<p align="center"><img alt="" src="docs/images/logo.svg" width="120"></p>

# Bubble

A design system that gives web pages and [Go UI](https://github.com/akrck02/go-ui) programs the
same look: tokens made in Go from your logo's colour, CSS parts that read them, and a Go module
that dresses Go UI in them.

[![Test](https://github.com/akrck02/bubble/actions/workflows/test.yaml/badge.svg)](https://github.com/akrck02/bubble/actions/workflows/test.yaml)
[![Licence](https://img.shields.io/github/license/akrck02/bubble)](LICENSE)

> **Status:** in development. Bubble 4 replaces Bubble 3, which is gone. The tokens, the CSS and
> the Go UI module work and are tested; Material's components for Go UI are on their way.

![The same screen in the browser and in Go UI](docs/images/bubble-paper-light.jpg)

```sh
go run ./cmd/bubble-css -accent "#4C9A5B" > tokens.css   # tokens for your colour
```

```html
<link rel="stylesheet" href="tokens.css">
<link rel="stylesheet" href="css/bubble.css">
<button class="button primary">Save</button>   <!-- green, with text that reads -->
```

```go
theme, _ := bubble.NewTheme(bubble.DefaultOptions())
context.Theme = goui.Theme(theme) // every Go UI control now wears Bubble
```

## Install

```sh
git clone -b feature/v4 https://github.com/akrck02/bubble   # main once Bubble 4 lands
```

For the web, copy `css/tokens.css` (or your own) and `css/bubble.css`. For Go, see
[getting started](docs/2-getting-started.md#3-the-same-look-in-a-go-window).

## Documentation

Read it at **<https://akrck02.org/bubble>**, or in [`docs/`](docs/1-overview.md):

- [Overview](docs/1-overview.md) · [Getting started](docs/2-getting-started.md) ·
  [Core ideas](docs/3-core-ideas.md) · 🏴‍☠️ [The quest](docs/4-the-quest/0-the-quest.md)
- Guides: [your logo's colour](docs/5-guides/1-the-accent.md) ·
  [light, dark and compact](docs/5-guides/2-mode-and-density.md) · [glass](docs/5-guides/3-glass.md) ·
  [Go UI](docs/5-guides/4-go-ui.md)
- Reference: [tokens](docs/6-reference/1-tokens.md) · [components](docs/6-reference/2-components.md) ·
  [Go](docs/6-reference/3-go.md)
- [Contributing](docs/8-contributing.md) · [How it was built](docs/9-journey/0-journey.md)

## AI disclosure

Parts of Bubble 4 (the Go code, the CSS and these docs) were written with an AI assistant and are
pending review by the author.
