// Package conf is a wrapper for viper.
package conf

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

var cfg *viper.Viper

func init() {
	cfg = viper.New()
	cfg.SetConfigName("config") // name of config file (without extension)
	cfg.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	cfg.AddConfigPath("./etc")  // optionally look for config in the working directory
	err := cfg.ReadInConfig()   // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Printf("fatal error config file: %s\n", err.Error())
		cfg = nil

		return
	}
}

// New provide translate the parsed viper object according to the given file.
func New(path string) (*viper.Viper, error) {
	cfg = viper.New()
	cfg.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	cfg.SetConfigFile(path)
	err := cfg.ReadInConfig() // Find and read the config file
	if err != nil {           // Handle errors reading the config file
		fmt.Printf("fatal error config file: %s\n", err.Error())
		cfg = nil

		return nil, err
	}

	return cfg, nil
}

// GetViper return the current viper object.
func GetViper() *viper.Viper {
	return cfg
}

// GetSubViper return the scope viper object
func GetSubViper(scope string) *viper.Viper {
	if cfg == nil {
		return nil
	}

	return cfg.Sub(scope)
}

// GetAllCfg return all config struct
func GetAllCfg[T any]() (*T, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	var entity T
	if err := cfg.Unmarshal(&entity); err != nil {
		fmt.Printf("Unable to decode into struct, %s\n", err.Error())

		return nil, err
	}

	return &entity, nil
}

// GetSubCfg return scope config struct
func GetSubCfg[T any](scope string) (*T, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	subv := cfg.Sub(scope)
	if subv == nil {
		fmt.Println("no '" + scope + "' key in config")

		return nil, errors.New("sub config is nil")
	}

	var entity T
	if err := subv.Unmarshal(&entity); err != nil {
		fmt.Printf("Unable to decode into struct, %s\n", err.Error())

		return nil, err
	}

	return &entity, nil
}
