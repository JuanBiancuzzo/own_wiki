package platforms

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"
)

type TerminalPlatform struct{}

func NewTerminalPlatform() *TerminalPlatform {
	return &TerminalPlatform{}
}

func (tp *TerminalPlatform) HandleInput(chan e.Event) {}

func (tp *TerminalPlatform) Render(v.SceneRepresentation) {}

func (tp *TerminalPlatform) Close() {}
