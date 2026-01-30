package views

import s "github.com/JuanBiancuzzo/own_wiki/core/scene"

type FnUpdate func(s.UnitRange)

type Tween struct {
	Time *Timer
}

func NewTween() *Tween {
	return &Tween{}
}

func (t *Tween) Tween(duration, dt float64, update FnUpdate) {
	if t.Time == nil {
		t.Time = NewTimer()
		t.Time.Start(duration)
	}

	if t.Time.HasEnded() {
		return
	}
	defer t.Time.Add(dt)

	update(t.Time.InvLerp())
}
