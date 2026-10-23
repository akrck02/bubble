# 🌗 Light, dark and compact

This guide is for someone choosing how their app follows the system and how much fits on screen.
At the end you will set the mode, the density and the style, on the web and in Go.

## On the web

Each choice is an attribute on `<html>`:

```html
<html lang="en" data-mode="dark" data-density="compact" data-style="material">
```

| Attribute | Values | Without it |
|---|---|---|
| `data-mode` | `light`, `dark` | follows the system (`prefers-color-scheme`) |
| `data-density` | `compact` | comfortable |
| `data-style` | `material` | paper |

To let people choose, set or remove the attribute from a control; the
[quest's level 5](../4-the-quest/5-night-and-density.md) has a complete switch. Remember the
choice yourself (in `localStorage`, for example): Bubble keeps no state.

## In Go

```go
options := bubble.Options{
	Style:   bubble.StyleMaterial,
	Mode:    bubble.ModeDark,
	Density: bubble.DensityCompact,
	Accent:  accent,
}
theme, err := bubble.NewTheme(options)
```

A wrong value is an error with its own code: `bubble.unknown-style`, `bubble.unknown-mode` or
`bubble.unknown-density`.

In Go UI, the window knows the system's mode: `window.Options.Dark`, which `-dark` also sets.
Make the theme again when it changes.

## Which density?

| Density | For | Controls |
|---|---|---|
| comfortable | web pages, phones, tablets | 44 pixels tall: a finger can hit them |
| compact | desktop tools used with a mouse | 28 pixels tall, with a 24-pixel target at least |

24 pixels is the smallest target the WCAG allows; a mouse hits it easily, a finger does not. On a
touch screen, stay comfortable.

## Which style?

Paper is Bubble's own look. Material is for apps that sit among Android apps. Both give the same
tokens, so switching costs one attribute or one field.
