//go:build wasm
// +build wasm

package graphics
import (
    "fmt"
    "math"
    "syscall/js"
)
type Image struct {
    TextureID uint32 // Can be unused in WASM
    Width     int32
    Height    int32
    filePath  string
    jsImage   js.Value // JavaScript image object (HTMLImageElement or canvas)
}
type DrawOptions struct {
    X, Y          float32 // Destination position
    Width, Height float32 // Destination size
    Rotation      float32 // Not used yet
    Tint          Color   // Tint color (RGBA)
    SrcX, SrcY    float32 // Source rect X, Y
    SrcW, SrcH    float32 // Source rect Width, Height
}
func loadImage(filePath string) (*Image, error) {
    img := &Image{
        filePath: filePath,
    }
    // Create a new HTMLImageElement
    img.jsImage = js.Global().Get("Image").New()
    img.jsImage.Set("src", filePath)
    // Wait for the image to load
    promise := img.jsImage.Call("addEventListener", "load", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        img.Width = int32(img.jsImage.Get("width").Float())
        img.Height = int32(img.jsImage.Get("height").Float())
        return nil
    }))
    if !promise.IsUndefined() {
        return nil, fmt.Errorf("failed to load image: %s", filePath)
    }
    return img, nil
}
func drawImage(img *Image, x, y float32) {
    ctx.Call("drawImage", img.jsImage, float64(x), float64(y))
}
func drawImageRotated(img *Image, x, y, angle float32) {
    cx := float64(x + float32(img.Width)/2)
    cy := float64(y + float32(img.Height)/2)
    radians := float64(angle) * math.Pi / 180
    ctx.Call("save")
    ctx.Call("translate", cx, cy)
    ctx.Call("rotate", radians)
    ctx.Call("drawImage", img.jsImage, -float64(img.Width)/2, -float64(img.Height)/2)
    ctx.Call("restore")
}
func drawImageScaled(img *Image, x, y, w, h float32) {
    ctx.Call("drawImage", img.jsImage, float64(x), float64(y), w *float32(img.Width), float32(img.Height)*h)
}

func drawImageTinted(img *Image, x, y float32, tint Color) {
    ctx.Call("save")
    ctx.Call("drawImage", img.jsImage, float64(x), float64(y))
    ctx.Set("globalAlpha", float64(tint.A))
    ctx.Set("globalCompositeOperation", "multiply")
    ctx.Set("fillStyle", colorToRGBA(tint))
    ctx.Call("fillRect", float64(x), float64(y), float64(img.Width), float64(img.Height))
    ctx.Call("restore")
}
func drawImageEx(img *Image, opts DrawOptions) {
    if opts.SrcW == 0 {
        opts.SrcW = float32(img.Width)
    }
    if opts.SrcH == 0 {
        opts.SrcH = float32(img.Height)
    }
    if opts.Width == 0 {
        opts.Width = float32(img.Width)
    }
    if opts.Height == 0 {
        opts.Height = float32(img.Height)
    }
    ctx.Call("save")
    // Move to center of image to apply rotation
    cx := float64(opts.X + opts.Width/2)
    cy := float64(opts.Y + opts.Height/2)
    radians := float64(opts.Rotation) * math.Pi / 180
    ctx.Call("translate", cx, cy)
    if opts.Rotation != 0 {
        ctx.Call("rotate", radians)
    }
    // Draw image with source cropping and destination scaling
    ctx.Call("drawImage",
        img.jsImage,
        float64(opts.SrcX), float64(opts.SrcY),
        float64(opts.SrcW), float64(opts.SrcH),
        -float64(opts.Width)/2, -float64(opts.Height)/2,
        float64(opts.Width), float64(opts.Height),
    )
    // Apply tint if alpha > 0
    if opts.Tint.A > 0 {
        ctx.Set("globalAlpha", float64(opts.Tint.A))
        ctx.Set("globalCompositeOperation", "multiply")
        ctx.Set("fillStyle", colorToRGBA(opts.Tint))
        ctx.Call("fillRect",
            -float64(opts.Width)/2, -float64(opts.Height)/2,
            float64(opts.Width), float64(opts.Height),
        )
    }
    ctx.Call("restore")
}
func deleteImage(img *Image) {
    if img.jsImage.IsUndefined() || img.jsImage.IsNull() {
        return
    }
    img.jsImage.Call("remove")
    img.jsImage = js.Undefined() // Clear the reference
}