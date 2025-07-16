//go:build wasm
// +build wasm

package graphics

import (
	"syscall/js"
)


func setupKeyboardEvents() {
	console := js.Global().Get("console")
	
	keyDownCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		keyCode := event.Get("code").String()
		// console.Call("log", "KeyDown:", keyCode) // Debug print

		if key := jsKeyToEngineKey(keyCode); key != -1 {
			keyStates[key] = true
			// console.Call("log", "Mapped to:", key) // Debug print
			event.Call("preventDefault")
		} else {
			// console.Call("log", "Unknown key:", keyCode)
		}
		return nil
	})

	keyUpCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		keyCode := event.Get("code").String()
		// console.Call("log", "KeyUp:", keyCode)

		if key := jsKeyToEngineKey(keyCode); key != -1 {
			keyStates[key] = false
			event.Call("preventDefault")
		}
		return nil
	})

	if canvas.IsUndefined() || canvas.IsNull() {
		canvas = js.Global().Get("document")
		console.Call("log", "Canvas is undefined, falling back to document")
	}

	canvas.Call("addEventListener", "keydown", keyDownCb)
	canvas.Call("addEventListener", "keyup", keyUpCb)
}


// Converts JS key string to your custom Key enum
func jsKeyToEngineKey(jsKey string) Key {
	switch jsKey {
	case "Space":
		return KeySpace
	case "Enter":
		return KeyEnter
	case "Escape":
		return KeyEscape
	case "ArrowLeft":
		return KeyLeft
	case "ArrowRight":
		return KeyRight
	case "ArrowUp":
		return KeyUp
	case "ArrowDown":
		return KeyDown
	// Letters A-Z
	default:
		if len(jsKey) == 4 && jsKey[:3] == "Key" {
			return Key(jsKey[3] - 'A') // Assuming KeyA = 0, KeyB = 1, etc.
		}
	}
	return -1
}

// Key input utility functions

func isKeyPressed(key Key) bool {
	return keyStates[key] 
}

func isKeyDown(key Key) bool {
	return keyStates[key]
}

func isKeyReleased(key Key) bool {
	return !keyStates[key] && prevKeyStates[key]
}

// Call this at the start of each frame
func beginDrawing() {
	for k := range keyStates {
		prevKeyStates[k] = keyStates[k]
	}
}