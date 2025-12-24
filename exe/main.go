package main

import (
	"fmt"
	"os"
	"sync"

	_ "embed"

	"github.com/JuanBiancuzzo/own_wiki/exe/platforms/htmx"
	"github.com/JuanBiancuzzo/own_wiki/exe/platforms/terminal"

	"github.com/JuanBiancuzzo/own_wiki/core"
	p "github.com/JuanBiancuzzo/own_wiki/core/platform"
	c "github.com/JuanBiancuzzo/own_wiki/core/system/configuration"
	log "github.com/JuanBiancuzzo/own_wiki/core/system/logger"
)

const USER_CONFIG_PATH string = "config/user_config.json"

//go:embed "config/system_config.json"
var systemConfigBytes []byte

func main() {
	if userConfigBytes, err := os.ReadFile(USER_CONFIG_PATH); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read User define configuration, with error: %v\n", err)

	} else if userConfig, err := c.NewUserConfig(userConfigBytes); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create User define configuration, with error: %v\n", err)

	} else if systemConfig, err := c.NewSystemConfig(systemConfigBytes); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create User define configuration, with error: %v\n", err)

	} else if err := log.CreateLogger("logs/logger.txt", log.VERBOSE); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)

	} else {
		defer log.Close()

		var platform p.Platform
		switch systemConfig.Platform {
		case "HTMX":
			platform = htmx.NewHTMX()

		case "Terminal":
			platform = terminal.NewTerminal()

		default:
			log.Error("Failed to asign platform, check configuration file, options are HTMX, Terminal")
			return
		}

		var waitGroup sync.WaitGroup

		core.Loop(userConfig, platform, &waitGroup)
		platform.Close()

		waitGroup.Wait()
	}
}
