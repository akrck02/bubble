package material

import (
	"fmt"
	"testing"
	"time"

	"github.com/akrck02/go-renderer/draw2d"
	"github.com/akrck02/go-renderer/input"
	ui "github.com/akrck02/go-ui"
)

// clocked gives the screen a clock that moves a frame at a time, as fingers' speeds need.
func (screen *screen) clocked() {
	now := time.Unix(1000, 0)
	screen.state.Now = func() time.Time { return now }
	frame := screen.build
	screen.build = func(look Look) {
		now = now.Add(16 * time.Millisecond)
		frame(look)
	}
}

// drag puts a finger down, moves it in ten frames, and lifts it when told.
func (screen *screen) drag(fromX, fromY, byX, byY float64, lift bool) {
	screen.state.Feed().SetTouch(1, fromX, fromY)
	screen.state.Feed().MoveCursor(fromX, fromY)
	screen.state.Feed().SetButton(input.MouseLeft, true)
	screen.frame()
	for step := 1; step <= 10; step++ {
		x, y := fromX+byX*float64(step)/10, fromY+byY*float64(step)/10
		screen.state.Feed().SetTouch(1, x, y)
		screen.state.Feed().MoveCursor(x, y)
		screen.frame()
	}
	if lift {
		screen.state.Feed().SetButton(input.MouseLeft, false)
		screen.state.Feed().LiftTouch(1)
		screen.frame()
	}
}

func rows(look Look, count int) {
	for index := range count {
		look.ListItem(fmt.Sprint("row ", index), Item{Headline: fmt.Sprint("Light ", index)})
	}
}

func TestPullingAListDownAndLettingGoRefreshesIt(t *testing.T) {
	refreshes := 0
	screen := newScreen(t, func(look Look) {
		if look.PullToRefresh("lights", false, ui.Layout{Width: ui.Grow(), Height: ui.Fixed(600)}, func() { rows(look, 30) }) {
			refreshes++
		}
	})
	screen.clocked()
	screen.frame()
	screen.frame()
	screen.drag(200, 100, 0, 60, true) // too short a pull
	if refreshes != 0 {
		t.Fatal("a short pull refreshed the list")
	}
	screen.drag(200, 100, 0, 300, true)
	screen.frame()
	if refreshes != 1 {
		t.Fatalf("%d refreshes after a long pull; want one", refreshes)
	}
}

func TestARowSwipedFarGoesAndOneSwipedLittleStays(t *testing.T) {
	dismissed := 0
	screen := newScreen(t, func(look Look) {
		if look.SwipeToDismiss("row", func() { look.ListItem("light", Item{Headline: "Harbour Light"}) }) {
			dismissed++
		}
	})
	screen.clocked()
	screen.frame()
	screen.frame()
	screen.drag(200, 30, 60, 2, true)
	for range 20 {
		screen.frame()
	}
	if dismissed != 0 {
		t.Fatal("a short swipe took the row away")
	}
	screen.drag(200, 30, -260, 2, true)
	if dismissed != 1 {
		t.Fatalf("%d dismissals after a long swipe; want one", dismissed)
	}
}

func TestPagesTurnWithAFingerAndWithTheirTabs(t *testing.T) {
	var pager Pager
	screen := newScreen(t, func(look Look) {
		look.TabPager("views", []string{"Map", "List", "Notes"}, &pager, func(index int) { rows(look, 3) })
	})
	screen.clocked()
	screen.frame()
	screen.frame()
	screen.drag(300, 300, -250, 5, true)
	if pager.Selected != 1 {
		t.Fatalf("page %d after swiping left; want the next one", pager.Selected)
	}
	screen.tap("views", "Notes")
	if pager.Selected != 2 {
		t.Fatalf("page %d after touching Notes; want it", pager.Selected)
	}
	for range 40 {
		screen.frame()
	}
	if pager.shown != 2 {
		t.Fatalf("pages shown at %v; want them settled on Notes", pager.shown)
	}
}

func TestACarouselSettlesOnAnItem(t *testing.T) {
	var carousel Carousel
	screen := newScreen(t, func(look Look) {
		look.Carousel("places", &carousel, 8, 160, func(index int) {
			look.CarouselItem(200, look.Roles.PrimaryContainer, "", fmt.Sprint("Place ", index))
		})
	})
	screen.clocked()
	screen.frame()
	screen.frame()
	screen.drag(300, 100, -250, 3, true)
	if carousel.rest < 1 || carousel.rest != float32(int(carousel.rest)) {
		t.Fatalf("rests at %v after a swipe; want a later item's edge", carousel.rest)
	}
}

func TestAMultiBrowseCarouselShrinksItemsTowardsItsEdges(t *testing.T) {
	var carousel Carousel
	widths := map[int]float32{}
	screen := newScreen(t, func(look Look) {
		clear(widths)
		look.MultiBrowseCarousel("bay", &carousel, 6, 200, 180, func(index int, width float32) { widths[index] = width })
	})
	screen.clocked()
	screen.frame()
	screen.frame()
	if _, past := widths[3]; past || widths[0] != 200 || widths[1] != 96 || widths[2] != smallItem {
		t.Fatalf("widths at rest %v; want a large item, a medium one filling what is left, and a small one", widths)
	}
	screen.drag(300, 90, -208, 0, false)
	for range 15 { // the finger stays still before it lifts, so it is no flick
		screen.frame()
	}
	screen.state.Feed().SetButton(input.MouseLeft, false)
	screen.state.Feed().LiftTouch(1)
	for range 30 {
		screen.frame()
	}
	if carousel.rest != 1 || widths[0] != smallItem || widths[1] != 200 || widths[2] != 96 {
		t.Fatalf("resting at %v with widths %v after sliding one item; want the first shrunk, the second large", carousel.rest, widths)
	}
}

func TestKeylinesKeepMoreLargeItemsOnAWideScreen(t *testing.T) {
	keylines := browseKeylines(860, 200)
	if len(keylines) != 5 || keylines[2] != 200 || keylines[3] != 172 || keylines[4] != smallItem {
		t.Fatalf("keylines %v on 860 pixels; want three large items, a medium of 172 and a small one", keylines)
	}
	keylines = browseKeylines(900, 200) // a medium of 4 pixels: the third large item takes its place
	if len(keylines) != 5 || keylines[3] != 200 {
		t.Fatalf("keylines %v on 900 pixels; want the medium item no narrower than the smallest", keylines)
	}
}

func TestACarouselPictureWritesItsTitleOverAVeil(t *testing.T) {
	photo := draw2d.NewBlankImage(40, 30)
	screen := newScreen(t, func(look Look) { look.CarouselPicture(photo, 180, "North Cape") })
	screen.frame()
	screen.frame()
	var image, veil, title bool
	for _, command := range screen.context.Drawing().Commands {
		switch {
		case command.Kind == draw2d.ImageKind && command.Picture == photo:
			image = true
		case command.Kind == draw2d.RectangleKind && command.Fill[3] > 0.3 && command.Fill[0] == 0 && command.HasClip:
			veil = image // drawn after the photo, over it
		case command.Kind == draw2d.TextKind && command.Text == "North Cape":
			title = veil
		}
	}
	if !image || !veil || !title {
		t.Fatalf("photo %v, veil over it %v, title over that %v; want all three in that order", image, veil, title)
	}
}
