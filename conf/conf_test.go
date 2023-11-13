package conf

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testConfig struct {
	Jacket   string        `json:"jacket"`
	Trousers string        `json:"trousers"`
	Times    time.Duration `json:"times"`
}

func TestGetViperNil(t *testing.T) {
	v := GetViper()
	assert.Nil(t, v)
}

func TestGetSubViperNil(t *testing.T) {
	v := GetSubViper("abc")
	assert.Nil(t, v)
}

func TestGetAllCfgNil(t *testing.T) {
	v, err := GetAllCfg[testConfig]()
	if err != nil {
		t.Logf("new failed: %v", err)
	}
	assert.Nil(t, v)
}

func TestGetSubCfgNil(t *testing.T) {
	v, err := GetSubCfg[testConfig]("abc")
	if err != nil {
		t.Logf("new failed: %v", err)
	}
	assert.Nil(t, v)
}

func TestNewSucc(t *testing.T) {
	v, err := New("./etc/config.yaml")
	if err != nil {
		t.Logf("new failed: %v", err)
	}
	assert.Nil(t, v)
}

func TestNewError(t *testing.T) {
	v, err := New("./etc/abc.yaml")
	if err != nil {
		t.Logf("new failed: %v", err)
	}
	assert.NotNil(t, v)
}

func TestGetViperSucc(t *testing.T) {
	v := GetViper()
	assert.NotNil(t, v)
}

func TestGetSubViper(t *testing.T) {
	v := GetSubViper("test")
	assert.Nil(t, v)
}

func TestGetAllCfg(t *testing.T) {
	var config *testConfig
	config, err := GetAllCfg[testConfig]()
	if err != nil {
		t.Fatalf("GetAllCfg failed: %v", err)
	}
	assert.NotNil(t, config)
}

func TestGetSubCfg(t *testing.T) {
	var config *testConfig
	config, err := GetSubCfg[testConfig]("invalid_scope")
	if err == nil {
		t.Fatalf("Expected error for invalid scope, but got none")
	}
	assert.Nil(t, config)

	config, err = GetSubCfg[testConfig]("clothing")
	if err != nil {
		t.Fatalf("GetSubCfg failed: %v", err)
	}
	assert.NotNil(t, config)
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
