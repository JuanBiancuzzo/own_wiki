package configuration

import (
	"encoding/json"
)

type UserConfig struct {
	TargetFrameRate         uint64 `json:"TargetFPS"`
	UserDefineDataDirectory string `json:"pluginDir"`
}

type SystemConfig struct {
	Platform string `json:"Platform"`
}

func NewUserConfig(bytes []byte) (userConfig UserConfig, err error) {
	if err = json.Unmarshal(bytes, &userConfig); err != nil {
		return userConfig, err

	} else {
		return userConfig, nil
	}
}

func NewSystemConfig(bytes []byte) (systemConfig SystemConfig, err error) {
	if err = json.Unmarshal(bytes, &systemConfig); err != nil {
		return systemConfig, err

	} else {
		return systemConfig, nil
	}
}
