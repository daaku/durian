package main

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"math"
	"time"

	"github.com/GeertJohan/go.rice"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/app/debug"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/gl"
	"golang.org/x/mobile/gl/glutil"
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
	startTime                    time.Time
	vertextBuffer, elementBuffer gl.Buffer
	textures                     [2]gl.Texture
	program                      gl.Program
	uniforms                     struct {
		fadeFactor gl.Uniform
		textures   [2]gl.Uniform
	}
	attribs struct {
		position gl.Attrib
	}
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

	g.program, err = glutil.CreateProgram(vertextShaderSource, fragmentShaderSource)
	if err != nil {
		panic(err)
	}

	g.uniforms.fadeFactor = gl.GetUniformLocation(g.program, "fade_factor")
	g.uniforms.textures[0] = gl.GetUniformLocation(g.program, "textures[0]")
	g.uniforms.textures[1] = gl.GetUniformLocation(g.program, "textures[1]")
	g.attribs.position = gl.GetAttribLocation(g.program, "position")
	g.startTime = time.Now()
}

func (g *game) stop() {
	fmt.Println("stop")
}

func (g *game) draw() {
	gl.UseProgram(g.program)

	fadeFactor := float32(math.Sin(time.Since(g.startTime).Seconds())*0.5 + 0.5)
	gl.Uniform1f(g.uniforms.fadeFactor, fadeFactor)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, g.textures[0])
	gl.Uniform1i(g.uniforms.textures[0], 0)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, g.textures[1])
	gl.Uniform1i(g.uniforms.textures[1], 1)

	gl.BindBuffer(gl.ARRAY_BUFFER, g.vertextBuffer)
	gl.EnableVertexAttribArray(g.attribs.position)
	gl.VertexAttribPointer(g.attribs.position, 2, gl.FLOAT, false, 0, 0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, g.elementBuffer)
	gl.DrawElements(gl.TRIANGLE_STRIP, gl.UNSIGNED_BYTE, 0, len(elementIndexes))

	gl.DisableVertexAttribArray(g.attribs.position)

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

const vertextShaderSource = `
#version 330

layout(location = 0) in vec2 position;
out vec2 texcoord;

void main() {
	gl_Position = vec4(position.x, -1*position.y, 0.0, 1.0);
	texcoord = position * vec2(0.5) + vec2(0.5);
}
`

const fragmentShaderSource = `
#version 330

uniform float fade_factor;
uniform sampler2D textures[2];

in vec2 texcoord;
out vec4 fragColor;

void main() {
	fragColor = mix(
		texture(textures[0], texcoord),
		texture(textures[1], texcoord),
		fade_factor
	);
}
`
