package pixu

import (
	"math"
)

// Vec2 represents a 2-dimensional vector with X and Y float32 components.
type Vec2 struct {
	X float32
	Y float32
}

// V2 is a convenient shorthand constructor for Vec2.
func V2(x, y float32) Vec2 {
	return Vec2{X: x, Y: y}
}

// Add returns the vector sum of v and other.
func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{X: v.X + other.X, Y: v.Y + other.Y}
}

// Sub returns the vector difference of v - other.
func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{X: v.X - other.X, Y: v.Y - other.Y}
}

// Mul returns the component-wise product of v and other.
func (v Vec2) Mul(other Vec2) Vec2 {
	return Vec2{X: v.X * other.X, Y: v.Y * other.Y}
}

// Div returns the component-wise division of v by other.
func (v Vec2) Div(other Vec2) Vec2 {
	return Vec2{X: v.X / other.X, Y: v.Y / other.Y}
}

// Scale multiplies vector v by a scalar factor.
func (v Vec2) Scale(factor float32) Vec2 {
	return Vec2{X: v.X * factor, Y: v.Y * factor}
}

// Length returns the magnitude (length) of vector v.
func (v Vec2) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

// LengthSqr returns the squared magnitude of vector v (faster than Length).
func (v Vec2) LengthSqr() float32 {
	return v.X*v.X + v.Y*v.Y
}

// Distance returns the Euclidean distance between v and other.
func (v Vec2) Distance(other Vec2) float32 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

// DistanceSqr returns the squared Euclidean distance between v and other.
func (v Vec2) DistanceSqr(other Vec2) float32 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return dx*dx + dy*dy
}

// Normalize returns a unit vector in the same direction as v.
// If v is a zero vector, it returns Vec2{0, 0}.
func (v Vec2) Normalize() Vec2 {
	l := v.Length()
	if l == 0 {
		return Vec2{0, 0}
	}
	return Vec2{X: v.X / l, Y: v.Y / l}
}

// Dot returns the dot product of v and other.
func (v Vec2) Dot(other Vec2) float32 {
	return v.X*other.X + v.Y*other.Y
}

// Angle returns the angle of the vector in radians relative to positive X-axis.
func (v Vec2) Angle() float32 {
	return float32(math.Atan2(float64(v.Y), float64(v.X)))
}

// Rotate rotates the vector by the specified angle in radians.
func (v Vec2) Rotate(radians float32) Vec2 {
	cos := float32(math.Cos(float64(radians)))
	sin := float32(math.Sin(float64(radians)))
	return Vec2{
		X: v.X*cos - v.Y*sin,
		Y: v.X*sin + v.Y*cos,
	}
}

// Lerp linearly interpolates between vector v and target by factor t [0.0, 1.0].
func (v Vec2) Lerp(target Vec2, t float32) Vec2 {
	return Vec2{
		X: v.X + (target.X-v.X)*t,
		Y: v.Y + (target.Y-v.Y)*t,
	}
}

// Rect represents a 2D rectangle defined by top-left (X, Y) position and size (Width, Height).
type Rect struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

// NewRect creates a new Rect with the given position and size.
func NewRect(x, y, width, height float32) Rect {
	return Rect{X: x, Y: y, Width: width, Height: height}
}

// Position returns the top-left position of the rectangle as a Vec2.
func (r Rect) Position() Vec2 {
	return Vec2{X: r.X, Y: r.Y}
}

// Size returns the width and height of the rectangle as a Vec2.
func (r Rect) Size() Vec2 {
	return Vec2{X: r.Width, Y: r.Height}
}

// Center returns the center point of the rectangle as a Vec2.
func (r Rect) Center() Vec2 {
	return Vec2{X: r.X + r.Width/2, Y: r.Y + r.Height/2}
}

// Contains returns true if point is inside the rectangle.
func (r Rect) Contains(point Vec2) bool {
	return point.X >= r.X && point.X <= r.X+r.Width &&
		point.Y >= r.Y && point.Y <= r.Y+r.Height
}

// ContainsXY returns true if (x, y) is inside the rectangle.
func (r Rect) ContainsXY(x, y float32) bool {
	return x >= r.X && x <= r.X+r.Width &&
		y >= r.Y && y <= r.Y+r.Height
}

// CheckCollisionRecs checks if two rectangles overlap.
func CheckCollisionRecs(r1, r2 Rect) bool {
	return r1.X < r2.X+r2.Width &&
		r1.X+r1.Width > r2.X &&
		r1.Y < r2.Y+r2.Height &&
		r1.Y+r1.Height > r2.Y
}

// CheckCollisionCircles checks if two circles overlap.
func CheckCollisionCircles(c1 Vec2, r1 float32, c2 Vec2, r2 float32) bool {
	distSqr := c1.DistanceSqr(c2)
	radiusSum := r1 + r2
	return distSqr <= radiusSum*radiusSum
}

// CheckCollisionCircleRec checks if a circle and a rectangle overlap.
func CheckCollisionCircleRec(center Vec2, radius float32, rec Rect) bool {
	closestX := Clamp(center.X, rec.X, rec.X+rec.Width)
	closestY := Clamp(center.Y, rec.Y, rec.Y+rec.Height)

	dx := center.X - closestX
	dy := center.Y - closestY

	return (dx*dx + dy*dy) <= radius*radius
}

// CheckCollisionPointRec checks if a point is inside a rectangle.
func CheckCollisionPointRec(point Vec2, rec Rect) bool {
	return rec.Contains(point)
}

// CheckCollisionPointCircle checks if a point is inside a circle.
func CheckCollisionPointCircle(point Vec2, center Vec2, radius float32) bool {
	return point.DistanceSqr(center) <= radius*radius
}

// CheckCollisionPointTriangle checks if a point is inside a triangle defined by p1, p2, p3.
func CheckCollisionPointTriangle(p, p1, p2, p3 Vec2) bool {
	alpha := ((p2.Y-p3.Y)*(p.X-p3.X) + (p3.X-p2.X)*(p.Y-p3.Y)) /
		((p2.Y-p3.Y)*(p1.X-p3.X) + (p3.X-p2.X)*(p1.Y-p3.Y))
	beta := ((p3.Y-p1.Y)*(p.X-p3.X) + (p1.X-p3.X)*(p.Y-p3.Y)) /
		((p2.Y-p3.Y)*(p1.X-p3.X) + (p3.X-p2.X)*(p1.Y-p3.Y))
	gamma := 1.0 - alpha - beta

	return alpha > 0 && beta > 0 && gamma > 0
}

// GetCollisionRec returns the intersection rectangle of two colliding rectangles.
// Returns a zero-area rectangle if they do not collide.
func GetCollisionRec(r1, r2 Rect) Rect {
	if !CheckCollisionRecs(r1, r2) {
		return Rect{0, 0, 0, 0}
	}

	left := float32(math.Max(float64(r1.X), float64(r2.X)))
	top := float32(math.Max(float64(r1.Y), float64(r2.Y)))
	right := float32(math.Min(float64(r1.X+r1.Width), float64(r2.X+r2.Width)))
	bottom := float32(math.Min(float64(r1.Y+r1.Height), float64(r2.Y+r2.Height)))

	return Rect{
		X:      left,
		Y:      top,
		Width:  right - left,
		Height: bottom - top,
	}
}

// Clamp restricts a value to be within the range [min, max].
func Clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Lerp linearly interpolates between start and end by amount t [0.0, 1.0].
func Lerp(start, end, amount float32) float32 {
	return start + (end-start)*amount
}

// NormalizeVal normalizes a value from the range [start, end] to [0.0, 1.0].
func NormalizeVal(value, start, end float32) float32 {
	if start == end {
		return 0
	}
	return (value - start) / (end - start)
}

// Remap maps a value from one range [inputStart, inputEnd] to another range [outputStart, outputEnd].
func Remap(value, inputStart, inputEnd, outputStart, outputEnd float32) float32 {
	return Lerp(outputStart, outputEnd, NormalizeVal(value, inputStart, inputEnd))
}

// Deg2Rad converts degrees to radians.
func Deg2Rad(deg float32) float32 {
	return deg * (math.Pi / 180.0)
}

// Rad2Deg converts radians to degrees.
func Rad2Deg(rad float32) float32 {
	return rad * (180.0 / math.Pi)
}
