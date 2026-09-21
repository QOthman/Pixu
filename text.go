package pixu

import (
	"fmt"
	"strings"
)

// Font represents a bitmap font atlas for text rendering.
type Font struct {
	Texture     *Image
	CharWidth   int
	CharHeight  int
	CharsPerRow int
	BaseSize    float32
}

var (
	defaultFont     *Font
	defaultBoldFont *Font
	activeFont      *Font
)

// initFontSystem initializes the embedded font system.
func initFontSystem() error {
	var err error

	// Load regular font atlas from embedded bytes
	if len(defaultFontAtlasData) > 0 {
		regImg, err := LoadImageFromBytes(defaultFontAtlasData)
		if err == nil {
			defaultFont = &Font{
				Texture:     regImg,
				CharWidth:   20,
				CharHeight:  24,
				CharsPerRow: 16,
				BaseSize:    24,
			}
		}
	}

	// Load bold font atlas from embedded bytes
	if len(defaultFontAtlasBoldData) > 0 {
		boldImg, err := LoadImageFromBytes(defaultFontAtlasBoldData)
		if err == nil {
			defaultBoldFont = &Font{
				Texture:     boldImg,
				CharWidth:   20,
				CharHeight:  24,
				CharsPerRow: 16,
				BaseSize:    24,
			}
		}
	}

	if defaultBoldFont != nil {
		activeFont = defaultBoldFont
	} else if defaultFont != nil {
		activeFont = defaultFont
	} else {
		return fmt.Errorf("embedded font atlases could not be initialized")
	}

	return err
}

// destroyFontSystem frees GPU memory used by font textures.
func destroyFontSystem() {
	if defaultFont != nil && defaultFont.Texture != nil {
		defaultFont.Texture.Delete()
		defaultFont = nil
	}
	if defaultBoldFont != nil && defaultBoldFont.Texture != nil {
		defaultBoldFont.Texture.Delete()
		defaultBoldFont = nil
	}
	activeFont = nil
}

// GetDefaultFont returns the embedded regular font.
func GetDefaultFont() *Font {
	return defaultFont
}

// GetDefaultBoldFont returns the embedded bold font.
func GetDefaultBoldFont() *Font {
	return defaultBoldFont
}

// GetFont returns the currently active font.
func GetFont() *Font {
	return activeFont
}

// SetFont sets the active font used for text rendering.
func SetFont(font *Font) {
	if font != nil && font.Texture != nil && font.Texture.IsValid() {
		activeFont = font
	}
}

// LoadFont loads a custom bitmap font atlas from a file.
func LoadFont(filePath string, charWidth, charHeight, charsPerRow int) (*Font, error) {
	img, err := LoadImage(filePath)
	if err != nil {
		return nil, err
	}
	return &Font{
		Texture:     img,
		CharWidth:   charWidth,
		CharHeight:  charHeight,
		CharsPerRow: charsPerRow,
		BaseSize:    float32(charHeight),
	}, nil
}

// LoadFontFromBytes loads a custom bitmap font atlas from byte slice.
func LoadFontFromBytes(data []byte, charWidth, charHeight, charsPerRow int) (*Font, error) {
	img, err := LoadImageFromBytes(data)
	if err != nil {
		return nil, err
	}
	return &Font{
		Texture:     img,
		CharWidth:   charWidth,
		CharHeight:  charHeight,
		CharsPerRow: charsPerRow,
		BaseSize:    float32(charHeight),
	}, nil
}

// MeasureText calculates the width and height in pixels of a text string when rendered at given size scale.
func MeasureText(text string, size float32) (width, height float32) {
	if activeFont == nil {
		return 0, 0
	}
	lines := strings.Split(text, "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}
	w := float32(maxLen) * float32(activeFont.CharWidth) * size * 0.5
	h := float32(len(lines)) * float32(activeFont.CharHeight) * size
	return w, h
}

// MeasureTextV returns the size of a text string as a Vec2.
func MeasureTextV(text string, size float32) Vec2 {
	w, h := MeasureText(text, size)
	return Vec2{X: w, Y: h}
}

// DrawText renders text at (x, y) with specified size scaling and color.
func DrawText(text string, x, y, size float32, color Color) {
	if activeFont == nil || activeFont.Texture == nil || !activeFont.Texture.IsValid() {
		return
	}

	charW := float32(activeFont.CharWidth)
	charH := float32(activeFont.CharHeight)
	atlasH := float32(activeFont.Texture.Height)

	lines := strings.Split(text, "\n")
	for lineIdx, line := range lines {
		lineY := y + float32(lineIdx)*charH*size
		for i, ch := range line {
			if ch < 32 || ch > 127 {
				continue
			}
			index := int(ch) - 32
			srcX := float32((index % activeFont.CharsPerRow) * activeFont.CharWidth)
			srcY := atlasH - float32(((index/activeFont.CharsPerRow)+1)*activeFont.CharHeight)

			drawX := x + float32(i)*charW*size*0.5

			DrawImageEx(activeFont.Texture, DrawOptions{
				X:        drawX,
				Y:        lineY,
				Width:    charW * size,
				Height:   charH * size,
				Rotation: 0,
				Tint:     color,
				SrcX:     srcX,
				SrcY:     srcY,
				SrcW:     charW,
				SrcH:     charH,
			})
		}
	}
}

// DrawTextV renders text at Vec2 position.
func DrawTextV(text string, pos Vec2, size float32, color Color) {
	DrawText(text, pos.X, pos.Y, size, color)
}

// DrawTextCentered renders text horizontally and vertically centered at (centerX, centerY).
func DrawTextCentered(text string, centerX, centerY, size float32, color Color) {
	w, h := MeasureText(text, size)
	DrawText(text, centerX-w/2.0, centerY-h/2.0, size, color)
}

// DrawTextWithBackground renders text on top of a solid colored background box.
func DrawTextWithBackground(text string, x, y, size float32, bgColor, textColor Color) {
	w, h := MeasureText(text, size)
	padding := float32(6.0) * size
	DrawRectangle(x-padding/2, y-padding/2, w+padding, h+padding, bgColor)
	DrawText(text, x, y, size, textColor)
}

// DrawTextOutline renders text with a contrasting border outline.
func DrawTextOutline(text string, x, y, size float32, textColor, outlineColor Color) {
	offsets := []struct{ dx, dy float32 }{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}

	for _, off := range offsets {
		DrawText(text, x+off.dx*size*0.5, y+off.dy*size*0.5, size, outlineColor)
	}

	DrawText(text, x, y, size, textColor)
}

// DrawTextPro renders text with full control over position, origin, rotation, spacing, and tint.
func DrawTextPro(font *Font, text string, pos Vec2, origin Vec2, rotation, size, spacing float32, tint Color) {
	if font == nil {
		font = activeFont
	}
	if font == nil || font.Texture == nil || !font.Texture.IsValid() {
		return
	}

	charW := float32(font.CharWidth)
	charH := float32(font.CharHeight)
	atlasH := float32(font.Texture.Height)

	lines := strings.Split(text, "\n")
	for lineIdx, line := range lines {
		lineY := float32(lineIdx) * charH * size
		for i, ch := range line {
			if ch < 32 || ch > 127 {
				continue
			}
			index := int(ch) - 32
			srcX := float32((index % font.CharsPerRow) * font.CharWidth)
			srcY := atlasH - float32(((index/font.CharsPerRow)+1)*font.CharHeight)

			charPosX := float32(i) * (charW*size*0.5 + spacing)

			DrawImageEx(font.Texture, DrawOptions{
				X:        pos.X + charPosX,
				Y:        pos.Y + lineY,
				Width:    charW * size,
				Height:   charH * size,
				Rotation: rotation,
				Origin:   origin,
				Tint:     tint,
				SrcX:     srcX,
				SrcY:     srcY,
				SrcW:     charW,
				SrcH:     charH,
			})
		}
	}
}
