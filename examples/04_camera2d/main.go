package main

import (
	"fmt"
	"log"

	"github.com/QOthman/Pixu"
)

func main() {
	if err := pixu.Init(900, 650, "Pixu - 2D Camera System Demo"); err != nil {
		log.Fatal("Failed to initialize Pixu:", err)
	}
	defer pixu.Close()

	playerPos := pixu.V2(0, 0)
	camera := pixu.Camera2D{
		Offset:   pixu.V2(450, 325), // Center of the 900x650 viewport
		Target:   playerPos,
		Rotation: 0,
		Zoom:     1.0,
	}

	// World objects
	buildings := []pixu.Rect{
		pixu.NewRect(-300, -200, 150, 150),
		pixu.NewRect(100, -250, 200, 120),
		pixu.NewRect(-250, 150, 180, 100),
		pixu.NewRect(150, 100, 140, 200),
		pixu.NewRect(-450, -400, 100, 800),
		pixu.NewRect(400, -400, 100, 800),
	}

	for pixu.ShouldContinue() {
		dt := pixu.GetDeltaTime()

		// Player Controls
		speed := float32(300.0) * dt
		if pixu.IsKeyDown(pixu.KeyW) || pixu.IsKeyDown(pixu.KeyUp) {
			playerPos.Y -= speed
		}
		if pixu.IsKeyDown(pixu.KeyS) || pixu.IsKeyDown(pixu.KeyDown) {
			playerPos.Y += speed
		}
		if pixu.IsKeyDown(pixu.KeyA) || pixu.IsKeyDown(pixu.KeyLeft) {
			playerPos.X -= speed
		}
		if pixu.IsKeyDown(pixu.KeyD) || pixu.IsKeyDown(pixu.KeyRight) {
			playerPos.X += speed
		}

		// Camera Zoom Controls (Mouse Wheel or +/- keys)
		wheel := pixu.GetMouseWheelMove()
		if wheel != 0 {
			camera.Zoom += wheel * 0.1
			if camera.Zoom < 0.25 {
				camera.Zoom = 0.25
			}
			if camera.Zoom > 3.0 {
				camera.Zoom = 3.0
			}
		}

		// Camera Rotation (Q / E)
		if pixu.IsKeyDown(pixu.KeyQ) {
			camera.Rotation -= 60.0 * dt
		}
		if pixu.IsKeyDown(pixu.KeyE) {
			camera.Rotation += 60.0 * dt
		}
		if pixu.IsKeyDown(pixu.KeyR) {
			camera.Rotation = 0
			camera.Zoom = 1.0
		}

		// Smooth camera follow (Lerp)
		camera.Target = camera.Target.Lerp(playerPos, 5.0*dt)

		// Convert screen mouse position to world coordinates
		mouseScreen := pixu.GetMousePositionV()
		mouseWorld := pixu.GetScreenToWorld2D(mouseScreen, camera)

		// RENDER
		pixu.ClearBackground(pixu.Hex(0x1a1a24FF))

		// ---------------------------------------------
		// 1. WORLD SPACE RENDERING (inside BeginMode2D)
		// ---------------------------------------------
		pixu.BeginMode2D(camera)

		// Draw World Grid
		for x := float32(-1000); x <= 1000; x += 100 {
			pixu.DrawLine(x, -1000, x, 1000, pixu.HexAlpha(0xFFFFFF, 0.08))
		}
		for y := float32(-1000); y <= 1000; y += 100 {
			pixu.DrawLine(-1000, y, 1000, y, pixu.HexAlpha(0xFFFFFF, 0.08))
		}

		// Draw World Boundary
		pixu.DrawRectangleOutlineThick(-500, -500, 1000, 1000, 4.0, pixu.Red)

		// Draw Buildings
		for i, b := range buildings {
			c := pixu.Blue
			if i%2 == 1 {
				c = pixu.Purple
			}
			pixu.DrawRectangleRoundedRec(b, 0.15, 8, c)
			pixu.DrawRectangleOutlineRec(b, 2.0, pixu.White.WithAlpha(0.5))
		}

		// Draw World Origin Marker
		pixu.DrawCircle(0, 0, 8, pixu.Gold)
		pixu.DrawText("(0,0)", 12, -8, 0.7, pixu.Gold)

		// Draw Player
		pixu.DrawCircle(playerPos.X, playerPos.Y, 20, pixu.Lime)
		pixu.DrawCircleOutlineThick(playerPos.X, playerPos.Y, 22, 3.0, pixu.White)

		// Draw World Cursor Indicator
		pixu.DrawCircleOutline(mouseWorld.X, mouseWorld.Y, 10, pixu.Orange)

		pixu.EndMode2D()

		// ---------------------------------------------
		// 2. SCREEN SPACE HUD (outside Camera2D)
		// ---------------------------------------------
		pixu.DrawRectangle(15, 15, 360, 130, pixu.Black.WithAlpha(0.7))
		pixu.DrawRectangleOutline(15, 15, 360, 130, pixu.White.WithAlpha(0.3))

		pixu.DrawText("CAMERA 2D CONTROLS", 25, 25, 0.9, pixu.Yellow)
		pixu.DrawText("WASD: Move Player", 25, 50, 0.75, pixu.White)
		pixu.DrawText("Mouse Wheel: Zoom In / Out", 25, 70, 0.75, pixu.White)
		pixu.DrawText("Q / E: Rotate Camera | R: Reset", 25, 90, 0.75, pixu.White)
		pixu.DrawText(fmt.Sprintf("Zoom: %.2fx | Rot: %.1f deg", camera.Zoom, camera.Rotation), 25, 115, 0.75, pixu.Cyan)

		pixu.DrawText(fmt.Sprintf("Player World Pos: (%.1f, %.1f)", playerPos.X, playerPos.Y), 25, 615, 0.8, pixu.LightGray)
		pixu.DrawText(fmt.Sprintf("Mouse World Pos: (%.1f, %.1f)", mouseWorld.X, mouseWorld.Y), 450, 615, 0.8, pixu.Orange)

		pixu.EndDrawing()
	}
}
