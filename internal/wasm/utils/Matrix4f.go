package utils

import "github.com/seqsense/webgl-go"

type Matrix4f struct {
	webgl.Mat4
	matrix [4][4]float32
}

//@formatter:off

func InitIdentity() Matrix4f {
	mat := Matrix4f{}

	mat.matrix[0][0] = 1
	mat.matrix[1][0] = 0
	mat.matrix[2][0] = 0
	mat.matrix[3][0] = 0
	mat.matrix[0][1] = 0
	mat.matrix[1][1] = 1
	mat.matrix[2][1] = 0
	mat.matrix[3][1] = 0
	mat.matrix[0][2] = 0
	mat.matrix[1][2] = 0
	mat.matrix[2][2] = 1
	mat.matrix[3][2] = 0
	mat.matrix[0][3] = 0
	mat.matrix[1][3] = 0
	mat.matrix[2][3] = 0
	mat.matrix[3][3] = 1

	return mat
}

func InitTranslationV(vec Vector3f) Matrix4f {
	return InitTranslation(vec.x, vec.y, vec.z)
}

func InitTranslation(x, y, z float32) Matrix4f {
	mat := Matrix4f{}

	mat.matrix[0][0] = 1
	mat.matrix[1][0] = 0
	mat.matrix[2][0] = 0
	mat.matrix[3][0] = x
	mat.matrix[0][1] = 0
	mat.matrix[1][1] = 1
	mat.matrix[2][0] = 0
	mat.matrix[3][1] = y
	mat.matrix[0][2] = 0
	mat.matrix[1][2] = 0
	mat.matrix[2][0] = 1
	mat.matrix[3][2] = z
	mat.matrix[0][3] = 0
	mat.matrix[1][3] = 0
	mat.matrix[2][0] = 0
	mat.matrix[3][3] = 1

	return mat
}

func (matrix Matrix4f) Mul(r Matrix4f) Matrix4f {
	res := Matrix4f{}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			res.matrix[i][j] = matrix.matrix[i][0]*r.matrix[0][j] +
				matrix.matrix[i][1]*r.matrix[1][j] +
				matrix.matrix[i][2]*r.matrix[2][j] +
				matrix.matrix[i][3]*r.matrix[3][j]
		}
	}

	return res
}

//@formatter:on
