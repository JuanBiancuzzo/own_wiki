package api

import (
	"fmt"
	"net/rpc"

	"github.com/hashicorp/go-plugin"

	"github.com/JuanBiancuzzo/own_wiki/src/core/ecv"
)

// Definimos un handshake para poder controlar versiones del plugin
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "USER_STRUCTURE_DATA_PLUGIN",
	MagicCookieValue: "enabled",
}

// Definimos un mapeo de plugins, para poder hacer dispatch del plugin
var PluginMap = map[string]plugin.Plugin{
	"userStructureData": &UserStructureDataPlugin{},
}

/*
La view tiene que ser basicamente una maquina de estados definina dinamicamente por el cambio de
view al final de cada una.

La view inicialmente va a recibir una escena, la entidad que necesita, y los eventos que suceden entre el
frame anterior y este.

De la view puede salir eventos, que si son del sistema seran manejados en el exe, y sino sera pasado a los
siguientes views. Los eventos pueden ser:
  - Sistema General:
  	  * Cerrar el programa
  - View externa al usuario
  	  * Cambiar la configuracion
  - Sistema de datos
  	  * Actualizar un componente
  	  * Eliminar un componente
*/

type ErrorLoadPath struct {
	HasError    bool
	ErrorReason string
}

func NoErrorLoadPath() ErrorLoadPath {
	return ErrorLoadPath{
		HasError: false,
	}
}

func NewErrorLoadPath(reason string, args ...any) ErrorLoadPath {
	return ErrorLoadPath{
		HasError:    true,
		ErrorReason: fmt.Sprintf(reason, args...),
	}
}

type ReturnRegisterStructure struct {
	HasError    bool
	ErrorReason string
	Ecv         *ecv.ECV
}

func NewErrorRegisterStructure(reason string, args ...any) ReturnRegisterStructure {
	return ReturnRegisterStructure{
		HasError:    true,
		ErrorReason: fmt.Sprintf(reason, args...),
	}
}

func ReturnStructure(system *ecv.ECV) ReturnRegisterStructure {
	if system == nil {
		return NewErrorRegisterStructure("No system register")
	}

	return ReturnRegisterStructure{
		HasError: false,
		Ecv:      system,
	}
}

type UserStructureData interface {
	// Carga el plugin definido por el usuario,
	LoadPlugin(path string) (ErrorLoadPath, error)

	// Manera de obtener una estructura general de plugin definido por el usuario
	RegisterStructures() (ReturnRegisterStructure, error)
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
type UserStructureDataPlugin struct {
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl UserStructureData
}

func (p *UserStructureDataPlugin) Server(*plugin.MuxBroker) (any, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (*UserStructureDataPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (any, error) {
	return &RPCClient{client: c}, nil
}
