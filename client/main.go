package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/JuanBiancuzzo/own_wiki/core/api"

	c "github.com/JuanBiancuzzo/own_wiki/core/systems/configuration"
	log "github.com/JuanBiancuzzo/own_wiki/core/systems/logger"
)

func main() {
	// Command arguments parsing
	pathConfig := flag.String("config_path", "", "This is the path to the config file")
	flag.Parse()

	configuration, err := c.NewConfiguration(*pathConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load the configuration at %q, with error: %v", *pathConfig, err)
		return
	}

	if err := log.CreateLogger(configuration.Logger); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create the logger with (%s), with error: %v", configuration.Logger.String(), err)
		return
	}
	defer log.Close()

	log.Debug("Creating UserInteraction serve, with config: %#v", configuration.UserInteraction)
	userServer, err := api.NewUserInteractionServer(configuration.UserInteraction)
	if err != nil {
		log.Error("Failed to create UserInteraction server, with error: %v", err)
		return
	}

	log.Info("Serving UserInteraction at %s:%d", configuration.UserInteraction.Ip, configuration.UserInteraction.Port)
	if err = userServer.Serve(); err != nil {
		log.Error("Failed to serve UserInteraction server, with error: %v", err)
		return
	}
}
