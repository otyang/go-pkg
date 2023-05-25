package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// type Example_ConfigStruct struct {
// 	AppName string `env:"APP_NAME" env-default:"Auth"`

// 	Redis struct {
// 		URL       string `env:"REDIS_URI" env-default:".." `
// 		EnableTLS bool   `env:"REDIS_ENABLE_TLS" env-default:"true"`
// 	}
// }

// MustLoad: takes a config file path, parses it to a struct of any type,
// then overwrite the struct with similar settings from environment.
//
// Take note: if configFile is empty, it skips to env. It exit the program
// if it encounters an error parsing the config.
//
// Supported Config file types are: YAML, JSON, TOML, .ENV
//
// Unlike most config loaders, here anytype can be parsed on the fly. It even gets updated env
func MustLoad(configFile string, configStructPtr any) {
	err := Load(configFile, configStructPtr)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

// Load unlike MustLoad, does return errors and not exists the system
func Load(configFile string, configStructPtr any) error {
	var err error

	if configFile != "" {
		err = cleanenv.ReadConfig(configFile, configStructPtr)
		if err != nil {
			return fmt.Errorf("unable read config from file: '%s' | %s", configFile, err.Error())
		}
	}

	// env overwriting config from file
	err = cleanenv.ReadEnv(configStructPtr)
	if err != nil {
		return fmt.Errorf("unable read config from env: %s", err.Error())
	}

	// updating env variables change via runtime
	err = cleanenv.UpdateEnv(configStructPtr)
	if err != nil {
		return fmt.Errorf("unable to update config from env: %s", err.Error())
	}

	return nil
}
