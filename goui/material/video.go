package material

import (
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
	"github.com/akrck02/go-ui/window"
)

// VideoPlayer shows a film (go-ui's window.Video) in a rounded box: whole, over black where their
// shapes differ, with a round play button over it while it is paused and a thin line at its foot
// of how far it has played. Touching it plays or pauses the film; it brings the film up to date
// itself (Video.Update). name tells players apart.
func (look Look) VideoPlayer(name string, video *window.Video, height float32) {
	video.Update(look.Context)
	layout := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(height)}
	looks := ui.Looks{} // the film is the look; a touch shows on the button
	touched := look.Context.Clickable(name, layout, looks, func() {
		look.Context.Paint(ui.Layout{Width: ui.Grow(), Height: ui.Grow()}, func(area draw2d.Rectangle, drawing *draw2d.List) {
			look.paintFilm(video, area, drawing)
		})
	})
	if !touched {
		return
	}
	if video.Playing() {
		video.Pause()
		return
	}
	video.Play()
}

const filmRadius = 16

// paintFilm draws the film's picture fitted in the area, the play button while it is paused, and
// how far it has played.
func (look Look) paintFilm(video *window.Video, area draw2d.Rectangle, drawing *draw2d.List) {
	scheme := look.Roles
	drawing.Rectangle(area, filmRadius, draw2d.Color{0, 0, 0, 1}, draw2d.Color{}, 0)
	if picture := video.Picture(); picture != nil {
		shown := fitted(picture, area)
		drawing.PushClip(area) // the corners the picture shares with the box stay round
		drawing.Image(picture, shown, filmRadius, draw2d.Color{1, 1, 1, 1})
		drawing.PopClip()
	}
	if !video.Playing() {
		drawing.Rectangle(area, filmRadius, draw2d.Color{0, 0, 0, 0.3}, draw2d.Color{}, 0)
		const button = 56
		middleX, middleY := area.X+area.Width/2, area.Y+area.Height/2
		drawing.Rectangle(draw2d.Rectangle{X: middleX - button/2, Y: middleY - button/2, Width: button, Height: button}, button/2,
			scheme.PrimaryContainer, draw2d.Color{}, 0)
		drawing.Icon("play", middleX-12, middleY-12, 24, scheme.OnPrimaryContainer)
	}
	if duration := video.Duration(); duration > 0 {
		played := min(float32(video.Position())/float32(duration), 1)
		line := draw2d.Rectangle{X: area.X + filmRadius, Y: area.Bottom() - 8, Width: (area.Width - 2*filmRadius) * played, Height: 3}
		drawing.Rectangle(line, 1.5, scheme.Primary, draw2d.Color{}, 0)
	}
}

// fitted is the largest area of the picture's shape inside an area, in its middle.
func fitted(picture *draw2d.Image, area draw2d.Rectangle) draw2d.Rectangle {
	width, height := picture.Size()
	pictureShape, areaShape := float32(width)/float32(height), area.Width/area.Height
	if pictureShape > areaShape {
		shownHeight := area.Width / pictureShape
		return draw2d.Rectangle{X: area.X, Y: area.Y + (area.Height-shownHeight)/2, Width: area.Width, Height: shownHeight}
	}
	shownWidth := area.Height * pictureShape
	return draw2d.Rectangle{X: area.X + (area.Width-shownWidth)/2, Y: area.Y, Width: shownWidth, Height: area.Height}
}
