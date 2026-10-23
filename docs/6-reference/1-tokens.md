# 🎨 Tokens

Every token Bubble's parts read, with its value in each theme for Bubble's own accent (`#2C5C9E`). A token is a CSS custom property (`var(--accent)`) and a field of `bubble.Theme` in Go.

This page is made by `docs/6-reference/generate.py` from the code; do not edit it by hand.

## Colours

| Token | Go | Paper | Paper, dark | Material | Material, dark | What it is |
|---|---|---|---|---|---|---|
| `--background` | `Colors.Background` | `#F7F5F0` | `#1A1D24` | `#F4F4F6` | `#121316` | the page or window |
| `--surface` | `Colors.Surface` | `#FDFCF9` | `#22262E` | `#FFFFFF` | `#1E1F24` | fields, raised cards, menus |
| `--ink` | `Colors.Ink` | `#1F2329` | `#E9E7E1` | `#1B1B1F` | `#E4E2E6` | body text |
| `--ink-soft` | `Colors.InkSoft` | `#555A62` | `#B0B3BB` | `#45474F` | `#C5C6D0` | secondary text and icons |
| `--ink-quiet` | `Colors.InkQuiet` | `#62666E` | `#8F939C` | `#5E616A` | `#9A9BA5` | hints and captions |
| `--line` | `Colors.Line` | `#DDD8CE` | `#2F3540` | `#D5D6DD` | `#34363C` | separators and the borders of cards |
| `--control-line` | `Colors.ControlLine` | `#8C8578` | `#6E7480` | `#74777F` | `#8E9099` | borders of fields and secondary buttons: strong enough to see the control |
| `--accent` | `Colors.Accent` | `#2C5C9E` | `#5B80B3` | `#2C5C9E` | `#5B80B3` | the main action's fill |
| `--on-accent` | `Colors.OnAccent` | `#FFFFFF` | `#0E1116` | `#FFFFFF` | `#0E1116` | text and icons on the accent |
| `--accent-edge` | `Colors.AccentEdge` | `#1F406F` | `#405A7D` | `#1F406F` | `#405A7D` | the block edge under a primary button: the logo's darker side, ×0.70 |
| `--accent-side` | `Colors.AccentSide` | `#254E86` | `#4D6D98` | `#254E86` | `#4D6D98` | the logo's lighter side, ×0.85, for pressed states |
| `--accent-ink` | `Colors.AccentInk` | `#2C5C9E` | `#7F9CC4` | `#2C5C9E` | `#7796C0` | accent-coloured text and icons on the background |
| `--accent-tint` | `Colors.AccentTint` | `#D2D9E1` | `#20334F` | `#CAD6E7` | `#1B2D46` | a wash of the accent: the selected item, a checkpoint, a Material container |
| `--positive` | `Colors.Positive` | `#1E7D4F` | `#4CC38A` | `#1E7D4F` | `#4CC38A` | success: saved, done, valid |
| `--warning` | `Colors.Warning` | `#8A5300` | `#E5A64A` | `#8A5300` | `#E5A64A` | something to look at, not yet wrong |
| `--danger` | `Colors.Danger` | `#B3261E` | `#EF7A6C` | `#B3261E` | `#EF7A6C` | errors and actions that lose something |
| `--focus-ring` | `Colors.FocusRing` | `#2C5C9E` | `#7F9CC4` | `#2C5C9E` | `#7796C0` | the ring around whatever the keyboard is on |
| `--selection` | `Colors.Selection` | `#2C5C9E4D` | `#5B80B34D` | `#2C5C9E4D` | `#5B80B34D` | behind selected text |
| `--shadow` | `Colors.Shadow` | `#1F232929` | `#00000066` | `#1F232929` | `#00000066` | the colour of shadows |
| `--scrim` | `Colors.Scrim` | `#1F232966` | `#0000008C` | `#1F232966` | `#0000008C` | over everything behind a sheet or a dialog |
| `--tile` | `Colors.Tile` | `#2F3540` | `#2F3540` | `#2F3540` | `#2F3540` | the dark tile every voxel logo stands on |
| `--color-scheme` |  | `light` | `dark` | `light` | `dark` | light or dark, for the browser's own controls and scroll bars |

## Shape

| Token | Go | Paper | Compact | Material | What it is |
|---|---|---|---|---|---|
| `--control-radius` | `Shape.ControlRadius` | `10px` | `6px` | `12px` | fields and selects |
| `--button-radius` | `Shape.ButtonRadius` | `10px` | `6px` | `999px` | buttons and chips |
| `--card-radius` | `Shape.CardRadius` | `12px` | `7.199999999999999px` | `20px` | cards and panes |
| `--panel-radius` | `Shape.PanelRadius` | `16px` | `9.6px` | `28px` | windows, sheets, dialogs |
| `--floating-radius` | `Shape.FloatingRadius` | `999px` | `999px` | `16px` | the floating action button |
| `--line-width` | `Shape.LineWidth` | `1px` | `1px` | `1px` | lines and borders |
| `--edge-width` | `Shape.EdgeWidth` | `2px` | `2px` | `0px` | the block edge under primary buttons; 0 when the style has none |
| `--focus-width` | `Shape.FocusWidth` | `2px` | `2px` | `2px` | the focus ring's thickness |
| `--focus-offset` | `Shape.FocusOffset` | `0px` | `0px` | `0px` | the gap between a control and its focus ring; 0, so the edge lights up, kept for go-ui |
| `--elevation` |  | `none` | `none` | `0 1px 3px var(--shadow), 0 4px 12px var(--shadow)` | the shadow of raised things: none on paper, two soft shadows on material |

## Type

| Token | Go | Paper | Compact | Material | What it is |
|---|---|---|---|---|---|
| `--font` | `Type.Family` | `system-ui, -apple-system, "Segoe UI", sans-serif` | `system-ui, -apple-system, "Segoe UI", sans-serif` | `same` | the system's own interface typeface |
| `--font-monospace` | `Type.MonospaceFamily` | `ui-monospace, "SF Mono", Menlo, Consolas, monospace` | `ui-monospace, "SF Mono", Menlo, Consolas, monospace` | `same` | code and keys |
| `--text-small` | `Type.SmallSize` | `13px` | `11px` | `same` | captions, hints, pills |
| `--text-body` | `Type.BodySize` | `16px` | `13px` | `same` | body text |
| `--text-heading` | `Type.HeadingSize` | `20px` | `15px` | `same` | section headings |
| `--text-title` | `Type.TitleSize` | `30px` | `20px` | `same` | a page's title |
| `--text-display` | `Type.DisplaySize` | `44px` | `28px` | `same` | a big number or a welcome |
| `--icon-size` | `Type.IconSize` | `20px` | `16px` | `same` | icons beside text |
| `--line-height` | `Type.LineHeight` | `1.5` | `1.4` | `same` | a multiple of the size |

## Space

| Token | Go | Paper | Compact | Material | What it is |
|---|---|---|---|---|---|
| `--gap` | `Space.Gap` | `12px` | `8px` | `same` | between controls |
| `--padding` | `Space.Padding` | `20px` | `12px` | `same` | inside cards and panels |
| `--control-height` | `Space.ControlHeight` | `44px` | `28px` | `same` | buttons, fields and selects |
| `--touch-target` | `Space.TouchTarget` | `44px` | `24px` | `same` | the smallest area a finger or pointer must be able to hit |

## Motion

| Token | Go | Paper | Compact | Material | What it is |
|---|---|---|---|---|---|
| `--motion-fast` | `Motion.Fast` | `0.25s` | `same` | `same` | a hover, a press |
| `--motion-medium` | `Motion.Medium` | `0.5s` | `same` | `same` | a panel opening |
| `--motion-slow` | `Motion.Slow` | `0.75s` | `same` | `same` | a page changing |

## Glass

Only under `data-glass` (`Options.Glass` in Go). On glass, `--sunken` takes the raised pane's value.

| Token | Go | Glass | Glass, dark | What it is |
|---|---|---|---|---|
| `--glass-pane` | `Glass.Pane` | `#FFFFFF8C` | `#FFFFFF0F` | the translucent fill of panels, cards and fields, with no outline |
| `--glass-raised` | `Glass.Raised` | `#FFFFFFD9` | `#FFFFFF1F` | a brighter fill, for what is hovered or chosen |
| `--glass-dialog` | `Glass.Dialog` | `#EFECE6E6` | `#16181DDB` | the denser fill of a dialog, which floats over other text |
| `--glass-glow-one` | `Glass.Glows` | `#FFFFFF` | `#24272D` | one of the two soft lights behind the panes |
| `--glass-glow-two` | `Glass.Glows` | `#CBC6BB` | `#181A1F` | one of the two soft lights behind the panes |
| `--glass-veil` | `Glass.Veil` | `#E6E3DCB8` | `#0B0C0FB8` | laid over a wallpaper, so any picture behind keeps the text readable |
| `--glass-dimming` | `Glass.Dimming` | `1` | `0.45` | how bright a wallpaper is kept before the veil, from 0 to 1 (CSS brightness) |
| `--glass-blur` | `Glass.Blur` | `24px` | `24px` | how much of what is behind shows through, in layout pixels |

## Choosing a theme

| Attribute on `<html>` | Go | Effect |
|---|---|---|
| (none) | `bubble.DefaultOptions()` | paper, light (or the system's dark), comfortable |
| `data-mode="dark"` / `"light"` | `Mode: bubble.ModeDark` | forces a mode |
| `data-style="material"` | `Style: bubble.StyleMaterial` | Material |
| `data-density="compact"` | `Density: bubble.DensityCompact` | compact |
| `data-glass` | `Glass: true` | glass (paper only) |
