package main

import (
	"Wasm-gl-trainning/internal/wasm/utils"
	_ "Wasm-gl-trainning/internal/wasm/utils"
	"github.com/go-gl/mathgl/mgl32"
	webgl "github.com/seqsense/webgl-go"
	"syscall/js"
	"unsafe"
)

const vsSource = `
attribute vec3 position;
uniform mat4 Pmatrix;
uniform mat4 Vmatrix;
uniform mat4 Mmatrix;
attribute vec3 color;
varying vec3 vColor;

void main(void) {
  gl_Position = Pmatrix*Vmatrix*Mmatrix*vec4(position, 1.);
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
var movMatrix mgl32.Mat4
var tmark float32
var rotation = float32(0)

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

func render(now float32) {
	width := gl.Canvas.ClientWidth()
	height := gl.Canvas.ClientHeight()
	gl.UseProgram(program)

	// Associate attributes to vertex shader
	PositionMatrix := gl.GetUniformLocation(program, "Pmatrix")
	ViewMatrix := gl.GetUniformLocation(program, "Vmatrix")
	ModelMatrix := gl.GetUniformLocation(program, "Mmatrix")

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

	ratio := float32(width / height)

	// Generate and apply projection matrix
	projMatrix := mgl32.Perspective(mgl32.DegToRad(45.0), ratio, 1, 100.0)
	var projMatrixBuffer *[16]float32
	projMatrixBuffer = (*[16]float32)(unsafe.Pointer(&projMatrix))
	projMatrixParsed := utils.NewMatrixFromArr(*projMatrixBuffer)
	gl.UniformMatrix4fv(PositionMatrix, false, projMatrixParsed)

	// Generate and apply view matrix
	viewMatrix := mgl32.LookAtV(mgl32.Vec3{3.0, 3.0, 3.0}, mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.0, 1.0, 0.0})
	var viewMatrixBuffer *[16]float32
	viewMatrixBuffer = (*[16]float32)(unsafe.Pointer(&viewMatrix))
	viewMatrixParsed := utils.NewMatrixFromArr(*viewMatrixBuffer)
	gl.UniformMatrix4fv(ViewMatrix, false, viewMatrixParsed)

	tdiff := now - tmark
	tmark = now
	rotation = rotation + float32(tdiff)/500

	// Do new model matrix calculations
	movMatrix = mgl32.HomogRotate3DX(0.5 * rotation)
	movMatrix = movMatrix.Mul4(mgl32.HomogRotate3DY(0.3 * rotation))
	movMatrix = movMatrix.Mul4(mgl32.HomogRotate3DZ(0.2 * rotation))

	// Convert model matrix to a JS TypedArray
	var modelMatrixBuffer *[16]float32
	modelMatrixBuffer = (*[16]float32)(unsafe.Pointer(&movMatrix))

	// Apply the model matrix
	modelMatrixParsed := utils.NewMatrixFromArr(*modelMatrixBuffer)
	gl.UniformMatrix4fv(ModelMatrix, false, modelMatrixParsed)

	gl.DrawArrays(gl.TRIANGLES, 0, len(vertices)/3)

}

func main() {
	setup()
	movMatrix = mgl32.Ident4()

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
		now := args[0].Float()
		//fmt.Println("Actual time ", now)
		render(float32(now))
		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	defer renderFrame.Release()

	// Start running
	js.Global().Call("requestAnimationFrame", renderFrame)
	select {}
}
