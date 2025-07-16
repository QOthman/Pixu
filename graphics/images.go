//go:build !wasm
// +build !wasm

package graphics

import (
	"image"
	_ "image/gif" // Support GIF
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Image struct represents a loaded image texture
// It contains the OpenGL texture ID, dimensions, and file path.
type Image struct {
	TextureID uint32
	Width     int32
	Height    int32
	filePath  string
}

type DrawOptions struct {
	X, Y          float32 // Destination position
	Width, Height float32 // Destination size
	Rotation      float32 // Not used yet
	Tint          Color   // Tint color (RGBA)
	SrcX, SrcY    float32 // Source rect X, Y
	SrcW, SrcH    float32 // Source rect Width, Height
}

// Global texture shader program
var textureShaderProgram uint32
var textureVAO, textureVBO uint32

// Initialize texture system
func initTextureSystem() {
	setupTextureShaders()
	setupTextureBuffers()
}

// Load image mn file - supports PNG, JPEG, GIF
func loadImage(filePath string) (*Image, error) {
	// Open l file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode l image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// Convert l RGBA
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, bounds.Max.Y-y-1, img.At(x, y))
		}
	}

	// Create OpenGL texture
	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	// Texture parameters 
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// Upload texture data
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		int32(width), int32(height), 0,
		gl.RGBA, gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)

	return &Image{
		TextureID: textureID,
		Width:     int32(width),
		Height:    int32(height),
		filePath:  filePath,
	}, nil
}

// Draw image at specific position
func drawImage(img *Image, x, y float32) {
	drawImageEx(img, DrawOptions{
		X:        x,
		Y:        y,
		Width:    float32(img.Width),
		Height:   float32(img.Height),
		Rotation: 0,
		Tint:     WHITE,
		SrcX:     0,
		SrcY:     0,
		SrcW:     float32(img.Width),
		SrcH:     float32(img.Height),
	})
}

// Draw image scaled
func drawImageScaled(img *Image, x, y, scaleX, scaleY float32) {
	drawImageEx(img, DrawOptions{
		X:        x,
		Y:        y,
		Width:    float32(img.Width) * scaleX,
		Height:   float32(img.Height) * scaleY,
		Rotation: 0,
		Tint:     WHITE,
		SrcX:     0,
		SrcY:     0,
		SrcW:     float32(img.Width),
		SrcH:     float32(img.Height),
	})
}

// Draw image rotated
func drawImageRotated(img *Image, x, y, rotation float32) {
	drawImageEx(img, DrawOptions{
		X:        x,
		Y:        y,
		Width:    float32(img.Width),
		Height:   float32(img.Height),
		Rotation: rotation,
		Tint:     WHITE,
		SrcX:     0,
		SrcY:     0,
		SrcW:     float32(img.Width),
		SrcH:     float32(img.Height),
	})
}

// Draw image with tint color
func drawImageTinted(img *Image, x, y float32, tint Color) {
	drawImageEx(img, DrawOptions{
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

// Draw image extended - kol l options
func drawImageEx(img *Image, opts DrawOptions) {
	if img == nil || img.TextureID == 0 {
		return // No image to draw
	}

	// Set default width/height if not provided
	if opts.Width == 0 {
		opts.Width = float32(img.Width)
	}
	if opts.Height == 0 {
		opts.Height = float32(img.Height)
	}
	if opts.SrcW <= 0 {
		opts.SrcW = float32(img.Width)
	}
	if opts.SrcH <= 0 {
		opts.SrcH = float32(img.Height)
	}
	if opts.Tint == (Color{}) {
		opts.Tint = WHITE
	}

	// Texture coordinates
	texX := opts.SrcX / float32(img.Width)
	texY := 1.0 - (opts.SrcY + opts.SrcH) / float32(img.Height)
	texW := opts.SrcW / float32(img.Width)
	texH := opts.SrcH / float32(img.Height)


	// Rotation angle in radians
	rotation := opts.Rotation * math.Pi / 180.0

	// Center of the image in screen (pixel) coordinates
	centerPx := opts.X + opts.Width/2
	centerPy := opts.Y + opts.Height/2

	// Helper function: rotate a point (x, y) around the center (in pixels)
	rotatePixel := func(x, y float32) (float32, float32) {
		translatedX := x - centerPx
		translatedY := y - centerPy

		sin, cos := float32(math.Sin(float64(rotation))), float32(math.Cos(float64(rotation)))

		rotatedX := translatedX*cos - translatedY*sin
		rotatedY := translatedX*sin + translatedY*cos

		return rotatedX + centerPx, rotatedY + centerPy
	}

	// Rotate all four corners in pixel space
	px_tlx, py_tly := rotatePixel(opts.X, opts.Y)                     // top-left
	px_trx, py_try := rotatePixel(opts.X+opts.Width, opts.Y)           // top-right
	px_brx, py_bry := rotatePixel(opts.X+opts.Width, opts.Y+opts.Height) // bottom-right
	px_blx, py_bly := rotatePixel(opts.X, opts.Y+opts.Height)          // bottom-left

	// Convert rotated pixel positions to OpenGL normalized coordinates
	gl_tlx, gl_tly := screenToGL(px_tlx, py_tly)
	gl_trx, gl_try := screenToGL(px_trx, py_try)
	gl_brx, gl_bry := screenToGL(px_brx, py_bry)
	gl_blx, gl_bly := screenToGL(px_blx, py_bly)

	// Vertex data with position (GL coords), texture coords, and tint color
	vertices := []float32{
		gl_tlx, gl_tly, texX, texY + texH, opts.Tint.R, opts.Tint.G, opts.Tint.B, opts.Tint.A, // top-left
		gl_trx, gl_try, texX + texW, texY + texH, opts.Tint.R, opts.Tint.G, opts.Tint.B, opts.Tint.A, // top-right
		gl_brx, gl_bry, texX + texW, texY, opts.Tint.R, opts.Tint.G, opts.Tint.B, opts.Tint.A, // bottom-right
		gl_blx, gl_bly, texX, texY, opts.Tint.R, opts.Tint.G, opts.Tint.B, opts.Tint.A, // bottom-left
	}

	indices := []uint32{0, 1, 2, 2, 3, 0}

	drawTexturedQuad(img, vertices, indices)
}

// Helper function to draw a textured quad
// This function sets up the vertex array object (VAO), vertex buffer object (VBO),
// and element buffer object (EBO) for rendering the textured quad.
func drawTexturedQuad(img *Image, vertices []float32, indices []uint32) {
	var EBO uint32
	gl.GenBuffers(1, &EBO)

	// Bind texture
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, img.TextureID)

	// Setup buffers
	gl.BindVertexArray(textureVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, textureVBO)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices[0]))*len(vertices), gl.Ptr(vertices), gl.DYNAMIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(unsafe.Sizeof(indices[0]))*len(indices), gl.Ptr(indices), gl.DYNAMIC_DRAW)

	// Enable blending for transparency
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// Draw
	gl.UseProgram(textureShaderProgram)
	gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)

	gl.Disable(gl.BLEND)
	gl.DeleteBuffers(1, &EBO)
}

// Delete image texture
// This function deletes the OpenGL texture associated with the image.
func deleteImage(img *Image) {
	gl.DeleteTextures(1, &img.TextureID)
}
