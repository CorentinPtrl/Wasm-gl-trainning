package utils

import (
	"github.com/seqsense/webgl-go"
	"math"
)

type Vector3f struct {
	webgl.Vec3
	x, y, z float32
}

func (vec *Vector3f) GetX() float32 {
	return vec.x
}

func (vec *Vector3f) SetX(x float32) {
	vec.x = x
}

func (vec *Vector3f) GetY() float32 {
	return vec.y
}

func (vec *Vector3f) SetY(y float32) {
	vec.y = y
}

func (vec *Vector3f) GetZ() float32 {
	return vec.z
}

func (vec *Vector3f) SetZ(z float32) {
	vec.z = z
}

func (vec *Vector3f) Length() float32 {
	return float32(math.Sqrt(float64(vec.x*vec.x + vec.y*vec.y + vec.z*vec.z)))
}

func (vec *Vector3f) Max() float32 {
	return float32(math.Max(float64(vec.x), math.Max(float64(vec.y), float64(vec.z))))
}

func (vec *Vector3f) Cross(r Vector3f) Vector3f {
	res := Vector3f{}

	res.x = vec.y*r.GetZ() - vec.z*r.GetY()
	res.y = vec.z*r.GetX() - vec.x*r.GetZ()
	res.z = vec.x*r.GetY() - vec.y*r.GetX()

	return res
}

func (vec *Vector3f) Normalized() Vector3f {
	length := vec.Length()
	res := Vector3f{}
	res.x = vec.x / length
	res.y = vec.y / length
	res.z = vec.z / length
	return res
}

func (vec *Vector3f) Lerp(dest Vector3f, lerpFactor float32) Vector3f {
	v := dest.SubV(*vec)
	v = v.Mul(lerpFactor)
	return v.AddV(*vec)
}

func (vec *Vector3f) AddV(r Vector3f) Vector3f {
	res := Vector3f{}
	res.x = vec.x + r.GetX()
	res.y = vec.y + r.GetY()
	res.z = vec.z + r.GetZ()

	return res
}

func (vec *Vector3f) Add(r float32) Vector3f {
	res := Vector3f{}
	res.x = vec.x + r
	res.y = vec.y + r
	res.z = vec.z + r

	return res
}

func (vec *Vector3f) SubV(r Vector3f) Vector3f {
	res := Vector3f{}
	res.x = vec.x - r.GetX()
	res.y = vec.y - r.GetY()
	res.z = vec.z - r.GetZ()

	return res
}

func (vec *Vector3f) Sub(r float32) Vector3f {
	res := Vector3f{}
	res.x = vec.x - r
	res.y = vec.y - r
	res.z = vec.z - r

	return res
}

func (vec *Vector3f) MulV(r Vector3f) Vector3f {
	res := Vector3f{}
	res.x = vec.x * r.GetX()
	res.y = vec.y * r.GetY()
	res.z = vec.z * r.GetZ()

	return res
}

func (vec *Vector3f) Mul(r float32) Vector3f {
	res := Vector3f{}
	res.x = vec.x * r
	res.y = vec.y * r
	res.z = vec.z * r

	return res
}

func (vec *Vector3f) DivV(r Vector3f) Vector3f {
	res := Vector3f{}
	res.x = vec.x / r.GetX()
	res.y = vec.y / r.GetY()
	res.z = vec.z / r.GetZ()

	return res
}

func (vec *Vector3f) Div(r float32) Vector3f {
	res := Vector3f{}
	res.x = vec.x / r
	res.y = vec.y / r
	res.z = vec.z / r

	return res
}

func (vec *Vector3f) Abs() Vector3f {
	res := Vector3f{}
	res.x = float32(math.Abs(float64(res.x)))
	res.y = float32(math.Abs(float64(res.y)))
	res.z = float32(math.Abs(float64(res.z)))
	return res
}

func (vec Vector3f) Floats() [3]float32 {
	return [3]float32{vec.x, vec.y, vec.z}
}
