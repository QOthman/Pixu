//go:build !wasm
// +build !wasm

package graphics

var fontAtlas *Image
var charWidth, charHeight = 20, 24
var charsPerRow = 16

func loadFontAtlas() error {
	var err error
	fontAtlas, err = loadImage("font/font_atlas_bold.png")
	return err
}

// DrawTextFromAtlas renders text using the font atlas (fixed-width).
func drawText(text string, x, y, size float32, color Color) {
	if fontAtlas == nil {
		return
	}
	scale := size
	for i, ch := range text {
		if ch < 32 || ch > 127 {
			continue // unsupported character
		}
		index := int(ch) - 32
		srcX := (index % charsPerRow) * charWidth // Adjust for atlas scale
		srcY := (index / charsPerRow) * charHeight

		drawX := x + float32(i)*float32(charWidth)*scale*0.5
		drawY := y // Adjust this if needed

		// Draw the character from the font atlas
		drawImageEx(fontAtlas, DrawOptions{
			X:        drawX,
			Y:        drawY,
			Width:    float32(charWidth) * scale,
			Height:   float32(charHeight) * scale,
			Rotation: 0,
			Tint:     color,
			SrcX:     float32(srcX),
			SrcY:     float32(srcY),
			SrcW:     float32(charWidth),
			SrcH:     float32(charHeight),
		})
	}
}

// DrawTextCentered renders centered text using the font atlas.
func drawTextCentered(text string, centerX, centerY, size float32, color Color) {
	if fontAtlas == nil {
		return
	}

	totalWidth := float32(len(text)) * float32(charWidth) * size * 0.5
	startX := centerX - totalWidth/2

	for i, ch := range text {
		if ch < 32 || ch > 127 {
			continue // unsupported character
		}
		index := int(ch) - 32
		srcX := (index % charsPerRow) * charWidth // Adjust for atlas scale
		srcY := (index / charsPerRow) * charHeight

		drawX := startX + float32(i)*float32(charWidth)*size*0.5
		drawY := centerY - float32(charHeight)*size/2 // Center vertically
		// Draw the character from the font atlas
		drawImageEx(fontAtlas, DrawOptions{
			X:        drawX,
			Y:        drawY,
			Width:    float32(charWidth) * size,
			Height:   float32(charHeight) * size,
			Rotation: 0,
			Tint:     color,
			SrcX:     float32(srcX),
			SrcY:     float32(srcY),
			SrcW:     float32(charWidth),
			SrcH:     float32(charHeight),
		})
	}
}

// DrawTextWithBackground renders text with a background rectangle.
func drawTextWithBackground(text string, x, y, size float32, bgColor, textColor Color, paddingX, paddingY float32) {
	if fontAtlas == nil {
		return
	}
	textWidth := float32(len(text)) * float32(charWidth) * size * 0.5
	textHeight := float32(charHeight) * size

	bgWidth := textWidth + 2*paddingX
	bgHeight := textHeight + 2*paddingY

	drawRectangle(x, y, bgWidth, bgHeight, bgColor)

	textX := x + paddingX
	textY := y + paddingY + (bgHeight-2*paddingY-textHeight)/2

	drawText(text, textX, textY, size, textColor)
}

// DrawTextOutline renders text with an outline effect.
func drawTextOutline(text string, x, y, size float32, textColor, outlineColor Color) {
	if fontAtlas == nil {
		return
	}

	// Draw outline by rendering text multiple times with slight offsets
	offsets := []struct{ dx, dy float32 }{
		{-1, -1}, {1, -1}, {-1, 1}, {1, 1},
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
	}

	for _, offset := range offsets {
		drawText(text, x+offset.dx*size*0.5, y+offset.dy*size/2, size, outlineColor)
	}

	// Draw the main text
	drawText(text, x, y+(float32(charHeight)*size-float32(charHeight)*size)/2, size, textColor)
}
