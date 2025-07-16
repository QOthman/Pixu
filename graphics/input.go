//go:build !wasm
// +build !wasm

package graphics

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Key constants - common keys
// const (
// 	KeySpace     = glfw.KeySpace
// 	KeyEscape    = glfw.KeyEscape
// 	KeyEnter     = glfw.KeyEnter
// 	KeyTab       = glfw.KeyTab
// 	KeyBackspace = glfw.KeyBackspace
// 	KeyDelete    = glfw.KeyDelete

// 	// Arrow keys
// 	KeyUp    = glfw.KeyUp
// 	KeyDown  = glfw.KeyDown
// 	KeyLeft  = glfw.KeyLeft
// 	KeyRight = glfw.KeyRight

// 	// WASD
// 	KeyW = glfw.KeyW
// 	KeyA = glfw.KeyA
// 	KeyS = glfw.KeyS
// 	KeyD = glfw.KeyD

// 	// Numbers
// 	Key0 = glfw.Key0
// 	Key1 = glfw.Key1
// 	Key2 = glfw.Key2
// 	Key3 = glfw.Key3
// 	Key4 = glfw.Key4
// 	Key5 = glfw.Key5
// 	Key6 = glfw.Key6
// 	Key7 = glfw.Key7
// 	Key8 = glfw.Key8
// 	Key9 = glfw.Key9

// 	// Mouse buttons
// 	MouseLeft   = glfw.MouseButton1
// 	MouseRight  = glfw.MouseButton2
// 	MouseMiddle = glfw.MouseButton3
// )

// Input state tracking
var (
	keyStates = make(map[Key]bool)
	prevKeyStates = make(map[Key]bool)

	keysPressed      = make(map[int]bool)
	keysJustPressed  = make(map[int]bool)
	keysJustReleased = make(map[int]bool)

	mousePressed      = make(map[int]bool)
	mouseJustPressed  = make(map[int]bool)
	mouseJustReleased = make(map[int]bool)

	mouseX, mouseY           float64
	mouseDeltaX, mouseDeltaY float64
	lastMouseX, lastMouseY   float64

	scrollX, scrollY           float64
	scrollDeltaX, scrollDeltaY float64

	lastWindowWidth, lastWindowHeight int
)

// Setup input callbacks
func setupInput(window *glfw.Window) {
	// Keyboard callbacks
	window.SetKeyCallback(keyCallback)

	// Mouse callbacks
	window.SetMouseButtonCallback(mouseCallback)
	window.SetCursorPosCallback(mousePosCallback)
	window.SetScrollCallback(scrollCallback)

	// Window size callback
	// window.SetSizeCallback(windowSizeCallback)
}

// // Update input state
// // Call this every frame to update input states
func updateInput() {
	// Clear "just pressed/released" states
	for key := range keysJustPressed {
		keysJustPressed[key] = false
	}
	for key := range keysJustReleased {
		keysJustReleased[key] = false
	}
	for btn := range mouseJustPressed {
		mouseJustPressed[btn] = false
	}
	for btn := range mouseJustReleased {
		mouseJustReleased[btn] = false
	}

	// Update mouse delta
	mouseDeltaX = mouseX - lastMouseX
	mouseDeltaY = mouseY - lastMouseY
	lastMouseX = mouseX
	lastMouseY = mouseY

	scrollDeltaX = 0
	scrollDeltaY = 0

	keyStates[KeySpace] = window.GetKey(glfw.KeySpace) == glfw.Press
	keyStates[KeyA] = window.GetKey(glfw.KeyA) == glfw.Press
	keyStates[KeyB] = window.GetKey(glfw.KeyB) == glfw.Press
	keyStates[KeyC] = window.GetKey(glfw.KeyC) == glfw.Press
	keyStates[KeyD] = window.GetKey(glfw.KeyD) == glfw.Press
	keyStates[KeyE] = window.GetKey(glfw.KeyE) == glfw.Press
	keyStates[KeyF] = window.GetKey(glfw.KeyF) == glfw.Press
	keyStates[KeyG] = window.GetKey(glfw.KeyG) == glfw.Press
	keyStates[KeyH] = window.GetKey(glfw.KeyH) == glfw.Press
	keyStates[KeyI] = window.GetKey(glfw.KeyI) == glfw.Press
	keyStates[KeyJ] = window.GetKey(glfw.KeyJ) == glfw.Press
	keyStates[KeyK] = window.GetKey(glfw.KeyK) == glfw.Press
	keyStates[KeyL] = window.GetKey(glfw.KeyL) == glfw.Press
	keyStates[KeyM] = window.GetKey(glfw.KeyM) == glfw.Press
	keyStates[KeyN] = window.GetKey(glfw.KeyN) == glfw.Press
	keyStates[KeyO] = window.GetKey(glfw.KeyO) == glfw.Press
	keyStates[KeyP] = window.GetKey(glfw.KeyP) == glfw.Press
	keyStates[KeyQ] = window.GetKey(glfw.KeyQ) == glfw.Press
	keyStates[KeyR] = window.GetKey(glfw.KeyR) == glfw.Press
	keyStates[KeyS] = window.GetKey(glfw.KeyS) == glfw.Press
	keyStates[KeyT] = window.GetKey(glfw.KeyT) == glfw.Press
	keyStates[KeyU] = window.GetKey(glfw.KeyU) == glfw.Press
	keyStates[KeyV] = window.GetKey(glfw.KeyV) == glfw.Press
	keyStates[KeyW] = window.GetKey(glfw.KeyW) == glfw.Press
	keyStates[KeyX] = window.GetKey(glfw.KeyX) == glfw.Press
	keyStates[KeyY] = window.GetKey(glfw.KeyY) == glfw.Press
	keyStates[KeyZ] = window.GetKey(glfw.KeyZ) == glfw.Press
	keyStates[KeyLeft] = window.GetKey(glfw.KeyLeft) == glfw.Press
	keyStates[KeyRight] = window.GetKey(glfw.KeyRight) == glfw.Press
	keyStates[KeyUp] = window.GetKey(glfw.KeyUp) == glfw.Press
	keyStates[KeyDown] = window.GetKey(glfw.KeyDown) == glfw.Press
	keyStates[KeyEnter] = window.GetKey(glfw.KeyEnter) == glfw.Press
	keyStates[KeyEscape] = window.GetKey(glfw.KeyEscape) == glfw.Press
}

// // ============= KEYBOARD FUNCTIONS =============

// Check if a key is pressed (continuously)
// Returns true if the key is currently pressed down
func isKeyPressed(key Key) bool {
	return keyStates[key] && !prevKeyStates[key]
}

// Check if a key is just pressed (one frame only)
// Returns true if the key was pressed this frame
func isKeyJustPressed(key Key) bool {
	return keyStates[key]
}

// Check if a key is just released (one frame only)
// Returns true if the key was released this frame
// func isKeyJustReleased(key glfw.Key) bool {
// 	return keysJustReleased[int(key)]
// }

// // ============= MOUSE FUNCTIONS =============

// // Check if a mouse button is pressed (continuously)
// // Returns true if the mouse button is currently pressed down
// func isMousePressed(button glfw.MouseButton) bool {
// 	return mousePressed[int(button)]
// }

// // Check if a mouse button is just pressed (one frame only)
// // Returns true if the mouse button was pressed this frame
// func isMouseJustPressed(button glfw.MouseButton) bool {
// 	return mouseJustPressed[int(button)]
// }

// // Check if a mouse button is just released (one frame only)
// // Returns true if the mouse button was released this frame
// func isMouseJustReleased(button glfw.MouseButton) bool {
// 	return mouseJustReleased[int(button)]
// }

// // Get mouse position
// func getMousePosition() (float32, float32) {
// 	return float32(mouseX), float32(mouseY)
// }

// // Get mouse delta (movement since last frame)
// func getMouseDelta() (float32, float32) {
// 	return float32(mouseDeltaX), float32(mouseDeltaY)
// }

// // Check if mouse moved
// func isMouseMoved() bool {
// 	return mouseDeltaX != 0 || mouseDeltaY != 0
// }

// // Get scroll delta
// func getScrollDelta() (float32, float32) {
// 	return float32(scrollDeltaX), float32(scrollDeltaY)
// }

// // Get scroll position
// func getScrollPosition() (float32, float32) {
// 	return float32(scrollX), float32(scrollY)
// }

// // Check if scroll wheel moved
// func isScrollMoved() bool {
// 	return scrollDeltaX != 0 || scrollDeltaY != 0
// }

// // ============= WINDOW FUNCTIONS =============

// // func IsWindowResized() bool {
// // 	if windowWidth != lastWindowWidth || windowHeight != lastWindowHeight {
// // 		lastWindowWidth = windowWidth
// // 		lastWindowHeight = windowHeight
// // 		return true
// // 	}
// // 	return false
// // }

// // func GetWindowSize() (int, int) {
// // 	return windowWidth, windowHeight
// // }

// // ============= CALLBACKS =============

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		keysPressed[int(key)] = true
		keysJustPressed[int(key)] = true
	case glfw.Release:
		keysPressed[int(key)] = false
		keysJustReleased[int(key)] = true
	}
}

func mouseCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		mousePressed[int(button)] = true
		mouseJustPressed[int(button)] = true
	case glfw.Release:
		mousePressed[int(button)] = false
		mouseJustReleased[int(button)] = true
	}
}

func mousePosCallback(w *glfw.Window, xpos, ypos float64) {
	mouseX = xpos
	mouseY = ypos
}

func scrollCallback(w *glfw.Window, xoffset, yoffset float64) {
	scrollDeltaX += xoffset
	scrollDeltaY += yoffset
}

// func windowSizeCallback(w *glfw.Window, width, height int) {
// 	windowWidth = width
// 	windowHeight = height
// 	gl.Viewport(0, 0, int32(width), int32(height)) // Update viewport
// }
