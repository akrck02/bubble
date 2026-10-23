module github.com/akrck02/bubble/goui

go 1.25.6

require (
	github.com/akrck02/bubble v0.0.0
	github.com/akrck02/go-renderer v0.2.1-0.20261022213249-8ecc924ab3e5
	github.com/akrck02/go-ui v0.0.0
)

require (
	github.com/go-gl/gl v0.0.0-20260331235117-4566fea9a276 // indirect
	github.com/go-gl/glfw/v3.4/glfw v0.1.0-pre.1 // indirect
	golang.org/x/image v0.45.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
)

// Bubble is the module one folder up. go-ui still points at go-renderer through a local replace,
// so build against both cloned next to Bubble, on their main.
replace (
	github.com/akrck02/bubble => ../
	github.com/akrck02/go-renderer => ../../go-renderer
	github.com/akrck02/go-ui => ../../go-ui
)
