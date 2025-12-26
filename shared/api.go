package shared

import (
	"net/rpc"

	// _ "github.com/JuanBiancuzzo/own_wiki/view"
	"github.com/hashicorp/go-plugin"
)

// Definimos un handshake para poder controlar versiones del plugin
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "USER_PLUGIN",
	MagicCookieValue: "enabled",
}

// Definimos un mapeo de plugins, para poder hacer dispatch del plugin
var PluginMap = map[string]plugin.Plugin{
	"userDefineData": &UserDefineDataPlugin{},
}

/*
Con esto podemos definir 3 funciones, que fuerzan al usuario a establecer
todo lo que deberia hacer, es la API/contrato, que necesitan cumplir para que
el sistema funcione. Estas son
  - Funcion para ingresar los struct que corresponden como componentes
  - Funcion que recibe un archivo de texto, con toda la metadata, y el sistema
    para ingresar las entidades
  - Funcion para ingresar el par (entidad, []view)
*/

/*
Esta es la interfaz que tiene que cumplir el plugin.
  - RegisterComponents: Esta permite definir los bloques minimos para crear entidades
*/
type UserDefineData interface {
	RegisterComponents() ([]*ComponentInformation, error)

	RegisterEntities() ([]*EntityInformation, error)

	RegisterViews() ([]*ViewInformation, error)
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
type UserDefineDataPlugin struct {
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl UserDefineData
}

func (p *UserDefineDataPlugin) Server(*plugin.MuxBroker) (any, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (*UserDefineDataPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (any, error) {
	return &RPCClient{client: c}, nil
}

// Server
// Here is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc
type RPCServer struct {
	// This is the real implementation
	Impl UserDefineData
}

func (m *RPCServer) RegisterComponents(args any, resp *[]*ComponentInformation) (err error) {
	*resp, err = m.Impl.RegisterComponents()
	return err
}

func (m *RPCServer) RegisterEntities(args any, resp *[]*EntityInformation) (err error) {
	*resp, err = m.Impl.RegisterEntities()
	return err
}

// Client
// RPCClient is an implementation of Plugin that talks over RPC.
type RPCClient struct{ client *rpc.Client }

func (m *RPCClient) RegisterComponents() ([]*ComponentInformation, error) {
	var resp []*ComponentInformation
	err := m.client.Call("Plugin.RegisterComponents", new(any), &resp)
	return resp, err
}

func (m *RPCClient) RegisterEntities() ([]*EntityInformation, error) {
	var resp []*EntityInformation
	err := m.client.Call("Plugin.RegisterEntities", new(any), &resp)
	return resp, err
}
