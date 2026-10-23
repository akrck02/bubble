package main

import (
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/akrck02/go-renderer/draw2d"
)

// The gallery's pictures are painted here rather than shipped: a coast with a lighthouse at
// different times of day, so the images pieces have something to show without photos of anyone's.

// daylight is a time of day's colours: the sky from its top to the horizon, the sun, the sea.
type daylight struct {
	skyTop, horizon, sun, seaNear, seaFar color.NRGBA
	sunHeight                             float64 // 0 at the horizon, 1 at the top of the picture
	night                                 bool    // the lamp shines
}

var daylights = []daylight{
	{skyTop: rgb(70, 110, 190), horizon: rgb(250, 180, 150), sun: rgb(255, 220, 160), seaNear: rgb(40, 60, 100), seaFar: rgb(170, 140, 160), sunHeight: 0.1},
	{skyTop: rgb(60, 140, 220), horizon: rgb(190, 225, 250), sun: rgb(255, 250, 220), seaNear: rgb(20, 80, 130), seaFar: rgb(90, 160, 200), sunHeight: 0.75},
	{skyTop: rgb(90, 60, 130), horizon: rgb(250, 140, 80), sun: rgb(255, 190, 110), seaNear: rgb(50, 40, 70), seaFar: rgb(200, 110, 90), sunHeight: 0.2},
	{skyTop: rgb(10, 15, 40), horizon: rgb(40, 60, 110), sun: rgb(230, 230, 240), seaNear: rgb(5, 10, 25), seaFar: rgb(30, 45, 80), sunHeight: 0.6, night: true},
	{skyTop: rgb(120, 150, 170), horizon: rgb(210, 215, 215), sun: rgb(240, 240, 235), seaNear: rgb(60, 80, 90), seaFar: rgb(150, 165, 170), sunHeight: 0.5},
	{skyTop: rgb(40, 90, 170), horizon: rgb(140, 200, 240), sun: rgb(255, 245, 200), seaNear: rgb(10, 70, 110), seaFar: rgb(60, 140, 180), sunHeight: 0.9},
}

func rgb(red, green, blue uint8) color.NRGBA { return color.NRGBA{R: red, G: green, B: blue, A: 255} }

// coast is the gallery's pictures, painted the first time they are needed.
var coast = sync.OnceValue(coastPictures)

// coastPictures paints one picture for each time of day.
func coastPictures() []*draw2d.Image {
	pictures := make([]*draw2d.Image, len(daylights))
	for index, light := range daylights {
		pictures[index] = draw2d.NewImage(paintCoast(light, index, 480, 360))
	}
	return pictures
}

// paintCoast paints the sky, the sun, the sea with its waves, a headland and a lighthouse on it.
// seed moves the headland from picture to picture.
func paintCoast(light daylight, seed, width, height int) *image.NRGBA {
	picture := image.NewNRGBA(image.Rect(0, 0, width, height))
	horizon := float64(height) * 0.58
	sunX, sunY := float64(width)*(0.25+0.1*float64(seed%3)), horizon*(1-light.sunHeight)
	headland := float64(width) * (0.62 + 0.06*float64(seed%2))
	for y := range height {
		for x := range width {
			shade := skyOrSea(light, float64(x), float64(y), horizon, height)
			if math.Hypot(float64(x)-sunX, float64(y)-sunY) < 22 && float64(y) < horizon {
				shade = light.sun
			}
			if onHeadland(float64(x), float64(y), headland, horizon, width) {
				shade = mixed(light.seaNear, color.NRGBA{A: 255}, 0.55)
			}
			if tower, lit := onLighthouse(float64(x), float64(y), headland, horizon, light.night); tower {
				shade = lit
			}
			picture.SetNRGBA(x, y, shade)
		}
	}
	return picture
}

// skyOrSea is the sky's gradient above the horizon, and the sea's below it with bands of waves.
func skyOrSea(light daylight, x, y, horizon float64, height int) color.NRGBA {
	if y < horizon {
		return mixed(light.skyTop, light.horizon, y/horizon)
	}
	depth := (y - horizon) / (float64(height) - horizon)
	sea := mixed(light.seaFar, light.seaNear, depth)
	wave := math.Sin(x/18+y*0.9) * math.Sin(y*0.35)
	if wave > 0.85 {
		sea = mixed(sea, light.horizon, 0.35*(1-depth))
	}
	return sea
}

// onHeadland is a low hill rising from the sea on the right.
func onHeadland(x, y, centre, horizon float64, width int) bool {
	spread := float64(width) * 0.32
	top := horizon - 46*math.Exp(-math.Pow((x-centre)/spread*2.2, 2))
	return y >= top && y < horizon+14 && math.Abs(x-centre) < spread
}

// onLighthouse is a white tower with a red band on the headland's top, and its lamp: bright at
// night.
func onLighthouse(x, y, centre, horizon float64, night bool) (bool, color.NRGBA) {
	base, top := horizon-44, horizon-118
	halfWidth := 7 + 4*(y-top)/(base-top)
	if y >= top && y < base && math.Abs(x-centre) < halfWidth {
		if band := (y - top) / (base - top); band > 0.35 && band < 0.55 {
			return true, rgb(200, 50, 50)
		}
		return true, rgb(240, 240, 235)
	}
	lamp := rgb(60, 60, 70)
	if night {
		lamp = rgb(255, 240, 170)
	}
	if y >= top-12 && y < top && math.Abs(x-centre) < 8 {
		return true, lamp
	}
	return false, color.NRGBA{}
}

func mixed(from, to color.NRGBA, amount float64) color.NRGBA {
	amount = min(max(amount, 0), 1)
	channel := func(first, second uint8) uint8 {
		return uint8(float64(first) + (float64(second)-float64(first))*amount)
	}
	return color.NRGBA{R: channel(from.R, to.R), G: channel(from.G, to.G), B: channel(from.B, to.B), A: 255}
}
