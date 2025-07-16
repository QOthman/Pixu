//go:build wasm
// +build wasm

package graphics

import (
	"syscall/js"
)

func runGameLoop(gameLoop func()) {
	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if true {
			beginDrawing()
			gameLoop()
			js.Global().Call("requestAnimationFrame", renderFrame)
		}
		return nil
	})
	js.Global().Call("requestAnimationFrame", renderFrame)
	select {}
}
