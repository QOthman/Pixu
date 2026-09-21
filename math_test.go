package pixu

import (
	"math"
	"testing"
)

func TestVec2Operations(t *testing.T) {
	v1 := V2(3, 4)
	if l := v1.Length(); l != 5.0 {
		t.Errorf("expected Length 5, got %f", l)
	}

	norm := v1.Normalize()
	if math.Abs(float64(norm.Length()-1.0)) > 0.001 {
		t.Errorf("expected normalized length 1, got %f", norm.Length())
	}

	v2 := V2(1, 2)
	sum := v1.Add(v2)
	if sum.X != 4 || sum.Y != 6 {
		t.Errorf("expected Add (4, 6), got (%f, %f)", sum.X, sum.Y)
	}
}

func TestCollisions(t *testing.T) {
	r1 := NewRect(0, 0, 50, 50)
	r2 := NewRect(25, 25, 50, 50)
	r3 := NewRect(100, 100, 50, 50)

	if !CheckCollisionRecs(r1, r2) {
		t.Errorf("expected r1 and r2 to collide")
	}
	if CheckCollisionRecs(r1, r3) {
		t.Errorf("expected r1 and r3 to NOT collide")
	}

	c1 := V2(0, 0)
	c2 := V2(30, 0)
	if !CheckCollisionCircles(c1, 20, c2, 20) {
		t.Errorf("expected circles to collide")
	}

	if !CheckCollisionCircleRec(V2(25, 25), 10, r1) {
		t.Errorf("expected circle inside rect to collide")
	}
}

func TestCamera2DTransforms(t *testing.T) {
	cam := Camera2D{
		Offset:   V2(400, 300),
		Target:   V2(100, 100),
		Rotation: 0,
		Zoom:     2.0,
	}

	// World target (100, 100) should project exactly to screen offset (400, 300)
	screenPos := GetWorldToScreen2D(V2(100, 100), cam)
	if math.Abs(float64(screenPos.X-400)) > 0.01 || math.Abs(float64(screenPos.Y-300)) > 0.01 {
		t.Errorf("expected screen pos (400, 300), got (%f, %f)", screenPos.X, screenPos.Y)
	}

	worldPos := GetScreenToWorld2D(screenPos, cam)
	if math.Abs(float64(worldPos.X-100)) > 0.01 || math.Abs(float64(worldPos.Y-100)) > 0.01 {
		t.Errorf("expected roundtrip world pos (100, 100), got (%f, %f)", worldPos.X, worldPos.Y)
	}
}
