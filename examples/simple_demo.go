package main

import (
	"fmt"
	"log"

	"github.com/QOthman/Pixu"
)

func main() {
	// Initialize window with 800x600 resolution and title
	if err := pixu.Init(800, 600, "Pixu - Simple Graphics Demo"); err != nil {
		log.Fatal("Failed to initialize Pixu graphics:", err)
	}
	defer pixu.Close()

	playerX, playerY := float32(400), float32(300)

	// Try loading optional background image (relative to current or examples directory)
	var backgroundImg *pixu.Image
	var err error
	for _, path := range []string{"image.jpg", "examples/image.jpg"} {
		backgroundImg, err = pixu.LoadImage(path)
		if err == nil {
			break
		}
	}
	if backgroundImg != nil {
		defer backgroundImg.Delete()
	}

	// Main game loop
	for pixu.ShouldContinue() {
		// Handle input
		moveSpeed := float32(200.0) * pixu.GetDeltaTime()
		if pixu.IsKeyPressed(pixu.KeyUp) || pixu.IsKeyPressed(pixu.KeyW) {
			playerY -= moveSpeed
		}
		if pixu.IsKeyPressed(pixu.KeyDown) || pixu.IsKeyPressed(pixu.KeyS) {
			playerY += moveSpeed
		}
		if pixu.IsKeyPressed(pixu.KeyLeft) || pixu.IsKeyPressed(pixu.KeyA) {
			playerX -= moveSpeed
		}
		if pixu.IsKeyPressed(pixu.KeyRight) || pixu.IsKeyPressed(pixu.KeyD) {
			playerX += moveSpeed
		}

		// Clear screen
		pixu.ClearBackground(pixu.BLACK)

		// Draw background image if available
		if backgroundImg != nil {
			pixu.DrawImage(backgroundImg, 0, 0)
		}

		// Draw shapes
		pixu.DrawLineThick(100, 100, 300, 200, 3.0, pixu.RED)
		pixu.DrawTriangle(400, 100, 350, 200, 450, 200, pixu.GREEN)
		pixu.DrawRectangleRounded(500, 150, 120, 80, 0.3, 8, pixu.BLUE)
		pixu.DrawCircle(playerX, playerY, 40, pixu.YELLOW)
		pixu.DrawCircleOutlineThick(playerX, playerY, 45, 3.0, pixu.WHITE)

		// Draw text
		pixu.DrawTextCentered("GAME OVER", 400, 200, 1.5, pixu.GREEN)
		pixu.DrawTextWithBackground("RESTART", 300, 400, 1.5, pixu.RED, pixu.WHITE)
		pixu.DrawTextOutline("EPIC!", 400, 500, 2.5, pixu.YELLOW, pixu.RED)

		// Draw HUD
		pixu.DrawText(fmt.Sprintf("FPS: %d", pixu.GetFps()), 15, 15, 1.0, pixu.WHITE)
		pixu.DrawText("Use WASD / Arrow Keys to move", 15, 45, 0.8, pixu.LIGHTGRAY)

		// Finish frame
		pixu.EndDrawing()
	}
}
