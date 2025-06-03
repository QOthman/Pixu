package graphics

// Color struct represents a color in RGBA format.
// Each component is a float32 in the range [0.0, 1.0].
type Color struct {
	R, G, B, A float32
}

// Predefined colors
// These colors are commonly used and can be referenced directly.
var (
	BLACK   = Color{0.0, 0.0, 0.0, 1.0}
	WHITE   = Color{1.0, 1.0, 1.0, 1.0}
	RED     = Color{1.0, 0.0, 0.0, 1.0}
	GREEN   = Color{0.0, 1.0, 0.0, 1.0}
	BLUE    = Color{0.0, 0.0, 1.0, 1.0}
	YELLOW  = Color{1.0, 1.0, 0.0, 1.0}
	CYAN    = Color{0.0, 1.0, 1.0, 1.0}
	MAGENTA = Color{1.0, 0.0, 1.0, 1.0}
	GRAY    = Color{0.5, 0.5, 0.5, 1.0}
)

// Create custom color
func NewColor(r, g, b, a float32) Color {
	return Color{r, g, b, a}
}

// RGB color (alpha = 1.0)
func RGB(r, g, b float32) Color {
	return Color{r, g, b, 1.0}
}