package pixu

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"os"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// TextureFilter defines texture min/mag filtering modes.
type TextureFilter int

const (
	FilterPoint   TextureFilter = iota // Nearest-neighbor (pixel art / retro look)
	FilterNearest                      // Alias for FilterPoint
	FilterLinear                       // Linear interpolation (smooth)
	FilterBilinear
)

// TextureWrap defines texture wrapping modes.
type TextureWrap int

const (
	WrapClamp TextureWrap = iota
	WrapRepeat
	WrapMirroredRepeat
)

// Image represents a 2D OpenGL texture loaded into GPU memory.
type Image struct {
	TextureID uint32
	Width     int32
	Height    int32
	filePath  string
}

// Texture2D is a type alias for Image to match standard game engine conventions.
type Texture2D = Image

// DrawOptions specifies the complete set of parameters for drawing a textured quad.
type DrawOptions struct {
	X, Y          float32 // Destination position (top-left by default, or transformed by Origin)
	Width, Height float32 // Destination width and height
	Rotation      float32 // Rotation angle in degrees (clockwise)
	Origin        Vec2    // Origin offset for rotation and positioning [0..Width, 0..Height]
	Tint          Color   // Multiplicative tint color (default White)
	SrcX, SrcY    float32 // Source sub-rectangle X and Y on texture
	SrcW, SrcH    float32 // Source sub-rectangle Width and Height on texture
}

// ID returns the OpenGL texture ID.
func (img *Image) ID() uint32 {
	if img == nil {
		return 0
	}
	return img.TextureID
}

// Size returns the dimensions as a Vec2.
func (img *Image) Size() Vec2 {
	if img == nil {
		return Vec2{0, 0}
	}
	return Vec2{X: float32(img.Width), Y: float32(img.Height)}
}

// Rect returns a Rect representing the full dimensions of the image.
func (img *Image) Rect() Rect {
	if img == nil {
		return Rect{}
	}
	return Rect{X: 0, Y: 0, Width: float32(img.Width), Height: float32(img.Height)}
}

// IsValid returns true if the image is loaded and has a valid OpenGL texture.
func (img *Image) IsValid() bool {
	return img != nil && img.TextureID != 0 && img.Width > 0 && img.Height > 0
}

// SetFilter sets the OpenGL minifying and magnifying filters for this texture.
func (img *Image) SetFilter(filter TextureFilter) {
	if !img.IsValid() {
		return
	}
	gl.BindTexture(gl.TEXTURE_2D, img.TextureID)
	switch filter {
	case FilterPoint, FilterNearest:
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	case FilterLinear, FilterBilinear:
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	}
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// SetWrap sets the texture coordinate wrapping mode.
func (img *Image) SetWrap(wrap TextureWrap) {
	if !img.IsValid() {
		return
	}
	gl.BindTexture(gl.TEXTURE_2D, img.TextureID)
	var glWrap int32 = gl.CLAMP_TO_EDGE
	switch wrap {
	case WrapClamp:
		glWrap = gl.CLAMP_TO_EDGE
	case WrapRepeat:
		glWrap = gl.REPEAT
	case WrapMirroredRepeat:
		glWrap = gl.MIRRORED_REPEAT
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, glWrap)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, glWrap)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// Delete frees the OpenGL texture from GPU memory.
func (img *Image) Delete() {
	if img != nil && img.TextureID != 0 {
		gl.DeleteTextures(1, &img.TextureID)
		img.TextureID = 0
	}
}

// Destroy is an alias for Delete.
func (img *Image) Destroy() {
	img.Delete()
}

// LoadImage loads an image file from disk (PNG, JPEG, GIF).
func LoadImage(filePath string) (*Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open image file %s: %w", filePath, err)
	}
	defer file.Close()

	img, err := LoadImageFromReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image %s: %w", filePath, err)
	}
	img.filePath = filePath
	return img, nil
}

// LoadImageFromBytes loads an image from a byte slice in memory.
func LoadImageFromBytes(data []byte) (*Image, error) {
	return LoadImageFromReader(bytes.NewReader(data))
}

// LoadImageFromReader decodes an image stream and uploads it to OpenGL.
func LoadImageFromReader(r io.Reader) (*Image, error) {
	decoded, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return LoadImageFromImage(decoded)
}

// LoadImageFromImage creates an OpenGL texture from an existing Go image.Image.
func LoadImageFromImage(img image.Image) (*Image, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid image dimensions: %dx%d", width, height)
	}

	// Convert and vertically flip for OpenGL coordinate compatibility
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcColor := img.At(bounds.Min.X+x, bounds.Min.Y+y)
			rgba.Set(x, height-y-1, srcColor)
		}
	}

	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		int32(width), int32(height), 0,
		gl.RGBA, gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return &Image{
		TextureID: textureID,
		Width:     int32(width),
		Height:    int32(height),
	}, nil
}

// NewImage creates a solid color image/texture of specified dimensions.
func NewImage(width, height int, fillColor Color) (*Image, error) {
	r, g, b, a := fillColor.ToRGBA()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	c := color.RGBA{R: r, G: g, B: b, A: a}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetRGBA(x, y, c)
		}
	}
	return LoadImageFromImage(img)
}

// DrawImage renders an image at position (x, y) with original dimensions.
func DrawImage(img *Image, x, y float32) {
	if !img.IsValid() {
		return
	}
	DrawImageEx(img, DrawOptions{
		X:        x,
		Y:        y,
		Width:    float32(img.Width),
		Height:   float32(img.Height),
		Rotation: 0,
		Tint:     White,
		SrcX:     0,
		SrcY:     0,
		SrcW:     float32(img.Width),
		SrcH:     float32(img.Height),
	})
}

// DrawImageV renders an image at a Vec2 position with a tint color.
func DrawImageV(img *Image, pos Vec2, tint Color) {
	if !img.IsValid() {
		return
	}
	DrawImageEx(img, DrawOptions{
		X:        pos.X,
		Y:        pos.Y,
		Width:    float32(img.Width),
		Height:   float32(img.Height),
		Rotation: 0,
		Tint:     tint,
		SrcX:     0,
		SrcY:     0,
		SrcW:     float32(img.Width),
		SrcH:     float32(img.Height),
	})
}

// DrawImageScaled renders an image scaled by scaleX and scaleY factors.
func DrawImageScaled(img *Image, x, y, scaleX, scaleY float32) {
	if !img.IsValid() {
		return
	}
	DrawImageEx(img, DrawOptions{
		X:        x,
		Y:        y,
		Width:    float32(img.Width) * scaleX,
		Height:   float32(img.Height) * scaleY,
		Rotation: 0,
		Tint:     White,
		SrcX:     0,
		SrcY:     0,
		SrcW:     float32(img.Width),
		SrcH:     float32(img.Height),
	})
}

// DrawImageRotated renders an image rotated by specified degrees around its center.
func DrawImageRotated(img *Image, x, y, rotation float32) {
	if !img.IsValid() {
		return
	}
	w := float32(img.Width)
	h := float32(img.Height)
	DrawImageEx(img, DrawOptions{
		X:        x,
		Y:        y,
		Width:    w,
		Height:   h,
		Rotation: rotation,
		Origin:   Vec2{X: w / 2, Y: h / 2},
		Tint:     White,
		SrcX:     0,
		SrcY:     0,
		SrcW:     w,
		SrcH:     h,
	})
}

// DrawImageTinted renders an image at (x, y) with a color tint.
func DrawImageTinted(img *Image, x, y float32, tint Color) {
	if !img.IsValid() {
		return
	}
	DrawImageEx(img, DrawOptions{
		X:        x,
		Y:        y,
		Width:    float32(img.Width),
		Height:   float32(img.Height),
		Rotation: 0,
		Tint:     tint,
		SrcX:     0,
		SrcY:     0,
		SrcW:     float32(img.Width),
		SrcH:     float32(img.Height),
	})
}

// DrawImageRec renders a rectangular slice of an image (sprite sheet) to a destination position.
func DrawImageRec(img *Image, srcRec Rect, destPos Vec2, tint Color) {
	if !img.IsValid() {
		return
	}
	DrawImageEx(img, DrawOptions{
		X:        destPos.X,
		Y:        destPos.Y,
		Width:    srcRec.Width,
		Height:   srcRec.Height,
		Rotation: 0,
		Tint:     tint,
		SrcX:     srcRec.X,
		SrcY:     srcRec.Y,
		SrcW:     srcRec.Width,
		SrcH:     srcRec.Height,
	})
}

// DrawImagePro renders an image with complete control over source rectangle, destination rectangle, origin, and rotation.
func DrawImagePro(img *Image, srcRec Rect, destRec Rect, origin Vec2, rotation float32, tint Color) {
	if !img.IsValid() {
		return
	}
	DrawImageEx(img, DrawOptions{
		X:        destRec.X,
		Y:        destRec.Y,
		Width:    destRec.Width,
		Height:   destRec.Height,
		Rotation: rotation,
		Origin:   origin,
		Tint:     tint,
		SrcX:     srcRec.X,
		SrcY:     srcRec.Y,
		SrcW:     srcRec.Width,
		SrcH:     srcRec.Height,
	})
}

// DrawImageEx renders an image with full custom DrawOptions.
func DrawImageEx(img *Image, opts DrawOptions) {
	if !img.IsValid() {
		return
	}

	if opts.Width == 0 {
		opts.Width = float32(img.Width)
	}
	if opts.Height == 0 {
		opts.Height = float32(img.Height)
	}
	if opts.SrcW == 0 {
		opts.SrcW = float32(img.Width)
	}
	if opts.SrcH == 0 {
		opts.SrcH = float32(img.Height)
	}
	if opts.Tint.Equals(Blank) {
		opts.Tint = White
	}

	// Texture UV coordinates [0.0, 1.0]
	texX := opts.SrcX / float32(img.Width)
	texY := opts.SrcY / float32(img.Height)
	texW := opts.SrcW / float32(img.Width)
	texH := opts.SrcH / float32(img.Height)

	// Destination quad corners relative to origin
	x0 := -opts.Origin.X
	y0 := -opts.Origin.Y
	x1 := x0 + opts.Width
	y1 := y0 + opts.Height

	var tlx, tly, trx, try, brx, bry, blx, bly float32

	if opts.Rotation != 0 {
		rad := float64(opts.Rotation) * math.Pi / 180.0
		cos := float32(math.Cos(rad))
		sin := float32(math.Sin(rad))

		rotX := func(x, y float32) float32 { return x*cos - y*sin + opts.X }
		rotY := func(x, y float32) float32 { return x*sin + y*cos + opts.Y }

		tlx, tly = rotX(x0, y0), rotY(x0, y0)
		trx, try = rotX(x1, y0), rotY(x1, y0)
		brx, bry = rotX(x1, y1), rotY(x1, y1)
		blx, bly = rotX(x0, y1), rotY(x0, y1)
	} else {
		tlx, tly = opts.X+x0, opts.Y+y0
		trx, try = opts.X+x1, opts.Y+y0
		brx, bry = opts.X+x1, opts.Y+y1
		blx, bly = opts.X+x0, opts.Y+y1
	}

	// 4 vertices (Pos: 2, TexCoord: 2, Color: 4)
	r, g, b, a := opts.Tint.R, opts.Tint.G, opts.Tint.B, opts.Tint.A
	vertices := []float32{
		tlx, tly, texX, texY + texH, r, g, b, a,       // Top-Left
		trx, try, texX + texW, texY + texH, r, g, b, a, // Top-Right
		brx, bry, texX + texW, texY, r, g, b, a,        // Bottom-Right
		blx, bly, texX, texY, r, g, b, a,               // Bottom-Left
	}

	drawTexturedQuad(img, vertices)
}

// drawTexturedQuad uploads vertex data and executes the draw call using pre-allocated buffers.
func drawTexturedQuad(img *Image, vertices []float32) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, img.TextureID)

	gl.UseProgram(textureShaderProgram)
	gl.BindVertexArray(textureVAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, textureVBO)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices[0]))*len(vertices), gl.Ptr(vertices), gl.DYNAMIC_DRAW)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}
