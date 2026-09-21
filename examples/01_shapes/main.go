package main

import (
	"fmt"
	"log"
	"math"

	"github.com/QOthman/Pixu"
)

func main() {
	cfg := pixu.DefaultConfig()
	cfg.Width = 900
	cfg.Height = 650
	cfg.Title = "Pixu - 2D Shapes & Geometry Gallery"
	cfg.MSAA = 4

	if err := pixu.InitWithConfig(cfg); err != nil {
		log.Fatal("Failed to initialize Pixu:", err)
	}
	defer pixu.Close()

	for pixu.ShouldContinue() {
		time := float32(pixu.GetTime())

		pixu.ClearBackground(pixu.Hex(0x181824FF))

		// Header
		pixu.DrawTextCentered("PIXU 2D SHAPES & GEOMETRY", 450, 30, 1.2, pixu.White)

		// 1. Lines & Curves
		pixu.DrawText("1. Lines & Beziers", 50, 70, 0.9, pixu.Yellow)
		pixu.DrawLineThick(50, 100, 200, 150, 4.0, pixu.Red)
		controlPoint := pixu.V2(125+float32(math.Sin(float64(time*2)))*40, 70)
		pixu.DrawLineBezier(pixu.V2(50, 170), pixu.V2(200, 170), controlPoint, 3.0, pixu.Cyan)

		// 2. Rectangles & Gradients
		pixu.DrawText("2. Rectangles & Gradients", 250, 70, 0.9, pixu.Yellow)
		pixu.DrawRectangle(250, 100, 80, 80, pixu.Blue)
		pixu.DrawRectangleOutlineThick(345, 100, 80, 80, 3.0, pixu.Lime)
		pixu.DrawRectangleGradientH(250, 195, 175, 40, pixu.Purple, pixu.Orange)

		// 3. Rounded Rectangles
		pixu.DrawText("3. Rounded Rectangles", 460, 70, 0.9, pixu.Yellow)
		pixu.DrawRectangleRounded(460, 100, 180, 60, 0.4, 16, pixu.Maroon)
		pixu.DrawRectangleRounded(460, 175, 180, 60, 0.25, 16, pixu.DarkGreen)

		// 4. Triangles & Polygons
		pixu.DrawText("4. Regular Polygons", 680, 70, 0.9, pixu.Yellow)
		pixu.DrawPolygon(730, 140, 5, 40, time*30, pixu.Gold)
		pixu.DrawPolygonOutline(820, 140, 6, 40, -time*20, 2.0, pixu.SkyBlue)

		// 5. Circles, Ellipses & Sectors
		pixu.DrawText("5. Circles & Ellipses", 50, 260, 0.9, pixu.Yellow)
		pixu.DrawCircle(100, 340, 40, pixu.Pink)
		pixu.DrawCircleOutlineThick(200, 340, 40, 4.0, pixu.Teal)
		pixu.DrawEllipse(100, 430, 50, 30, pixu.Orange)
		pixu.DrawEllipseOutline(200, 430, 50, 30, pixu.Violet)

		// 6. Rings & Sectors
		pixu.DrawText("6. Rings & Sectors", 280, 260, 0.9, pixu.Yellow)
		sweepAngle := float32(math.Mod(float64(time*90), 360))
		pixu.DrawCircleSector(340, 340, 45, 0, sweepAngle, 32, pixu.Green)
		pixu.DrawRing(450, 340, 25, 45, 0, 270, 32, pixu.Magenta)

		// 7. Interactive Mouse Follower
		mx, my := pixu.GetMousePosition()
		pixu.DrawText("7. Interactive Ring", 580, 260, 0.9, pixu.Yellow)
		pixu.DrawRing(660, 340, 20+float32(math.Sin(float64(time*4)))*5, 45, 0, 360, 36, pixu.Gold)

		// Draw subtle cursor highlight
		pixu.DrawCircleOutline(mx, my, 15, pixu.White.WithAlpha(0.6))

		// Footer FPS Info
		pixu.DrawText(fmt.Sprintf("FPS: %d | Time: %.1fs", pixu.GetFPS(), pixu.GetTime()), 20, 620, 0.8, pixu.Gray)

		pixu.EndDrawing()
	}
}
