# 🎨 Tokens

Every token Bubble's parts read, with its value in each theme for Bubble's own accent (`#2C5C9E`). A token is a CSS custom property (`var(--accent)`) and a field of `bubble.Theme` in Go.

This page is made by `docs/6-reference/generate.py` from the code; do not edit it by hand.

## Colours

| Token | Go | Paper | Paper, dark | Material | Material, dark | What it is |
|---|---|---|---|---|---|---|
| `--background` | `Colors.Background` | `#F7F5F0` | `#1A1D24` | `#F6F9FE` | `#0E1115` | the page or window |
| `--surface` | `Colors.Surface` | `#FDFCF9` | `#22262E` | `#EEF4FD` | `#171C23` | fields, raised cards, menus |
| `--ink` | `Colors.Ink` | `#1F2329` | `#E9E7E1` | `#171C23` | `#DCE3ED` | body text |
| `--ink-soft` | `Colors.InkSoft` | `#555A62` | `#B0B3BB` | `#3E4755` | `#BCC7D8` | secondary text and icons |
| `--ink-quiet` | `Colors.InkQuiet` | `#62666E` | `#8F939C` | `#3E4755` | `#BCC7D8` | hints and captions |
| `--line` | `Colors.Line` | `#DDD8CE` | `#2F3540` | `#BCC7D8` | `#3E4755` | separators and the borders of cards |
| `--control-line` | `Colors.ControlLine` | `#8C8578` | `#6E7480` | `#6D7887` | `#8792A1` | borders of fields and secondary buttons: strong enough to see the control |
| `--accent` | `Colors.Accent` | `#2C5C9E` | `#5B80B3` | `#305F9F` | `#A3C9FF` | the main action's fill |
| `--on-accent` | `Colors.OnAccent` | `#FFFFFF` | `#0E1116` | `#FFFFFF` | `#002E68` | text and icons on the accent |
| `--accent-edge` | `Colors.AccentEdge` | `#1F406F` | `#405A7D` | `#305F9F` | `#A3C9FF` | the logo's darker side, ×0.70: the edge of a checked box and of a raised button |
| `--accent-side` | `Colors.AccentSide` | `#254E86` | `#4D6D98` | `#4972AB` | `#8FB6ED` | the logo's lighter side, ×0.85, for pressed states |
| `--accent-ink` | `Colors.AccentInk` | `#2C5C9E` | `#7F9CC4` | `#305F9F` | `#A3C9FF` | accent-coloured text and icons on the background |
| `--accent-tint` | `Colors.AccentTint` | `#D2D9E1` | `#20334F` | `#D1E4FF` | `#354862` | a wash of the accent: the selected item, a checkpoint, a Material container |
| `--positive` | `Colors.Positive` | `#1E7D4F` | `#4CC38A` | `#1E7D4F` | `#4CC38A` | success: saved, done, valid |
| `--warning` | `Colors.Warning` | `#8A5300` | `#E5A64A` | `#8A5300` | `#E5A64A` | something to look at, not yet wrong |
| `--danger` | `Colors.Danger` | `#B3261E` | `#EF7A6C` | `#B3261E` | `#F2B8B5` | errors and actions that lose something |
| `--focus-ring` | `Colors.FocusRing` | `#2C5C9E` | `#7F9CC4` | `#305F9F` | `#A3C9FF` | the ring around whatever the keyboard is on |
| `--selection` | `Colors.Selection` | `#2C5C9E4D` | `#5B80B34D` | `#305F9F4D` | `#A3C9FF4D` | behind selected text |
| `--shadow` | `Colors.Shadow` | `#1F232929` | `#00000066` | `#1F232929` | `#00000066` | the colour of shadows |
| `--scrim` | `Colors.Scrim` | `#1F232966` | `#0000008C` | `#00000052` | `#00000052` | over everything behind a sheet or a dialog |
| `--tile` | `Colors.Tile` | `#2F3540` | `#2F3540` | `#2F3540` | `#2F3540` | the dark tile every voxel logo stands on |
| `--container-low` | `Colors.ContainerLow` | `#FDFCF9` | `#22262E` | `#EEF4FD` | `#171C23` |  |
| `--container` | `Colors.Container` | `#F8F6F1` | `#1E2129` | `#E7EEF9` | `#1B2027` |  |
| `--container-high` | `Colors.ContainerHigh` | `#F4F0E8` | `#191D23` | `#E2E9F3` | `#252B32` |  |
| `--container-highest` | `Colors.ContainerHighest` | `#EFEAE0` | `#15181E` | `#DCE3ED` | `#2F353D` |  |
| `--on-accent-tint` | `Colors.OnAccentTint` | `#1F2329` | `#E9E7E1` | `#0B1C33` | `#D1E4FF` | text and icons on the accent tint |
| `--accent-container` | `Colors.AccentContainer` | `#D2D9E1` | `#20334F` | `#D1E4FF` | `#184684` | a stronger wash of the accent: a chosen chip, a floating button |
| `--on-accent-container` | `Colors.OnAccentContainer` | `#1F2329` | `#E9E7E1` | `#001A41` | `#D1E4FF` |  |
| `--on-danger` | `Colors.OnDanger` | `#FFFFFF` | `#0E1116` | `#FFFFFF` | `#601410` | text on a danger fill |
| `--danger-container` | `Colors.DangerContainer` | `#EDD6D1` | `#3A2B2F` | `#F9DEDC` | `#8C1D18` | a wash of danger: an error banner |
| `--on-danger-container` | `Colors.OnDangerContainer` | `#1F2329` | `#E9E7E1` | `#410E0B` | `#F9DEDC` |  |
| `--inverse` | `Colors.Inverse` | `#1F2329` | `#E9E7E1` | `#2B3139` | `#DCE3ED` | what stands out against the page: a snackbar |
| `--on-inverse` | `Colors.OnInverse` | `#F7F5F0` | `#1A1D24` | `#EAF1FC` | `#2B3139` |  |
| `--inverse-accent` | `Colors.InverseAccent` | `#8BA5CA` | `#324662` | `#A3C9FF` | `#305F9F` | the accent on the inverse surface: a snackbar's action |
| `--color-scheme` |  | `light` | `dark` | `light` | `dark` | light or dark, for the browser's own controls and scroll bars |

## Shape

| Token | Go | Paper | Compact | Material | What it is |
|---|---|---|---|---|---|
| `--control-radius` | `Shape.ControlRadius` | `10px` | `6px` | `20px` | fields and selects |
| `--button-radius` | `Shape.ButtonRadius` | `10px` | `6px` | `999px` | buttons and chips |
| `--card-radius` | `Shape.CardRadius` | `12px` | `7.2px` | `12px` | cards and panes |
| `--panel-radius` | `Shape.PanelRadius` | `16px` | `9.6px` | `28px` | windows, sheets, dialogs |
| `--floating-radius` | `Shape.FloatingRadius` | `999px` | `999px` | `16px` | the floating action button |
| `--line-width` | `Shape.LineWidth` | `1px` | `1px` | `1px` | lines and borders |
| `--edge-width` | `Shape.EdgeWidth` | `0px` | `0px` | `0px` | a block edge under the main button; 0 in every style today, which draws it flat |
| `--focus-width` | `Shape.FocusWidth` | `2px` | `2px` | `2px` | the focus ring's thickness |
| `--focus-offset` | `Shape.FocusOffset` | `0px` | `0px` | `0px` | the gap between a control and its focus ring; 0, so the edge lights up, kept for go-ui |
| `--elevation` |  | `none` | `none` | `0 1px 3px var(--shadow), 0 4px 12px var(--shadow)` | the shadow of raised things: none on paper, two soft shadows on material |

## Type

| Token | Go | Paper | Compact | Material | What it is |
|---|---|---|---|---|---|
| `--font` | `Type.Family` | `system-ui, -apple-system, "Segoe UI", sans-serif` | `system-ui, -apple-system, "Segoe UI", sans-serif` | `Roboto, system-ui, -apple-system, "Segoe UI", sans-serif` | the system's own interface typeface |
| `--font-monospace` | `Type.MonospaceFamily` | `ui-monospace, "SF Mono", Menlo, Consolas, monospace` | `ui-monospace, "SF Mono", Menlo, Consolas, monospace` | `ui-monospace, "SF Mono", Menlo, Consolas, monospace` | code and keys |
| `--text-small` | `Type.SmallSize` | `13px` | `11px` | `14px` | captions, hints, pills |
| `--text-body` | `Type.BodySize` | `16px` | `13px` | `16px` | body text |
| `--text-heading` | `Type.HeadingSize` | `20px` | `15px` | `22px` | section headings |
| `--text-title` | `Type.TitleSize` | `30px` | `20px` | `32px` | a page's title |
| `--text-display` | `Type.DisplaySize` | `44px` | `28px` | `45px` | a big number or a welcome |
| `--icon-size` | `Type.IconSize` | `20px` | `16px` | `24px` | icons beside text |
| `--line-height` | `Type.LineHeight` | `1.5` | `1.4` | `1.5` | a multiple of the size |

## Space

| Token | Go | Paper | Compact | Material | What it is |
|---|---|---|---|---|---|
| `--gap` | `Space.Gap` | `12px` | `8px` | `12px` | between controls |
| `--padding` | `Space.Padding` | `20px` | `12px` | `16px` | inside cards and panels |
| `--control-height` | `Space.ControlHeight` | `44px` | `28px` | `48px` | buttons, fields and selects |
| `--touch-target` | `Space.TouchTarget` | `44px` | `24px` | `48px` | the smallest area a finger or pointer must be able to hit |

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
