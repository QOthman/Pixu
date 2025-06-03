package graphics

import (
	"runtime"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Global states for the graphics system
var (
	window *glfw.Window
	shaderProgram uint32
	currentVAO uint32
	currentVBO uint32
	windowWidth, windowHeight int
)

func init() {
	runtime.LockOSThread()
}

// Initialize for graphics system
func Init(width, height int, title string) error {
	windowWidth, windowHeight = width, height
	
	if err := glfw.Init(); err != nil {
		return err
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	var err error
	window, err = glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return err
	}
	window.MakeContextCurrent()
	setupInput(window)
	if err := gl.Init(); err != nil {
		return err
	}

	// Setup viewport and OpenGL settings
	gl.Viewport(0, 0, int32(width), int32(height))
	setupShaders()
	setupBuffers()
	
	initTextureSystem()
	loadFontAtlas()
	InitFps(60)
	return nil
}

// Check if window should continue running
func ShouldContinue() bool {
	return !window.ShouldClose()
}

// Clear the background with a color
func ClearBackground(color Color) {
	gl.ClearColor(color.R, color.G, color.B, color.A)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

// Present the current frame
// This swaps the buffers and polls events
// It should be called after all drawing operations are done.
func Present() {
	window.SwapBuffers()
	glfw.PollEvents()
}

// Close the graphics system
// This should be called when the application is done with graphics
func Close() {
	glfw.Terminate()
}
