package user

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/JuanBiancuzzo/own_wiki/src/core/api"
	"github.com/JuanBiancuzzo/own_wiki/src/core/ecv"
	plugin "github.com/hashicorp/go-plugin"
)

type UserPlugin struct {
	client *plugin.Client
	Plugin api.UserStructureData
}

func GetUserDefineData(pluginPath string) (*UserPlugin, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: api.Handshake,
		Plugins:         api.PluginMap,
		// Se puede usar air para que cada vez que se cambia el código, se
		//   buildea de nuevo
		// Mover a configuracion
		Cmd: exec.Command("plugin/plugin.so"),
	})

	if rpcClient, err := client.Client(); err != nil { // Connect via RPC
		client.Kill()
		return nil, err

	} else if raw, err := rpcClient.Dispense("UserDefineStructure"); err != nil { // Request the plugin
		client.Kill()
		return nil, err

	} else if plugin, ok := raw.(api.UserStructureData); !ok {
		client.Kill()
		return nil, errors.New("Plugin doesnt satisfy UserStructureData interface")

	} else if info, err := plugin.LoadPlugin(pluginPath); err != nil {
		client.Kill()
		return nil, fmt.Errorf("Faield to connect to the plugin to initialize it, with error: %v", err)

	} else if info.HasError {
		client.Kill()
		return nil, errors.New(info.ErrorReason)

	} else {
		return &UserPlugin{
			client: client,
			Plugin: plugin,
		}, nil
	}
}

func (up *UserPlugin) RegisterStructures() (*ecv.ECV, error) {
	if info, err := up.Plugin.RegisterStructures(); err != nil {
		return nil, fmt.Errorf("Error in the connection with the plugin, with error: %v", err)

	} else if info.HasError {
		return nil, errors.New(info.ErrorReason)

	} else {
		return info.Ecv, nil
	}
}

// Implementar importar archivos al programa

// Implementar las views

func (up *UserPlugin) Close() {
	up.client.Kill()
}
