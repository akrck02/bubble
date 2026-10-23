# 🧊 Glass

This guide is for someone who wants frosted panes, with or without a wallpaper behind them. At
the end you will turn glass on and know what it promises.

## On the web

```html
<html lang="en" data-glass style="--wallpaper: url(garden.jpg)">
<body>
<div class="wallpaper"></div>
```

- `data-glass` turns cards, panes, the index and fields into translucent panes.
- `--wallpaper` is optional: any `url()` or gradient. Without it, two soft lights glow behind.
- `.wallpaper` draws it, blurred, dimmed and veiled. Put it first in `<body>`.

Glass works with `data-mode="dark"` too.

## In Go

```go
options := bubble.DefaultOptions()
options.Glass = true
theme, err := bubble.NewTheme(options)
// …
ground := draw2d.NewImage(goui.Ground(theme, width, height))
```

`goui.Ground` paints the two lights as the browser does. Draw it first, then the interface in an
overlay over it; the [quest's level 7](../4-the-quest/7-to-the-desktop.md) does exactly this.

## What glass promises

- **Text reads.** The tests check every text over the ground, over each light, over a dialog
  above any page, and over a wallpaper of any brightness.
- **What is chosen is brighter, never coloured.** The accent stays for the main action.
- **Less transparency means none.** When the system asks for less transparency
  (`prefers-reduced-transparency`), every pane turns opaque.

## Limits

- Glass is a finish of **paper** only. With material, `Glass` is ignored.
- Blur costs: on slow devices, prefer paper.
