//go:build !wasm
// +build !wasm

package graphics

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

func runGameLoop(loopFunc func()) {

	for !window.ShouldClose() {
		glfw.PollEvents()
		updateInput()
		loopFunc()
		Wait()
		window.SwapBuffers()
	}
}
