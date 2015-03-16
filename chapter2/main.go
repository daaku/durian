package main

import (
	"encoding/binary"
	"fmt"

	"image"
	"image/draw"
	_ "image/png"

	"github.com/GeertJohan/go.rice"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/app/debug"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/gl"
)

var (
	vertexCoords = f32.Bytes(binary.LittleEndian,
		-1.0, -1.0,
		+1.0, -1.0,
		-1.0, +1.0,
		+1.0, +1.0,
	)
	elementIndexes = []byte{0, 1, 2, 3}
)

func roundToPower2(x int) int {
	xa := 1
	for xa < x {
		xa *= 2
	}
	return xa
}

func toRGBA(in image.Image) *image.RGBA {
	switch i := in.(type) {
	case *image.RGBA:
		return i
	default:
		copy := image.NewRGBA(i.Bounds())
		draw.Draw(copy, i.Bounds(), i, image.Pt(0, 0), draw.Src)
		return copy
	}
}

type game struct {
	box                          *rice.Box
	vertextBuffer, elementBuffer gl.Buffer
	textures                     [2]gl.Texture
}

func (g *game) loadImage(name string) (gl.Texture, error) {
	r, err := g.box.Open(name)
	if err != nil {
		return gl.Texture{}, err
	}
	defer r.Close()

	oi, _, err := image.Decode(r)
	if err != nil {
		return gl.Texture{}, err
	}
	i := toRGBA(oi)
	size := i.Bounds().Size()

	t := gl.GenTexture()
	gl.BindTexture(gl.TEXTURE_2D, t)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(gl.TEXTURE_2D, 0, size.X, size.Y, gl.RGBA, gl.UNSIGNED_BYTE, i.Pix)
	return t, nil
}

func (g *game) start() {
	g.vertextBuffer = gl.GenBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, g.vertextBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.STATIC_DRAW, vertexCoords)

	g.elementBuffer = gl.GenBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, g.elementBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, gl.STATIC_DRAW, elementIndexes)

	var err error
	g.textures[0], err = g.loadImage("gl2-hello-1.png")
	if err != nil {
		panic(err)
	}

	g.textures[1], err = g.loadImage("gl2-hello-2.png")
	if err != nil {
		panic(err)
	}
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
	g := game{box: rice.MustFindBox("../resources")}
	app.Run(app.Callbacks{
		Start: g.start,
		Stop:  g.stop,
		Draw:  g.draw,
		Touch: g.touch,
	})
}
