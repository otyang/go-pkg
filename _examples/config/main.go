package main

import (
	"fmt"

	"github.com/otyang/go-starter/pkg/config"
)

func main() {
	type ConfigSampleStruct struct {
		Name string `env:"APP_NAME" json:"name" toml:"name" env-default:"Auth"`
	}

	cfgStruct := &ConfigSampleStruct{}      // config struct (must be a pointer)
	pathToConfigFile := "file/location.env" // file could be .env or .json or .toml or .yaml

	/* MustLoad */
	config.MustLoad(pathToConfigFile, cfgStruct) // loads the config and panics on error
	fmt.Println(cfgStruct.Name)

	/* Load */
	err := config.Load(pathToConfigFile, cfgStruct) // loads the config and returns error
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfgStruct.Name)
}
