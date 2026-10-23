// Package material draws Material 3's pieces with go-ui, in a Bubble theme (with Material 3
// Expressive's newer pieces). They came from Go UI for Android. The pieces are made of go-ui
// boxes, texts and clickables, as any program can make controls of its own, so they draw the same
// on a phone, on a desktop and in a browser. They read the theme through Roles, Material's names
// for its colours, so they draw in every Bubble style: Material's own, paper and glass. A program
// makes a Look for each frame and calls its pieces:
//
//	theme, _ := bubble.NewTheme(bubble.Options{Style: bubble.StyleMaterial, Mode: bubble.ModeLight,
//		Density: bubble.DensityComfortable, Accent: accent})
//	context.Theme = goui.Theme(theme)
//	look := material.New(context, theme)
//	look.TopAppBar(material.TopBar{Title: "Inbox"})
//	if look.FilledButton("plus", "New message") { … }
//
// On a phone, Go for Android reads the wallpaper's palettes, which go in bubble.Options.Palettes.
// The files follow Material's catalogue:
//
//	look.go         the Look, text and paragraphs
//	roles.go        a Bubble theme's colours by Material's names
//	navigation.go   top bars, the navigation bar, rail and drawer, tabs and swiping pages, the back
//	actions.go      buttons, floating buttons and their menu, split buttons, button groups, toolbars
//	selection.go    switches, check boxes, radio buttons, segmented buttons, sliders
//	chips.go        filter, assist and input chips, and tags
//	inputs.go       text fields and the search bar
//	containment.go  cards, list items, avatars, banners, section titles, dividers
//	overlays.go     dialogs, bottom sheets, snackbars, menus
//	pickers.go      date and time pickers
//	progress.go     loading indicators, circular, linear and wavy progress, skeletons
//	gestures.go     pull to refresh, rows swiped away, carousels
//	stick.go        an on-screen stick
package material
