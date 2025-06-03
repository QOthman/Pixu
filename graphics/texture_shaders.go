package graphics

import (
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

const textureVertexShaderSource = `
#version 330 core
layout (location = 0) in vec2 aPos;
layout (location = 1) in vec2 aTexCoord;
layout (location = 2) in vec4 aColor;

out vec2 TexCoord;
out vec4 Color;

void main() {
    gl_Position = vec4(aPos, 0.0, 1.0);
    TexCoord = aTexCoord;
    Color = aColor;
}
` + "\x00"

const textureFragmentShaderSource = `
#version 330 core
in vec2 TexCoord;
in vec4 Color;
out vec4 FragColor;

uniform sampler2D ourTexture;

void main() {
    vec4 texColor = texture(ourTexture, TexCoord);
    FragColor = texColor * Color;
}
` + "\x00"

func setupTextureShaders() {
	vertexShader := compileShader(textureVertexShaderSource, gl.VERTEX_SHADER)
	fragmentShader := compileShader(textureFragmentShaderSource, gl.FRAGMENT_SHADER)

	textureShaderProgram = gl.CreateProgram()
	gl.AttachShader(textureShaderProgram, vertexShader)
	gl.AttachShader(textureShaderProgram, fragmentShader)
	gl.LinkProgram(textureShaderProgram)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	// Set texture uniform
	gl.UseProgram(textureShaderProgram)
	textureUniform := gl.GetUniformLocation(textureShaderProgram, gl.Str("ourTexture\x00"))
	gl.Uniform1i(textureUniform, 0) // Texture unit 0
}

func setupTextureBuffers() {
	gl.GenVertexArrays(1, &textureVAO)
	gl.GenBuffers(1, &textureVBO)

	gl.BindVertexArray(textureVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, textureVBO)

	// Position attribute (location 0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 8*4, nil)
	gl.EnableVertexAttribArray(0)

	// Texture coordinate attribute (location 1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 8*4, unsafe.Pointer(uintptr(2*4)))
	gl.EnableVertexAttribArray(1)

	// Color attribute (location 2)
	gl.VertexAttribPointer(2, 4, gl.FLOAT, false, 8*4, unsafe.Pointer(uintptr(4*4)))
	gl.EnableVertexAttribArray(2)
}
