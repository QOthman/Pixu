package pixu

import (
	"math"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// ==========================================
// Basic Shape Rendering Helpers
// ==========================================

func drawShapeVertices(vertices []float32, mode uint32) {
	if len(vertices) == 0 {
		return
	}
	gl.UseProgram(shapeShaderProgram)
	gl.BindVertexArray(shapeVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, shapeVBO)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices[0]))*len(vertices), gl.Ptr(vertices), gl.DYNAMIC_DRAW)
	gl.DrawArrays(mode, 0, int32(len(vertices)/6))
	gl.BindVertexArray(0)
}

func drawShapeQuad(vertices []float32) {
	if len(vertices) < 24 {
		return
	}
	gl.UseProgram(shapeShaderProgram)
	gl.BindVertexArray(shapeVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, shapeVBO)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices[0]))*len(vertices), gl.Ptr(vertices), gl.DYNAMIC_DRAW)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}

// ==========================================
// Pixel & Lines
// ==========================================

// DrawPixel renders a single pixel at (x, y).
func DrawPixel(x, y float32, color Color) {
	DrawRectangle(x, y, 1, 1, color)
}

// DrawPixelV renders a single pixel at position Vec2.
func DrawPixelV(pos Vec2, color Color) {
	DrawRectangle(pos.X, pos.Y, 1, 1, color)
}

// DrawLine draws a 1-pixel wide line between (x1, y1) and (x2, y2).
func DrawLine(x1, y1, x2, y2 float32, color Color) {
	vertices := []float32{
		x1, y1, color.R, color.G, color.B, color.A,
		x2, y2, color.R, color.G, color.B, color.A,
	}
	drawShapeVertices(vertices, gl.LINES)
}

// DrawLineV draws a 1-pixel wide line between start and end vectors.
func DrawLineV(start, end Vec2, color Color) {
	DrawLine(start.X, start.Y, end.X, end.Y, color)
}

// DrawLineThick draws a line with custom thickness.
func DrawLineThick(x1, y1, x2, y2, thickness float32, color Color) {
	if thickness <= 1.0 {
		DrawLine(x1, y1, x2, y2, color)
		return
	}

	dx := x2 - x1
	dy := y2 - y1
	length := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	if length == 0 {
		return
	}

	// Normal vector scaled by half-thickness
	nx := (-dy / length) * (thickness / 2.0)
	ny := (dx / length) * (thickness / 2.0)

	vertices := []float32{
		x1 + nx, y1 + ny, color.R, color.G, color.B, color.A,
		x2 + nx, y2 + ny, color.R, color.G, color.B, color.A,
		x2 - nx, y2 - ny, color.R, color.G, color.B, color.A,
		x1 - nx, y1 - ny, color.R, color.G, color.B, color.A,
	}
	drawShapeQuad(vertices)
}

// DrawLineThickV draws a line with custom thickness between vectors.
func DrawLineThickV(start, end Vec2, thickness float32, color Color) {
	DrawLineThick(start.X, start.Y, end.X, end.Y, thickness, color)
}

// DrawLineBezier draws a smooth cubic/quadratic bezier line curve.
func DrawLineBezier(start, end, control Vec2, thickness float32, color Color) {
	const segments = 32
	prev := start
	for i := 1; i <= segments; i++ {
		t := float32(i) / float32(segments)
		u := 1.0 - t
		// Quadratic bezier: B(t) = (1-t)^2 * P0 + 2*(1-t)*t * P1 + t^2 * P2
		curr := Vec2{
			X: u*u*start.X + 2*u*t*control.X + t*t*end.X,
			Y: u*u*start.Y + 2*u*t*control.Y + t*t*end.Y,
		}
		if thickness > 1.0 {
			DrawLineThickV(prev, curr, thickness, color)
		} else {
			DrawLineV(prev, curr, color)
		}
		prev = curr
	}
}

// DrawPolyline draws connected line segments through an array of points.
func DrawPolyline(points []Vec2, color Color) {
	if len(points) < 2 {
		return
	}
	vertices := make([]float32, 0, len(points)*6)
	for _, p := range points {
		vertices = append(vertices, p.X, p.Y, color.R, color.G, color.B, color.A)
	}
	drawShapeVertices(vertices, gl.LINE_STRIP)
}

// ==========================================
// Rectangles
// ==========================================

// DrawRectangle draws a solid filled rectangle.
func DrawRectangle(x, y, width, height float32, color Color) {
	vertices := []float32{
		x, y, color.R, color.G, color.B, color.A,                         // Top-Left
		x + width, y, color.R, color.G, color.B, color.A,                 // Top-Right
		x + width, y + height, color.R, color.G, color.B, color.A,        // Bottom-Right
		x, y + height, color.R, color.G, color.B, color.A,                // Bottom-Left
	}
	drawShapeQuad(vertices)
}

// DrawRectangleV draws a rectangle using position and size vectors.
func DrawRectangleV(pos, size Vec2, color Color) {
	DrawRectangle(pos.X, pos.Y, size.X, size.Y, color)
}

// DrawRectangleRec draws a rectangle using a Rect struct.
func DrawRectangleRec(rec Rect, color Color) {
	DrawRectangle(rec.X, rec.Y, rec.Width, rec.Height, color)
}

// DrawRectangleOutline draws a 1-pixel outline around a rectangle.
func DrawRectangleOutline(x, y, width, height float32, color Color) {
	vertices := []float32{
		x, y, color.R, color.G, color.B, color.A,
		x + width, y, color.R, color.G, color.B, color.A,
		x + width, y + height, color.R, color.G, color.B, color.A,
		x, y + height, color.R, color.G, color.B, color.A,
		x, y, color.R, color.G, color.B, color.A,
	}
	drawShapeVertices(vertices, gl.LINE_STRIP)
}

// DrawRectangleOutlineThick draws a rectangle outline with custom thickness.
func DrawRectangleOutlineThick(x, y, width, height, thickness float32, color Color) {
	if thickness <= 1.0 {
		DrawRectangleOutline(x, y, width, height, color)
		return
	}
	// Top
	DrawRectangle(x, y, width, thickness, color)
	// Bottom
	DrawRectangle(x, y+height-thickness, width, thickness, color)
	// Left
	DrawRectangle(x, y+thickness, thickness, height-2*thickness, color)
	// Right
	DrawRectangle(x+width-thickness, y+thickness, thickness, height-2*thickness, color)
}

// DrawRectangleOutlineRec draws a rectangle outline using a Rect struct.
func DrawRectangleOutlineRec(rec Rect, thickness float32, color Color) {
	DrawRectangleOutlineThick(rec.X, rec.Y, rec.Width, rec.Height, thickness, color)
}

// DrawRectangleGradientH draws a horizontal color gradient across a rectangle.
func DrawRectangleGradientH(x, y, width, height float32, left, right Color) {
	vertices := []float32{
		x, y, left.R, left.G, left.B, left.A,
		x + width, y, right.R, right.G, right.B, right.A,
		x + width, y + height, right.R, right.G, right.B, right.A,
		x, y + height, left.R, left.G, left.B, left.A,
	}
	drawShapeQuad(vertices)
}

// DrawRectangleGradientV draws a vertical color gradient down a rectangle.
func DrawRectangleGradientV(x, y, width, height float32, top, bottom Color) {
	vertices := []float32{
		x, y, top.R, top.G, top.B, top.A,
		x + width, y, top.R, top.G, top.B, top.A,
		x + width, y + height, bottom.R, bottom.G, bottom.B, bottom.A,
		x, y + height, bottom.R, bottom.G, bottom.B, bottom.A,
	}
	drawShapeQuad(vertices)
}

// DrawRectangleGradientEx draws a 4-corner color gradient rectangle.
func DrawRectangleGradientEx(x, y, width, height float32, topLeft, bottomLeft, topRight, bottomRight Color) {
	vertices := []float32{
		x, y, topLeft.R, topLeft.G, topLeft.B, topLeft.A,
		x + width, y, topRight.R, topRight.G, topRight.B, topRight.A,
		x + width, y + height, bottomRight.R, bottomRight.G, bottomRight.B, bottomRight.A,
		x, y + height, bottomLeft.R, bottomLeft.G, bottomLeft.B, bottomLeft.A,
	}
	drawShapeQuad(vertices)
}

// DrawRectangleRounded draws a rounded rectangle with corner radii.
func DrawRectangleRounded(x, y, width, height, roundness float32, segments int, color Color) {
	if roundness <= 0 {
		DrawRectangle(x, y, width, height, color)
		return
	}
	if segments < 4 {
		segments = 8
	}

	maxRadius := float32(math.Min(float64(width), float64(height))) / 2.0
	radius := roundness * maxRadius

	// Central rectangles
	DrawRectangle(x+radius, y, width-2*radius, height, color)
	DrawRectangle(x, y+radius, radius, height-2*radius, color)
	DrawRectangle(x+width-radius, y+radius, radius, height-2*radius, color)

	// Corner fan sectors
	drawSectorFan(x+radius, y+radius, radius, 180, 270, segments, color)
	drawSectorFan(x+width-radius, y+radius, radius, 270, 360, segments, color)
	drawSectorFan(x+width-radius, y+height-radius, radius, 0, 90, segments, color)
	drawSectorFan(x+radius, y+height-radius, radius, 90, 180, segments, color)
}

// DrawRectangleRoundedRec draws a rounded rectangle from a Rect.
func DrawRectangleRoundedRec(rec Rect, roundness float32, segments int, color Color) {
	DrawRectangleRounded(rec.X, rec.Y, rec.Width, rec.Height, roundness, segments, color)
}

// ==========================================
// Triangles
// ==========================================

// DrawTriangle draws a solid filled triangle between three 2D points.
func DrawTriangle(x1, y1, x2, y2, x3, y3 float32, color Color) {
	vertices := []float32{
		x1, y1, color.R, color.G, color.B, color.A,
		x2, y2, color.R, color.G, color.B, color.A,
		x3, y3, color.R, color.G, color.B, color.A,
	}
	drawShapeVertices(vertices, gl.TRIANGLES)
}

// DrawTriangleV draws a solid triangle between three Vec2 vertices.
func DrawTriangleV(v1, v2, v3 Vec2, color Color) {
	DrawTriangle(v1.X, v1.Y, v2.X, v2.Y, v3.X, v3.Y, color)
}

// DrawTriangleOutline draws the outline of a triangle.
func DrawTriangleOutline(x1, y1, x2, y2, x3, y3 float32, color Color) {
	vertices := []float32{
		x1, y1, color.R, color.G, color.B, color.A,
		x2, y2, color.R, color.G, color.B, color.A,
		x3, y3, color.R, color.G, color.B, color.A,
		x1, y1, color.R, color.G, color.B, color.A,
	}
	drawShapeVertices(vertices, gl.LINE_STRIP)
}

// DrawTriangleOutlineV draws the outline of a triangle from vectors.
func DrawTriangleOutlineV(v1, v2, v3 Vec2, color Color) {
	DrawTriangleOutline(v1.X, v1.Y, v2.X, v2.Y, v3.X, v3.Y, color)
}

// ==========================================
// Circles, Ellipses, Sectors & Rings
// ==========================================

// DrawCircle draws a solid filled circle.
func DrawCircle(centerX, centerY, radius float32, color Color) {
	if radius <= 0 {
		return
	}
	const segments = 36
	vertices := make([]float32, 0, (segments+2)*6)

	// Center vertex
	vertices = append(vertices, centerX, centerY, color.R, color.G, color.B, color.A)

	for i := 0; i <= segments; i++ {
		angle := float32(i) * (2.0 * math.Pi / segments)
		px := centerX + radius*float32(math.Cos(float64(angle)))
		py := centerY + radius*float32(math.Sin(float64(angle)))
		vertices = append(vertices, px, py, color.R, color.G, color.B, color.A)
	}

	drawShapeVertices(vertices, gl.TRIANGLE_FAN)
}

// DrawCircleV draws a circle using center vector.
func DrawCircleV(center Vec2, radius float32, color Color) {
	DrawCircle(center.X, center.Y, radius, color)
}

// DrawCircleOutline draws the 1-pixel boundary of a circle.
func DrawCircleOutline(centerX, centerY, radius float32, color Color) {
	if radius <= 0 {
		return
	}
	const segments = 36
	vertices := make([]float32, 0, (segments+1)*6)

	for i := 0; i <= segments; i++ {
		angle := float32(i) * (2.0 * math.Pi / segments)
		px := centerX + radius*float32(math.Cos(float64(angle)))
		py := centerY + radius*float32(math.Sin(float64(angle)))
		vertices = append(vertices, px, py, color.R, color.G, color.B, color.A)
	}

	drawShapeVertices(vertices, gl.LINE_STRIP)
}

// DrawCircleOutlineThick draws a circle outline with custom thickness.
func DrawCircleOutlineThick(centerX, centerY, radius, thickness float32, color Color) {
	DrawRing(centerX, centerY, radius-thickness/2.0, radius+thickness/2.0, 0, 360, 36, color)
}

// DrawCircleSector draws a piece of pie (sector) of a circle.
func DrawCircleSector(centerX, centerY, radius, startAngle, endAngle float32, segments int, color Color) {
	drawSectorFan(centerX, centerY, radius, startAngle, endAngle, segments, color)
}

func drawSectorFan(centerX, centerY, radius, startAngle, endAngle float32, segments int, color Color) {
	if segments < 4 {
		segments = 16
	}
	startRad := Deg2Rad(startAngle)
	endRad := Deg2Rad(endAngle)
	step := (endRad - startRad) / float32(segments)

	vertices := make([]float32, 0, (segments+2)*6)
	vertices = append(vertices, centerX, centerY, color.R, color.G, color.B, color.A)

	for i := 0; i <= segments; i++ {
		angle := startRad + float32(i)*step
		px := centerX + radius*float32(math.Cos(float64(angle)))
		py := centerY + radius*float32(math.Sin(float64(angle)))
		vertices = append(vertices, px, py, color.R, color.G, color.B, color.A)
	}

	drawShapeVertices(vertices, gl.TRIANGLE_FAN)
}

// DrawEllipse draws a solid filled ellipse.
func DrawEllipse(centerX, centerY, radiusH, radiusV float32, color Color) {
	if radiusH <= 0 || radiusV <= 0 {
		return
	}
	const segments = 36
	vertices := make([]float32, 0, (segments+2)*6)

	vertices = append(vertices, centerX, centerY, color.R, color.G, color.B, color.A)
	for i := 0; i <= segments; i++ {
		angle := float32(i) * (2.0 * math.Pi / segments)
		px := centerX + radiusH*float32(math.Cos(float64(angle)))
		py := centerY + radiusV*float32(math.Sin(float64(angle)))
		vertices = append(vertices, px, py, color.R, color.G, color.B, color.A)
	}
	drawShapeVertices(vertices, gl.TRIANGLE_FAN)
}

// DrawEllipseOutline draws the outline of an ellipse.
func DrawEllipseOutline(centerX, centerY, radiusH, radiusV float32, color Color) {
	if radiusH <= 0 || radiusV <= 0 {
		return
	}
	const segments = 36
	vertices := make([]float32, 0, (segments+1)*6)

	for i := 0; i <= segments; i++ {
		angle := float32(i) * (2.0 * math.Pi / segments)
		px := centerX + radiusH*float32(math.Cos(float64(angle)))
		py := centerY + radiusV*float32(math.Sin(float64(angle)))
		vertices = append(vertices, px, py, color.R, color.G, color.B, color.A)
	}
	drawShapeVertices(vertices, gl.LINE_STRIP)
}

// DrawRing draws a filled ring (donut shape) or partial ring between startAngle and endAngle.
func DrawRing(centerX, centerY, innerRadius, outerRadius, startAngle, endAngle float32, segments int, color Color) {
	if segments < 4 {
		segments = 36
	}
	if innerRadius < 0 {
		innerRadius = 0
	}
	if outerRadius <= innerRadius {
		return
	}

	startRad := Deg2Rad(startAngle)
	endRad := Deg2Rad(endAngle)
	step := (endRad - startRad) / float32(segments)

	// Triangle strip (2 vertices per segment step)
	vertices := make([]float32, 0, (segments+1)*2*6)

	for i := 0; i <= segments; i++ {
		angle := startRad + float32(i)*step
		cos := float32(math.Cos(float64(angle)))
		sin := float32(math.Sin(float64(angle)))

		// Outer vertex
		ox := centerX + outerRadius*cos
		oy := centerY + outerRadius*sin
		vertices = append(vertices, ox, oy, color.R, color.G, color.B, color.A)

		// Inner vertex
		ix := centerX + innerRadius*cos
		iy := centerY + innerRadius*sin
		vertices = append(vertices, ix, iy, color.R, color.G, color.B, color.A)
	}

	drawShapeVertices(vertices, gl.TRIANGLE_STRIP)
}

// DrawPolygon draws a regular N-sided polygon.
func DrawPolygon(centerX, centerY float32, sides int, radius, rotation float32, color Color) {
	if sides < 3 || radius <= 0 {
		return
	}
	rotRad := Deg2Rad(rotation)
	step := (2.0 * math.Pi) / float32(sides)

	vertices := make([]float32, 0, (sides+2)*6)
	vertices = append(vertices, centerX, centerY, color.R, color.G, color.B, color.A)

	for i := 0; i <= sides; i++ {
		angle := rotRad + float32(i)*step
		px := centerX + radius*float32(math.Cos(float64(angle)))
		py := centerY + radius*float32(math.Sin(float64(angle)))
		vertices = append(vertices, px, py, color.R, color.G, color.B, color.A)
	}

	drawShapeVertices(vertices, gl.TRIANGLE_FAN)
}

// DrawPolygonOutline draws the outline of a regular polygon.
func DrawPolygonOutline(centerX, centerY float32, sides int, radius, rotation float32, color Color) {
	if sides < 3 || radius <= 0 {
		return
	}
	rotRad := Deg2Rad(rotation)
	step := (2.0 * math.Pi) / float32(sides)

	vertices := make([]float32, 0, (sides+1)*6)
	for i := 0; i <= sides; i++ {
		angle := rotRad + float32(i)*step
		px := centerX + radius*float32(math.Cos(float64(angle)))
		py := centerY + radius*float32(math.Sin(float64(angle)))
		vertices = append(vertices, px, py, color.R, color.G, color.B, color.A)
	}

	drawShapeVertices(vertices, gl.LINE_STRIP)
}
