//go:build !wasm
// +build !wasm

package graphics

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"unsafe"
)

// Convert screen coordinates to OpenGL coordinates
func screenToGL(x, y float32) (float32, float32) {
	glX := (x / float32(windowWidth)) * 2.0 - 1.0
	glY := 1.0 - (y / float32(windowHeight)) * 2.0
	return glX, glY
}

// Draw vertices using OpenGL
func drawVertices(vertices []float32, count int32, mode uint32) {
	gl.BindVertexArray(currentVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, currentVBO)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices[0]))*len(vertices), gl.Ptr(vertices), gl.DYNAMIC_DRAW)
	
	gl.UseProgram(shaderProgram)
	gl.DrawArrays(mode, 0, count)
}

// Draw with indices 
func drawIndexed(vertices []float32, indices []uint32) {
	var EBO uint32
	gl.GenBuffers(1, &EBO)
	
	gl.BindVertexArray(currentVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, currentVBO)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(vertices[0]))*len(vertices), gl.Ptr(vertices), gl.DYNAMIC_DRAW)
	
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(unsafe.Sizeof(indices[0]))*len(indices), gl.Ptr(indices), gl.DYNAMIC_DRAW)
	
	gl.UseProgram(shaderProgram)
	gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)
	
	gl.DeleteBuffers(1, &EBO)
}
