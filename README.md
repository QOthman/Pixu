# 🎮 Pixu

[![Go Reference](https://pkg.go.dev/badge/github.com/QOthman/Pixu.svg)](https://pkg.go.dev/github.com/QOthman/Pixu)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/QOthman/Pixu)](https://golang.org)

**Pixu** is a lightweight, fast, and modern 2D graphics and game development library for Go, powered by OpenGL 3.3 and GLFW. Inspired by simple immediate-mode drawing APIs like Raylib, Pixu makes 2D rendering, math, inputs, camera transforms, and collision detection effortless.

---

## ✨ Features

- ⚡ **High Performance Rendering**: Persistent GPU buffers (VAO/VBO/EBO) with zero per-frame buffer allocation overhead.
- 🎨 **Rich 2D Geometry**: Filled & outlined shapes, thick lines, bezier curves, rounded rectangles, gradient fills, regular polygons, circles, rings, and ellipses.
- 🖼️ **Textures & Sprites**: PNG, JPEG, and GIF support with nearest/linear filtering, spritesheet sub-rectangles, custom pivot origins, rotation, scaling, and color tinting.
- ✍️ **Embedded Font Engine**: Built-in bitmap font atlases embedded via `go:embed` — text works anywhere immediately without missing asset crashes.
- 📷 **2D Camera System**: Full camera transforms (target tracking, viewport offset, zoom, rotation) and screen-to-world / world-to-screen matrix transformations.
- 💥 **2D Collision Engine**: Fast collision algorithms for AABB rectangles, circles, points, and circle-to-box intersections.
- 🎮 **Thread-Safe Input**: Clean keyboard, mouse buttons, cursor coordinates, delta movement, mouse wheel scrolling, and cursor lock modes.
- ⏱️ **Frame Timing**: High-precision frame limiter, delta time calculations, and real-time FPS counter.
- 📸 **Built-in Utilities**: Screenshot capture to PNG, RGBA / Hex color constructors, and vector mathematics.

---

## 📦 Installation

```bash
go get github.com/QOthman/Pixu
```

### Prerequisites
Pixu requires OpenGL 3.3+ and GLFW dependencies:
- **Windows**: Works out of the box with MinGW-w64 / CGO.
- **Linux**: Requires X11 / Wayland and OpenGL development libraries:
  ```bash
  sudo apt-get install libgl1-mesa-dev xorg-dev
  ```
- **macOS**: Requires Xcode Command Line Tools.

---

## 🚀 Quick Start

Here is a minimal, complete application:

```go
package main

import (
	"log"

	"github.com/QOthman/Pixu"
)

func main() {
	// Initialize an 800x600 window
	if err := pixu.Init(800, 600, "Pixu - Hello World"); err != nil {
		log.Fatal(err)
	}
	defer pixu.Close()

	for pixu.ShouldContinue() {
		// Clear background
		pixu.ClearBackground(pixu.Black)

		// Draw shapes
		pixu.DrawCircle(400, 300, 60, pixu.Gold)
		pixu.DrawCircleOutlineThick(400, 300, 65, 3.0, pixu.White)

		// Draw text
		pixu.DrawTextCentered("Welcome to Pixu!", 400, 300, 1.2, pixu.White)

		// End frame (pacing + input polling + buffer swap)
		pixu.EndDrawing()
	}
}
```

---

## 📚 Examples Gallery

Explore ready-to-run interactive examples in the [`examples/`](./examples) directory:

| Example | Description |
| :--- | :--- |
| **[`01_shapes`](./examples/01_shapes/main.go)** | Full gallery of lines, beziers, rounded rects, gradients, polygons, rings, and ellipses. |
| **[`02_sprites_and_textures`](./examples/02_sprites_and_textures/main.go)** | Texture loading, scaling, rotations, tinting, alpha channels, and custom pivot points. |
| **[`03_input_and_movement`](./examples/03_input_and_movement/main.go)** | WASD/arrow keyboard navigation, laser mouse aiming, shooting mechanics, and interactive UI buttons. |
| **[`04_camera2d`](./examples/04_camera2d/main.go)** | 2D Camera world viewport with zoom, rotation, smooth player follow (lerp), and world coordinate mapping. |
| **[`05_collisions`](./examples/05_collisions/main.go)** | Interactive circle-to-box, circle-to-circle, point-in-rect collision queries, and bouncing physics. |
| **[`simple_demo`](./examples/simple_demo.go)** | Upgraded all-in-one simple demonstration. |

Run any example with:
```bash
go run ./examples/01_shapes
```

---

## 🛠️ API Overview

### 1. Window & Lifecycle
```go
pixu.Init(width, height int, title string) error
pixu.InitWithConfig(cfg WindowConfig) error
pixu.ShouldContinue() bool
pixu.ShouldClose() bool
pixu.Close()
pixu.BeginDrawing()
pixu.EndDrawing()
pixu.ClearBackground(color Color)
pixu.SetWindowTitle(title string)
pixu.SetWindowSize(w, h int)
pixu.TakeScreenshot("capture.png")
```

### 2. Shapes & Geometry
```go
pixu.DrawLine(x1, y1, x2, y2 float32, color Color)
pixu.DrawLineThick(x1, y1, x2, y2, thickness float32, color Color)
pixu.DrawLineBezier(start, end, control Vec2, thickness float32, color Color)
pixu.DrawRectangle(x, y, w, h float32, color Color)
pixu.DrawRectangleOutline(x, y, w, h float32, color Color)
pixu.DrawRectangleRounded(x, y, w, h, roundness float32, segments int, color Color)
pixu.DrawRectangleGradientH(x, y, w, h float32, left, right Color)
pixu.DrawTriangle(x1, y1, x2, y2, x3, y3 float32, color Color)
pixu.DrawCircle(cx, cy, radius float32, color Color)
pixu.DrawCircleOutlineThick(cx, cy, radius, thickness float32, color Color)
pixu.DrawCircleSector(cx, cy, radius, startAngle, endAngle float32, segments int, color Color)
pixu.DrawRing(cx, cy, innerR, outerR, startAngle, endAngle float32, segments int, color Color)
pixu.DrawPolygon(cx, cy float32, sides int, radius, rotation float32, color Color)
```

### 3. Textures & Sprites
```go
img, err := pixu.LoadImage("player.png")
defer img.Delete()

pixu.DrawImage(img, x, y)
pixu.DrawImageScaled(img, x, y, scaleX, scaleY)
pixu.DrawImageRotated(img, x, y, rotation)
pixu.DrawImageTinted(img, x, y, pixu.Red)
pixu.DrawImageRec(img, srcRect, destPos, tint)
pixu.DrawImagePro(img, srcRect, destRect, origin, rotation, tint)
```

### 4. Text & Fonts
```go
pixu.DrawText(text string, x, y, size float32, color Color)
pixu.DrawTextCentered(text string, cx, cy, size float32, color Color)
pixu.DrawTextWithBackground(text string, x, y, size float32, bgColor, textColor Color)
pixu.DrawTextOutline(text string, x, y, size float32, textColor, outlineColor Color)
w, h := pixu.MeasureText(text, size)
```

### 5. Input Management
```go
if pixu.IsKeyDown(pixu.KeySpace) { ... }
if pixu.IsKeyJustPressed(pixu.KeyEscape) { ... }
if pixu.IsMouseButtonJustPressed(pixu.MouseLeft) { ... }

mousePos := pixu.GetMousePositionV()
dx, dy := pixu.GetMouseDelta()
wheel := pixu.GetMouseWheelMove()
```

### 6. 2D Camera
```go
camera := pixu.Camera2D{
    Offset: pixu.V2(400, 300),
    Target: playerPos,
    Zoom:   1.0,
}

pixu.BeginMode2D(camera)
    // Draw all world-space objects here
pixu.EndMode2D()

// Screen to world conversion
worldPos := pixu.GetScreenToWorld2D(mousePos, camera)
```

### 7. 2D Math & Collisions
```go
v := pixu.V2(100, 200).Add(pixu.V2(10, 0)).Normalize()
rect := pixu.NewRect(0, 0, 100, 100)

hit := pixu.CheckCollisionRecs(rect1, rect2)
hit = pixu.CheckCollisionCircles(center1, radius1, center2, radius2)
hit = pixu.CheckCollisionCircleRec(circleCenter, radius, rect)
```

---

## 📄 License

Pixu is licensed under the [MIT License](./LICENSE).
