package main

import (
	"Wasm-gl-trainning/internal/wasm/utils"
	_ "Wasm-gl-trainning/internal/wasm/utils"
	webgl "github.com/seqsense/webgl-go"
	"syscall/js"
)

const vsSource = `
attribute vec3 position;
uniform vec3 offset;
attribute vec3 color;
varying vec3 vColor;

void main(void) {
  gl_Position = vec4(position, 1.0) + vec4(offset, 1.0);
  //gl_Position.x = gl_Position.x + 1.0;
  vColor = color;
}
`

const fsSource = `
precision mediump float;
varying vec3 vColor;
void main(void) {
  gl_FragColor = vec4(vColor, 1.);
}
`

var vertices = []float32{
	-0.5, -0.5, 0,
	0.5, -0.5, 0,
	0, 0.5, 0,
}

var colors = []float32{
	1, 0, 0,
	0, 0, 255,
	0, 0, 0,
}

var gl *webgl.WebGL
var program webgl.Program
var colorBuffer webgl.Buffer
var vertexBuffer webgl.Buffer
var mousePos [2]float64
var deltaMouse [2]float64
var uniform_fetched bool
var cube_offset utils.Vector3f
var offset webgl.Location

func setup() {
	canvas := js.Global().Get("document").Call("getElementById", "glcanvas")
	var err error
	gl, err = webgl.New(canvas)
	if err != nil {
		panic(err)
	}

	vertexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, webgl.Float32ArrayBuffer(vertices), gl.STATIC_DRAW)

	colorBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, colorBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, webgl.Float32ArrayBuffer(colors), gl.STATIC_DRAW)

	var vs, fs webgl.Shader
	if vs, err = initVertexShader(gl, vsSource); err != nil {
		panic(err)
	}

	if fs, err = initFragmentShader(gl, fsSource); err != nil {
		panic(err)
	}

	program, err = linkShaders(gl, nil, vs, fs)
	if err != nil {
		panic(err)
	}
}

func render() {
	width := gl.Canvas.ClientWidth()
	height := gl.Canvas.ClientHeight()
	gl.GetUniformLocation(program, "offset")
	gl.UseProgram(program)
	if uniform_fetched == false {
		offset = gl.GetUniformLocation(program, "offset")
		uniform_fetched = true
	}
	gl.Uniform3fv(offset, cube_offset)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	position := gl.GetAttribLocation(program, "position")
	gl.VertexAttribPointer(position, 3, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(position)

	gl.BindBuffer(gl.ARRAY_BUFFER, colorBuffer)
	color := gl.GetAttribLocation(program, "color")
	gl.VertexAttribPointer(color, 3, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(color)

	gl.ClearColor(0.5, 0.5, 0.5, 0.9)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)
	gl.Viewport(0, 0, width, height)
	gl.DrawArrays(gl.TRIANGLES, 0, len(vertices)/3)

}

func main() {
	setup()
	cube_offset = utils.Vector3f{}
	doc := js.Global().Get("document")
	mouseMoveEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		deltaMouse[0] = e.Get("clientX").Float() - deltaMouse[0]
		deltaMouse[1] = e.Get("clientY").Float() - deltaMouse[1]
		if deltaMouse[0]*0.01 > 10 {
			cube_offset.SetX(cube_offset.GetX() + 0.01*(float32(deltaMouse[0])*0.001))
		} else if deltaMouse[0]*0.01 > 10 {
			cube_offset.SetX(cube_offset.GetX() - 0.01*(float32(deltaMouse[0])*0.001))
		}
		if deltaMouse[1]*0.01 < 10 {
			cube_offset.SetY(cube_offset.GetY() + 0.01*(float32(deltaMouse[1])*0.001))
		} else if deltaMouse[1]*0.01 > 10 {
			cube_offset.SetY(cube_offset.GetY() - 0.01*(float32(deltaMouse[1])*0.001))
		}
		return nil
	})
	defer mouseMoveEvt.Release()

	doc.Call("addEventListener", "mousemove", mouseMoveEvt)
	var renderFrame js.Func

	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		//now := args[0].Float()
		//fmt.Println("Actual time ", now)
		render()
		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	defer renderFrame.Release()

	// Start running
	js.Global().Call("requestAnimationFrame", renderFrame)
	select {}
}
