package graphics

import (
	"time"
)

var (
	frameCount     int
	fps            int
	fpsTimer       = time.Now()
	frameDuration  time.Duration
	lastFrameTime  = time.Now()
	deltaTime      float64
	totalStartTime = time.Now()
)

// GetFps returns current FPS, updates once per second
func GetFps() int {
	frameCount++
	if time.Since(fpsTimer) >= time.Second {
		fps = frameCount
		frameCount = 0
		fpsTimer = time.Now()
	}
	return fps
}

// InitFps initializes frame timing for target FPS
func InitFps(targetFps int) {
	frameDuration = time.Second / time.Duration(targetFps)
	lastFrameTime = time.Now()
	deltaTime = 0
}

// Wait should be called at the end of each frame to maintain frame rate
func Wait() {
	now := time.Now()
	elapsed := now.Sub(lastFrameTime)

	if elapsed < frameDuration {
		time.Sleep(frameDuration - elapsed)
		now = time.Now() // Update now after sleeping
		elapsed = now.Sub(lastFrameTime)
	}

	deltaTime = elapsed.Seconds()

	lastFrameTime = lastFrameTime.Add(frameDuration)

	// Prevent spiral of death by resetting timer if too much behind
	if time.Since(lastFrameTime) > frameDuration {
		lastFrameTime = time.Now()
	}
}

// GetDeltaTime returns elapsed time in seconds since last frame
func GetDeltaTime() float64 {
	return deltaTime
}

// GetTotalTime returns total time since the start of the application
func ResetTimer() {
	totalStartTime = time.Now()
}

func GetTime() float64 {
	return time.Since(totalStartTime).Seconds()
}
