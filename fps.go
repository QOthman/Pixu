package pixu

import (
	"time"
)

var (
	frameCount     int
	fps            int
	fpsTimer       = time.Now()
	targetFPS      = 60
	frameDuration  = time.Second / 60
	lastFrameTime  = time.Now()
	deltaTime      float32
	totalStartTime = time.Now()
)

// InitFps initializes the target frame rate (default 60 FPS).
func InitFps(target int) {
	SetTargetFPS(target)
}

// SetTargetFPS sets the target frame rate limit in frames per second.
// Pass 0 or negative value to disable frame limiting.
func SetTargetFPS(target int) {
	targetFPS = target
	if target > 0 {
		frameDuration = time.Second / time.Duration(target)
	} else {
		frameDuration = 0
	}
	lastFrameTime = time.Now()
	deltaTime = 0
}

// GetTargetFPS returns the current target FPS limit.
func GetTargetFPS() int {
	return targetFPS
}

// GetFps returns the measured frames per second (updated once every second).
func GetFps() int {
	return fps
}

// GetFPS is an alias for GetFps.
func GetFPS() int {
	return fps
}

// GetDeltaTime returns the elapsed time in seconds since the last frame.
func GetDeltaTime() float32 {
	return deltaTime
}

// GetFrameTime is an alias for GetDeltaTime.
func GetFrameTime() float32 {
	return deltaTime
}

// GetTime returns the total elapsed time in seconds since the application started.
func GetTime() float64 {
	return time.Since(totalStartTime).Seconds()
}

// ResetTimer resets the application start timer to now.
func ResetTimer() {
	totalStartTime = time.Now()
}

// Wait paces the frame execution to maintain the target FPS.
// Call this once per frame (or use EndDrawing which calls it automatically).
func Wait() {
	now := time.Now()
	elapsed := now.Sub(lastFrameTime)

	if frameDuration > 0 && elapsed < frameDuration {
		time.Sleep(frameDuration - elapsed)
		now = time.Now()
		elapsed = now.Sub(lastFrameTime)
	}

	deltaTime = float32(elapsed.Seconds())
	lastFrameTime = now

	// Track FPS
	frameCount++
	if now.Sub(fpsTimer) >= time.Second {
		fps = frameCount
		frameCount = 0
		fpsTimer = now
	}
}
