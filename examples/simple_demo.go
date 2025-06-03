package main

import (
	"fmt"
	"log"
	"main/graphics"
)

func main() {

	if err := graphics.Init(800, 600, "Simple Graphics"); err != nil {
		log.Fatal("Failed to initialize graphics:", err)
	}
	defer graphics.Close()

	playerX, playerY := float32(400), float32(300)

	backgroundImg, err := graphics.LoadImage("image.jpg")
	if err != nil {
		log.Printf("Failed to load background: %v", err)
		return
	}
	defer backgroundImg.Delete()
	
	// Main loop 
	for graphics.ShouldContinue() {
		graphics.ClearBackground(graphics.BLACK)

		// Draw background image
		graphics.DrawImage(backgroundImg,0,0)

		// shapes
		graphics.DrawLine(100, 100, 300, 200, graphics.RED)
		graphics.DrawTriangle(400, 100, 350, 200, 450, 200, graphics.GREEN)
		graphics.DrawRectangle(500, 150, 100, 80, graphics.BLUE)
		graphics.DrawCircle(playerX, playerY, 50, graphics.YELLOW)

		// text
		graphics.DrawTextCentered("GAME OVER", 400, 200, 1.5, graphics.GREEN)
		graphics.DrawTextWithBackground("RESTART", 300, 400, 1.5, graphics.RED, graphics.GREEN)
		graphics.DrawTextOutline("EPIC!", 400, 500, 2.5, graphics.YELLOW, graphics.RED)

		// // Handle input
		if graphics.IsKeyPressed(graphics.KeyUp) {
			playerY -= 2 // Move up
		}
		if graphics.IsKeyPressed(graphics.KeyDown) {
			playerY += 2 // Move down
		}
		if graphics.IsKeyPressed(graphics.KeyLeft) {
			playerX -= 2 // Move left
		}
		if graphics.IsKeyPressed(graphics.KeyRight) {
			playerX += 2 // Move right
		}

		graphics.DrawText(fmt.Sprintf("FPS: %d", graphics.GetFps()), 10, 10, 1.0, graphics.RED)
		
		graphics.Wait()
		graphics.UpdateInput()
		graphics.Present()
	}
}
