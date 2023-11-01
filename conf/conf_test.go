package conf

import (
	"log"
	"testing"
	"time"
)

type databaseConfig struct {
	Jacket   string        `json:"jacket"`
	Trousers string        `json:"trousers"`
	Times    time.Duration `json:"times"`
}

func TestSubViber(t *testing.T) {

	subv := cfg.Sub("clothing")
	if subv == nil {
		log.Fatalf("No 'database' key in config")
	}

	var data databaseConfig

	if err := subv.Unmarshal(&data); err != nil {
		log.Fatalf("Unable to decode into struct, %s", err.Error())
	} else {
		log.Printf("config %v ", data)
	}
}

func TestSub(t *testing.T) {
	data, err := GetSubCfg[databaseConfig]("clothing")
	if err != nil {
		log.Fatalf("get sub cfg  error, %s", err.Error())
	} else {
		log.Printf("config %v ", data)
	}
}

func TestGetAnyField(t *testing.T) {
	viper := GetViper()
	log.Printf("clothing.jacket : %v ", viper.GetString("clothing.jacket"))
	log.Printf("clothing.times : %v ", viper.GetDuration("clothing.times"))
}
