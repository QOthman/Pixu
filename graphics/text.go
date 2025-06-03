package graphics

var fontAtlas *Image
var charWidth, charHeight = 20, 24
var charsPerRow = 16

func loadFontAtlas() error {
	var err error
	fontAtlas, err = LoadImage("font/font_atlas_bold.png")
	return err
}

// DrawTextFromAtlas renders text using the font atlas (fixed-width).
func DrawText(text string, x, y, size float32, color Color) {
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
		srcY := int(fontAtlas.Height) - ((index / charsPerRow) + 1) * charHeight

		drawX := x + float32(i)*float32(charWidth)*scale * 0.5 
		drawY := y // Adjust this if needed
		
		// Draw the character from the font atlas
		DrawImageEx(fontAtlas, DrawOptions{
			X:      drawX,
			Y:      drawY,
			Width:  float32(charWidth) * scale,
			Height: float32(charHeight) * scale,
			Rotation: 0,
			Tint:   color,
			SrcX:   float32(srcX),
			SrcY:   float32(srcY),
			SrcW:   float32(charWidth),
			SrcH:   float32(charHeight),
		})
	}
}


// DrawTextCentered renders centered text using the font atlas.
func DrawTextCentered(text string, centerX, centerY, size float32, color Color) {
	if fontAtlas == nil {
		return
	}

	totalWidth := float32(len(text)) * float32(charWidth) * size * 0.5
	startX := centerX - totalWidth / 2

	for i, ch := range text {
		if ch < 32 || ch > 127 {
			continue // unsupported character
		}
		index := int(ch) - 32
		srcX := (index % charsPerRow) * charWidth // Adjust for atlas scale
		srcY := int(fontAtlas.Height) - ((index / charsPerRow) + 1) * charHeight

		drawX := startX + float32(i)*float32(charWidth)*size * 0.5 
		drawY := centerY - float32(charHeight)*size / 2 // Center vertically
		// Draw the character from the font atlas
		DrawImageEx(fontAtlas, DrawOptions{
			X:      drawX,
			Y:      drawY,
			Width:  float32(charWidth) * size,
			Height: float32(charHeight) * size,
			Rotation: 0,
			Tint:   color,
			SrcX:   float32(srcX),
			SrcY:   float32(srcY),
			SrcW:   float32(charWidth),
			SrcH:   float32(charHeight),
		})
	}
}

// DrawTextWithBackground renders text with a background rectangle.
func DrawTextWithBackground(text string, x, y, size float32, bgColor, textColor Color) {
	if fontAtlas == nil {
		return
	}

	// Calculate total width and height of the text
	totalWidth := float32(len(text)) * float32(charWidth) * size * 0.5
	totalHeight := float32(charHeight) * size

	// Draw background rectangle
	DrawRectangle(x, y, totalWidth, totalHeight, bgColor)

	// Draw the text on top
	DrawText(text, x, y + (totalHeight - float32(charHeight)*size) / 2, size, textColor)
}

// DrawTextOutline renders text with an outline effect.
func DrawTextOutline(text string, x, y, size float32, textColor, outlineColor Color) {
	if fontAtlas == nil {
		return
	}

	// Draw outline by rendering text multiple times with slight offsets
	offsets := []struct{ dx, dy float32 }{
		{-1, -1}, {1, -1}, {-1, 1}, {1, 1},
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
	}

	for _, offset := range offsets {
		DrawText(text, x + offset.dx * size * 0.5, y + offset.dy * size / 2, size, outlineColor)
	}

	// Draw the main text
	DrawText(text, x, y + (float32(charHeight)*size - float32(charHeight)*size) / 2, size, textColor)
}
