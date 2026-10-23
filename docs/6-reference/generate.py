#!/usr/bin/env python3
"""Writes 1-tokens.md from the code: the token names from css.go, their meanings from the comments
of theme.go, and their values from css/tokens.css (Bubble's own accent).

    python3 docs/6-reference/generate.py      from the repository's root
"""
import re

css_go = open("css.go").read()
theme_go = open("theme.go").read()
tokens = open("css/tokens.css").read()

# what each Go field is for, from the comment beside it in theme.go
meaning = {}
for struct, body in re.findall(r"type (\w+) struct \{(.*?)\n\}", theme_go, re.S):
    for line in body.split("\n"):
        field = re.match(r"\s+([A-Z]\w*)(?:, ([A-Z]\w*))?\s+[\w.\[\]]+\s*(?://\s*(.*))?$", line)
        if field:
            for name in filter(None, field.group(1, 2)):
                meaning[(struct, name)] = (field.group(3) or "").strip()

# each token's Go field, from css.go's {"name", group.Field…} pairs
groups = {"colors": "Colors", "shape": "Shape", "kind": "Type", "space": "Space", "motion": "Motion", "glass": "Glass"}
field_of = {}
for name, group, field in re.findall(r'\{"([a-z-]+)", (?:pixels\()?(?:fmt\.Sprintf\("[^"]*", )?(\w+)\.(\w+)', css_go):
    if group in groups:
        field_of[name] = (groups[group], field.replace("Glows", "Glows"))

# the values of each block of tokens.css (the media copies repeat the dark blocks)
blocks = {}
for selector, body in re.findall(r"\n(:root[^{]*) \{\n(.*?)\n\}", tokens, re.S):
    blocks[selector] = dict(re.findall(r"--([a-z-]+): ([^;]+);", body))
def value(selector, name):
    return blocks.get(selector, {}).get(name, "")

LIGHT, DARK = ":root", ':root[data-mode="dark"]'
MATERIAL, MATERIAL_DARK = ':root[data-style="material"]', ':root[data-style="material"][data-mode="dark"]'
COMPACT, GLASS, GLASS_DARK = ':root[data-density="compact"]', ":root[data-glass]", ':root[data-glass][data-mode="dark"]'

# fields theme.go leaves without a comment, because their name says it
PLAIN = {
    "Positive": "success: saved, done, valid", "Warning": "something to look at, not yet wrong",
    "Danger": "errors and actions that lose something", "Shadow": "the colour of shadows",
    "CardRadius": "cards and panes", "LineWidth": "lines and borders",
    "FocusWidth": "the focus ring's thickness", "MonospaceFamily": "code and keys",
    "SmallSize": "captions, hints, pills", "BodySize": "body text", "HeadingSize": "section headings",
    "TitleSize": "a page's title", "DisplaySize": "a big number or a welcome", "IconSize": "icons beside text",
    "ControlHeight": "buttons, fields and selects", "Fast": "a hover, a press",
    "Medium": "a panel opening", "Slow": "a page changing",
}

def describe(name):
    struct, field = field_of.get(name, ("", ""))
    go = f"`{struct}.{field}`" if field else ""
    return go, meaning.get((struct, field)) or PLAIN.get(field, "")

def swatch(colour):
    return f"`{colour}`" if colour else "—"

out = ["# 🎨 Tokens", "",
       "Every token Bubble's parts read, with its value in each theme for Bubble's own accent "
       "(`#2C5C9E`). A token is a CSS custom property (`var(--accent)`) and a field of `bubble.Theme` in Go.",
       "", "This page is made by `docs/6-reference/generate.py` from the code; do not edit it by hand.", ""]

out += ["## Colours", "", "| Token | Go | Paper | Paper, dark | Material | Material, dark | What it is |", "|---|---|---|---|---|---|---|"]
for name in blocks[LIGHT]:
    if field_of.get(name, ("",))[0] != "Colors" and name != "color-scheme":
        continue
    go, text = describe(name)
    if name == "color-scheme":
        text = "light or dark, for the browser's own controls and scroll bars"
    out.append(f"| `--{name}` | {go} | {swatch(value(LIGHT, name))} | {swatch(value(DARK, name))} | "
               f"{swatch(value(MATERIAL, name))} | {swatch(value(MATERIAL_DARK, name))} | {text} |")

def sizes(title, struct, extra=()):
    rows = ["", f"## {title}", "", "| Token | Go | Paper | Compact | Material | What it is |", "|---|---|---|---|---|---|"]
    for name in blocks[LIGHT]:
        if field_of.get(name, ("",))[0] != struct and name not in extra:
            continue
        go, text = describe(name)
        if name == "elevation":
            text = "the shadow of raised things: none on paper, two soft shadows on material"
        material = value(MATERIAL, name) or "same"
        compact = value(COMPACT, name) or "same"
        rows.append(f"| `--{name}` | {go} | `{value(LIGHT, name)}` | `{compact}` | `{material}` | {text} |")
    return rows

out += sizes("Shape", "Shape", ("elevation",))
out += sizes("Type", "Type")
out += sizes("Space", "Space")
out += sizes("Motion", "Motion")

out += ["", "## Glass", "", "Only under `data-glass` (`Options.Glass` in Go). On glass, `--sunken` takes the raised pane's value.",
        "", "| Token | Go | Glass | Glass, dark | What it is |", "|---|---|---|---|---|"]
for name in blocks[GLASS]:
    if not name.startswith("glass-"):
        continue
    go, text = describe(name)
    if name.startswith("glass-glow"):
        go, text = "`Glass.Glows`", "one of the two soft lights behind the panes"
    out.append(f"| `--{name}` | {go} | `{value(GLASS, name)}` | `{value(GLASS_DARK, name) or value(GLASS, name)}` | {text} |")

out += ["", "## Choosing a theme", "",
        "| Attribute on `<html>` | Go | Effect |", "|---|---|---|",
        "| (none) | `bubble.DefaultOptions()` | paper, light (or the system's dark), comfortable |",
        '| `data-mode="dark"` / `"light"` | `Mode: bubble.ModeDark` | forces a mode |',
        '| `data-style="material"` | `Style: bubble.StyleMaterial` | Material |',
        '| `data-density="compact"` | `Density: bubble.DensityCompact` | compact |',
        "| `data-glass` | `Glass: true` | glass (paper only) |", ""]
open("docs/6-reference/1-tokens.md", "w").write("\n".join(out))
print("docs/6-reference/1-tokens.md")

# 3-go.md: the Go packages, as go doc prints them
import subprocess
def godoc(folder, package):
    return subprocess.run(["go", "doc", "-all", package], cwd=folder, capture_output=True, text=True, check=True).stdout.rstrip()

go = ["# 🐹 Go", "",
      "The Go API of Bubble's two modules, as `go doc -all` prints it. This page is made by "
      "`docs/6-reference/generate.py`; do not edit it by hand.", "",
      "| You want to… | Use |", "|---|---|",
      "| make a theme | `bubble.NewTheme(bubble.Options{…})`, or start from `bubble.DefaultOptions()` |",
      "| read a colour | `bubble.ParseColor(\"#4C9A5B\")` |",
      "| write the CSS tokens | `bubble.CSS(writer, accent)`, or the `bubble-css` command |",
      "| check two colours | `bubble.Contrast(a, b)` |",
      "| dress Go UI | `goui.Theme(theme)`, and `goui.Ground` for glass |", "",
      "## Package bubble", "", "`github.com/akrck02/bubble`", "", "```", godoc(".", "."), "```", "",
      "## Package goui", "", "`github.com/akrck02/bubble/goui`, a module of its own (`goui/go.mod`).", "",
      "```", godoc("goui", "."), "```", "",
      "## Command bubble-css", "", "```", godoc(".", "./cmd/bubble-css"), "```", ""]
open("docs/6-reference/3-go.md", "w").write("\n".join(go))
print("docs/6-reference/3-go.md")
