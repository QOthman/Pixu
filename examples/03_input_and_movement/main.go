package main

import (
	"fmt"
	"log"
	"math"

	"github.com/QOthman/Pixu"
)

type Bullet struct {
	Pos pixu.Vec2
	Vel pixu.Vec2
}

func main() {
	if err := pixu.Init(800, 600, "Pixu - Input & Movement Demo"); err != nil {
		log.Fatal("Failed to initialize Pixu:", err)
	}
	defer pixu.Close()

	playerPos := pixu.V2(400, 300)
	playerSpeed := float32(250.0)
	var bullets []Bullet

	buttonRect := pixu.NewRect(20, 20, 140, 40)
	buttonColor := pixu.Blue

	for pixu.ShouldContinue() {
		dt := pixu.GetDeltaTime()
		mousePos := pixu.GetMousePositionV()

		// 1. Keyboard Movement (WASD / Arrows)
		moveDir := pixu.V2(0, 0)
		if pixu.IsKeyDown(pixu.KeyW) || pixu.IsKeyDown(pixu.KeyUp) {
			moveDir.Y -= 1
		}
		if pixu.IsKeyDown(pixu.KeyS) || pixu.IsKeyDown(pixu.KeyDown) {
			moveDir.Y += 1
		}
		if pixu.IsKeyDown(pixu.KeyA) || pixu.IsKeyDown(pixu.KeyLeft) {
			moveDir.X -= 1
		}
		if pixu.IsKeyDown(pixu.KeyD) || pixu.IsKeyDown(pixu.KeyRight) {
			moveDir.X += 1
		}

		if moveDir.LengthSqr() > 0 {
			playerPos = playerPos.Add(moveDir.Normalize().Scale(playerSpeed * dt))
		}

		// Clamp player inside window
		playerPos.X = pixu.Clamp(playerPos.X, 20, 780)
		playerPos.Y = pixu.Clamp(playerPos.Y, 20, 580)

		// 2. Mouse Click Shooting
		if pixu.IsMouseButtonJustPressed(pixu.MouseLeft) {
			aimDir := mousePos.Sub(playerPos).Normalize()
			bullets = append(bullets, Bullet{
				Pos: playerPos,
				Vel: aimDir.Scale(500.0),
			})
		}

		// 3. UI Button Click Test
		buttonHovered := buttonRect.Contains(mousePos)
		if buttonHovered {
			if pixu.IsMouseButtonDown(pixu.MouseLeft) {
				buttonColor = pixu.Red
			} else {
				buttonColor = pixu.SkyBlue
			}
			if pixu.IsMouseButtonJustPressed(pixu.MouseLeft) {
				// Reset position button action
				playerPos = pixu.V2(400, 300)
			}
		} else {
			buttonColor = pixu.Blue
		}

		// Update bullets
		var activeBullets []Bullet
		for _, b := range bullets {
			b.Pos = b.Pos.Add(b.Vel.Scale(dt))
			if b.Pos.X >= 0 && b.Pos.X <= 800 && b.Pos.Y >= 0 && b.Pos.Y <= 600 {
				activeBullets = append(activeBullets, b)
			}
		}
		bullets = activeBullets

		// RENDER
		pixu.ClearBackground(pixu.Hex(0x14161DFF))

		// Draw Grid Lines
		for x := float32(0); x <= 800; x += 50 {
			pixu.DrawLine(x, 0, x, 600, pixu.HexAlpha(0xFFFFFF, 0.05))
		}
		for y := float32(0); y <= 600; y += 50 {
			pixu.DrawLine(0, y, 800, y, pixu.HexAlpha(0xFFFFFF, 0.05))
		}

		// Draw Aiming Laser Line from player to mouse
		pixu.DrawLine(playerPos.X, playerPos.Y, mousePos.X, mousePos.Y, pixu.HexAlpha(0xFF3333, 0.3))

		// Draw Bullets
		for _, b := range bullets {
			pixu.DrawCircle(b.Pos.X, b.Pos.Y, 5, pixu.Yellow)
		}

		// Draw Player Ship (Rotating Triangle aimed at mouse)
		angle := float32(math.Atan2(float64(mousePos.Y-playerPos.Y), float64(mousePos.X-playerPos.X)))
		p1 := playerPos.Add(pixu.V2(float32(math.Cos(float64(angle)))*25, float32(math.Sin(float64(angle)))*25))
		p2 := playerPos.Add(pixu.V2(float32(math.Cos(float64(angle)+2.4))*18, float32(math.Sin(float64(angle)+2.4))*18))
		p3 := playerPos.Add(pixu.V2(float32(math.Cos(float64(angle)-2.4))*18, float32(math.Sin(float64(angle)-2.4))*18))

		pixu.DrawTriangle(p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y, pixu.Lime)
		pixu.DrawTriangleOutline(p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y, pixu.White)

		// Draw Mouse Reticle
		pixu.DrawCircleOutline(mousePos.X, mousePos.Y, 12, pixu.Red)
		pixu.DrawLine(mousePos.X-16, mousePos.Y, mousePos.X+16, mousePos.Y, pixu.Red)
		pixu.DrawLine(mousePos.X, mousePos.Y-16, mousePos.X, mousePos.Y+16, pixu.Red)

		// Draw Reset UI Button
		pixu.DrawRectangleRoundedRec(buttonRect, 0.3, 8, buttonColor)
		pixu.DrawRectangleRoundedOutline(buttonRect.X, buttonRect.Y, buttonRect.Width, buttonRect.Height, 0.3, 2.0, 8, pixu.White)
		pixu.DrawTextCentered("Reset Pos", buttonRect.X+buttonRect.Width/2, buttonRect.Y+buttonRect.Height/2, 0.8, pixu.White)

		// Draw Instructions
		pixu.DrawText("Controls: [WASD / Arrows] Move | [Left Click] Shoot", 180, 30, 0.8, pixu.White)
		pixu.DrawText(fmt.Sprintf("Player: (%.0f, %.0f) | Bullets: %d | FPS: %d", playerPos.X, playerPos.Y, len(bullets), pixu.GetFPS()), 20, 560, 0.8, pixu.Gray)

		pixu.EndDrawing()
	}
}
