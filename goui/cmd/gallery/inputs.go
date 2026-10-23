package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/akrck02/bubble/goui/material"
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// inputsState is what the inputs screen remembers.
type inputsState struct {
	search       material.Search
	name, reach  string
	password     string
	revealed     bool
	light        material.Dropdown
	kind, view   int
	nightly      bool
	filters      [3]bool
	tags         []string
	volume       float32
	date         time.Time
	datePicker   material.DatePicker
	timePicker   material.TimePicker
	choosingDate bool
	choosingTime bool
}

func newInputs() inputsState {
	return inputsState{date: time.Now(), tags: []string{"North shore", "Cliffs", "Harbour"}, volume: 0.4}
}

// draw shows the fields, the choices and the pickers.
func (inputs *inputsState) draw(program *gallery, look material.Look) {
	context, scheme := look.Context, look.Roles
	program.largeTopBar(look, "Inputs", "Fields, choices and pickers", "inputs")
	page(context, "inputs", func() {
		inputs.searchLighthouses(program, look)
		look.Heading("Text fields")
		reachError := ""
		for _, character := range inputs.reach {
			if character < '0' || character > '9' {
				reachError = "Nautical miles, in digits"
				break
			}
		}
		look.TextField("name", &inputs.name, material.Field{Label: "Name", Placeholder: "Harbour Light", Supporting: "As the chart writes it",
			MaxLength: 40, Counter: true, Filled: true})
		look.TextField("reach", &inputs.reach, material.Field{Label: "Reach", Placeholder: "18", Supporting: "How far out it is seen",
			Error: reachError, MaxLength: 3, Suffix: "miles"})
		look.Dropdown("light", "Kind of light", []string{"Fixed", "Flashing", "Occulting", "Isophase"}, &inputs.light)
		look.TextField("password", &inputs.password, material.Field{Label: "Keeper's password", Secret: true, Reveal: &inputs.revealed})
		look.Heading("Choices")
		look.Checkbox("Lit every night", &inputs.nightly)
		for index, kind := range []string{"Fixed light", "Flashing", "Occulting"} {
			if look.RadioButton(kind, index == inputs.kind) {
				inputs.kind = index
			}
		}
		dark := program.isDark()
		look.ListItem("dark theme", material.Item{Icon: "moon", Headline: "Dark theme", Supporting: "Follows the system until chosen here",
			Trailing: func() {
				if look.Switch("dark switch", &dark) {
					program.setDark(dark)
				}
			}})
		program.colourChoice(look)
		look.SegmentedButtons("view", []string{"Day", "Week", "Month"}, &inputs.view)
		row(context, func() {
			for index, label := range []string{"Lit", "Dark", "Far out"} {
				if look.FilterChip(label, inputs.filters[index]) {
					inputs.filters[index] = !inputs.filters[index]
				}
			}
		})
		row(context, func() {
			for index := 0; index < len(inputs.tags); index++ {
				if _, removed := look.InputChip(inputs.tags[index]); removed {
					inputs.tags = append(inputs.tags[:index], inputs.tags[index+1:]...)
					index--
				}
			}
		})
		look.Context.Box(uiRow(12), func() {
			context.IconSized("volume-2", 24, scheme.OnSurfaceVariant)
			look.Slider("volume", &inputs.volume, 0, 1, 0)
		})
		look.Heading("Pickers")
		look.Text("Last inspection", 14, draw2d.Medium, scheme.OnSurfaceVariant)
		row(context, func() {
			if look.OutlinedButton("calendar", inputs.date.Format("2 Jan 2006")) {
				inputs.choosingDate = true
			}
			if look.OutlinedButton("clock", inputs.date.Format("15:04")) {
				inputs.choosingTime = true
			}
		})
	})
	if inputs.choosingDate {
		inputs.choosingDate = look.DatePicker("inspection day", &inputs.datePicker, &inputs.date) == material.Open
	}
	if inputs.choosingTime {
		inputs.choosingTime = look.TimePicker("inspection time", &inputs.timePicker, &inputs.date) == material.Open
	}
}

// colourChoice chooses where the theme's colours come from: the wallpaper, or one accent that
// Bubble makes into every colour.
func (program *gallery) colourChoice(look material.Look) {
	context, scheme := look.Context, look.Roles
	context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.PadXY(16, 4), Gap: 8, AlignY: ui.Center}, func() {
		context.Box(ui.Layout{Width: ui.Grow()}, func() { look.Text("Colours", 16, draw2d.Regular, scheme.OnSurface) })
		for index := range len(seeds) + 1 {
			chosen := program.colour == index
			fill, icon := scheme.SurfaceContainerHighest, "image" // the wallpaper's
			if index > 0 {
				fill, icon = draw2d.MustHex(seeds[index-1].Hex()), ""
			}
			if chosen {
				icon = "check"
			}
			ring := ui.Look{Fill: fill, Radius: 20, Border: scheme.Outline, BorderWidth: 1}
			if chosen {
				ring.Border, ring.BorderWidth = scheme.Primary, 3
			}
			layout := ui.Layout{Width: ui.Fixed(40), Height: ui.Fixed(40), AlignX: ui.Center, AlignY: ui.Center}
			if context.Clickable(fmt.Sprint("colour ", index), layout, ui.Looks{Normal: ring, Hovered: ring, Pressed: ring}, func() {
				if icon != "" {
					context.IconSized(icon, 20, readable(fill))
				}
			}) {
				program.colour = index
			}
		}
	})
}

// readable is black or white, whichever shows better on a colour.
func readable(colour draw2d.Color) draw2d.Color {
	if 0.299*colour[0]+0.587*colour[1]+0.114*colour[2] > 0.6 {
		return draw2d.Color{0, 0, 0, 0.87}
	}
	return draw2d.Color{1, 1, 1, 1}
}

// searchLighthouses is the search at the top of the inputs: it opens over the screen and finds
// lighthouses by name or place as they are typed; touching one goes to the lists.
func (inputs *inputsState) searchLighthouses(program *gallery, look material.Look) {
	query := &inputs.search
	look.SearchView("search", "Search the lighthouses", query, func() {
		if query.Query == "" {
			look.SectionTitle("All the lighthouses")
		}
		wanted := strings.ToLower(query.Query)
		for index, light := range program.lists.lighthouses {
			if !strings.Contains(strings.ToLower(light.name+" "+light.place), wanted) {
				continue
			}
			picture := coast()[index%len(coast())]
			if look.ListItem(light.name, material.Item{Picture: picture, Headline: light.name, Supporting: light.place}) {
				query.Active, program.screen = false, listsScreen
			}
		}
	})
}
