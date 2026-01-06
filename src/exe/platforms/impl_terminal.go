//go:build terminal
// +build terminal

package platforms

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
	p "github.com/JuanBiancuzzo/own_wiki/src/core/platform"
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"
)

type TerminalPlatform struct{}

func GetPlatformImplementation() p.Platform {
	return &TerminalPlatform{}
}

func (tp *TerminalPlatform) HandleInput(chan e.Event) {}

func (tp *TerminalPlatform) Render(v.SceneRepresentation) {}

func (tp *TerminalPlatform) Close() {}
