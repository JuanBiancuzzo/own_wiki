package views

import (
	"math"

	s "github.com/JuanBiancuzzo/own_wiki/core/scene"
)

type LerpFollow[V s.Value, LV s.LerpValue[V]] struct {
	LerpFunction s.FnLerp[LV]
}

func NewLerpFollow[V s.Value, LV s.LerpValue[V]](lerpFunction s.FnLerp[LV]) LerpFollow[V, LV] {
	return LerpFollow[V, LV]{
		LerpFunction: lerpFunction,
	}
}

func (lf LerpFollow[_, LV]) Follow(K, dt float64, current, target LV) LV {
	t := 1 - math.Pow(K, -dt)
	return lf.LerpFunction(s.UnitRange(t), current, target)
}

// This function is the result of a second order system, given by:
// y + k1 * y' + k2 * y” = x + k3 * x'
// Where (f = frequencyHz, zeta = damping, r = initialResponse):
// k1 = zeta / (pi * f)
// k2 = 1 / (2 * pi * f)^2
// k3 = (r * zeta) / (2 * pi * f)
type SecondOrderFollow[V s.Value, LV s.LerpValue[V]] struct {
	K1, K2, K3    float64
	PreviosTarget LV
}

func NewSecondOrderFollow[V s.Value, LV s.LerpValue[V]](frequencyHz, damping, initialResponse float64) SecondOrderFollow[V, LV] {
	return SecondOrderFollow[V, LV]{
		K1: damping / (math.Pi * frequencyHz),
		K2: 1 / math.Pow(2*math.Pi*frequencyHz, 2),
		K3: (initialResponse * damping) / (2 * math.Pi * frequencyHz),
	}
}

func (sof SecondOrderFollow[_, LV]) Follow(dt float64, current, target LV) LV {
	panic("TODO")
}
