package user

import (
	"os/exec"
	"reflect"

	"github.com/JuanBiancuzzo/own_wiki/shared"
	plugin "github.com/hashicorp/go-plugin"
)

type UserPlugin struct {
	client *plugin.Client
	plugin shared.UserDefineData
}

func GetUserDefineData(pluginPath string) (*UserPlugin, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("sh", "-c", pluginPath),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC,
		},
	})

	if rpcClient, err := client.Client(); err != nil { // Connect via RPC
		client.Kill()
		return nil, err

	} else if raw, err := rpcClient.Dispense("userDefineData"); err != nil { // Request the plugin
		client.Kill()
		return nil, err

	} else {
		return &UserPlugin{
			client: client,
			plugin: raw.(shared.UserDefineData),
		}, nil
	}
}

func (up *UserPlugin) RegisterComponents() ([]reflect.Type, error) {
	return up.plugin.RegisterComponents()
}

func (up *UserPlugin) Close() {
	up.client.Kill()
}
