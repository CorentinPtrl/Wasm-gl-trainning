package utils

import "github.com/seqsense/webgl-go"

type Matrix4f struct {
	webgl.Mat4
	matrix [16]float32
}

func NewMatrixFromArr(matrix [16]float32) Matrix4f {
	return Matrix4f{matrix: matrix}
}

func (mat Matrix4f) Floats() [16]float32 {
	return mat.matrix
}
