package views

import s "github.com/JuanBiancuzzo/own_wiki/core/scene"

type Counter struct {
	Counter, Min, Max int
}

func NewCounter() *Counter {
	return &Counter{
		Counter: 0,
		Min:     0,
		Max:     0,
	}
}

func (c *Counter) Start(min, max int) {
	if min > max {
		min, max = max, min
	}

	c.Counter = min
	c.Min = min
	c.Max = max
}

func (c *Counter) Restart() {
	c.Counter = c.Min
}

func (c *Counter) Add(deltaCount uint) {
	c.Counter += int(deltaCount)
	if c.Counter > c.Max {
		c.Counter = c.Max
	}
}

func (c Counter) InvLerp() s.UnitRange {
	return s.InvLerpInt(c.Counter, c.Min, c.Max)
}

func EaseCounter[Out any](c Counter, lerp s.FnLerp[Out], ease FnEase, outStart, outStop Out) Out {
	t := c.InvLerp()
	return lerp(ease(t), outStart, outStop)
}

func (c Counter) HasEnded() bool {
	return c.Counter >= c.Max
}
