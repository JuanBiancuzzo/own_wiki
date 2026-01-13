package configuration

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/JuanBiancuzzo/own_wiki/core/systems/logger"
)

type Configuration struct {
	Logger log.LoggerConfiguration

	SystemInteraction SystemInteractionConfig `json:"system_interaction"`
	UserInteraction   UserInteractionConfig   `json:"user_interaction"`
}

func NewConfiguration(path string) (config Configuration, err error) {
	file, err := os.Open(path)
	if err != nil {
		return config, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return config, fmt.Errorf("failed to decode config file: %v", err)
	}

	return config, nil
}
