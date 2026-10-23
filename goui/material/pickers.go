package material

import (
	"fmt"
	"math"
	"time"

	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// DateNames are the words a date picker shows, for a program in another language than English.
type DateNames struct {
	Title    string     // over the date chosen
	Months   [12]string // January first
	Weekdays [7]string  // one or two letters each, Monday first
	Confirm  string     // the action that keeps the date
	Dismiss  string     // the one that leaves it
}

// EnglishDates are the words a date picker shows by default.
var EnglishDates = DateNames{Title: "Select date",
	Months:   [12]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
	Weekdays: [7]string{"M", "T", "W", "T", "F", "S", "S"}, Confirm: "OK", Dismiss: "Cancel"}

// DatePicker is what a date picker remembers while shown: the month on show and the day chosen.
type DatePicker struct {
	Names DateNames // the zero value shows EnglishDates
	month time.Time // the first of the month on show
	day   time.Time // the day chosen, not yet kept
	open  bool
}

// DatePicker asks for a day in a dialog: a month at a time, with arrows to the months around it.
// It starts at the date given; on Confirmed it writes the day chosen there. The program shows it
// with an if while it wants an answer, as Dialog.
func (look Look) DatePicker(name string, picker *DatePicker, date *time.Time) Answer {
	names := picker.Names
	if names.Confirm == "" {
		names = EnglishDates
	}
	if !picker.open {
		picker.open, picker.day = true, dayOf(*date)
		picker.month = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	}
	roles := look.Roles
	answer := look.Dialog(name, Dialog{Title: names.Title, Confirm: names.Confirm, Dismiss: names.Dismiss, Content: func() {
		look.Text(fmt.Sprintf("%d %s %d", picker.day.Day(), names.Months[picker.day.Month()-1], picker.day.Year()), 28, draw2d.Regular, roles.OnSurface)
		look.Divider()
		look.monthBar(picker, names)
		look.monthGrid(picker, names)
	}})
	if answer == Confirmed {
		*date = time.Date(picker.day.Year(), picker.day.Month(), picker.day.Day(), date.Hour(), date.Minute(), date.Second(), 0, date.Location())
	}
	if answer != Open {
		picker.open = false
	}
	return answer
}

// monthBar names the month on show, between arrows to the one before and the one after.
func (look Look) monthBar(picker *DatePicker, names DateNames) {
	look.Context.Box(ui.Layout{Width: ui.Grow(), AlignY: ui.Center}, func() {
		look.Context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.Padding{Left: 12}}, func() {
			look.Text(fmt.Sprintf("%s %d", names.Months[picker.month.Month()-1], picker.month.Year()), 14, draw2d.Medium, look.Roles.OnSurfaceVariant)
		})
		if look.NamedIconButton("chevron-left", look.words().PreviousMonth) {
			picker.month = picker.month.AddDate(0, -1, 0)
		}
		if look.NamedIconButton("chevron-right", look.words().NextMonth) {
			picker.month = picker.month.AddDate(0, 1, 0)
		}
	})
}

// monthGrid lays the month's days under the weekdays, Monday first, the chosen one filled and
// today ringed.
func (look Look) monthGrid(picker *DatePicker, names DateNames) {
	roles := look.Roles
	today := dayOf(time.Now())
	week := func(content func()) { look.Context.Box(ui.Layout{Width: ui.Grow()}, content) }
	cell := ui.Layout{Width: ui.Share(1), Height: ui.Fixed(44), AlignX: ui.Center, AlignY: ui.Center}
	week(func() {
		for _, weekday := range names.Weekdays {
			look.Context.Box(cell, func() { look.Text(weekday, 14, draw2d.Regular, roles.OnSurface) })
		}
	})
	blanks := (int(picker.month.Weekday()) + 6) % 7 // Monday first
	days := picker.month.AddDate(0, 1, -1).Day()
	for start := 1 - blanks; start <= days; start += 7 {
		week(func() {
			for day := start; day < start+7; day++ {
				if day < 1 || day > days {
					look.Context.Box(cell, nil)
					continue
				}
				date := picker.month.AddDate(0, 0, day-1)
				if look.dayCell(date, date.Equal(picker.day), date.Equal(today), cell) {
					picker.day = date
				}
			}
		})
	}
}

func (look Look) dayCell(date time.Time, chosen, today bool, cell ui.Layout) bool {
	roles := look.Roles
	circle, ink := ui.Look{Radius: 20}, roles.OnSurface
	switch {
	case chosen:
		circle, ink = ui.Look{Fill: roles.Primary, Radius: 20}, roles.OnPrimary
	case today:
		circle, ink = ui.Look{Border: roles.Primary, BorderWidth: 1, Radius: 20}, roles.Primary
	}
	touched := false
	look.Context.Box(cell, func() {
		layout := ui.Layout{Width: ui.Fixed(40), Height: ui.Fixed(40), AlignX: ui.Center, AlignY: ui.Center}
		touched = look.Context.Clickable(date.Format("day 2006-01-02"), layout, pressedOnly(circle, stateLayer(roles.OnSurface, 0.12)), func() {
			look.Text(fmt.Sprint(date.Day()), 14, draw2d.Regular, ink)
		})
		look.Context.Describe(ui.Description{Selected: chosen})
	})
	return touched
}

// dayOf is a date's day, at midnight, in its own zone.
func dayOf(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

// TimeNames are the words a time picker shows, for a program in another language than English.
type TimeNames struct {
	Title   string
	Confirm string
	Dismiss string
}

// EnglishTimes are the words a time picker shows by default.
var EnglishTimes = TimeNames{Title: "Select time", Confirm: "OK", Dismiss: "Cancel"}

// TimePicker is what a time picker remembers while shown: the time chosen, not yet kept, and
// whether the dial chooses the hour or the minute.
type TimePicker struct {
	Names   TimeNames // the zero value shows EnglishTimes
	hour    int
	minute  int
	minutes bool // the dial shows minutes
	open    bool
	dial    draw2d.Rectangle // where the dial was painted in the last frame
}

// TimePicker asks for a time of day in a dialog, on Material's dial: the hour first, on a ring of
// 1 to 12 with 13 to 00 inside it, then the minute, by touching or dragging round the dial. The
// hour and the minute at the top choose which the dial shows. It starts at the time given; on
// Confirmed it writes the time chosen there, keeping the day.
func (look Look) TimePicker(name string, picker *TimePicker, date *time.Time) Answer {
	names := picker.Names
	if names.Confirm == "" {
		names = EnglishTimes
	}
	if !picker.open {
		picker.open, picker.hour, picker.minute, picker.minutes = true, date.Hour(), date.Minute(), false
	}
	answer := look.Dialog(name, Dialog{Title: names.Title, Confirm: names.Confirm, Dismiss: names.Dismiss, Content: func() {
		look.timeHeader(picker)
		look.dial(picker)
	}})
	if answer == Confirmed {
		*date = time.Date(date.Year(), date.Month(), date.Day(), picker.hour, picker.minute, 0, 0, date.Location())
	}
	if answer != Open {
		picker.open = false
	}
	return answer
}

// timeHeader shows the hour and the minute chosen, in large boxes that choose what the dial shows.
func (look Look) timeHeader(picker *TimePicker) {
	roles := look.Roles
	box := func(text string, chosen bool) bool {
		normal, ink := ui.Look{Fill: roles.SurfaceContainerHighest, Radius: 8}, roles.OnSurface
		if chosen {
			normal, ink = ui.Look{Fill: roles.PrimaryContainer, Radius: 8}, roles.OnPrimaryContainer
		}
		layout := ui.Layout{Width: ui.Fixed(96), Height: ui.Fixed(80), AlignX: ui.Center, AlignY: ui.Center}
		touched := look.Context.Clickable("time "+text, layout, pressedOnly(normal, stateLayer(ink, 0.12)), func() {
			look.Text(text, 48, draw2d.Regular, ink)
		})
		look.Context.Describe(ui.Description{Selected: chosen})
		return touched
	}
	look.Context.Box(ui.Layout{Width: ui.Grow(), Gap: 8, AlignX: ui.Center, AlignY: ui.Center}, func() {
		if box(fmt.Sprintf("%02d", picker.hour), !picker.minutes) {
			picker.minutes = false
		}
		look.Text(":", 48, draw2d.Regular, roles.OnSurface)
		if box(fmt.Sprintf("%02d", picker.minute), picker.minutes) {
			picker.minutes = true
		}
	})
}

// dialSize is the dial's side.
const dialSize = 256

// dial is the clock face: numbers round it, a hand to the value chosen, and a finger anywhere on it
// choosing. Letting go of an hour goes on to the minutes.
func (look Look) dial(picker *TimePicker) {
	context := look.Context
	layout := ui.Layout{Width: ui.Fixed(dialSize), Height: ui.Fixed(dialSize)}
	context.Box(ui.Layout{Width: ui.Grow(), AlignX: ui.Center}, func() {
		touched := context.Clickable("dial", layout, ui.Looks{}, func() {
			context.Paint(ui.Layout{Width: ui.Grow(), Height: ui.Grow()}, func(area draw2d.Rectangle, drawing *draw2d.List) {
				look.paintDial(area, drawing, picker)
			})
		})
		// a finger on the dial chooses as it moves; letting go of an hour goes on to the minutes
		pointerX, pointerY := context.Pointer()
		if picker.dial.Width > 0 && picker.dial.Contains(pointerX, pointerY) && (context.Pressing() || touched) {
			picker.choose(picker.dial, pointerX, pointerY)
		}
		if touched && !picker.minutes {
			picker.minutes = true
		}
	})
}

// choose sets the hour or the minute from where a finger is on the dial.
func (picker *TimePicker) choose(area draw2d.Rectangle, x, y float32) {
	middleX, middleY := float64(area.X+area.Width/2), float64(area.Y+area.Height/2)
	angle := math.Atan2(float64(x)-middleX, middleY-float64(y)) // 0 at twelve o'clock, clockwise
	if angle < 0 {
		angle += 2 * math.Pi
	}
	if picker.minutes {
		picker.minute = int(math.Round(angle/(2*math.Pi)*60)) % 60
		return
	}
	position := int(math.Round(angle/(2*math.Pi)*12)) % 12
	inner := math.Hypot(float64(x)-middleX, float64(y)-middleY) < float64(area.Width)*0.3
	switch {
	case inner && position == 0:
		picker.hour = 0
	case inner:
		picker.hour = position + 12
	case position == 0:
		picker.hour = 12
	default:
		picker.hour = position
	}
}

// paintDial draws the face, the numbers and the hand.
func (look Look) paintDial(area draw2d.Rectangle, drawing *draw2d.List, picker *TimePicker) {
	roles := look.Roles
	picker.dial = area
	drawing.Rectangle(area, area.Width/2, roles.SurfaceContainerHighest, draw2d.Transparent, 0)
	middleX, middleY := float64(area.X+area.Width/2), float64(area.Y+area.Height/2)
	outer, inner := float64(area.Width)*0.4, float64(area.Width)*0.26
	// the hand: from the middle to the value chosen, ending in a circle round it
	var angle, reach float64
	if picker.minutes {
		angle, reach = float64(picker.minute)/60*2*math.Pi, outer
	} else {
		angle, reach = float64(picker.hour%12)/12*2*math.Pi, outer
		if picker.hour == 0 || picker.hour > 12 {
			reach = inner
		}
	}
	handX, handY := middleX+reach*math.Sin(angle), middleY-reach*math.Cos(angle)
	for step := 0.0; step <= 1; step += 0.02 {
		dot(drawing, middleX+(handX-middleX)*step, middleY+(handY-middleY)*step, 2, roles.Primary)
	}
	dot(drawing, middleX, middleY, 8, roles.Primary)
	dot(drawing, handX, handY, 48, roles.Primary)
	number := func(text string, angle, distance float64, chosen bool) {
		x, y := middleX+distance*math.Sin(angle), middleY-distance*math.Cos(angle)
		ink := roles.OnSurface
		if chosen {
			ink = roles.OnPrimary
		}
		size := float32(16)
		if distance == inner {
			size = 14
		}
		width := float32(len(text)) * size * 0.55 // close enough to centre digits
		drawing.Text(text, float32(x)-width/2, float32(y)+size*0.35, size, draw2d.Regular, ink)
	}
	if picker.minutes {
		for minute := 0; minute < 60; minute += 5 {
			number(fmt.Sprintf("%02d", minute), float64(minute)/60*2*math.Pi, outer, minute == picker.minute)
		}
		return
	}
	for hour := 1; hour <= 12; hour++ {
		number(fmt.Sprint(hour), float64(hour%12)/12*2*math.Pi, outer, hour == picker.hour)
		innerHour := (hour + 12) % 24
		number(fmt.Sprintf("%02d", innerHour), float64(hour%12)/12*2*math.Pi, inner, innerHour == picker.hour)
	}
}
