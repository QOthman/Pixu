package pixu

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Global window and graphics state
var (
	window            *glfw.Window
	windowWidth       int
	windowHeight      int
	windowResizedFlag bool
	isFullscreen      bool
	savedWindowX      int
	savedWindowY      int
	savedWindowW      int
	savedWindowH      int
)

func init() {
	// GLFW and OpenGL require execution locked to the main OS thread
	runtime.LockOSThread()
}

// WindowConfig configures window properties during initialization.
type WindowConfig struct {
	Title       string // Window title text
	Width       int    // Window initial width in pixels
	Height      int    // Window initial height in pixels
	TargetFPS   int    // Target frame rate (e.g. 60)
	Resizable   bool   // Allow window resizing
	Fullscreen  bool   // Start in borderless fullscreen mode
	VSync       bool   // Enable vertical sync
	MSAA        int    // Multi-sample anti-aliasing sample count (0 to disable, 4 recommended)
	Decorated   bool   // Window borders and titlebar (default true)
	AlwaysOnTop bool   // Window stays on top of other windows
}

// DefaultConfig returns the recommended default window configuration.
func DefaultConfig() WindowConfig {
	return WindowConfig{
		Title:       "Pixu Application",
		Width:       800,
		Height:      600,
		TargetFPS:   60,
		Resizable:   true,
		Fullscreen:  false,
		VSync:       true,
		MSAA:        4,
		Decorated:   true,
		AlwaysOnTop: false,
	}
}

// Init initializes the graphics engine and creates a window with default settings.
func Init(width, height int, title string) error {
	cfg := DefaultConfig()
	cfg.Width = width
	cfg.Height = height
	cfg.Title = title
	return InitWithConfig(cfg)
}

// InitWindow is an alias for Init to follow common 2D game framework conventions.
func InitWindow(width, height int, title string) error {
	return Init(width, height, title)
}

// InitWithConfig initializes the window and OpenGL context according to custom config.
func InitWithConfig(cfg WindowConfig) error {
	if cfg.Width <= 0 {
		cfg.Width = 800
	}
	if cfg.Height <= 0 {
		cfg.Height = 600
	}
	windowWidth = cfg.Width
	windowHeight = cfg.Height

	fmt.Println("Pixu: starting GLFW initialization")
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("failed to initialize GLFW: %w", err)
	}
	fmt.Println("Pixu: GLFW initialized")

	// OpenGL 3.3 Core Profile hints
	fmt.Println("Pixu: setting OpenGL window hints")
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Window properties hints
	if cfg.Resizable {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	} else {
		glfw.WindowHint(glfw.Resizable, glfw.False)
	}

	if cfg.Decorated {
		glfw.WindowHint(glfw.Decorated, glfw.True)
	} else {
		glfw.WindowHint(glfw.Decorated, glfw.False)
	}

	if cfg.AlwaysOnTop {
		glfw.WindowHint(glfw.Floating, glfw.True)
	} else {
		glfw.WindowHint(glfw.Floating, glfw.False)
	}

	if cfg.MSAA > 0 {
		glfw.WindowHint(glfw.Samples, cfg.MSAA)
	}

	var monitor *glfw.Monitor
	if cfg.Fullscreen {
		monitor = glfw.GetPrimaryMonitor()
		mode := monitor.GetVideoMode()
		windowWidth = mode.Width
		windowHeight = mode.Height
	}

	var err error
	fmt.Println("Pixu: creating GLFW window")
	window, err = glfw.CreateWindow(windowWidth, windowHeight, cfg.Title, monitor, nil)
	if err != nil {
		glfw.Terminate()
		return fmt.Errorf("failed to create GLFW window: %w", err)
	}
	if window == nil {
		glfw.Terminate()
		return fmt.Errorf("failed to create GLFW window: returned nil window")
	}
	fmt.Println("Pixu: GLFW window created")

	fmt.Println("Pixu: binding OpenGL context")
	window.MakeContextCurrent()
	fmt.Println("Pixu: OpenGL context bound")

	// Initialize OpenGL function pointers
	fmt.Println("Pixu: initializing OpenGL bindings")
	if err := gl.Init(); err != nil {
		window.Destroy()
		glfw.Terminate()
		return fmt.Errorf("failed to initialize OpenGL bindings: %w", err)
	}
	fmt.Println("Pixu: OpenGL bindings initialized")

	if cfg.VSync {
		glfw.SwapInterval(1)
	} else {
		glfw.SwapInterval(0)
	}

	// Setup callbacks
	setupInput(window)
	window.SetSizeCallback(func(w *glfw.Window, width, height int) {
		if width > 0 && height > 0 {
			windowWidth = width
			windowHeight = height
			windowResizedFlag = true
			gl.Viewport(0, 0, int32(width), int32(height))
			updateProjectionMatrix()
		}
	})

	// Setup viewport & GL capabilities
	gl.Viewport(0, 0, int32(windowWidth), int32(windowHeight))
	if cfg.MSAA > 0 {
		gl.Enable(gl.MULTISAMPLE)
	}

	// Initialize rendering pipeline & shaders
	fmt.Println("Pixu: creating renderer")
	if err := initRenderer(); err != nil {
		Close()
		return fmt.Errorf("failed to initialize renderer: %w", err)
	}
	fmt.Println("Pixu: renderer created")

	// Initialize embedded fonts
	fmt.Println("Pixu: initializing font system")
	if err := initFontSystem(); err != nil {
		Close()
		return fmt.Errorf("failed to initialize font system: %w", err)
	}
	fmt.Println("Pixu: font system ready")

	// Initialize timing
	InitFps(cfg.TargetFPS)

	return nil
}

// ShouldClose checks if the application window has received a close request.
func ShouldClose() bool {
	if window == nil {
		return true
	}
	return window.ShouldClose()
}

// ShouldContinue is an alias for !ShouldClose().
func ShouldContinue() bool {
	return !ShouldClose()
}

// ClearBackground clears the entire screen with the specified color.
func ClearBackground(color Color) {
	gl.ClearColor(color.R, color.G, color.B, color.A)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

// Clear is an alias for ClearBackground.
func Clear(color Color) {
	ClearBackground(color)
}

// BeginDrawing sets up the frame for drawing.
func BeginDrawing() {
	windowResizedFlag = false
}

// EndDrawing finalizes the current frame: swaps buffers, processes input, and paces frame rate.
func EndDrawing() {
	Present()
	UpdateInput()
	Wait()
}

// Present swaps the front and back buffers and polls GLFW event loop.
func Present() {
	if window != nil {
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

// Close gracefully terminates the OpenGL context and GLFW window.
func Close() {
	destroyFontSystem()
	destroyRenderer()

	if window != nil {
		window.Destroy()
		window = nil
	}
	glfw.Terminate()
}

// Terminate is an alias for Close.
func Terminate() {
	Close()
}

// ==========================================
// Window Manipulation API
// ==========================================

// SetWindowTitle updates the title string of the window.
func SetWindowTitle(title string) {
	if window != nil {
		window.SetTitle(title)
	}
}

// SetWindowSize changes the dimensions of the window.
func SetWindowSize(width, height int) {
	if window != nil && width > 0 && height > 0 {
		window.SetSize(width, height)
		windowWidth = width
		windowHeight = height
		gl.Viewport(0, 0, int32(width), int32(height))
		updateProjectionMatrix()
	}
}

// GetWindowSize returns current window width and height.
func GetWindowSize() (int, int) {
	return windowWidth, windowHeight
}

// GetWindowWidth returns the current window width in pixels.
func GetWindowWidth() int {
	return windowWidth
}

// GetWindowHeight returns the current window height in pixels.
func GetWindowHeight() int {
	return windowHeight
}

// GetScreenWidth returns the current window width in pixels.
func GetScreenWidth() int {
	return windowWidth
}

// GetScreenHeight returns the current window height in pixels.
func GetScreenHeight() int {
	return windowHeight
}

// GetMonitorSize returns the primary monitor screen resolution.
func GetMonitorSize() (int, int) {
	monitor := glfw.GetPrimaryMonitor()
	if monitor != nil {
		mode := monitor.GetVideoMode()
		return mode.Width, mode.Height
	}
	return windowWidth, windowHeight
}

// SetWindowPosition moves the window to (x, y) screen coordinates.
func SetWindowPosition(x, y int) {
	if window != nil {
		window.SetPos(x, y)
	}
}

// GetWindowPosition returns the window's (x, y) coordinates on the desktop.
func GetWindowPosition() (int, int) {
	if window != nil {
		return window.GetPos()
	}
	return 0, 0
}

// SetWindowMinSize sets the minimum allowed window dimensions.
func SetWindowMinSize(minW, minH int) {
	if window != nil {
		window.SetSizeLimits(minW, minH, glfw.DontCare, glfw.DontCare)
	}
}

// SetWindowMaxSize sets the maximum allowed window dimensions.
func SetWindowMaxSize(maxW, maxH int) {
	if window != nil {
		window.SetSizeLimits(glfw.DontCare, glfw.DontCare, maxW, maxH)
	}
}

// SetWindowFullscreen enables or disables borderless fullscreen mode.
func SetWindowFullscreen(fullscreen bool) {
	if window == nil || isFullscreen == fullscreen {
		return
	}

	monitor := glfw.GetPrimaryMonitor()
	if monitor == nil {
		return
	}

	if fullscreen {
		savedWindowX, savedWindowY = window.GetPos()
		savedWindowW, savedWindowH = window.GetSize()
		mode := monitor.GetVideoMode()
		window.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
		isFullscreen = true
	} else {
		window.SetMonitor(nil, savedWindowX, savedWindowY, savedWindowW, savedWindowH, 0)
		isFullscreen = false
	}
}

// ToggleFullscreen switches between windowed and fullscreen modes.
func ToggleFullscreen() {
	SetWindowFullscreen(!isFullscreen)
}

// IsWindowFullscreen returns true if the window is currently in fullscreen mode.
func IsWindowFullscreen() bool {
	return isFullscreen
}

// IsWindowFocused returns true if the window currently has keyboard/mouse focus.
func IsWindowFocused() bool {
	if window == nil {
		return false
	}
	return window.GetAttrib(glfw.Focused) == glfw.True
}

// IsWindowResized returns true if the window was resized during the last frame.
func IsWindowResized() bool {
	return windowResizedFlag
}

// SetVSync enables or disables vertical sync (1 = vsync on, 0 = vsync off).
func SetVSync(enabled bool) {
	if enabled {
		glfw.SwapInterval(1)
	} else {
		glfw.SwapInterval(0)
	}
}

// TakeScreenshot captures the current frame buffer and saves it as a PNG file.
func TakeScreenshot(filePath string) error {
	w, h := windowWidth, windowHeight
	pixels := make([]byte, w*h*4)

	gl.ReadPixels(0, 0, int32(w), int32(h), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			srcIdx := ((h-1-y)*w + x) * 4
			img.Pix[(y*w+x)*4+0] = pixels[srcIdx+0]
			img.Pix[(y*w+x)*4+1] = pixels[srcIdx+1]
			img.Pix[(y*w+x)*4+2] = pixels[srcIdx+2]
			img.Pix[(y*w+x)*4+3] = pixels[srcIdx+3]
		}
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create screenshot file %s: %w", filePath, err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("failed to encode screenshot png: %w", err)
	}

	return nil
}
