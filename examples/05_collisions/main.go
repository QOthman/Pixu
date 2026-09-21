package main

import (
	"fmt"
	"log"
	"math"

	"github.com/QOthman/Pixu"
)

func main() {
	if err := pixu.Init(900, 650, "Pixu - 2D Collision Detection Demo"); err != nil {
		log.Fatal("Failed to initialize Pixu:", err)
	}
	defer pixu.Close()

	// Static target shapes
	targetBox := pixu.NewRect(150, 180, 140, 120)
	targetCirclePos := pixu.V2(550, 240)
	targetCircleRadius := float32(70.0)

	// Bouncing ball in bottom area
	ballPos := pixu.V2(450, 500)
	ballVel := pixu.V2(220, 160)
	ballRadius := float32(20.0)
	arenaRect := pixu.NewRect(100, 420, 700, 180)

	for pixu.ShouldContinue() {
		dt := pixu.GetDeltaTime()
		mousePos := pixu.GetMousePositionV()

		// Update bouncing ball
		ballPos = ballPos.Add(ballVel.Scale(dt))
		if ballPos.X-ballRadius < arenaRect.X {
			ballPos.X = arenaRect.X + ballRadius
			ballVel.X = -ballVel.X
		}
		if ballPos.X+ballRadius > arenaRect.X+arenaRect.Width {
			ballPos.X = arenaRect.X + arenaRect.Width - ballRadius
			ballVel.X = -ballVel.X
		}
		if ballPos.Y-ballRadius < arenaRect.Y {
			ballPos.Y = arenaRect.Y + ballRadius
			ballVel.Y = -ballVel.Y
		}
		if ballPos.Y+ballRadius > arenaRect.Y+arenaRect.Height {
			ballPos.Y = arenaRect.Y + arenaRect.Height - ballRadius
			ballVel.Y = -ballVel.Y
		}

		// Interactive cursor shapes
		cursorCircleRadius := float32(40.0)
		cursorBox := pixu.NewRect(mousePos.X-40, mousePos.Y-40, 80, 80)

		// Test 1: Point vs Rect
		pointInBox := targetBox.Contains(mousePos)

		// Test 2: Circle (Mouse) vs Rect (Box)
		circleBoxColliding := pixu.CheckCollisionCircleRec(mousePos, cursorCircleRadius, targetBox)

		// Test 3: Circle (Mouse) vs Circle (Target)
		circleCircleColliding := pixu.CheckCollisionCircles(mousePos, cursorCircleRadius, targetCirclePos, targetCircleRadius)

		// Test 4: Rect (Mouse Box) vs Rect (Target Box)
		boxBoxColliding := pixu.CheckCollisionRecs(cursorBox, targetBox)

		// RENDER
		pixu.ClearBackground(pixu.Hex(0x131722FF))

		pixu.DrawTextCentered("2D COLLISION DETECTION SYSTEM", 450, 30, 1.2, pixu.White)
		pixu.DrawText("Move your mouse over the shapes to test various collision algorithms", 160, 60, 0.75, pixu.LightGray)

		// 1. Draw Target Box
		boxColor := pixu.Blue
		if circleBoxColliding || pointInBox {
			boxColor = pixu.Red
		}
		pixu.DrawRectangleRec(targetBox, boxColor.WithAlpha(0.6))
		pixu.DrawRectangleOutlineRec(targetBox, 3.0, pixu.White)
		pixu.DrawTextCentered("Target Box", targetBox.Center().X, targetBox.Center().Y, 0.8, pixu.White)

		// 2. Draw Target Circle
		circleColor := pixu.Purple
		if circleCircleColliding {
			circleColor = pixu.Red
		}
		pixu.DrawCircleV(targetCirclePos, targetCircleRadius, circleColor.WithAlpha(0.6))
		pixu.DrawCircleOutlineThick(targetCirclePos.X, targetCirclePos.Y, targetCircleRadius, 3.0, pixu.White)
		pixu.DrawTextCentered("Target Circle", targetCirclePos.X, targetCirclePos.Y, 0.8, pixu.White)

		// 3. Draw Cursor Circle / Tester
		cursorColor := pixu.Lime
		if circleBoxColliding || circleCircleColliding {
			cursorColor = pixu.Yellow
		}
		pixu.DrawCircleOutlineThick(mousePos.X, mousePos.Y, cursorCircleRadius, 2.5, cursorColor)
		pixu.DrawCircle(mousePos.X, mousePos.Y, 4, pixu.White) // Center point

		// 4. Draw Collision Status HUD
		statusY := float32(140)
		pixu.DrawText("Collision Tests:", 730, statusY, 0.8, pixu.Gold)

		drawStatus := func(label string, hit bool, y float32) {
			col := pixu.Gray
			status := "NO"
			if hit {
				col = pixu.Lime
				status = "YES"
			}
			pixu.DrawText(fmt.Sprintf("%s: %s", label, status), 730, y, 0.7, col)
		}

		drawStatus("Point -> Box", pointInBox, statusY+30)
		drawStatus("Circle -> Box", circleBoxColliding, statusY+55)
		drawStatus("Circle -> Circle", circleCircleColliding, statusY+80)
		drawStatus("Box -> Box", boxBoxColliding, statusY+105)

		// 5. Draw Bouncing Arena
		pixu.DrawRectangleRec(arenaRect, pixu.Black.WithAlpha(0.4))
		pixu.DrawRectangleOutlineRec(arenaRect, 2.0, pixu.Cyan)
		pixu.DrawText("Dynamic Bouncing Particle Arena", arenaRect.X+15, arenaRect.Y+12, 0.75, pixu.Cyan)

		// Draw bouncing ball with speed trail
		trailDir := ballVel.Normalize().Scale(-15)
		pixu.DrawLineThick(ballPos.X, ballPos.Y, ballPos.X+trailDir.X, ballPos.Y+trailDir.Y, 4.0, pixu.Orange.WithAlpha(0.5))
		pixu.DrawCircleV(ballPos, ballRadius, pixu.Orange)
		pixu.DrawCircleOutlineThick(ballPos.X, ballPos.Y, ballRadius, 2.0, pixu.White)

		// Wave oscillation indicator
		waveY := 400 + float32(math.Sin(float64(pixu.GetTime()*4)))*15
		pixu.DrawCircle(60, waveY, 10, pixu.Gold)

		pixu.DrawText(fmt.Sprintf("FPS: %d", pixu.GetFPS()), 20, 620, 0.8, pixu.Gray)

		pixu.EndDrawing()
	}
}
