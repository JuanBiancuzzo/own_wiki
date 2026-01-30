package views

import s "github.com/JuanBiancuzzo/own_wiki/core/scene"

type Timer struct {
	Time, Max float64
}

func NewTimer() *Timer {
	return &Timer{
		Time: 0,
		Max:  0,
	}
}

func (t *Timer) Start(max float64) {
	if max < 0 {
		max = 0
	}

	t.Time = 0
	t.Max = max
}

func (t *Timer) Restart() {
	t.Time = 0
}

func (t *Timer) Add(dt float64) {
	if dt < 0 {
		dt = 0
	}

	t.Time += dt
	if t.Time > t.Max {
		t.Time = t.Max
	}
}

func (t Timer) InvLerp() s.UnitRange {
	return s.InvLerpFloat(t.Time, 0, t.Max)
}

func EaseTimer[Out any](t Timer, lerp s.FnLerp[Out], ease FnEase, outStart, outStop Out) Out {
	return lerp(ease(t.InvLerp()), outStart, outStop)
}

func (t Timer) HasEnded() bool {
	return t.Time >= t.Max
}
