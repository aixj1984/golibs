package conf

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

var (
	cfg *viper.Viper = nil
)

func init() {
	cfg = viper.New()
	cfg.SetConfigName("config") // name of config file (without extension)
	cfg.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	cfg.AddConfigPath("./etc")  // optionally look for config in the working directory
	fmt.Println("init")
	err := cfg.ReadInConfig() // Find and read the config file
	if err != nil {           // Handle errors reading the config file
		fmt.Printf("fatal error config file: %s\n", err.Error())
	}
}

func GetViper() *viper.Viper {
	return cfg
}

func GetSubViper(scope string) *viper.Viper {
	if cfg == nil {
		return nil
	}
	return cfg.Sub(scope)
}

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
