//go:build !wasm
// +build !wasm

package graphics

import (
	"log"
	"github.com/go-gl/gl/v3.3-core/gl"
	"unsafe"
)

const vertexShaderSource = `
#version 330 core
layout (location = 0) in vec2 aPos;
layout (location = 1) in vec3 aColor;

out vec3 vertexColor;

void main() {
    gl_Position = vec4(aPos, 0.0, 1.0);
    vertexColor = aColor;
}
` + "\x00"

const fragmentShaderSource = `
#version 330 core
in vec3 vertexColor;
out vec4 FragColor;

void main() {
    FragColor = vec4(vertexColor, 1.0);
}
` + "\x00"

func setupShaders() {
	vertexShader := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragmentShader := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	
	shaderProgram = gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)
	
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
}

func setupBuffers() {
	gl.GenVertexArrays(1, &currentVAO)
	gl.GenBuffers(1, &currentVBO)
	
	gl.BindVertexArray(currentVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, currentVBO)
	
	// Position attribute
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)
	
	// Color attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 5*4, unsafe.Pointer(uintptr(2*4)))
	gl.EnableVertexAttribArray(1)
}

func compileShader(source string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)
	
	var success int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		logText := make([]byte, logLength)
		gl.GetShaderInfoLog(shader, logLength, nil, &logText[0])
		log.Printf("Shader compilation error: %s", string(logText))
	}
	return shader
}