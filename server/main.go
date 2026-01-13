package main

import (
	"flag"
	"fmt"
	"os"

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

	/*  1. Initialize the app view
	 	2. Load all the program information, as in that projects are known
		3. Select a known project o create a new one
		4. Create client process (if possible create a hot reloadable process)
		5. Build userPlugin as plugin
		6. Create UserInteractionClient
			1. LoadPlugin(path) Components
			2. Set database base in component data
			3. Create SystemInteractionServer
			if user selects import files, and selects a path:
				4. Import files
			5. Set up user plugin render loop
	*/
}
