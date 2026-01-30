package views

import (
	"math"

	s "github.com/JuanBiancuzzo/own_wiki/core/scene"
)

type LerpFollow[V s.Value, VecOrVal s.Vector[V] | s.Value] struct {
	LerpFunction s.FnLerp[VecOrVal]
}

func NewLerpFollow[V s.Value, VecOrVal s.Vector[V] | s.Value](lerpFunction s.FnLerp[VecOrVal]) LerpFollow[V, VecOrVal] {
	return LerpFollow[V, VecOrVal]{
		LerpFunction: lerpFunction,
	}
}

func (f LerpFollow[_, VecOrVal]) Follow(K, dt float64, current, target VecOrVal) VecOrVal {
	t := 1 - math.Pow(K, -dt)
	return f.LerpFunction(s.UnitRange(t), current, target)
}

type SecondOrderFollowVec2[V s.Value] struct {
	FollowX, FollowY SecondOrderFollow[V]
}

func NewSecondOrderFollowVec2[V s.Value](frequencyHz, damping, initialResponse float64, initialTargetPosition s.Vector2[V]) SecondOrderFollowVec2[V] {
	return SecondOrderFollowVec2[V]{
		FollowX: NewSecondOrderFollow(frequencyHz, damping, initialResponse, V(initialTargetPosition.X)),
		FollowY: NewSecondOrderFollow(frequencyHz, damping, initialResponse, V(initialTargetPosition.Y)),
	}
}

func (f SecondOrderFollowVec2[V]) Follow(dt float64, target s.Vector2[V]) s.Vector2[V] {
	return s.NewVector2(
		f.FollowX.Follow(dt, target.X),
		f.FollowY.Follow(dt, target.Y),
	)
}

type SecondOrderFollowVec3[V s.Value] struct {
	FollowX, FollowY, FollowZ SecondOrderFollow[V]
}

func NewSecondOrderFollowVec3[V s.Value](frequencyHz, damping, initialResponse float64, initialTargetPosition s.Vector3[V]) SecondOrderFollowVec3[V] {
	return SecondOrderFollowVec3[V]{
		FollowX: NewSecondOrderFollow(frequencyHz, damping, initialResponse, V(initialTargetPosition.X)),
		FollowY: NewSecondOrderFollow(frequencyHz, damping, initialResponse, V(initialTargetPosition.Y)),
		FollowZ: NewSecondOrderFollow(frequencyHz, damping, initialResponse, V(initialTargetPosition.Z)),
	}
}

func (f SecondOrderFollowVec3[V]) Follow(dt float64, target s.Vector3[V]) s.Vector3[V] {
	return s.NewVector3(
		f.FollowX.Follow(dt, target.X),
		f.FollowY.Follow(dt, target.Y),
		f.FollowZ.Follow(dt, target.Z),
	)
}

type SecondOrderFollowVec4[V s.Value] struct {
	FollowX, FollowY, FollowZ, FollowW SecondOrderFollow[V]
}

func NewSecondOrderFollowVec4[V s.Value](frequencyHz, damping, initialResponse float64, initialTargetPosition s.Vector4[V]) SecondOrderFollowVec4[V] {
	return SecondOrderFollowVec4[V]{
		FollowX: NewSecondOrderFollow(frequencyHz, damping, initialResponse, V(initialTargetPosition.X)),
		FollowY: NewSecondOrderFollow(frequencyHz, damping, initialResponse, V(initialTargetPosition.Y)),
		FollowZ: NewSecondOrderFollow(frequencyHz, damping, initialResponse, V(initialTargetPosition.Z)),
		FollowW: NewSecondOrderFollow(frequencyHz, damping, initialResponse, V(initialTargetPosition.W)),
	}
}

func (f SecondOrderFollowVec4[V]) Follow(dt float64, target s.Vector4[V]) s.Vector4[V] {
	return s.NewVector4(
		f.FollowX.Follow(dt, target.X),
		f.FollowY.Follow(dt, target.Y),
		f.FollowZ.Follow(dt, target.Z),
		f.FollowW.Follow(dt, target.W),
	)
}

// This function is the result of a second order system, given by:
// y + k1 * y' + k2 * y” = x + k3 * x'
// Where (f = frequencyHz, zeta = damping, r = initialResponse):
// k1 = zeta / (pi * f)
// k2 = 1 / (2 * pi * f)^2
// k3 = (r * zeta) / (2 * pi * f)
type SecondOrderFollow[V s.Value] struct {
	K1, K2, K3      float64
	PreviosTarget   V
	CurrentPosition V
	CurrentVelocity V
}

func NewSecondOrderFollow[V s.Value](frequencyHz, damping, initialResponse float64, initialTargetPosition V) SecondOrderFollow[V] {
	return SecondOrderFollow[V]{
		K1:              damping / (math.Pi * frequencyHz),
		K2:              1 / math.Pow(2*math.Pi*frequencyHz, 2),
		K3:              (initialResponse * damping) / (2 * math.Pi * frequencyHz),
		PreviosTarget:   initialTargetPosition,
		CurrentPosition: initialTargetPosition,
		// The velocity is by default 0
	}
}

func (f SecondOrderFollow[V]) Follow(dt float64, target V) V {
	// Estimate the target velocity
	targetVelocity := float64(target-f.PreviosTarget) / dt
	f.PreviosTarget = target

	// Update position an velocity
	nextPosition := float64(f.CurrentPosition) + dt*float64(f.CurrentVelocity)
	f.CurrentVelocity += V(dt * (float64(target) + f.K3*targetVelocity - nextPosition - f.K1*float64(f.CurrentVelocity)) / f.K2)
	f.CurrentPosition = V(nextPosition)

	return f.CurrentPosition
}
