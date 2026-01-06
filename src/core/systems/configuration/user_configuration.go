package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

var UserConfig UserConfiguration

type UserConfiguration struct {
	TargetFrameRate uint64 `json:"TargetFPS"`
}

func LoadUserConfiguration(userConfigurationPath string) error {
	if userConfigBytes, err := os.ReadFile(userConfigurationPath); err != nil {
		return fmt.Errorf("failed to read User define configuration, with error: %v\n", err)

	} else if err = json.Unmarshal(userConfigBytes, &UserConfig); err != nil {
		return err
	}

	return nil
}
