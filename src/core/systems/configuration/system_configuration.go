package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

var SystemConfig SystemConfiguration

type SystemConfiguration struct {
	TargetFrameRate uint64 `json:"TargetFPS"`
}

func LoadSystemConfiguration(userConfigurationPath string) error {
	if userConfigBytes, err := os.ReadFile(userConfigurationPath); err != nil {
		return fmt.Errorf("failed to read System define configuration, with error: %v\n", err)

	} else if err = json.Unmarshal(userConfigBytes, &UserConfig); err != nil {
		return err
	}

	return nil
}
