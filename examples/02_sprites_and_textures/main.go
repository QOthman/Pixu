package main

import (
	"fmt"
	"log"
	"math"

	"github.com/QOthman/Pixu"
)

func main() {
	if err := pixu.Init(800, 600, "Pixu - Textures & Sprites Demo"); err != nil {
		log.Fatal("Failed to initialize Pixu:", err)
	}
	defer pixu.Close()

	// Try loading sample image from examples folder or disk
	var sprite *pixu.Image
	for _, path := range []string{"examples/image.jpg", "image.jpg"} {
		var err error
		sprite, err = pixu.LoadImage(path)
		if err == nil {
			break
		}
	}

	// If no image on disk, create a procedural 64x64 texture
	if sprite == nil {
		var err error
		sprite, err = pixu.NewImage(64, 64, pixu.Cyan)
		if err != nil {
			log.Fatal("Failed to create procedural texture:", err)
		}
	}
	defer sprite.Delete()

	for pixu.ShouldContinue() {
		time := float32(pixu.GetTime())

		pixu.ClearBackground(pixu.Hex(0x101520FF))

		pixu.DrawTextCentered("TEXTURES & SPRITE RENDERING", 400, 30, 1.2, pixu.White)

		// 1. Standard Draw
		pixu.DrawText("Original", 100, 80, 0.8, pixu.Yellow)
		pixu.DrawImage(sprite, 80, 110)

		// 2. Scaled Draw
		scale := 1.0 + float32(math.Sin(float64(time*3)))*0.3
		pixu.DrawText("Dynamic Scale", 300, 80, 0.8, pixu.Yellow)
		pixu.DrawImageScaled(sprite, 280, 110, scale, scale)

		// 3. Rotated around Center Origin
		rot := time * 45.0
		pixu.DrawText("Rotating Sprite", 550, 80, 0.8, pixu.Yellow)
		pixu.DrawImageRotated(sprite, 580, 140, rot)

		// 4. Color Tinting & Alpha
		pixu.DrawText("Color Tinting & Opacity", 100, 320, 0.8, pixu.Yellow)
		pixu.DrawImageTinted(sprite, 80, 350, pixu.Red)
		pixu.DrawImageTinted(sprite, 200, 350, pixu.Green)
		pixu.DrawImageTinted(sprite, 320, 350, pixu.Gold)
		pixu.DrawImageTinted(sprite, 440, 350, pixu.White.WithAlpha(0.4))

		// 5. Transform Pro with Origin
		pixu.DrawText("Custom Pivot / Origin", 600, 320, 0.8, pixu.Yellow)
		w, h := float32(sprite.Width), float32(sprite.Height)
		pixu.DrawImagePro(
			sprite,
			pixu.NewRect(0, 0, w, h),
			pixu.NewRect(650, 400, w*1.5, h*1.5),
			pixu.V2(0, 0), // Top-left pivot
			time*30.0,
			pixu.SkyBlue,
		)

		pixu.DrawText(fmt.Sprintf("FPS: %d", pixu.GetFPS()), 20, 560, 0.8, pixu.Gray)

		pixu.EndDrawing()
	}
}
