//go:build wasm
// +build wasm


package graphics

import (
	"fmt"
)

func (c Color) ToRGBA() string {
	// Convert the Color struct to an "rgba(r, g, b, a)" format for Canvas context
	return fmt.Sprintf("rgba(%d, %d, %d, %.2f)", 
		int(c.R*255), 
		int(c.G*255), 
		int(c.B*255), 
		c.A)
}



func drawText(text string, x, y float32, fontSize float32, color Color) {
    // Ensure the font is set with the correct scaling factor
    ctx.Set("font", fmt.Sprintf("%.0fpx sans-serif", fontSize*16))
    ctx.Set("fillStyle", color.ToRGBA())
	ctx.Set("textBaseline", "top")
    ctx.Call("fillText", text, x, y)
}


func drawTextCentered(text string, x, y float32, fontSize float32, color Color) {
	ctx.Set("font", fmt.Sprintf("%.0fpx sans-serif", fontSize*16))
	ctx.Set("fillStyle", color.ToRGBA())
	ctx.Set("textAlign", "center")
	ctx.Set("textBaseline", "middle")
	ctx.Call("fillText", text, x, y)

	// Reset
	ctx.Set("textAlign", "start")
	ctx.Set("textBaseline", "alphabetic")
}

func drawTextWithBackground(text string, x, y float32, fontSize float32, textColor, bgColor Color, paddingX, paddingY float32) {
	font := fmt.Sprintf("%.0fpx sans-serif", fontSize*16)
	ctx.Set("font", font)
	ctx.Set("textBaseline", "top")

	metrics := ctx.Call("measureText", text)
	textWidth := metrics.Get("width").Float()
	textHeight := fontSize * 16 // Approximation

	// Draw background
	ctx.Set("fillStyle", bgColor.ToRGBA())
	ctx.Call("fillRect", x-paddingX, y-paddingY, float32(textWidth)+2*paddingX, textHeight+2*paddingY)

	// Draw text
	ctx.Set("fillStyle", textColor.ToRGBA())
	ctx.Call("fillText", text, x, y)

	// Reset baseline if needed
	ctx.Set("textBaseline", "alphabetic")
}


func drawTextOutline(text string, x, y float32, fontSize float32, textColor, outlineColor Color) {
	font := fmt.Sprintf("%.0fpx sans-serif", fontSize*16)
	ctx.Set("font", font)
	ctx.Set("lineWidth", 3)
	ctx.Set("strokeStyle", outlineColor.ToRGBA())
	ctx.Set("fillStyle", textColor.ToRGBA())
	ctx.Set("textBaseline", "top")

	ctx.Call("strokeText", text, x, y)
	ctx.Call("fillText", text, x, y)

	ctx.Set("textBaseline", "alphabetic")
}
