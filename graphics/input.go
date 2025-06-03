package graphics

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Key constants - common keys
const (
	KeySpace     = glfw.KeySpace
	KeyEscape    = glfw.KeyEscape
	KeyEnter     = glfw.KeyEnter
	KeyTab       = glfw.KeyTab
	KeyBackspace = glfw.KeyBackspace
	KeyDelete    = glfw.KeyDelete

	// Arrow keys
	KeyUp    = glfw.KeyUp
	KeyDown  = glfw.KeyDown
	KeyLeft  = glfw.KeyLeft
	KeyRight = glfw.KeyRight

	// WASD
	KeyW = glfw.KeyW
	KeyA = glfw.KeyA
	KeyS = glfw.KeyS
	KeyD = glfw.KeyD

	// Numbers
	Key0 = glfw.Key0
	Key1 = glfw.Key1
	Key2 = glfw.Key2
	Key3 = glfw.Key3
	Key4 = glfw.Key4
	Key5 = glfw.Key5
	Key6 = glfw.Key6
	Key7 = glfw.Key7
	Key8 = glfw.Key8
	Key9 = glfw.Key9

	// Mouse buttons
	MouseLeft   = glfw.MouseButton1
	MouseRight  = glfw.MouseButton2
	MouseMiddle = glfw.MouseButton3
)

// Input state tracking
var (
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
	window.SetSizeCallback(windowSizeCallback)
}

// Update input state
// Call this every frame to update input states
func UpdateInput() {
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
}

// ============= KEYBOARD FUNCTIONS =============

// Check if a key is pressed (continuously)
// Returns true if the key is currently pressed down
func IsKeyPressed(key glfw.Key) bool {
	return keysPressed[int(key)]
}

// Check if a key is just pressed (one frame only)
// Returns true if the key was pressed this frame
func IsKeyJustPressed(key glfw.Key) bool {
	return keysJustPressed[int(key)]
}

// Check if a key is just released (one frame only)
// Returns true if the key was released this frame
func IsKeyJustReleased(key glfw.Key) bool {
	return keysJustReleased[int(key)]
}

// ============= MOUSE FUNCTIONS =============

// Check if a mouse button is pressed (continuously)
// Returns true if the mouse button is currently pressed down
func IsMousePressed(button glfw.MouseButton) bool {
	return mousePressed[int(button)]
}

// Check if a mouse button is just pressed (one frame only)
// Returns true if the mouse button was pressed this frame
func IsMouseJustPressed(button glfw.MouseButton) bool {
	return mouseJustPressed[int(button)]
}

// Check if a mouse button is just released (one frame only)
// Returns true if the mouse button was released this frame
func IsMouseJustReleased(button glfw.MouseButton) bool {
	return mouseJustReleased[int(button)]
}

// Get mouse position
func GetMousePosition() (float32, float32) {
	return float32(mouseX), float32(mouseY)
}

// Get mouse delta (movement since last frame)
func GetMouseDelta() (float32, float32) {
	return float32(mouseDeltaX), float32(mouseDeltaY)
}

// Check if mouse moved
func IsMouseMoved() bool {
	return mouseDeltaX != 0 || mouseDeltaY != 0
}

// Get scroll delta
func GetScrollDelta() (float32, float32) {
	return float32(scrollDeltaX), float32(scrollDeltaY)
}

// Get scroll position
func GetScrollPosition() (float32, float32) {
	return float32(scrollX), float32(scrollY)
}

// Check if scroll wheel moved
func IsScrollMoved() bool {
	return scrollDeltaX != 0 || scrollDeltaY != 0
}

// ============= WINDOW FUNCTIONS =============

func IsWindowResized() bool {
	if windowWidth != lastWindowWidth || windowHeight != lastWindowHeight {
		lastWindowWidth = windowWidth
		lastWindowHeight = windowHeight
		return true
	}
	return false
}

func GetWindowSize() (int, int) {
	return windowWidth, windowHeight
}

// ============= CALLBACKS =============

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

func windowSizeCallback(w *glfw.Window, width, height int) {
	windowWidth = width
	windowHeight = height
	gl.Viewport(0, 0, int32(width), int32(height)) // Update viewport
}
