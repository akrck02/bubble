# 3. Material into Bubble

> Status: building · Author's approval: proposal ☑ plan ☑ (CSS included)

## Why

Go UI for Android's `material` package draws Material 3 the way the author wants it to look:
tonal surfaces tinted from a seed colour, filled fields, the pill behind the rail's icon, every kind
of button, Roboto at Material's scale. Bubble's own material style is a first step: paper with
pill buttons. The colours and the theme are tokens, which Bubble already makes. The parts are
controls that Paper and Glass apps want too. So both move into Bubble, and Go UI for Android keeps
only what is Android and becomes Go for Android.

## What

| Id | Requirement | Check |
|---|---|---|
| R1 | `StyleMaterial` takes its colours from Material's scheme, made from the accent as a seed (tonal palettes, in OKLCH), in light and dark. Paper and Glass don't change | Test that Bubble's material roles equal `material.FromColour` for a set of seeds; the existing tests of Paper and Glass pass unchanged |
| R2 | Every style gets a ladder of surfaces from low to highest, the containers Material's parts sit on. Material's ladder comes from its scheme; Paper's and Glass's from their own colours | Contrast tests: ink reads on every step of every style |
| R3 | The Material style has Material's sizes: 48 px controls, 16/14/22 px text, 24 px icons, its corners. Roboto comes with `goui` | Test of the shape and type; the gallery's captures |
| R4 | The parts of `material/` live in `bubble/goui` as controls that read a Bubble theme: navigation, buttons, selection, chips, fields, containment, dialogs and sheets, pickers, progress, gestures | Their tests, moved and passing |
| R5 | Bubble's material gallery looks like Go UI for Android's: the same screens, compared side by side | Captures of both galleries, every screen, light and dark, and their pixel difference |
| R6 | Bubble stays free of Android. A phone's wallpaper colours reach Bubble as palettes (`bubble.Palettes`), which Go for Android reads from the system | `go list -deps` of `bubble` and `bubble/goui` names no Android package |
| R7 | Bubble's CSS follows the same Material roles, so a web page and a Go window look alike in Material | `tokens.css` holds the ladder and Material's colours; the showcase pair in Material |

## Out of scope

- Renaming Go UI for Android to Go for Android, and deleting what moved from it: that session does
  it once this lands.
- Moving `go-renderer/android` and `go-ui/window/*_android.go` into Go for Android: later.
- The gallery's sound and film screens: they stay in Go for Android, which has go-sound.
- Valhalla's Material: it follows the CSS later, in its own change.

## Design

**Colours (root package, no go-ui).**

- `material.go` brings the seed maths (OKLCH, tonal palettes at Android's shades) and the scheme's
  roles from `seed.go` and `dynamic.go`.
- `NewTheme` with `StyleMaterial` builds the scheme from `Options.Accent`, or from
  `Options.Palettes` when a phone gave them.
- Each role lands on a Bubble colour. Surface becomes background, SurfaceContainerLow becomes
  surface, SurfaceContainerHighest becomes sunken, OnSurface becomes ink, OnSurfaceVariant becomes
  soft ink, OutlineVariant becomes line, Outline becomes the control line, Primary becomes the
  accent, and SecondaryContainer becomes the accent tint.
- The baseline scheme's hand-written colours go: the seed makes them.

**The ladder.** `Colors` gains `SurfaceLow`, `SurfaceHigh` and `SurfaceHighest` around `Surface`
and `Sunken`:

- In Material, they are the scheme's containers.
- In Paper, they are steps of its ground's tone.
- In Glass, they are panes of more white.

This is the surface ladder Bubble 3 had, now with a reason.

**Shape and type.** `shapeOf` and the type scale get Material's sizes. The 44 px touch target stays
the floor.

**Parts (`bubble/goui`).** `material/` moves to `bubble/goui/material`, keeping its `Look` and its
pieces. `Look` reads a `bubble.Theme` instead of a `material.Scheme`: each role it used maps to a
Bubble colour or a step of the ladder. The pieces' code otherwise moves as it is, with its tests. The
gallery goes to `bubble/goui/cmd/gallery`, without its sound and film screens.

**The phone.** `bubble.Palettes` (five tonal palettes of 13 shades) mirrors `window.Palettes`, so
Bubble needs no go-ui to accept them. Go for Android turns the system's palettes into it.

**Checking the look.** `goui/cmd/gallery/compare.sh` captures both galleries, screen by screen, in
light and dark, side by side, and prints their difference, as `showcase/compare.sh` does for the
browser.

**Not taken:**

- Keeping `material/` as a package of its own beside Bubble: two sources of the same tokens would
  drift.
- Making Material's parts Material-only: read through the ladder, they serve Paper and Glass apps.

**A question for the author:** the CSS (R7) widens the work. Everything else can land without it.

## Tasks

- [x] T1 — Seed maths and Material's scheme in the root package · `material.go`, `palettes.go` · R1 R6 · test: `TestMaterialRolesComeFromTheSeedAsGoUIForAndroidMadeThem`
- [x] T2 — The ladder of surfaces in every style, with contrast tests · `theme.go`, `palettes.go`, `accent.go` · R2 · test: `TestInkReadsOnEveryStepOfTheLadder`
- [x] T3 — Material's shape and type, Roboto in `goui` · `scale.go`, `goui/fonts/` · R3 · test of sizes
- [x] T4 — `goui.Theme` maps the ladder and Material's sizes · `goui/theme.go` · R1 R3 · goui tests
- [ ] T5 — Move `material/` into `bubble/goui/material`, reading a Bubble theme; its tests pass · R4 · moved tests
- [ ] T6 — The gallery in `bubble/goui/cmd/gallery`, and `compare.sh` against Go UI for Android's · R5 · captures and their difference
- [ ] T7 — The CSS: the ladder and Material's roles in `tokens.css`, Material's parts in `bubble.css`, if the author wants it now · R7 · CSS tests, showcase pair
- [ ] T8 — Docs: `5-guides/5-material-components.md` with the real API, the tokens reference, the journey · all

## Result

Not built yet.
