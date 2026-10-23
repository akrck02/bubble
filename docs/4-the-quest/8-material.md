# 8. 🤖 Material

**Goal:** the herbarium in Material, the look Android users know, without touching its parts.

## On the web

Add one attribute:

```html
<html lang="en" data-style="material">
```

The finished page also takes it from the address: open `6-glass.html?style=material`.

![The herbarium in Material, in the browser](../images/quest/8-material-web.jpg)

What changed, and what did not:

| Changed | The same |
|---|---|
| every colour is made from the green, as Android makes its own | every class in the HTML |
| cards, panes and rows sit on tones of the page, with no lines | the green as the main action |
| fields are filled, with a line under them | the order of the page |
| buttons are pills with medium labels; the index is a rail of icons | the contrast Bubble checks |

## In Go

```sh
go run ./examples/herbarium -style material
```

![The herbarium in Material, in Go](../images/quest/8-material-go.jpg)

In the program it is one field: `Style: bubble.StyleMaterial`, and Roboto for the window
(`goui.FontsFor`). The herbarium's own panes read the style too: on Material they sit on the low
container, as Bubble's CSS draws them.

## Material's own parts

Android apps also share their parts: a navigation bar at the bottom, snackbars, dialogs,
carousels. Bubble has them for Go UI in `goui/material`, the parts Go UI for Android wrote. They
read Bubble's theme, so they work in paper and glass too. [Material's components](../5-guides/5-material-components.md)
shows them, next to the gallery they came from.

## 🏁 The end of the quest

You built one app, in two technologies and four looks, from one set of tokens. From here:

- [The guides](../5-guides/1-the-accent.md) go deeper into each choice.
- [The reference](../6-reference/2-components.md) lists every part and every token.
- [Contributing](../8-contributing.md) shows how to add a part of your own to Bubble.
