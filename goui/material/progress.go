package material

import (
	"math"
	"time"

	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// Loading indicators as Material 3 Expressive draws them: a filled shape that turns and morphs
// from one form to the next, circular and linear progress that can run as a wave, and skeletons
// standing in for content still coming. They are painted with rectangles (go-ui's Paint): a shape
// is filled row by row and edged with small dots, a wave is a line of dots. While shown they ask
// for frames; a screen without them sleeps.

// frameDelay is how soon a moving indicator asks for its next frame: about sixty a second.
const frameDelay = 16 * time.Millisecond

// shape is a closed form around a middle, as its distance from the middle at an angle, for a
// radius of 1.
type shape func(angle float64) float64

// lobed is a form with lobes all around: a cookie, a sun, a burst. depth is how far in the gaps go.
func lobed(lobes int, depth float64) shape {
	return func(angle float64) float64 {
		return 1 - depth*(0.5-0.5*math.Cos(float64(lobes)*angle))
	}
}

// oval is an ellipse as wide as width, for a height of 1.
func oval(width float64) shape {
	return func(angle float64) float64 {
		cosine, sine := math.Cos(angle), math.Sin(angle)
		return width / math.Sqrt(cosine*cosine+width*width*sine*sine)
	}
}

// indicatorShapes are the forms the loading indicator goes through, in Material's order: a soft
// burst, a cookie, a pentagon, a pill, a sun, a four-lobed cookie and an oval.
var indicatorShapes = []shape{lobed(10, 0.22), lobed(9, 0.14), lobed(5, 0.12), oval(0.62), lobed(8, 0.1), lobed(4, 0.18), oval(0.8)}

// LoadingIndicator is Material's indicator for a short wait of unknown length: a filled shape in
// the primary colour, turning and morphing from one form to the next. size is its side, 48 by
// default in Material.
func (look Look) LoadingIndicator(size float32) {
	look.loading(size, false)
}

// ContainedLoadingIndicator is the indicator on a round container, to stand out over content (at
// the top of a list being refreshed).
func (look Look) ContainedLoadingIndicator(size float32) {
	look.loading(size, true)
}

func (look Look) loading(size float32, contained bool) {
	scheme := look.Roles
	elapsed := time.Since(loadingEpoch).Seconds()
	// each form holds for a moment, then eases into the next; the whole turns meanwhile
	const formSeconds = 0.65
	step := elapsed / formSeconds
	index := int(step) % len(indicatorShapes)
	morph := easeInOut(math.Min((step-math.Floor(step))*1.6, 1))
	from, to := indicatorShapes[index], indicatorShapes[(index+1)%len(indicatorShapes)]
	turn := elapsed*2*math.Pi/4.5 + (float64(int(step))+morph)*math.Pi/7
	fill, scale := scheme.Primary, 0.38
	if contained {
		fill, scale = scheme.OnPrimaryContainer, 0.3
	}
	look.Context.Paint(ui.Layout{Width: ui.Fixed(size), Height: ui.Fixed(size)}, func(area draw2d.Rectangle, drawing *draw2d.List) {
		if contained {
			drawing.Rectangle(area, size/2, scheme.PrimaryContainer, draw2d.Transparent, 0)
		}
		middleX, middleY := float64(area.X+area.Width/2), float64(area.Y+area.Height/2)
		radius := float64(size) * scale
		form := func(angle float64) float64 {
			return radius * ((1-morph)*from(angle-turn) + morph*to(angle-turn))
		}
		fillShape(drawing, middleX, middleY, form, fill)
	})
	look.Context.RequestFrame(frameDelay)
}

// loadingEpoch keeps every indicator on one clock, so two on a screen move together.
var loadingEpoch = time.Now()

// fillShape fills a form around a middle row by row, and edges it with small dots so its outline
// reads smooth rather than stepped.
func fillShape(drawing *draw2d.List, middleX, middleY float64, form shape, colour draw2d.Color) {
	const points = 120
	var outline [points][2]float64
	top, bottom := math.Inf(1), math.Inf(-1)
	for index := range outline {
		angle := 2 * math.Pi * float64(index) / points
		distance := form(angle)
		x, y := middleX+distance*math.Cos(angle), middleY+distance*math.Sin(angle)
		outline[index] = [2]float64{x, y}
		top, bottom = math.Min(top, y), math.Max(bottom, y)
	}
	const row = 0.5
	for y := top; y < bottom; y += row {
		middle := y + row/2
		crossings := crossingsAt(outline[:], middle)
		for pair := 0; pair+1 < len(crossings); pair += 2 {
			left, right := crossings[pair], crossings[pair+1]
			drawing.Rectangle(draw2d.Rectangle{X: float32(left), Y: float32(y), Width: float32(right - left), Height: row + 0.1}, 0, colour, draw2d.Transparent, 0)
		}
	}
	for _, point := range outline {
		drawing.Rectangle(draw2d.Rectangle{X: float32(point[0] - 0.6), Y: float32(point[1] - 0.6), Width: 1.2, Height: 1.2}, 0.6, colour, draw2d.Transparent, 0)
	}
}

// crossingsAt is where a closed outline crosses a row, from left to right.
func crossingsAt(outline [][2]float64, y float64) []float64 {
	var crossings []float64
	for index := range outline {
		start, end := outline[index], outline[(index+1)%len(outline)]
		if (start[1] <= y) == (end[1] <= y) {
			continue
		}
		share := (y - start[1]) / (end[1] - start[1])
		crossings = append(crossings, start[0]+share*(end[0]-start[0]))
	}
	for index := 1; index < len(crossings); index++ { // a few, nearly sorted: insertion sort
		for back := index; back > 0 && crossings[back] < crossings[back-1]; back-- {
			crossings[back], crossings[back-1] = crossings[back-1], crossings[back]
		}
	}
	return crossings
}

func easeInOut(progress float64) float64 {
	if progress < 0.5 {
		return 4 * progress * progress * progress
	}
	left := -2*progress + 2
	return 1 - left*left*left/2
}

// dot paints a round dot of a diameter centred on a point.
func dot(drawing *draw2d.List, x, y, diameter float64, colour draw2d.Color) {
	radius := diameter / 2
	drawing.Rectangle(draw2d.Rectangle{X: float32(x - radius), Y: float32(y - radius), Width: float32(diameter), Height: float32(diameter)},
		float32(radius), colour, draw2d.Transparent, 0)
}

// Progress is how an indicator of known length shows its work: a fraction from 0 to 1, or
// Indeterminate for work of unknown length.
const Indeterminate = -1

// CircularProgress is a ring that fills as work gets done (fraction from 0 to 1), or, with
// Indeterminate, an arc chasing round. wavy draws the filled part as a wave, as Material 3
// Expressive does. size is its side, 48 by default in Material.
func (look Look) CircularProgress(fraction float32, size float32, wavy bool) {
	scheme := look.Roles
	elapsed := time.Since(loadingEpoch).Seconds()
	start, sweep := -math.Pi/2, 2*math.Pi*float64(min(max(fraction, 0), 1))
	if fraction < 0 {
		// the arc grows and shrinks as it goes round
		cycle := math.Mod(elapsed/1.333, 1)
		sweep = 0.1*math.Pi + 1.4*math.Pi*(0.5-0.5*math.Cos(2*math.Pi*cycle))
		start = elapsed*2*math.Pi/1.6 - math.Pi/2 + cycle*math.Pi
		look.Context.RequestFrame(frameDelay)
	}
	look.Context.Paint(ui.Layout{Width: ui.Fixed(size), Height: ui.Fixed(size)}, func(area draw2d.Rectangle, drawing *draw2d.List) {
		middleX, middleY := float64(area.X+area.Width/2), float64(area.Y+area.Height/2)
		const stroke = 4
		radius := float64(size)/2 - stroke - 2
		amplitude := 0.0
		if wavy {
			amplitude = 1.6
		}
		// the track: what is left to do, after a gap, when the work's length is known
		if fraction >= 0 && fraction < 1 {
			gap := 6 / radius
			arc(drawing, middleX, middleY, radius, start+sweep+gap, 2*math.Pi-sweep-2*gap, 0, elapsed, stroke, scheme.SecondaryContainer)
		}
		arc(drawing, middleX, middleY, radius, start, sweep, amplitude, elapsed, stroke, scheme.Primary)
	})
}

// arc paints an arc as a line of dots, waving out and in by amplitude when it is more than 0.
func arc(drawing *draw2d.List, middleX, middleY, radius, start, sweep, amplitude, elapsed, stroke float64, colour draw2d.Color) {
	if sweep <= 0 {
		return
	}
	steps := int(math.Ceil(radius * sweep / 1.2))
	for step := 0; step <= steps; step++ {
		angle := start + sweep*float64(step)/float64(steps)
		distance := radius + amplitude*math.Sin(10*angle-elapsed*6)
		dot(drawing, middleX+distance*math.Cos(angle), middleY+distance*math.Sin(angle), stroke, colour)
	}
}

// WavyProgress is a line that fills as work gets done (fraction from 0 to 1), the filled part a
// travelling wave and the rest a flat track ending in a stop dot, as Material 3 Expressive draws
// it; with Indeterminate, a wave runs along it.
func (look Look) WavyProgress(fraction float32) {
	scheme := look.Roles
	elapsed := time.Since(loadingEpoch).Seconds()
	look.Context.Paint(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(12)}, func(area draw2d.Rectangle, drawing *draw2d.List) {
		const stroke, amplitude, wavelength = 4.0, 3.0, 40.0
		left, right := float64(area.X)+stroke/2, float64(area.Right())-stroke/2
		middle := float64(area.Y + area.Height/2)
		wave := func(from, to float64) {
			for x := from; x <= to; x += 1 {
				dot(drawing, x, middle+amplitude*math.Sin(2*math.Pi*(x-left)/wavelength-elapsed*5), stroke, scheme.Primary)
			}
		}
		flat := func(from, to float64) {
			if to > from {
				drawing.Rectangle(draw2d.Rectangle{X: float32(from - stroke/2), Y: float32(middle - stroke/2), Width: float32(to - from + stroke), Height: stroke},
					stroke/2, scheme.SecondaryContainer, draw2d.Transparent, 0)
			}
		}
		length := right - left
		if fraction < 0 {
			// a stretch of wave runs from left to right, the track around it
			cycle := math.Mod(elapsed/1.75, 1)
			head := left + length*math.Min(cycle*1.4, 1)
			tail := left + length*math.Max(cycle*1.4-0.5, 0)
			flat(left, tail-8)
			wave(tail, head)
			flat(head+8, right)
			return
		}
		end := left + length*float64(min(max(fraction, 0), 1))
		wave(left, end)
		flat(end+8, right)
		if fraction < 1 {
			dot(drawing, right, middle, stroke, scheme.Primary) // the stop: where the work ends
		}
	})
	look.Context.RequestFrame(frameDelay)
}

// Skeleton stands in for content still coming (a picture, a line of text): a quiet block with a
// light passing across it. A row of them in the content's shape tells what will be there.
func (look Look) Skeleton(width, height ui.Sizing, radius float32) {
	scheme := look.Roles
	elapsed := time.Since(loadingEpoch).Seconds()
	base := scheme.SurfaceContainerHighest
	light := mix(base, scheme.Surface, 0.6)
	look.Context.Paint(ui.Layout{Width: width, Height: height}, func(area draw2d.Rectangle, drawing *draw2d.List) {
		drawing.Rectangle(area, radius, base, draw2d.Transparent, 0)
		// a soft band of light crossing every 1.5 seconds, made of strips fading out at its sides
		drawing.PushClip(area)
		defer drawing.PopClip()
		const band, strips = 80.0, 8
		window := float64(area.Width) + band*2
		position := float64(area.X) - band + math.Mod(elapsed/1.5, 1)*window
		for strip := 0; strip < strips; strip++ {
			share := float64(strip) / strips
			strength := float32(math.Sin(math.Pi * (share + 0.5/strips)))
			colour := draw2d.Color{light[0], light[1], light[2], strength * 0.7}
			drawing.Rectangle(draw2d.Rectangle{X: float32(position + share*band), Y: area.Y, Width: band/strips + 0.5, Height: area.Height}, 0, colour, draw2d.Transparent, 0)
		}
	})
	look.Context.RequestFrame(frameDelay)
}

// LinearProgress is a bar across that fills as work gets done: fraction from 0 to 1. A negative
// fraction is work of unknown length: a piece of the bar runs along it, and it asks for frames
// while shown.
func (look Look) LinearProgress(fraction float32) {
	scheme := look.Roles
	track := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(4)}
	if fraction >= 0 {
		look.Context.Frame(track, ui.Look{Fill: scheme.SecondaryContainer, Radius: 2}, func() {
			look.Context.Frame(ui.Layout{Width: ui.Share(max(min(fraction, 1), 0.001)), Height: ui.Grow()}, ui.Look{Fill: scheme.Primary, Radius: 2}, nil)
			look.Context.Box(ui.Layout{Width: ui.Share(max(1-fraction, 0.001))}, nil)
		})
		return
	}
	// a third of the bar sweeps across it every 1.5 seconds
	phase := float32(math.Mod(float64(time.Now().UnixMilli())/1500, 1))
	before := max(phase*1.3-0.3, 0)
	length := min(phase*1.3, 1) - before
	look.Context.Frame(track, ui.Look{Fill: scheme.SecondaryContainer, Radius: 2}, func() {
		look.Context.Box(ui.Layout{Width: ui.Share(before + 0.001)}, nil)
		look.Context.Frame(ui.Layout{Width: ui.Share(length + 0.001), Height: ui.Grow()}, ui.Look{Fill: scheme.Primary, Radius: 2}, nil)
		look.Context.Box(ui.Layout{Width: ui.Share(max(1-before-length, 0) + 0.001)}, nil)
	})
	look.Context.RequestFrame(16 * time.Millisecond)
}
