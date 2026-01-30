package scene

import (
	"math"
	"math/cmplx"

	c "golang.org/x/exp/constraints"
)

const EPSILON = 1e-5

type IVec2[N c.Signed] vector2[N]
type IVec3[N c.Signed] vector3[N]
type IVec4[N c.Signed] vector4[N]

type UVec2[N c.Unsigned] vector2[N]
type UVec3[N c.Unsigned] vector3[N]
type UVec4[N c.Unsigned] vector4[N]

type Vec2[F c.Float] vector2[F]
type Vec3[F c.Float] vector3[F]
type Vec4[F c.Float] vector4[F]

type UvVector[F c.Float] vector2[F]
type ColorRGB[F c.Float] vector3[F]
type ColorRGBA[F c.Float] vector4[F]
type Quaternion[F c.Float] vector4[F]

type Mat2[F c.Float] Matrix2x2[F]
type Mat3[F c.Float] Matrix3x3[F]
type Mat4[F c.Float] Matrix4x4[F]

type UnitRange float64

type Value interface {
	c.Signed | c.Unsigned | c.Float
}

// ---+--- Lerp and InverseLerp ---+---
type LerpValue[V Value] interface {
	Value | vector2[V] | vector3[V] | vector4[V]
}

type FnInvLerp[T any] func(in, start, end T) UnitRange
type FnLerp[T any] func(t UnitRange, start, end T) T

// ---+--- Lerp & InvLerp of Value types ---+---

func InvLerpInt[I c.Integer](in, start, end I) UnitRange {
	return UnitRange(in-start) / UnitRange(end-start)
}

func LerpInt[I c.Integer](t UnitRange, start, end I) I {
	return I(t*UnitRange(end-start)) + start
}

func InvLerpFloat[F c.Float](in, start, end F) UnitRange {
	return UnitRange((in - start) / (end - start))
}

func LerpFloat[F c.Float](t UnitRange, start, end F) F {
	return F(t*UnitRange(end-start)) + start
}

func InvLerpComplex[C c.Complex](in, start, end C) UnitRange {
	return UnitRange(
		cmplx.Abs(
			complex128((in - start) / (end - start)),
		),
	)
}

func LerpComplex[C c.Complex](t UnitRange, start, end C) C {
	return C(complex(t, 0)*complex128(end-start)) + start
}

// ---+--- Lerp & InvLerp of Vectors ---+---
func InvLerpVec2[V Value](in, start, end vector2[V]) UnitRange {
	endToStart := end.Sub(start)
	startToValue := in.Sub(start)
	return UnitRange(endToStart.Dot(startToValue) / endToStart.Dot(endToStart))
}

func LerpVec2[V Value](t UnitRange, start, end vector2[V]) vector2[V] {
	endToStart := end.Sub(start)
	return start.Sum(endToStart.Mult(float64(t)))
}

func SlerpVec2[V Value](t UnitRange, start, end vector2[V]) vector2[V] {
	cosAngle := start.Normalize().Dot(end.Normalize())
	if -EPSILON < cosAngle && cosAngle < EPSILON {
		return LerpVec2(t, start, end)
	}

	scaleStart := start.Mult(math.Sin(float64(1-t)*cosAngle) / math.Sin(cosAngle))
	scaleEnd := end.Mult(math.Sin(float64(t)*cosAngle) / math.Sin(cosAngle))

	return scaleStart.Sum(scaleEnd)
}

func InvLerpVec3[V Value](in, start, end vector3[V]) UnitRange {
	endToStart := end.Sub(start)
	startToValue := in.Sub(start)
	return UnitRange(endToStart.Dot(startToValue) / endToStart.Dot(endToStart))
}

func LerpVec3[V Value](t UnitRange, start, end vector3[V]) vector3[V] {
	endToStart := end.Sub(start)
	return start.Sum(endToStart.Mult(float64(t)))
}

func SlerpVec3[V Value](t UnitRange, start, end vector3[V]) vector3[V] {
	cosAngle := start.Normalize().Dot(end.Normalize())
	if -EPSILON < cosAngle && cosAngle < EPSILON {
		return LerpVec3(t, start, end)
	}

	scaleStart := start.Mult(math.Sin(float64(1-t)*cosAngle) / math.Sin(cosAngle))
	scaleEnd := end.Mult(math.Sin(float64(t)*cosAngle) / math.Sin(cosAngle))

	return scaleStart.Sum(scaleEnd)
}

func InvLerpVec4[V Value](in, start, end vector4[V]) UnitRange {
	endToStart := end.Sub(start)
	startToValue := in.Sub(start)
	return UnitRange(endToStart.Dot(startToValue) / endToStart.Dot(endToStart))
}

func LerpVec4[V Value](t UnitRange, start, end vector4[V]) vector4[V] {
	endToStart := end.Sub(start)
	return start.Sum(endToStart.Mult(float64(t)))
}

func SlerpVec4[V Value](t UnitRange, start, end vector4[V]) vector4[V] {
	cosAngle := start.Normalize().Dot(end.Normalize())
	if -EPSILON < cosAngle && cosAngle < EPSILON {
		return LerpVec4(t, start, end)
	}

	scaleStart := start.Mult(math.Sin(float64(1-t)*cosAngle) / math.Sin(cosAngle))
	scaleEnd := end.Mult(math.Sin(float64(t)*cosAngle) / math.Sin(cosAngle))

	return scaleStart.Sum(scaleEnd)
}

// Custom types
func InvLerpUv[F c.Float](in, start, end UvVector[F]) UnitRange {
	return InvLerpVec2(vector2[F](in), vector2[F](start), vector2[F](end))
}

func LerpUv[F c.Float](t UnitRange, start, end UvVector[F]) UvVector[F] {
	return UvVector[F](LerpVec2(t, vector2[F](start), vector2[F](end)))
}

func SlerpUv[F c.Float](t UnitRange, start, end UvVector[F]) UvVector[F] {
	return UvVector[F](SlerpVec2(t, vector2[F](start), vector2[F](end)))
}

func InvLerpColorRGB[F c.Float](in, start, end ColorRGB[F]) UnitRange {
	return InvLerpVec3(vector3[F](in), vector3[F](start), vector3[F](end))
}

func LerpColorRGB[F c.Float](t UnitRange, start, end ColorRGB[F]) ColorRGB[F] {
	return ColorRGB[F](LerpVec3(t, vector3[F](start), vector3[F](end)))
}

func SlerpColorRGB[F c.Float](t UnitRange, start, end ColorRGB[F]) ColorRGB[F] {
	return ColorRGB[F](SlerpVec3(t, vector3[F](start), vector3[F](end)))
}

func InvLerpColorRGBA[F c.Float](in, start, end ColorRGBA[F]) UnitRange {
	return InvLerpVec4(vector4[F](in), vector4[F](start), vector4[F](end))
}

func LerpColorRGBA[F c.Float](t UnitRange, start, end ColorRGBA[F]) ColorRGBA[F] {
	return ColorRGBA[F](LerpVec4(t, vector4[F](start), vector4[F](end)))
}

func SlerpColorRGBA[F c.Float](t UnitRange, start, end ColorRGBA[F]) ColorRGBA[F] {
	return ColorRGBA[F](SlerpVec4(t, vector4[F](start), vector4[F](end)))
}

func InvLerpQuaternion[F c.Float](in, start, end Quaternion[F]) UnitRange {
	return InvLerpVec4(vector4[F](in), vector4[F](start), vector4[F](end))
}

func LerpQuaternion[F c.Float](t UnitRange, start, end Quaternion[F]) Quaternion[F] {
	return Quaternion[F](LerpVec4(t, vector4[F](start), vector4[F](end)))
}

func SlerpQuaternion[F c.Float](t UnitRange, start, end Quaternion[F]) Quaternion[F] {
	return Quaternion[F](SlerpVec4(t, vector4[F](start), vector4[F](end)))
}

// ---+--- definitions ---+---

type vector2[V Value] struct {
	X, Y V
}

func (v2 vector2[V]) Sum(other vector2[V]) vector2[V] {
	return vector2[V]{
		X: v2.X + other.X,
		Y: v2.Y + other.Y,
	}
}

func (v2 vector2[V]) Sub(other vector2[V]) vector2[V] {
	return vector2[V]{
		X: v2.X - other.X,
		Y: v2.Y - other.Y,
	}
}

func (v2 vector2[V]) Dot(other vector2[V]) float64 {
	return float64(v2.X*other.X + v2.Y*other.Y)
}

func (v2 vector2[V]) Mult(scalar float64) vector2[V] {
	return vector2[V]{
		X: V(float64(v2.X) * scalar),
		Y: V(float64(v2.Y) * scalar),
	}
}

func (v2 vector2[V]) Normalize() vector2[V] {
	return v2.Mult(1 / v2.Magnitud())
}

func (v2 vector2[V]) Magnitud() float64 {
	return math.Sqrt(v2.Dot(v2))
}

type vector3[V Value] struct {
	X, Y, Z V
}

func (v3 vector3[V]) Sum(other vector3[V]) vector3[V] {
	return vector3[V]{
		X: v3.X + other.X,
		Y: v3.Y + other.Y,
		Z: v3.Z + other.Z,
	}
}

func (v3 vector3[V]) Sub(other vector3[V]) vector3[V] {
	return vector3[V]{
		X: v3.X - other.X,
		Y: v3.Y - other.Y,
		Z: v3.Z - other.Z,
	}
}

func (v3 vector3[V]) Dot(other vector3[V]) float64 {
	return float64(v3.X*other.X + v3.Y*other.Y + v3.Z*other.Z)
}

func (v3 vector3[V]) Mult(scalar float64) vector3[V] {
	return vector3[V]{
		X: V(float64(v3.X) * scalar),
		Y: V(float64(v3.Y) * scalar),
		Z: V(float64(v3.Z) * scalar),
	}
}

func (v3 vector3[V]) Normalize() vector3[V] {
	return v3.Mult(1 / v3.Magnitud())
}

func (v3 vector3[V]) Magnitud() float64 {
	return math.Sqrt(v3.Dot(v3))
}

type vector4[V Value] struct {
	X, Y, Z, W V
}

func (v4 vector4[V]) Sum(other vector4[V]) vector4[V] {
	return vector4[V]{
		X: v4.X + other.X,
		Y: v4.Y + other.Y,
		Z: v4.Z + other.Z,
		W: v4.W + other.W,
	}
}

func (v4 vector4[V]) Sub(other vector4[V]) vector4[V] {
	return vector4[V]{
		X: v4.X - other.X,
		Y: v4.Y - other.Y,
		Z: v4.Z - other.Z,
		W: v4.W - other.W,
	}
}

func (v4 vector4[V]) Mult(scalar float64) vector4[V] {
	return vector4[V]{
		X: V(float64(v4.X) * scalar),
		Y: V(float64(v4.Y) * scalar),
		Z: V(float64(v4.Z) * scalar),
		W: V(float64(v4.W) * scalar),
	}
}

func (v4 vector4[V]) Dot(other vector4[V]) float64 {
	return float64(v4.X*other.X + v4.Y*other.Y + v4.Z*other.Z + v4.W*other.W)
}

func (v4 vector4[V]) Normalize() vector4[V] {
	return v4.Mult(1 / v4.Magnitud())
}

func (v4 vector4[V]) Magnitud() float64 {
	return math.Sqrt(v4.Dot(v4))
}

type Matrix2x2[V Value] struct {
	A11, A12 V
	A21, A22 V
}

type Matrix3x3[V Value] struct {
	A11, A12, A13 V
	A21, A22, A23 V
	A31, A32, A33 V
}

type Matrix4x4[V Value] struct {
	A11, A12, A13, A14 V
	A21, A22, A23, A24 V
	A31, A32, A33, A34 V
	A41, A42, A43, A44 V
}
