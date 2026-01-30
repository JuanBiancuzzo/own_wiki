package views

import (
	"math"

	s "github.com/JuanBiancuzzo/own_wiki/core/scene"
)

type EaseType uint

const (
	LINEAR_EASE = iota
	QUAD_EASE
	CUBIC_EASE
	QUART_EASE
	QUINT_EASE
	SINE_EASE
	EXPO_EASE
	CIRC_EASE
)

type EaseTimeApply uint

const (
	EA_BOTH = iota
	EA_IN
	EA_OUT
)

// ---+--- Easing functions ---+---
type FnEase func(in s.UnitRange) (out s.UnitRange)

/*
Example of use:

	func test() {
		in := 15
		out := Ease(s.InvLerpInt, s.LerpFloat, GetCircularEase(EA_BOTH), in, 10, 20, 0.5, 0.7)
		_ = out
	}
*/
func Ease[In any, Out any](invLerp s.FnInvLerp[In], lerp s.FnLerp[Out], ease FnEase, in, inStart, inStop In, outStart, outStop Out) Out {
	t := invLerp(in, inStart, inStop)
	return lerp(ease(t), outStart, outStop)
}

func GetLinearEase() FnEase {
	return easeLinear
}

func easeLinear(in s.UnitRange) s.UnitRange {
	return in
}

func GetPolinomialEase(power uint, timeApply EaseTimeApply) FnEase {
	return func(in s.UnitRange) s.UnitRange {
		return easePolynomial(in, power, timeApply)
	}
}

func toPolynomial(in s.UnitRange, power uint) (out s.UnitRange) {
	out = 1
	for range power {
		out *= in
	}
	return out
}

func easePolynomial(in s.UnitRange, power uint, timeApply EaseTimeApply) (out s.UnitRange) {
	switch timeApply {
	case EA_IN:
		out = toPolynomial(in, power)

	case EA_OUT:
		out = 1 - toPolynomial(1-in, power)

	case EA_BOTH:
		if in < 0.5 {
			out = toPolynomial(2*in, power)

		} else {
			out = 1 - toPolynomial(2*(1-in), power)
		}
	}
	return out
}

func GetSineEase(timeApply EaseTimeApply) FnEase {
	return func(in s.UnitRange) s.UnitRange {
		return s.UnitRange(easeSine(float64(in), timeApply))
	}
}

func easeSine(in float64, timeApply EaseTimeApply) (out float64) {
	switch timeApply {
	case EA_IN:
		out = 1 - math.Cos(0.5*math.Pi*in)

	case EA_OUT:
		out = math.Sin(0.5 * math.Pi * in)

	case EA_BOTH:
		out = 0.5 * (1 - math.Cos(0.5*math.Pi*in))
	}

	return out
}

func GetExponentialEase(power s.UnitRange, exp FnExp, timeApply EaseTimeApply) FnEase {
	return func(in s.UnitRange) s.UnitRange {
		return easeExponential(in, power, exp, timeApply)
	}
}

func easeExponential(in, power s.UnitRange, exp FnExp, timeApply EaseTimeApply) (out s.UnitRange) {
	switch timeApply {
	case EA_IN:
		out = exp(power * (in - 1))

	case EA_OUT:
		out = 1 - exp(-power*in)

	case EA_BOTH:
		if in < 0.5 {
			out = exp(power*(2*in-1) - 1)

		} else {
			out = 1 - exp(power*(1-2*in)-1)
		}
	}
	return out
}

func GetCircularEase(timeApply EaseTimeApply) FnEase {
	return func(in s.UnitRange) s.UnitRange {
		return s.UnitRange(easeCircular(float64(in), timeApply))
	}
}

func easeCircular(in float64, timeApply EaseTimeApply) (out float64) {
	switch timeApply {
	case EA_IN:
		out = 1 - math.Sqrt(1-in*in)

	case EA_OUT:
		out = math.Sqrt(1 - (in-2)*(in-2))

	case EA_BOTH:
		if in < 0.5 {
			out = 0.5 * (1 - math.Sqrt(1-4*in*in))

		} else {
			out = 0.5 * (math.Sqrt(1-4*(1-in)*(1-in)) - 1)
		}
	}
	return out
}

func GetSplineEase() FnEase {
	return easeSpline
}

// We need to implement a spline elements
func easeSpline(in s.UnitRange) (out s.UnitRange) {
	panic("TODO")
}

// ---+--- Exponential functions ---+---
type FnExp func(s.UnitRange) s.UnitRange

func Exp(power s.UnitRange) s.UnitRange {
	return s.UnitRange(math.Exp(float64(power)))
}

func Exp2(power s.UnitRange) s.UnitRange {
	return s.UnitRange(math.Exp2(float64(power)))
}

func ExpNum(base float64) FnExp {
	return func(power s.UnitRange) s.UnitRange {
		return s.UnitRange(math.Pow(base, float64(power)))
	}
}
