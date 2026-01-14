package configuration

import (
	"encoding/json"
	"reflect"
)

const DEFAULT_FILE_PATH_CAPACITY uint16 = 50
const DEFAULT_COMPONENTS_CAPACITY uint16 = 25
const DEFAULT_FILE_WORKERS uint8 = 10

type UserInteractionConfig struct {
	Protocol string `json:"protocol"`
	Ip       string `json:"ip"`
	Port     uint8  `json:"port"`

	FilePathCapacity  uint16 `json:"file_path_capacity,omitempty"`
	ComponentCapacity uint16 `json:"component_capacity,omitempty"`
	AmountFileWorkers uint8  `json:"file_workers,omitempty"`
}

func (uc *UserInteractionConfig) UnmarshalJSON(data []byte) error {
	// We unmarshal it as default
	if err := json.Unmarshal(data, uc); err != nil {
		return err
	}

	checkDefault(&uc.FilePathCapacity, DEFAULT_FILE_PATH_CAPACITY)
	checkDefault(&uc.ComponentCapacity, DEFAULT_COMPONENTS_CAPACITY)
	checkDefault(&uc.AmountFileWorkers, DEFAULT_FILE_WORKERS)

	return nil
}

type SystemInteractionConfig struct {
	Protocol string `json:"protocol"`
	Ip       string `json:"ip"`
	Port     uint8  `json:"port"`
}

func checkDefault[T any](element *T, defaultValue T) {
	var empty T
	if reflect.DeepEqual(*element, empty) {
		*element = defaultValue
	}
}
