package graphics

import (
	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Draw line between two points
func DrawLine(x1, y1, x2, y2 float32, color Color) {
	x1GL, y1GL := screenToGL(x1, y1)
	x2GL, y2GL := screenToGL(x2, y2)

	vertices := []float32{
		x1GL, y1GL, color.R, color.G, color.B,
		x2GL, y2GL, color.R, color.G, color.B,
	}

	drawVertices(vertices, 2, gl.LINES)
}

// Draw triangle between three points
func DrawTriangle(x1, y1, x2, y2, x3, y3 float32, color Color) {
	glX1, glY1 := screenToGL(x1, y1)
	glX2, glY2 := screenToGL(x2, y2)
	glX3, glY3 := screenToGL(x3, y3)

	vertices := []float32{
		glX1, glY1, color.R, color.G, color.B,
		glX2, glY2, color.R, color.G, color.B,
		glX3, glY3, color.R, color.G, color.B,
	}

	drawVertices(vertices, 3, gl.TRIANGLES)
}

// Draw rectangle 
func DrawRectangle(x, y, width, height float32, color Color) {
	glX, glY := screenToGL(x, y)
	glW := width / float32(windowWidth) * 2.0
	glH := height / float32(windowHeight) * 2.0

	vertices := []float32{
		glX, glY, color.R, color.G, color.B, // top-left
		glX + glW, glY, color.R, color.G, color.B, // top-right
		glX + glW, glY - glH, color.R, color.G, color.B, // bottom-right
		glX, glY - glH, color.R, color.G, color.B, // bottom-left
	}

	indices := []uint32{0, 1, 2, 2, 3, 0}
	drawIndexed(vertices, indices)
}

// Draw circle by center and radius
// Circle is drawn using triangle fan with segments
func DrawCircle(centerX, centerY, radius float32, color Color) {
	const segments = 32
	vertices := make([]float32, 0, (segments+2)*5) // center + segments + first point again

	// Center point
	glX, glY := screenToGL(centerX, centerY)
	vertices = append(vertices, glX, glY, color.R, color.G, color.B)

	// Circle points
	for i := 0; i <= segments; i++ {
		angle := float32(i) * 2.0 * math.Pi / segments
		x := centerX + radius*float32(math.Cos(float64(angle)))
		y := centerY + radius*float32(math.Sin(float64(angle)))
		glPx, glPy := screenToGL(x, y)
		vertices = append(vertices, glPx, glPy, color.R, color.G, color.B)
	}

	drawVertices(vertices, int32(len(vertices)/5), gl.TRIANGLE_FAN)
}

// Draw rectangle outline
func DrawRectangleOutline(x, y, width, height float32, color Color) {
	glX, glY := screenToGL(x, y)
	glW := width / float32(windowWidth) * 2.0
	glH := height / float32(windowHeight) * 2.0

	vertices := []float32{
		glX, glY, color.R, color.G, color.B, // top-left
		glX + glW, glY, color.R, color.G, color.B, // top-right
		glX + glW, glY - glH, color.R, color.G, color.B, // bottom-right
		glX, glY - glH, color.R, color.G, color.B, // bottom-left
		glX, glY, color.R, color.G, color.B, // back to start
	}

	drawVertices(vertices, 5, gl.LINE_STRIP)
}
