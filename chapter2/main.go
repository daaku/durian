package main

import (
	"encoding/binary"
	"fmt"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/app/debug"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/gl"
)

var (
	vertexCoords = f32.Bytes(binary.LittleEndian,
		-1.0, -1.0,
		1.0, -1.0,
		-1.0, 1.0,
		1.0, 1.0,
	)
	elementIndexes = []byte{0, 1, 2, 3}
)

type game struct {
	vertextBuffer, elementBuffer gl.Buffer
	textures                     [2]gl.Texture
}

func (g *game) start() {
	g.vertextBuffer = gl.GenBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, g.vertextBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.STATIC_DRAW, vertexCoords)

	g.elementBuffer = gl.GenBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, g.elementBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, gl.STATIC_DRAW, elementIndexes)
}

func (g *game) stop() {
	fmt.Println("stop")
}

func (g *game) draw() {
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	debug.DrawFPS()
}

func (g *game) touch(e event.Touch) {
	fmt.Println("touch", e)
}

func main() {
	var g game
	app.Run(app.Callbacks{
		Start: g.start,
		Stop:  g.stop,
		Draw:  g.draw,
		Touch: g.touch,
	})
}
