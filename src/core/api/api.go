package api

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"

	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
	q "github.com/JuanBiancuzzo/own_wiki/src/core/query"
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"
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

type UploadEntity interface {
	Upload(entity any) error
}

type OWData interface {
	Query(q.QueryRequest) (any, error)

	SendEvent(e.Event) error
}

type UserStructureData interface {
	// ---+--- Register ---+---
	// Carga el plugin definido por el usuario,
	LoadPlugin(path string) (ErrorLoadPath, error)

	// Manera de obtener una estructura general de plugin definido por el usuario
	RegisterStructures() (ReturnRegisterStructure, error)

	// ---+--- Importing ---+---
	// Inicializar el proceso de importacion de archivos
	InitializeImport(uploader UploadEntity) error

	// Recibe la informacion de un path a un archivo importado
	ProcessFile(file string) error

	// Define el fin del proceso de importar archivos
	FinishImporing() error

	// ---+--- View Management ---+---
	// La view initial esta llena con la información default esperada de no tener
	// datos incluidos en esa view
	InitializeView(initialView string, world *v.World, data OWData) error

	// Avanza la escena al siguiente frame, pidiendo una nueva view si es necesario
	WalkScene(events []e.Event) error

	// Renderiza el mundo definido por el usuario
	RenderScene() (v.SceneRepresentation, error)
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
