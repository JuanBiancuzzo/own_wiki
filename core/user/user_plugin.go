package user

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/JuanBiancuzzo/own_wiki/shared"
	plugin "github.com/hashicorp/go-plugin"
)

type UserPlugin struct {
	client *plugin.Client
	plugin shared.UserDefineData
}

func GetUserDefineData(pluginPath string) (*UserPlugin, error) {
	pluginFilePath := filepath.Join(strings.Split(pluginPath, "/")...)

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command(pluginFilePath),
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

func (up *UserPlugin) RegisterComponents() ([]*shared.ComponentInformation, error) {
	return up.plugin.RegisterComponents()
}

func (up *UserPlugin) RegisterEntities() ([]*shared.EntityInformation, error) {
	return up.plugin.RegisterEntities()
}

func (up *UserPlugin) Close() {
	up.client.Kill()
}
