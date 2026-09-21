package pixu

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Color represents an RGBA color with float32 components in the range [0.0, 1.0].
type Color struct {
	R float32 // Red component [0.0, 1.0]
	G float32 // Green component [0.0, 1.0]
	B float32 // Blue component [0.0, 1.0]
	A float32 // Alpha component [0.0, 1.0]
}

// Predefined palette colors (TitleCase)
var (
	Blank       = Color{0.0, 0.0, 0.0, 0.0}
	White       = Color{1.0, 1.0, 1.0, 1.0}
	Black       = Color{0.0, 0.0, 0.0, 1.0}
	Red         = Color{0.9, 0.16, 0.22, 1.0}
	Green       = Color{0.0, 0.89, 0.31, 1.0}
	Blue        = Color{0.0, 0.47, 0.95, 1.0}
	Yellow      = Color{0.99, 0.86, 0.0, 1.0}
	Orange      = Color{1.0, 0.63, 0.0, 1.0}
	Pink        = Color{1.0, 0.43, 0.78, 1.0}
	Purple      = Color{0.78, 0.31, 0.88, 1.0}
	Violet      = Color{0.53, 0.24, 0.75, 1.0}
	DarkPurple  = Color{0.44, 0.12, 0.49, 1.0}
	Cyan        = Color{0.0, 0.89, 0.89, 1.0}
	Magenta     = Color{1.0, 0.0, 1.0, 1.0}
	Lime        = Color{0.0, 0.62, 0.18, 1.0}
	DarkGreen   = Color{0.0, 0.46, 0.17, 1.0}
	SkyBlue     = Color{0.4, 0.75, 1.0, 1.0}
	DarkBlue    = Color{0.0, 0.32, 0.67, 1.0}
	Teal        = Color{0.0, 0.5, 0.5, 1.0}
	Gold        = Color{1.0, 0.8, 0.0, 1.0}
	Maroon      = Color{0.75, 0.13, 0.22, 1.0}
	Beige       = Color{0.83, 0.69, 0.51, 1.0}
	Brown       = Color{0.5, 0.3, 0.16, 1.0}
	DarkBrown   = Color{0.3, 0.17, 0.09, 1.0}
	LightGray   = Color{0.78, 0.78, 0.78, 1.0}
	Gray        = Color{0.51, 0.51, 0.51, 1.0}
	DarkGray    = Color{0.31, 0.31, 0.31, 1.0}
)

// Predefined palette colors (UPPERCASE for backward compatibility)
var (
	BLANK      = Blank
	WHITE      = White
	BLACK      = Black
	RED        = Red
	GREEN      = Green
	BLUE       = Blue
	YELLOW     = Yellow
	ORANGE     = Orange
	PINK       = Pink
	PURPLE     = Purple
	VIOLET     = Violet
	DARKPURPLE = DarkPurple
	CYAN       = Cyan
	MAGENTA    = Magenta
	LIME       = Lime
	DARKGREEN  = DarkGreen
	SKYBLUE    = SkyBlue
	DARKBLUE   = DarkBlue
	TEAL       = Teal
	GOLD       = Gold
	MAROON     = Maroon
	BEIGE      = Beige
	BROWN      = Brown
	DARKBROWN  = DarkBrown
	LIGHTGRAY  = LightGray
	GRAY       = Gray
	DARKGRAY   = DarkGray
)

// NewColor creates a new Color from float32 RGBA values in range [0.0, 1.0].
func NewColor(r, g, b, a float32) Color {
	return Color{r, g, b, a}
}

// RGB creates a new Color with alpha set to 1.0.
func RGB(r, g, b float32) Color {
	return Color{r, g, b, 1.0}
}

// RGBA creates a new Color from uint8 values in range [0, 255].
func RGBA(r, g, b, a uint8) Color {
	return Color{
		R: float32(r) / 255.0,
		G: float32(g) / 255.0,
		B: float32(b) / 255.0,
		A: float32(a) / 255.0,
	}
}

// Hex creates a Color from a 32-bit hex integer (0xRRGGBBAA or 0xRRGGBB).
// If hex <= 0xFFFFFF, alpha defaults to 1.0 (255).
func Hex(hex uint32) Color {
	if hex <= 0xFFFFFF {
		r := uint8((hex >> 16) & 0xFF)
		g := uint8((hex >> 8) & 0xFF)
		b := uint8(hex & 0xFF)
		return RGBA(r, g, b, 255)
	}
	r := uint8((hex >> 24) & 0xFF)
	g := uint8((hex >> 16) & 0xFF)
	b := uint8((hex >> 8) & 0xFF)
	a := uint8(hex & 0xFF)
	return RGBA(r, g, b, a)
}

// HexAlpha creates a Color from a 24-bit hex integer (0xRRGGBB) and a float32 alpha [0.0, 1.0].
func HexAlpha(hex uint32, alpha float32) Color {
	r := float32((hex>>16)&0xFF) / 255.0
	g := float32((hex>>8)&0xFF) / 255.0
	b := float32(hex&0xFF) / 255.0
	return Color{R: r, G: g, B: b, A: alpha}
}

// ColorFromHex parses a hex color string like "#FF5733", "FF5733", "#FF5733AA", or "FF5733AA".
func ColorFromHex(hexStr string) (Color, error) {
	hexStr = strings.TrimPrefix(hexStr, "#")
	if len(hexStr) == 6 {
		val, err := strconv.ParseUint(hexStr, 16, 32)
		if err != nil {
			return Blank, err
		}
		return Hex(uint32(val)), nil
	} else if len(hexStr) == 8 {
		val, err := strconv.ParseUint(hexStr, 16, 32)
		if err != nil {
			return Blank, err
		}
		return Hex(uint32(val)), nil
	}
	return Blank, fmt.Errorf("invalid hex color length: %s", hexStr)
}

// WithAlpha returns a copy of the color with a modified alpha value.
func (c Color) WithAlpha(a float32) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: a}
}

// Lerp linearly interpolates between this color and a target color by factor t [0.0, 1.0].
func (c Color) Lerp(target Color, t float32) Color {
	t = float32(math.Max(0.0, math.Min(1.0, float64(t))))
	return Color{
		R: c.R + (target.R-c.R)*t,
		G: c.G + (target.G-c.G)*t,
		B: c.B + (target.B-c.B)*t,
		A: c.A + (target.A-c.A)*t,
	}
}

// Brightness scales the RGB components by the given factor while keeping Alpha unchanged.
func (c Color) Brightness(factor float32) Color {
	return Color{
		R: float32(math.Max(0.0, math.Min(1.0, float64(c.R*factor)))),
		G: float32(math.Max(0.0, math.Min(1.0, float64(c.G*factor)))),
		B: float32(math.Max(0.0, math.Min(1.0, float64(c.B*factor)))),
		A: c.A,
	}
}

// ToRGBA converts the Color to 8-bit unsigned integer channels [0, 255].
func (c Color) ToRGBA() (r, g, b, a uint8) {
	return uint8(c.R * 255.0), uint8(c.G * 255.0), uint8(c.B * 255.0), uint8(c.A * 255.0)
}

// ToHex returns the 32-bit RGBA integer representation (0xRRGGBBAA).
func (c Color) ToHex() uint32 {
	r, g, b, a := c.ToRGBA()
	return (uint32(r) << 24) | (uint32(g) << 16) | (uint32(b) << 8) | uint32(a)
}

// Equals checks if two colors have identical RGBA values.
func (c Color) Equals(other Color) bool {
	return c.R == other.R && c.G == other.G && c.B == other.B && c.A == other.A
}
