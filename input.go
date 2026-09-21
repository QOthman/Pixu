package pixu

import (
	"sync"

	"github.com/go-gl/glfw/v3.3/glfw"
)

// Key type definition for keyboard keys.
type Key = glfw.Key

// MouseButton type definition for mouse buttons.
type MouseButton = glfw.MouseButton

// Keyboard key constants
const (
	// Control & whitespace
	KeySpace        Key = glfw.KeySpace
	KeyApostrophe   Key = glfw.KeyApostrophe
	KeyComma        Key = glfw.KeyComma
	KeyMinus        Key = glfw.KeyMinus
	KeyPeriod       Key = glfw.KeyPeriod
	KeySlash        Key = glfw.KeySlash
	KeySemicolon    Key = glfw.KeySemicolon
	KeyEqual        Key = glfw.KeyEqual
	KeyLeftBracket  Key = glfw.KeyLeftBracket
	KeyBackslash    Key = glfw.KeyBackslash
	KeyRightBracket Key = glfw.KeyRightBracket
	KeyGraveAccent  Key = glfw.KeyGraveAccent
	KeyEscape       Key = glfw.KeyEscape
	KeyEnter        Key = glfw.KeyEnter
	KeyTab          Key = glfw.KeyTab
	KeyBackspace    Key = glfw.KeyBackspace
	KeyInsert       Key = glfw.KeyInsert
	KeyDelete       Key = glfw.KeyDelete

	// Numbers
	Key0 Key = glfw.Key0
	Key1 Key = glfw.Key1
	Key2 Key = glfw.Key2
	Key3 Key = glfw.Key3
	Key4 Key = glfw.Key4
	Key5 Key = glfw.Key5
	Key6 Key = glfw.Key6
	Key7 Key = glfw.Key7
	Key8 Key = glfw.Key8
	Key9 Key = glfw.Key9

	// Alphabet
	KeyA Key = glfw.KeyA
	KeyB Key = glfw.KeyB
	KeyC Key = glfw.KeyC
	KeyD Key = glfw.KeyD
	KeyE Key = glfw.KeyE
	KeyF Key = glfw.KeyF
	KeyG Key = glfw.KeyG
	KeyH Key = glfw.KeyH
	KeyI Key = glfw.KeyI
	KeyJ Key = glfw.KeyJ
	KeyK Key = glfw.KeyK
	KeyL Key = glfw.KeyL
	KeyM Key = glfw.KeyM
	KeyN Key = glfw.KeyN
	KeyO Key = glfw.KeyO
	KeyP Key = glfw.KeyP
	KeyQ Key = glfw.KeyQ
	KeyR Key = glfw.KeyR
	KeyS Key = glfw.KeyS
	KeyT Key = glfw.KeyT
	KeyU Key = glfw.KeyU
	KeyV Key = glfw.KeyV
	KeyW Key = glfw.KeyW
	KeyX Key = glfw.KeyX
	KeyY Key = glfw.KeyY
	KeyZ Key = glfw.KeyZ

	// Arrows & navigation
	KeyRight    Key = glfw.KeyRight
	KeyLeft     Key = glfw.KeyLeft
	KeyDown     Key = glfw.KeyDown
	KeyUp       Key = glfw.KeyUp
	KeyPageUp   Key = glfw.KeyPageUp
	KeyPageDown Key = glfw.KeyPageDown
	KeyHome     Key = glfw.KeyHome
	KeyEnd      Key = glfw.KeyEnd

	// Function keys
	KeyF1  Key = glfw.KeyF1
	KeyF2  Key = glfw.KeyF2
	KeyF3  Key = glfw.KeyF3
	KeyF4  Key = glfw.KeyF4
	KeyF5  Key = glfw.KeyF5
	KeyF6  Key = glfw.KeyF6
	KeyF7  Key = glfw.KeyF7
	KeyF8  Key = glfw.KeyF8
	KeyF9  Key = glfw.KeyF9
	KeyF10 Key = glfw.KeyF10
	KeyF11 Key = glfw.KeyF11
	KeyF12 Key = glfw.KeyF12

	// Modifiers
	KeyLeftShift    Key = glfw.KeyLeftShift
	KeyLeftControl  Key = glfw.KeyLeftControl
	KeyLeftAlt      Key = glfw.KeyLeftAlt
	KeyLeftSuper    Key = glfw.KeyLeftSuper
	KeyRightShift   Key = glfw.KeyRightShift
	KeyRightControl Key = glfw.KeyRightControl
	KeyRightAlt     Key = glfw.KeyRightAlt
	KeyRightSuper   Key = glfw.KeyRightSuper
	KeyMenu         Key = glfw.KeyMenu
)

// Mouse button constants
const (
	MouseLeft    MouseButton = glfw.MouseButton1
	MouseRight   MouseButton = glfw.MouseButton2
	MouseMiddle  MouseButton = glfw.MouseButton3
	MouseButton1 MouseButton = glfw.MouseButton1
	MouseButton2 MouseButton = glfw.MouseButton2
	MouseButton3 MouseButton = glfw.MouseButton3
	MouseButton4 MouseButton = glfw.MouseButton4
	MouseButton5 MouseButton = glfw.MouseButton5
	MouseButton6 MouseButton = glfw.MouseButton6
	MouseButton7 MouseButton = glfw.MouseButton7
	MouseButton8 MouseButton = glfw.MouseButton8
)

// Internal input state tracker
var (
	inputMutex sync.RWMutex

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
)

// setupInput registers all GLFW input callbacks.
func setupInput(w *glfw.Window) {
	w.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		inputMutex.Lock()
		defer inputMutex.Unlock()
		k := int(key)
		switch action {
		case glfw.Press:
			keysPressed[k] = true
			keysJustPressed[k] = true
		case glfw.Release:
			keysPressed[k] = false
			keysJustReleased[k] = true
		}
	})

	w.SetMouseButtonCallback(func(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		inputMutex.Lock()
		defer inputMutex.Unlock()
		btn := int(button)
		switch action {
		case glfw.Press:
			mousePressed[btn] = true
			mouseJustPressed[btn] = true
		case glfw.Release:
			mousePressed[btn] = false
			mouseJustReleased[btn] = true
		}
	})

	w.SetCursorPosCallback(func(window *glfw.Window, xpos, ypos float64) {
		inputMutex.Lock()
		defer inputMutex.Unlock()
		mouseX = xpos
		mouseY = ypos
	})

	w.SetScrollCallback(func(window *glfw.Window, xoffset, yoffset float64) {
		inputMutex.Lock()
		defer inputMutex.Unlock()
		scrollDeltaX += xoffset
		scrollDeltaY += yoffset
		scrollX += xoffset
		scrollY += yoffset
	})
}

// UpdateInput clears transient per-frame inputs and recalculates deltas.
func UpdateInput() {
	inputMutex.Lock()
	defer inputMutex.Unlock()

	for k := range keysJustPressed {
		delete(keysJustPressed, k)
	}
	for k := range keysJustReleased {
		delete(keysJustReleased, k)
	}
	for b := range mouseJustPressed {
		delete(mouseJustPressed, b)
	}
	for b := range mouseJustReleased {
		delete(mouseJustReleased, b)
	}

	mouseDeltaX = mouseX - lastMouseX
	mouseDeltaY = mouseY - lastMouseY
	lastMouseX = mouseX
	lastMouseY = mouseY

	scrollDeltaX = 0
	scrollDeltaY = 0
}

// ==========================================
// Keyboard API
// ==========================================

// IsKeyPressed returns true if the key is currently held down.
func IsKeyPressed(key Key) bool {
	inputMutex.RLock()
	defer inputMutex.RUnlock()
	return keysPressed[int(key)]
}

// IsKeyDown is an alias for IsKeyPressed.
func IsKeyDown(key Key) bool {
	return IsKeyPressed(key)
}

// IsKeyJustPressed returns true if the key was pressed down in the current frame.
func IsKeyJustPressed(key Key) bool {
	inputMutex.RLock()
	defer inputMutex.RUnlock()
	return keysJustPressed[int(key)]
}

// IsKeyTriggered is an alias for IsKeyJustPressed.
func IsKeyTriggered(key Key) bool {
	return IsKeyJustPressed(key)
}

// IsKeyJustReleased returns true if the key was released in the current frame.
func IsKeyJustReleased(key Key) bool {
	inputMutex.RLock()
	defer inputMutex.RUnlock()
	return keysJustReleased[int(key)]
}

// IsKeyUp returns true if the key is currently NOT pressed.
func IsKeyUp(key Key) bool {
	return !IsKeyPressed(key)
}

// ==========================================
// Mouse API
// ==========================================

// IsMouseButtonPressed returns true if the mouse button is currently held down.
func IsMouseButtonPressed(button MouseButton) bool {
	inputMutex.RLock()
	defer inputMutex.RUnlock()
	return mousePressed[int(button)]
}

// IsMouseButtonDown is an alias for IsMouseButtonPressed.
func IsMouseButtonDown(button MouseButton) bool {
	return IsMouseButtonPressed(button)
}

// IsMouseButtonJustPressed returns true if the mouse button was pressed this frame.
func IsMouseButtonJustPressed(button MouseButton) bool {
	inputMutex.RLock()
	defer inputMutex.RUnlock()
	return mouseJustPressed[int(button)]
}

// IsMouseButtonJustReleased returns true if the mouse button was released this frame.
func IsMouseButtonJustReleased(button MouseButton) bool {
	inputMutex.RLock()
	defer inputMutex.RUnlock()
	return mouseJustReleased[int(button)]
}

// IsMouseButtonUp returns true if the mouse button is currently NOT pressed.
func IsMouseButtonUp(button MouseButton) bool {
	return !IsMouseButtonPressed(button)
}

// GetMousePosition returns the mouse cursor coordinates (x, y) relative to top-left of window.
func GetMousePosition() (float32, float32) {
	inputMutex.RLock()
	defer inputMutex.RUnlock()
	return float32(mouseX), float32(mouseY)
}

// GetMousePositionV returns the mouse cursor position as a Vec2.
func GetMousePositionV() Vec2 {
	x, y := GetMousePosition()
	return Vec2{X: x, Y: y}
}

// GetMouseX returns the mouse cursor X position.
func GetMouseX() float32 {
	inputMutex.RLock()
	defer inputMutex.RUnlock()
	return float32(mouseX)
}

// GetMouseY returns the mouse cursor Y position.
func GetMouseY() float32 {
	inputMutex.RLock()
	defer inputMutex.RUnlock()
	return float32(mouseY)
}

// GetMouseDelta returns the mouse displacement (dx, dy) since the last frame.
func GetMouseDelta() (float32, float32) {
	inputMutex.RLock()
	defer inputMutex.RUnlock()
	return float32(mouseDeltaX), float32(mouseDeltaY)
}

// GetMouseDeltaV returns the mouse displacement since last frame as a Vec2.
func GetMouseDeltaV() Vec2 {
	dx, dy := GetMouseDelta()
	return Vec2{X: dx, Y: dy}
}

// IsMouseMoved returns true if the mouse moved since the last frame.
func IsMouseMoved() bool {
	inputMutex.RLock()
	defer inputMutex.RUnlock()
	return mouseDeltaX != 0 || mouseDeltaY != 0
}

// GetScrollDelta returns the mouse wheel scroll movement (dx, dy) in this frame.
func GetScrollDelta() (float32, float32) {
	inputMutex.RLock()
	defer inputMutex.RUnlock()
	return float32(scrollDeltaX), float32(scrollDeltaY)
}

// GetMouseWheelMove returns the vertical scroll wheel movement in this frame.
func GetMouseWheelMove() float32 {
	_, dy := GetScrollDelta()
	return dy
}

// GetScrollPosition returns the cumulative scroll position.
func GetScrollPosition() (float32, float32) {
	inputMutex.RLock()
	defer inputMutex.RUnlock()
	return float32(scrollX), float32(scrollY)
}

// SetCursorVisible shows or hides the mouse cursor over the window.
func SetCursorVisible(visible bool) {
	if window != nil {
		if visible {
			window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		} else {
			window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
		}
	}
}

// SetCursorDisabled locks and hides the cursor to capture unbounded relative mouse movement.
func SetCursorDisabled(disabled bool) {
	if window != nil {
		if disabled {
			window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
		} else {
			window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		}
	}
}
