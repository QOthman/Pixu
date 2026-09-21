package pixu

import (
	"testing"
)

func TestColorRGBA(t *testing.T) {
	c := RGBA(255, 128, 0, 255)
	if c.R < 0.99 || c.R > 1.01 {
		t.Errorf("expected R close to 1.0, got %f", c.R)
	}
	if c.G < 0.49 || c.G > 0.51 {
		t.Errorf("expected G close to 0.5, got %f", c.G)
	}
	if c.B != 0.0 {
		t.Errorf("expected B == 0, got %f", c.B)
	}
	if c.A != 1.0 {
		t.Errorf("expected A == 1.0, got %f", c.A)
	}

	r, g, b, a := c.ToRGBA()
	if r != 255 || g != 128 || b != 0 || a != 255 {
		t.Errorf("unexpected ToRGBA output: (%d, %d, %d, %d)", r, g, b, a)
	}
}

func TestColorHex(t *testing.T) {
	c := Hex(0xFF0000FF)
	if c.R < 0.99 || c.G != 0.0 || c.B != 0.0 || c.A != 1.0 {
		t.Errorf("Hex(0xFF0000FF) parsed incorrectly: %+v", c)
	}

	cParsed, err := ColorFromHex("#00FF00")
	if err != nil {
		t.Fatalf("failed to parse hex #00FF00: %v", err)
	}
	if cParsed.G < 0.99 || cParsed.R != 0 || cParsed.B != 0 {
		t.Errorf("ColorFromHex(#00FF00) incorrect: %+v", cParsed)
	}
}

func TestColorLerp(t *testing.T) {
	c1 := Black
	c2 := White
	mid := c1.Lerp(c2, 0.5)

	if mid.R < 0.49 || mid.R > 0.51 {
		t.Errorf("expected lerp R 0.5, got %f", mid.R)
	}
}
