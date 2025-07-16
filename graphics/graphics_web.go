//go:build wasm
// +build wasm

package graphics

import (
	"syscall/js"
)

var (
	canvas     js.Value
	ctx        js.Value
	width      int
	height     int
	gameLoopFn js.Func
	running    bool
	document   js.Value
	keyStates     map[Key]bool
	prevKeyStates map[Key]bool
	keyDownCb     js.Func
	keyUpCb       js.Func
)

func initPlatform(w, h int, title string) error {
	running = true

	document = js.Global().Get("document")

	// Dynamically create a new canvas element
	canvas = js.Global().Get("document").Call("createElement", "canvas")

	// Set canvas width and height to 100% of the window size
	updateCanvasSize()

	// Apply styles to remove any margins/paddings and ensure full screen
	style := canvas.Get("style")
	style.Set("position", "absolute")
	style.Set("top", "0")
	style.Set("left", "0")
	style.Set("margin", "0")
	style.Set("padding", "0")
	style.Set("width", "100%")
	style.Set("height", "100%")
	canvas.Set("tabIndex", 0)
	// Append canvas to the body
	body := document.Get("body")
	body.Call("appendChild", canvas)

	// Get the 2D drawing context
	ctx = canvas.Call("getContext", "2d")

	// Set the title of the document
	document.Set("title", title)

	// Set up keyboard event listeners
	setupKeyboardEvents()

	// Attach the keydown and keyup event listeners to the canvas element
	canvas.Call("addEventListener", "keydown", keyDownCb)
	canvas.Call("addEventListener", "keyup", keyUpCb)

	// Set focus on the canvas to ensure it receives keyboard input
	canvas.Call("focus")

	// Initialize key state maps
	keyStates = make(map[Key]bool)
	prevKeyStates = make(map[Key]bool)

	// Add resize listener to update canvas size on window resize
	js.Global().Call("addEventListener", "resize", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		updateCanvasSize()
		return nil
	}))

	return nil
}

// updateCanvasSize updates the canvas to fit the window's dimensions.
func updateCanvasSize() {
	width = js.Global().Get("innerWidth").Int()
	height = js.Global().Get("innerHeight").Int()

	// Update canvas width and height
	canvas.Set("width", width)
	canvas.Set("height", height)
}

func clearBackground(c Color) {
	ctx.Set("fillStyle", colorToRGBA(c))
	ctx.Call("fillRect", 0, 0, width, height)
}

func close() {
	running = false
	if !canvas.IsUndefined() {
		canvas.Call("remove")
	}
	if !gameLoopFn.IsUndefined() {
		gameLoopFn.Release()
	}
}
