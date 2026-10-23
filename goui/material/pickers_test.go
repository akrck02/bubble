package material

import (
	"testing"
	"time"

	"github.com/akrck02/go-renderer/input"
)

func TestADatePickerKeepsTheDayChosenOnlyWhenConfirmed(t *testing.T) {
	date := time.Date(2026, time.March, 10, 9, 30, 0, 0, time.UTC)
	var picker DatePicker
	answer := Open
	shown := true
	screen := newScreen(t, func(look Look) {
		if shown {
			answer = look.DatePicker("when", &picker, &date)
			shown = answer == Open
		}
	})
	screen.frame()
	screen.tap("when", "icon chevron-right") // to April
	screen.tap("when", "day 2026-04-21")
	if !date.Equal(time.Date(2026, time.March, 10, 9, 30, 0, 0, time.UTC)) {
		t.Fatal("the date changed before it was confirmed")
	}
	screen.tap("when", "button OK")
	if answer != Confirmed || !date.Equal(time.Date(2026, time.April, 21, 9, 30, 0, 0, time.UTC)) {
		t.Fatalf("answer %v, date %v: want Confirmed and 21 April, keeping the time of day", answer, date)
	}
}

func TestATimeIsChosenOnTheDial(t *testing.T) {
	date := time.Date(2026, time.March, 10, 9, 30, 0, 0, time.UTC)
	var picker TimePicker
	answer := Open
	screen := newScreen(t, func(look Look) {
		if answer == Open {
			answer = look.TimePicker("when", &picker, &date)
		}
	})
	screen.frame()
	screen.frame()
	dial := picker.dial
	middleX, middleY := float64(dial.X+dial.Width/2), float64(dial.Y+dial.Height/2)
	tapAt := func(x, y float64) {
		screen.state.Feed().MoveCursor(x, y)
		screen.state.Feed().SetButton(input.MouseLeft, true)
		screen.frame()
		screen.state.Feed().SetButton(input.MouseLeft, false)
		screen.frame()
	}
	tapAt(middleX+float64(dial.Width)*0.4, middleY) // three o'clock, outer ring
	if picker.hour != 3 || !picker.minutes {
		t.Fatalf("hour %d, choosing minutes %v; want 3, and on to the minutes", picker.hour, picker.minutes)
	}
	tapAt(middleX, middleY+float64(dial.Width)*0.4) // half past
	if picker.minute != 30 {
		t.Fatalf("minute %d; want 30", picker.minute)
	}
	screen.tap("when", "button OK")
	if answer != Confirmed || date.Hour() != 3 || date.Minute() != 30 || date.Day() != 10 {
		t.Fatalf("answer %v, date %v: want 03:30 the same day", answer, date)
	}
}
