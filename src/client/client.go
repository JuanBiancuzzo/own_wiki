package main

import (
	"reflect"

	"github.com/hashicorp/go-plugin"

	"github.com/JuanBiancuzzo/own_wiki/src/core/api"
	"github.com/JuanBiancuzzo/own_wiki/src/core/ecv"

	"github.com/JuanBiancuzzo/own_wiki/src/shared"
)

type OwnWikiUserStructure struct {
	Plugin shared.UserDefineStructure

	Components map[string]reflect.Type
	Entities   map[string]reflect.Type
}

func NewOwnWiki() *OwnWikiUserStructure {
	return &OwnWikiUserStructure{
		Plugin: nil,

		Components: make(map[string]reflect.Type),
		Entities:   make(map[string]reflect.Type),
	}
}

func (o *OwnWikiUserStructure) LoadPlugin(path string) error {
	return nil
}

func (o *OwnWikiUserStructure) RegisterStructures() (*ecv.ECV, error) {
	if o.Plugin == nil {
		// Si definimos ecv como nil estamos diciendo que hubo un error
		return nil, nil
	}

	return nil, nil
}

func (o *OwnWikiUserStructure) Close() {}

func main() {
	ownWiki := NewOwnWiki()
	defer ownWiki.Close()

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: api.Handshake,
		Plugins: map[string]plugin.Plugin{
			"userStructureData": &api.UserStructureDataPlugin{Impl: ownWiki},
		},
	})
}
