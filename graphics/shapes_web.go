//go:build wasm
// +build wasm

package graphics

import (
	"fmt"
)


func drawLine(x1, y1, x2, y2 float32, c Color) {
	ctx.Set("strokeStyle", colorToRGBA(c))
	ctx.Call("beginPath")
	ctx.Call("moveTo", float64(x1), float64(y1))
	ctx.Call("lineTo", float64(x2), float64(y2))
	ctx.Call("stroke")
}

func drawTriangle(x1, y1, x2, y2, x3, y3 float32, c Color) {
	ctx.Set("fillStyle", colorToRGBA(c))
	ctx.Call("beginPath")
	ctx.Call("moveTo", float64(x1), float64(y1))
	ctx.Call("lineTo", float64(x2), float64(y2))
	ctx.Call("lineTo", float64(x3), float64(y3))
	ctx.Call("closePath")
	ctx.Call("fill")
}

func drawCircle(cx, cy, radius float32, c Color) {
	ctx.Set("fillStyle", colorToRGBA(c))
	ctx.Call("beginPath")
	ctx.Call("arc", float64(cx), float64(cy), float64(radius), 0, 2*3.14159) // Full circle
	ctx.Call("fill")
}


func drawRectangle(x, y, w, h float32, c Color) {

	ctx.Set("fillStyle", colorToRGBA(c))
	ctx.Call("fillRect", float64(x), float64(y), float64(w), float64(h))
}

func drawRectangleOutline(x, y, w, h float32, c Color) {
	ctx.Set("strokeStyle", colorToRGBA(c))
	ctx.Call("strokeRect", float64(x), float64(y), float64(w), float64(h))
}



func colorToRGBA(c Color) string {
	r := float64(c.R) * 257.0
	g := float64(c.G) * 257.0
	b := float64(c.B) * 257.0
	a := float64(c.A) * 65535.0
	return fmt.Sprintf("rgba(%d,%d,%d,%.2f)", int(r), int(g), int(b), a)
}