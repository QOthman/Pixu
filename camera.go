package pixu

import (
	"math"
)

// Camera2D represents a 2D camera with offset, target, rotation, and zoom.
type Camera2D struct {
	Offset   Vec2    // Camera viewport offset (e.g., screen center: V2(width/2, height/2))
	Target   Vec2    // Camera target position in world space
	Rotation float32 // Camera rotation in degrees
	Zoom     float32 // Camera zoom level (1.0 = standard 1:1 scale)
}

// DefaultCamera2D creates a default 2D camera with zoom 1.0 and zero offset/target.
func DefaultCamera2D() Camera2D {
	return Camera2D{
		Offset:   Vec2{X: 0, Y: 0},
		Target:   Vec2{X: 0, Y: 0},
		Rotation: 0,
		Zoom:     1.0,
	}
}

// activeCamera holds the currently active 2D camera mode, if any.
var (
	activeCamera Camera2D
	isCamera2D   bool
)

// BeginMode2D activates 2D camera transformation for subsequent draw calls.
func BeginMode2D(camera Camera2D) {
	if camera.Zoom <= 0 {
		camera.Zoom = 1.0
	}
	activeCamera = camera
	isCamera2D = true
	updateProjectionMatrix()
}

// EndMode2D deactivates 2D camera mode and returns to standard screen-space rendering.
func EndMode2D() {
	isCamera2D = false
	updateProjectionMatrix()
}

// GetScreenToWorld2D converts a screen position to world coordinates based on the camera.
func GetScreenToWorld2D(position Vec2, camera Camera2D) Vec2 {
	zoom := camera.Zoom
	if zoom <= 0 {
		zoom = 1.0
	}

	// 1. Subtract offset
	dx := position.X - camera.Offset.X
	dy := position.Y - camera.Offset.Y

	// 2. Rotate in opposite direction
	rad := float64(-camera.Rotation) * math.Pi / 180.0
	cos := float32(math.Cos(rad))
	sin := float32(math.Sin(rad))

	rx := dx*cos - dy*sin
	ry := dx*sin + dy*cos

	// 3. Unscale
	rx /= zoom
	ry /= zoom

	// 4. Add target
	return Vec2{
		X: rx + camera.Target.X,
		Y: ry + camera.Target.Y,
	}
}

// GetWorldToScreen2D converts a world position to screen coordinates based on the camera.
func GetWorldToScreen2D(position Vec2, camera Camera2D) Vec2 {
	zoom := camera.Zoom
	if zoom <= 0 {
		zoom = 1.0
	}

	// 1. Subtract target
	dx := position.X - camera.Target.X
	dy := position.Y - camera.Target.Y

	// 2. Scale
	dx *= zoom
	dy *= zoom

	// 3. Rotate
	rad := float64(camera.Rotation) * math.Pi / 180.0
	cos := float32(math.Cos(rad))
	sin := float32(math.Sin(rad))

	rx := dx*cos - dy*sin
	ry := dx*sin + dy*cos

	// 4. Add offset
	return Vec2{
		X: rx + camera.Offset.X,
		Y: ry + camera.Offset.Y,
	}
}

// computeCameraMatrix4 returns the 4x4 transformation matrix for Camera2D in column-major order.
func computeCameraMatrix4(camera Camera2D) [16]float32 {
	zoom := camera.Zoom
	if zoom <= 0 {
		zoom = 1.0
	}

	rad := float64(camera.Rotation) * math.Pi / 180.0
	cos := float32(math.Cos(rad))
	sin := float32(math.Sin(rad))

	// Combined affine 2D transform into 4x4 matrix:
	// translation(offset) * rotation(rot) * scale(zoom) * translation(-target)
	//
	// [ cos*zoom, -sin*zoom, 0, offset.x - (target.x*cos - target.y*sin)*zoom ]
	// [ sin*zoom,  cos*zoom, 0, offset.y - (target.x*sin + target.y*cos)*zoom ]
	// [ 0,         0,        1, 0                                            ]
	// [ 0,         0,        0, 1                                            ]

	tx := camera.Offset.X - (camera.Target.X*cos-camera.Target.Y*sin)*zoom
	ty := camera.Offset.Y - (camera.Target.X*sin+camera.Target.Y*cos)*zoom

	// Column-major array:
	return [16]float32{
		cos * zoom, sin * zoom, 0, 0,
		-sin * zoom, cos * zoom, 0, 0,
		0, 0, 1, 0,
		tx, ty, 0, 1,
	}
}
