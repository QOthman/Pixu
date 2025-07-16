package graphics


type Key int

const (
	KeySpace Key = iota
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	KeyLeft
	KeyRight
	KeyUp
	KeyDown
	KeyEnter
	KeyEscape
)


func Init(width, height int, title string) error {
	return initPlatform(width, height, title)
}

func ClearBackground(color Color) {
	clearBackground(color)
}

func DrawLine(x1, y1, x2, y2 float32, color Color) {
	drawLine(x1, y1, x2, y2, color)
}

func DrawTriangle(x1, y1, x2, y2, x3, y3 float32, color Color) {
	drawTriangle(x1, y1, x2, y2, x3, y3, color)
}

func DrawRectangle(x, y, width, height float32, color Color) {
	drawRectangle(x, y, width, height, color)
}

func DrawCircle(centerX, centerY, radius float32, color Color) {
	drawCircle(centerX, centerY, radius, color)
}

func DrawRectangleOutline(x, y, width, height float32, color Color) {
	drawRectangleOutline(x, y, width, height, color)
}

func DrawText(text string, x, y float32, size float32, color Color) {
	drawText(text, x, y, size, color)
}

func DrawTextCentered(text string, centerX, centerY, size float32, color Color) {
	drawTextCentered(text, centerX, centerY, size, color)
}

func DrawTextWithBackground(text string, x, y, size float32, bgColor, textColor Color, paddingX, paddingY float32) {
	drawTextWithBackground(text, x, y, size, bgColor, textColor,paddingX, paddingY)
}

func DrawTextOutline(text string, x, y, size float32, color, outlineColor Color) {
	drawTextOutline(text, x, y, size, color, outlineColor)
}

func LoadImage(filePath string) (*Image, error) {
	return loadImage(filePath)
}

func DrawImage(img *Image, x, y float32) {
	drawImage(img, x, y)
}

func DrawImageScaled(img *Image, x, y, scaleX, scaleY float32) {
	drawImageScaled(img, x, y, scaleX, scaleY)
}

func DrawImageRotated(img *Image, x, y, rotation float32) {
	drawImageRotated(img, x, y, rotation)
}

func DrawImageTinted(img *Image, x, y float32, tint Color) {
	drawImageTinted(img, x, y, tint)
}

func DrawImageEx(img *Image, options DrawOptions) {
	drawImageEx(img, options)
}

func DeleteImage(img *Image) {
	deleteImage(img)
}

func RunGameLoop(loopFunc func()) {
	runGameLoop(loopFunc)
}

func Close() {
	close()
}


func IsKeyPressed(key Key) bool {
	return isKeyPressed(key)
}
