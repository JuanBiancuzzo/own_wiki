package user

import (
	"os/exec"
	"path/filepath"
	"strings"

	v "github.com/JuanBiancuzzo/own_wiki/view"

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

/*
Aca deberiamos agarrar esa informacion y traducirla, generando una estructura que
represente todo lo que puede hacer el usuario
*/
func (up *UserPlugin) RegisterComponents() ([]*shared.ComponentInformation, error) {
	return up.plugin.RegisterComponents()
}

func (up *UserPlugin) RegisterEntities() ([]*shared.EntityInformation, error) {
	return up.plugin.RegisterEntities()
}

func (up *UserPlugin) RegisterViews() ([]*shared.ViewInformation, error) {
	return up.plugin.RegisterViews()
}

func (up *UserPlugin) CreateView(viewName string, scene *v.Scene) error {
	return up.plugin.CreateView(shared.SceneInformation{
		ViewName:   viewName,
		EntityName: "",
		Scene:      scene,
	})
}

func (up *UserPlugin) AvanzarView(events []v.Event) (*v.SceneOperation, error) {
	return up.plugin.AvanzarView(events)
}

func (up *UserPlugin) Close() {
	up.client.Kill()
}
