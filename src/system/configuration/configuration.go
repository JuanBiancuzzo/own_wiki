package configuration

import (
	"encoding/json"
	"os"
)

type UserConfig struct {
	TargetFrameRate uint64 `json:"TargetFPS"`
}

type SystemConfig struct {
	Platform string `json:"Platform"`
}

func NewUserConfig(path string) (userConfig *UserConfig, err error) {
	if bytes, err := os.ReadFile(path); err != nil {
		return nil, err

	} else if err = json.Unmarshal(bytes, userConfig); err != nil {
		return nil, err

	} else {
		return userConfig, nil
	}
}

func NewSystemConfig(path string) (systemConfig *SystemConfig, err error) {
	if bytes, err := os.ReadFile(path); err != nil {
		return nil, err

	} else if err = json.Unmarshal(bytes, systemConfig); err != nil {
		return nil, err

	} else {
		return systemConfig, nil
	}
}
