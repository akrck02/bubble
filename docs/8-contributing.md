# 🛠️ Contributing

This page is for someone who wants to change Bubble: fix a part, add one, or change a token. At
the end you will know where everything is, how to check a change and how to propose it.

## Where things are

```
theme.go  palettes.go  scale.go     the theme: options, colours, shapes, type, space
accent.go  color.go                 the accent's steps, colours and contrast
css.go                              the tokens as CSS, one block per choice
errors.go                           the error codes (bubble.unknown-style…)
cmd/bubble-css/                     the command that writes tokens.css
css/tokens.css                      the tokens for Bubble's own accent (generated)
css/bubble.css                      the parts, reading only tokens
goui/                               the Go UI module: Theme, Ground, the showcase and examples
examples/herbarium/                 the quest's web pages, one per level
docs/                               these pages; 6-reference/ has the generators
```

## Build and test

```sh
go test ./...                    # the theme, the contrasts, the CSS rules
cd goui && go vet ./... && cd ..   # the Go UI module (needs go-ui and go-renderer next to Bubble)
```

The tests are the promises Bubble makes: text reads in every theme, a focus ring touches its
control, glass turns opaque when asked. A change that breaks one breaks a promise.

## Changing a token

Tokens are code. Change the Go (`palettes.go`, `scale.go`, `accent.go`), then make the CSS again:

```sh
go run ./cmd/bubble-css > css/tokens.css
python3 docs/6-reference/generate.py
```

`TestTheTokensFileIsUpToDate` fails if you forget the first line.

## Adding a part

1. Write its rules in `css/bubble.css`, in the section where it belongs. Read only tokens: no raw
   colours, no sizes that should follow the density.
2. Mark state with ARIA (`aria-current`, `aria-pressed`, `:checked`), not with classes, so the look
   and the screen reader agree.
3. Check it on glass (`data-glass`), on Material, at night and compact.
4. Add it to `PARTS` in `docs/6-reference/components.py`, and run it to draw its picture.
5. If Go UI has no control for it, show how to draw it from the theme, as the herbarium's `pane`.

## Comparing the browser and Go

```sh
goui/cmd/showcase/compare.sh /tmp/compare
```

It draws the same screen with Bubble's CSS in Chromium and with Go UI, in every look, and prints
how many pixels differ. Today that is 4 to 6 %. A jump means one side drifted. It needs Chromium,
ImageMagick, and Go UI and Go Renderer cloned next to Bubble.

## The docs

- Pages are numbered in reading order; papiro publishes them on every push to `main`.
- `docs/6-reference/generate.py` writes the tokens and Go pages; `components.py` the components.
- The quest's pictures come from `examples/herbarium` and `goui/examples/herbarium`.

## Proposing a change

Open an issue first for anything bigger than a fix: say what an app needed and could not do.
Then a pull request against `main`, with the tests passing and a picture of what changed.
