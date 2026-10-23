package material

import (
	"testing"

	"github.com/akrck02/go-renderer/draw2d"
)

func TestAFilmIsFittedInItsBox(t *testing.T) {
	picture := draw2d.NewBlankImage(160, 90) // sixteen by nine
	shown := fitted(picture, draw2d.Rectangle{X: 0, Y: 0, Width: 320, Height: 320})
	if shown != (draw2d.Rectangle{X: 0, Y: 70, Width: 320, Height: 180}) {
		t.Fatalf("a wide film in a square shows at %+v; want across it, in its middle", shown)
	}
	shown = fitted(picture, draw2d.Rectangle{X: 10, Y: 0, Width: 400, Height: 90})
	if shown != (draw2d.Rectangle{X: 130, Y: 0, Width: 160, Height: 90}) {
		t.Fatalf("a film in a wider box shows at %+v; want its full height, in the middle", shown)
	}
}
