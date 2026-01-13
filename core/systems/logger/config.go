package logger

import (
	"encoding/json"
	"fmt"
)

type Verbosity string

const (
	MUTE    = "mute"
	NORMAL  = "normal"
	VERBOSE = "verbose"
)

const DEFAULT_CAPACITY = 25

type LoggerConfiguration struct {
	LogPath         string    `json:"log_path,omitempty"`
	Verbosity       Verbosity `json:"verbosity"`
	MessageCapacity uint      `json:"message_capacity,omitempty"`
}

func (lc *LoggerConfiguration) UnmarshalJSON(data []byte) error {
	// We unmarshal it as default
	if err := json.Unmarshal(data, lc); err != nil {
		return err
	}

	// Now we check if the capacity is valid
	if lc.MessageCapacity == 0 {
		lc.MessageCapacity = DEFAULT_CAPACITY
	}

	return nil
}

func (lc LoggerConfiguration) String() string {
	return fmt.Sprintf("Logger path: %q, Verbosity: %s and Message capacity of: %d", lc.LogPath, lc.Verbosity, lc.MessageCapacity)
}
