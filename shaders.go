package pixu

import (
	"fmt"
	"log"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Internal OpenGL shader programs and buffer IDs
var (
	shapeShaderProgram   uint32
	shapeVAO             uint32
	shapeVBO             uint32
	shapeEBO             uint32
	shapeProjLoc         int32

	textureShaderProgram uint32
	textureVAO           uint32
	textureVBO           uint32
	textureEBO           uint32
	textureProjLoc       int32
	textureSamplerLoc    int32

	// Current combined projection matrix (column-major)
	currentProjection [16]float32
)

// Default quad indices
var quadIndices = []uint32{0, 1, 2, 2, 3, 0}

const shapeVertexShaderSrc = `#version 330 core
layout (location = 0) in vec2 aPos;
layout (location = 1) in vec4 aColor;

uniform mat4 uProjection;

out vec4 vertexColor;

void main() {
    gl_Position = uProjection * vec4(aPos, 0.0, 1.0);
    vertexColor = aColor;
}
` + "\x00"

const shapeFragmentShaderSrc = `#version 330 core
in vec4 vertexColor;
out vec4 FragColor;

void main() {
    FragColor = vertexColor;
}
` + "\x00"

const textureVertexShaderSrc = `#version 330 core
layout (location = 0) in vec2 aPos;
layout (location = 1) in vec2 aTexCoord;
layout (location = 2) in vec4 aColor;

uniform mat4 uProjection;

out vec2 TexCoord;
out vec4 Color;

void main() {
    gl_Position = uProjection * vec4(aPos, 0.0, 1.0);
    TexCoord = aTexCoord;
    Color = aColor;
}
` + "\x00"

const textureFragmentShaderSrc = `#version 330 core
in vec2 TexCoord;
in vec4 Color;
out vec4 FragColor;

uniform sampler2D ourTexture;

void main() {
    FragColor = texture(ourTexture, TexCoord) * Color;
}
` + "\x00"

// initRenderer initializes all shader programs, buffers, and default render states.
func initRenderer() error {
	var err error

	// 1. Build Shape Shader Program
	shapeShaderProgram, err = createShaderProgram(shapeVertexShaderSrc, shapeFragmentShaderSrc)
	if err != nil {
		return fmt.Errorf("failed to create shape shader: %w", err)
	}
	shapeProjLoc = gl.GetUniformLocation(shapeShaderProgram, gl.Str("uProjection\x00"))

	// 2. Build Texture Shader Program
	textureShaderProgram, err = createShaderProgram(textureVertexShaderSrc, textureFragmentShaderSrc)
	if err != nil {
		return fmt.Errorf("failed to create texture shader: %w", err)
	}
	textureProjLoc = gl.GetUniformLocation(textureShaderProgram, gl.Str("uProjection\x00"))
	textureSamplerLoc = gl.GetUniformLocation(textureShaderProgram, gl.Str("ourTexture\x00"))

	// Set sampler uniform once
	gl.UseProgram(textureShaderProgram)
	gl.Uniform1i(textureSamplerLoc, 0)

	// 3. Setup Shape Buffers
	gl.GenVertexArrays(1, &shapeVAO)
	gl.GenBuffers(1, &shapeVBO)
	gl.GenBuffers(1, &shapeEBO)

	gl.BindVertexArray(shapeVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, shapeVBO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, shapeEBO)

	// Populate quad indices in EBO
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(quadIndices)*4, gl.Ptr(quadIndices), gl.STATIC_DRAW)

	// Shape layout: aPos (vec2 = 2 floats), aColor (vec4 = 4 floats) -> stride = 6 * 4 bytes
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 6*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 6*4, unsafe.Pointer(uintptr(2*4)))
	gl.EnableVertexAttribArray(1)

	// 4. Setup Texture Buffers
	gl.GenVertexArrays(1, &textureVAO)
	gl.GenBuffers(1, &textureVBO)
	gl.GenBuffers(1, &textureEBO)

	gl.BindVertexArray(textureVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, textureVBO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, textureEBO)

	// Populate quad indices in EBO
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(quadIndices)*4, gl.Ptr(quadIndices), gl.STATIC_DRAW)

	// Texture layout: aPos (vec2 = 2 floats), aTexCoord (vec2 = 2 floats), aColor (vec4 = 4 floats) -> stride = 8 * 4 bytes
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 8*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 8*4, unsafe.Pointer(uintptr(2*4)))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(2, 4, gl.FLOAT, false, 8*4, unsafe.Pointer(uintptr(4*4)))
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0)

	// 5. Global OpenGL configuration
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	updateProjectionMatrix()
	return nil
}

// destroyRenderer frees GPU resources when closing Pixu.
func destroyRenderer() {
	if shapeShaderProgram != 0 {
		gl.DeleteProgram(shapeShaderProgram)
		shapeShaderProgram = 0
	}
	if textureShaderProgram != 0 {
		gl.DeleteProgram(textureShaderProgram)
		textureShaderProgram = 0
	}
	if shapeVAO != 0 {
		gl.DeleteVertexArrays(1, &shapeVAO)
		shapeVAO = 0
	}
	if shapeVBO != 0 {
		gl.DeleteBuffers(1, &shapeVBO)
		shapeVBO = 0
	}
	if shapeEBO != 0 {
		gl.DeleteBuffers(1, &shapeEBO)
		shapeEBO = 0
	}
	if textureVAO != 0 {
		gl.DeleteVertexArrays(1, &textureVAO)
		textureVAO = 0
	}
	if textureVBO != 0 {
		gl.DeleteBuffers(1, &textureVBO)
		textureVBO = 0
	}
	if textureEBO != 0 {
		gl.DeleteBuffers(1, &textureEBO)
		textureEBO = 0
	}
}

// computeOrthoMatrix returns a 4x4 orthographic projection matrix in column-major order.
// Maps screen coords (0, 0) top-left to (width, height) bottom-right to NDC [-1, 1].
func computeOrthoMatrix(width, height float32) [16]float32 {
	if width <= 0 {
		width = 1.0
	}
	if height <= 0 {
		height = 1.0
	}

	return [16]float32{
		2.0 / width, 0, 0, 0,
		0, -2.0 / height, 0, 0,
		0, 0, -1.0, 0,
		-1.0, 1.0, 0, 1.0,
	}
}

// multiplyMat4 multiplies two 4x4 matrices in column-major order: out = a * b.
func multiplyMat4(a, b [16]float32) [16]float32 {
	var out [16]float32
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			sum := float32(0.0)
			for i := 0; i < 4; i++ {
				sum += a[i*4+row] * b[col*4+i]
			}
			out[col*4+row] = sum
		}
	}
	return out
}

// updateProjectionMatrix recalculates and uploads the projection matrix to active shaders.
func updateProjectionMatrix() {
	ortho := computeOrthoMatrix(float32(windowWidth), float32(windowHeight))
	if isCamera2D {
		camMat := computeCameraMatrix4(activeCamera)
		currentProjection = multiplyMat4(ortho, camMat)
	} else {
		currentProjection = ortho
	}

	if shapeShaderProgram != 0 {
		gl.UseProgram(shapeShaderProgram)
		gl.UniformMatrix4fv(shapeProjLoc, 1, false, &currentProjection[0])
	}
	if textureShaderProgram != 0 {
		gl.UseProgram(textureShaderProgram)
		gl.UniformMatrix4fv(textureProjLoc, 1, false, &currentProjection[0])
	}
}

// createShaderProgram compiles and links vertex and fragment shaders.
func createShaderProgram(vSrc, fSrc string) (uint32, error) {
	vShader, err := compileShader(vSrc, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(vShader)

	fShader, err := compileShader(fSrc, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(fShader)

	program := gl.CreateProgram()
	gl.AttachShader(program, vShader)
	gl.AttachShader(program, fShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLen)
		logStr := strings.Repeat("\x00", int(logLen+1))
		gl.GetProgramInfoLog(program, logLen, nil, gl.Str(logStr))
		gl.DeleteProgram(program)
		return 0, fmt.Errorf("failed to link shader program: %s", logStr)
	}

	return program, nil
}

// compileShader compiles an individual OpenGL shader.
func compileShader(src string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(src)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLen)
		logStr := strings.Repeat("\x00", int(logLen+1))
		gl.GetShaderInfoLog(shader, logLen, nil, gl.Str(logStr))
		gl.DeleteShader(shader)
		log.Printf("Shader compilation error (%d): %s", shaderType, logStr)
		return 0, fmt.Errorf("failed to compile shader: %s", logStr)
	}

	return shader, nil
}
