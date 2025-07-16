package main

import (
	"fmt"
	"log"
	"main/graphics"
)

func gameLoop(playerX, playerY float32, img *graphics.Image) func() {
	return func() {
		graphics.ClearBackground(graphics.GRAY)
		graphics.DrawImage(img, 0, 0)
		graphics.DrawImageRotated(img, 100, 100, 180)
		graphics.DrawImageScaled(img, 200, 200, 2, 2)
		graphics.DrawImageTinted(img, 300, 300, graphics.RED)
		graphics.DrawImageEx(img, graphics.DrawOptions{
			X:        400,
			Y:        400,
			Width:    100,
			Height:   100,
			Rotation: 45,
			Tint:     graphics.BLUE,
			SrcX:     0,
			SrcY:     0,
			SrcW:     100,
			SrcH:     100,
		})
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
		graphics.DrawLine(100, 100, 300, 200, graphics.RED)
		graphics.DrawTriangle(400, 100, 350, 200, 450, 200, graphics.GREEN)
		graphics.DrawRectangle(500, 150, 100, 80, graphics.BLUE)
		graphics.DrawCircle(playerX, playerY, 50, graphics.YELLOW)

		graphics.DrawText(fmt.Sprintf("FPS: %d", graphics.GetFps()), 10, 10, 1.0, graphics.RED)
		graphics.DrawTextCentered("GAME OVER", 400, 200, 1, graphics.GREEN)
		graphics.DrawTextWithBackground("RESTART", 300, 400, 1, graphics.RED, graphics.GREEN,5,5)
		graphics.DrawTextOutline("EPIC!", 400, 500, 1, graphics.YELLOW, graphics.RED)
	}

}

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
	defer graphics.DeleteImage(backgroundImg)
	graphics.RunGameLoop(gameLoop(playerX, playerY, backgroundImg))

	// Main loop
	// for graphics.ShouldContinue() {
	// graphics.ClearBackground(graphics.BLACK)
	// Draw background image
	// v , h := graphics.GetWindowSize()
	// graphics.DrawImageEx(backgroundImg, graphics.DrawOptions{
	// 	Width:  float32(v),
	// 	Height: float32(h),
	// })

	// // shapes
	// graphics.DrawLine(100, 100, 300, 200, graphics.RED)
	// graphics.DrawTriangle(400, 100, 350, 200, 450, 200, graphics.GREEN)
	// graphics.DrawRectangle(500, 150, 100, 80, graphics.BLUE)
	// graphics.DrawCircle(playerX, playerY, 50, graphics.YELLOW)

	// // text
	// graphics.DrawTextCentered("GAME OVER", 400, 200, 1.5, graphics.GREEN)
	// graphics.DrawTextWithBackground("RESTART", 300, 400, 1.5, graphics.RED, graphics.GREEN)
	// graphics.DrawTextOutline("EPIC!", 400, 500, 2.5, graphics.YELLOW, graphics.RED)

	// // Handle input

	// if graphics.IsMouseMoved() {
	// 	mouseX, mouseY := graphics.GetMousePosition()
	// 	mouseDeltaX, mouseDeltaY := graphics.GetMouseDelta()
	// 	fmt.Println("Mouse Position:", mouseX, mouseY)
	// 	fmt.Println("Mouse Delta:", mouseDeltaX, mouseDeltaY)
	// }

	// graphics.DrawText(fmt.Sprintf("FPS: %d ! \" ' $", graphics.GetFps()), 10, 10, 1.0, graphics.RED)

	// graphics.Wait()
	// graphics.UpdateInput()
	// graphics.Present()
	// }
}
