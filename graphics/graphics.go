package graphics

import "github.com/QOthman/Pixu"

// Type aliases for backward compatibility
type Color = pixu.Color
type Image = pixu.Image
type DrawOptions = pixu.DrawOptions
type Vec2 = pixu.Vec2
type Rect = pixu.Rect
type Camera2D = pixu.Camera2D
type Font = pixu.Font
type Key = pixu.Key
type MouseButton = pixu.MouseButton

// Predefined colors
var (
	BLACK      = pixu.BLACK
	WHITE      = pixu.WHITE
	RED        = pixu.RED
	GREEN      = pixu.GREEN
	BLUE       = pixu.BLUE
	YELLOW     = pixu.YELLOW
	CYAN       = pixu.CYAN
	MAGENTA    = pixu.MAGENTA
	GRAY       = pixu.GRAY
	DARKGRAY   = pixu.DARKGRAY
	LIGHTGRAY  = pixu.LIGHTGRAY
	ORANGE     = pixu.ORANGE
	PINK       = pixu.PINK
	PURPLE     = pixu.PURPLE
	VIOLET     = pixu.VIOLET
	DARKPURPLE = pixu.DARKPURPLE
	LIME       = pixu.LIME
	DARKGREEN  = pixu.DARKGREEN
	SKYBLUE    = pixu.SKYBLUE
	DARKBLUE   = pixu.DARKBLUE
	TEAL       = pixu.TEAL
	GOLD       = pixu.GOLD
	MAROON     = pixu.MAROON
	BEIGE      = pixu.BEIGE
	BROWN      = pixu.BROWN
	DARKBROWN  = pixu.DARKBROWN
	BLANK      = pixu.BLANK
)

// Key constants
const (
	KeySpace     = pixu.KeySpace
	KeyEscape    = pixu.KeyEscape
	KeyEnter     = pixu.KeyEnter
	KeyTab       = pixu.KeyTab
	KeyBackspace = pixu.KeyBackspace
	KeyDelete    = pixu.KeyDelete

	KeyUp    = pixu.KeyUp
	KeyDown  = pixu.KeyDown
	KeyLeft  = pixu.KeyLeft
	KeyRight = pixu.KeyRight

	KeyW = pixu.KeyW
	KeyA = pixu.KeyA
	KeyS = pixu.KeyS
	KeyD = pixu.KeyD

	Key0 = pixu.Key0
	Key1 = pixu.Key1
	Key2 = pixu.Key2
	Key3 = pixu.Key3
	Key4 = pixu.Key4
	Key5 = pixu.Key5
	Key6 = pixu.Key6
	Key7 = pixu.Key7
	Key8 = pixu.Key8
	Key9 = pixu.Key9

	MouseLeft   = pixu.MouseLeft
	MouseRight  = pixu.MouseRight
	MouseMiddle = pixu.MouseMiddle
)

// Engine, Rendering & Lifecycle API
var (
	Init            = pixu.Init
	InitWindow      = pixu.InitWindow
	Close           = pixu.Close
	Terminate       = pixu.Terminate
	ShouldContinue  = pixu.ShouldContinue
	ShouldClose     = pixu.ShouldClose
	ClearBackground = pixu.ClearBackground
	Clear           = pixu.Clear
	Present         = pixu.Present
	BeginDrawing    = pixu.BeginDrawing
	EndDrawing      = pixu.EndDrawing
)

// Color API
var (
	NewColor      = pixu.NewColor
	RGB           = pixu.RGB
	RGBA          = pixu.RGBA
	Hex           = pixu.Hex
	HexAlpha      = pixu.HexAlpha
	ColorFromHex  = pixu.ColorFromHex
)

// Image & Texture API
var (
	LoadImage          = pixu.LoadImage
	LoadImageFromBytes = pixu.LoadImageFromBytes
	DrawImage          = pixu.DrawImage
	DrawImageV         = pixu.DrawImageV
	DrawImageScaled    = pixu.DrawImageScaled
	DrawImageRotated   = pixu.DrawImageRotated
	DrawImageTinted    = pixu.DrawImageTinted
	DrawImageRec       = pixu.DrawImageRec
	DrawImagePro       = pixu.DrawImagePro
	DrawImageEx        = pixu.DrawImageEx
)

// Shapes API
var (
	DrawPixel              = pixu.DrawPixel
	DrawLine               = pixu.DrawLine
	DrawLineV              = pixu.DrawLineV
	DrawLineThick          = pixu.DrawLineThick
	DrawLineThickV         = pixu.DrawLineThickV
	DrawLineBezier         = pixu.DrawLineBezier
	DrawPolyline           = pixu.DrawPolyline
	DrawRectangle          = pixu.DrawRectangle
	DrawRectangleV         = pixu.DrawRectangleV
	DrawRectangleRec       = pixu.DrawRectangleRec
	DrawRectangleOutline   = pixu.DrawRectangleOutline
	DrawRectangleRounded   = pixu.DrawRectangleRounded
	DrawRectangleGradientH = pixu.DrawRectangleGradientH
	DrawRectangleGradientV = pixu.DrawRectangleGradientV
	DrawRectangleGradientEx= pixu.DrawRectangleGradientEx
	DrawTriangle           = pixu.DrawTriangle
	DrawTriangleV          = pixu.DrawTriangleV
	DrawTriangleOutline    = pixu.DrawTriangleOutline
	DrawCircle             = pixu.DrawCircle
	DrawCircleV            = pixu.DrawCircleV
	DrawCircleOutline      = pixu.DrawCircleOutline
	DrawCircleSector       = pixu.DrawCircleSector
	DrawEllipse            = pixu.DrawEllipse
	DrawRing               = pixu.DrawRing
	DrawPolygon            = pixu.DrawPolygon
	DrawPolygonOutline     = pixu.DrawPolygonOutline
)

// Text & Font API
var (
	DrawText               = pixu.DrawText
	DrawTextV              = pixu.DrawTextV
	DrawTextCentered       = pixu.DrawTextCentered
	DrawTextWithBackground = pixu.DrawTextWithBackground
	DrawTextOutline        = pixu.DrawTextOutline
	MeasureText            = pixu.MeasureText
	MeasureTextV           = pixu.MeasureTextV
	SetFont                = pixu.SetFont
	GetFont                = pixu.GetFont
	GetDefaultFont         = pixu.GetDefaultFont
	GetDefaultBoldFont     = pixu.GetDefaultBoldFont
)

// Input API
var (
	IsKeyPressed        = pixu.IsKeyPressed
	IsKeyDown           = pixu.IsKeyDown
	IsKeyJustPressed    = pixu.IsKeyJustPressed
	IsKeyTriggered      = pixu.IsKeyTriggered
	IsKeyJustReleased   = pixu.IsKeyJustReleased
	IsKeyUp             = pixu.IsKeyUp
	IsMousePressed      = pixu.IsMouseButtonPressed
	IsMouseButtonPressed= pixu.IsMouseButtonPressed
	IsMouseJustPressed  = pixu.IsMouseButtonJustPressed
	IsMouseJustReleased = pixu.IsMouseButtonJustReleased
	GetMousePosition    = pixu.GetMousePosition
	GetMousePositionV   = pixu.GetMousePositionV
	GetMouseX           = pixu.GetMouseX
	GetMouseY           = pixu.GetMouseY
	GetMouseDelta       = pixu.GetMouseDelta
	GetMouseDeltaV      = pixu.GetMouseDeltaV
	IsMouseMoved        = pixu.IsMouseMoved
	GetScrollDelta      = pixu.GetScrollDelta
	GetScrollPosition   = pixu.GetScrollPosition
	UpdateInput         = pixu.UpdateInput
	SetCursorVisible    = pixu.SetCursorVisible
	SetCursorDisabled   = pixu.SetCursorDisabled
)

// FPS, Time & Window API
var (
	GetFps          = pixu.GetFps
	GetFPS          = pixu.GetFPS
	InitFps         = pixu.InitFps
	SetTargetFPS    = pixu.SetTargetFPS
	Wait            = pixu.Wait
	GetDeltaTime    = pixu.GetDeltaTime
	GetFrameTime    = pixu.GetFrameTime
	GetTime         = pixu.GetTime
	ResetTimer      = pixu.ResetTimer
	GetWindowSize   = pixu.GetWindowSize
	GetWindowWidth  = pixu.GetWindowWidth
	GetWindowHeight = pixu.GetWindowHeight
	IsWindowResized = pixu.IsWindowResized
	SetWindowTitle  = pixu.SetWindowTitle
	SetWindowSize   = pixu.SetWindowSize
	TakeScreenshot  = pixu.TakeScreenshot
)
