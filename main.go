package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/labstack/echo"
	"github.com/techievee/opensimsim/config"
	"github.com/techievee/opensimsim/server"
	shift "github.com/techievee/opensimsim/shifts"
	"github.com/techievee/opensimsim/storage"
	"os"
	"path/filepath"
)

func main() {

	//Configurations are read from the Config file,
	// config.json or the file passed as argument during the runtime
	var cfg config.Config
	err := readConfig(&cfg)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	//For Modularity, the codes are segregated according to logic,
	//Backend modules are embedded into all the other modules which requires database functionality
	//Server modules are embedded to all the modules that need to be exposed as API Endpoints
	//Shift modules contains all the shifts functionality with backend and server injected into it
	backend := storage.NewStorageModule(&cfg.MySql)
	server := server.NewServerModule(&cfg.Server)
	shiftmodules := shift.NewShiftModule(server, backend)

	defer backend.Close()

	//Setup routes for all the modules
	shiftmodules.SetupRoutes()

	//Start the Server
	server.StartServer()

}

func readConfig(cfg *config.Config) error {
	configFileName := "config.json"
	if len(os.Args) > 1 {
		configFileName = os.Args[1]
	}
	configFileName, _ = filepath.Abs(configFileName)
	fmt.Printf("Loading config: %v", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		return err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&cfg); err != nil {
		return err
	}
	return nil
}
