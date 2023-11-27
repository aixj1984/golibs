package conf

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
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
	v, err := GetAllCfg[TestConfig]()
	if err != nil {
		t.Logf("new failed: %v", err)
	}
	assert.Nil(t, v)
}

func TestGetSubCfgNil(t *testing.T) {
	v, err := GetSubCfg[TestConfig]("abc")
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
	type LogConfig struct {
		LogPath    string `json:"logPath"`
		AppName    string `json:"appName"`
		Debug      bool   `json:"debug"`
		Level      int8   `json:"level"`
		MaxSize    int
		MaxAge     int  `json:"maxAge"`
		MaxBackups int  `json:"maxBackups"`
		Compress   bool `json:"compress"`
	}

	type testConfigA struct {
		Log  LogConfig  `mapstructure:"log"`
		Test TestConfig `mapstructure:"clothing"`
	}

	var config *testConfigA
	config, err := GetAllCfg[testConfigA]()
	if err != nil {
		t.Fatalf("GetAllCfg failed: %v", err)
	}
	fmt.Printf("config %+v", config)
	assert.NotNil(t, config)
}

func TestGetSubCfg(t *testing.T) {
	var config *TestConfig
	config, err := GetSubCfg[TestConfig]("invalid_scope")
	if err == nil {
		t.Fatalf("Expected error for invalid scope, but got none")
	}
	assert.Nil(t, config)

	config, err = GetSubCfg[TestConfig]("clothing")
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

	var data TestConfig

	if err := subv.Unmarshal(&data); err != nil {
		log.Fatalf("Unable to decode into struct, %s", err.Error())
	} else {
		log.Printf("config %v ", data)
	}
}

func TestSub(t *testing.T) {
	data, err := GetSubCfg[TestConfig]("clothing")
	if err != nil {
		log.Fatalf("get sub cfg  error, %s", err.Error())
	} else {
		log.Printf("config %v ", data)
	}
}

func TestAllFieldErr(t *testing.T) {
	type testField struct {
		A string `json:"a"`
		B int    `json:"b"`
	}
	type allField struct {
		ErrField testField `json:"errField"`
	}
	data, err := GetAllCfg[allField]()
	if err != nil {
		log.Printf("get all cfg  error, %s", err.Error())
	} else {
		log.Printf("config %v ", data)
	}
	assert.Nil(t, data)
}

func TestSubFieldErr(t *testing.T) {
	type testField struct {
		A string `json:"a"`
		B int    `json:"b"`
	}
	data, err := GetSubCfg[testField]("errField")
	if err != nil {
		log.Printf("get sub cfg  error, %s", err.Error())
	} else {
		log.Printf("config %v ", data)
	}
	assert.Nil(t, data)
}

func TestGetAnyField(t *testing.T) {
	viper := GetViper()
	log.Printf("clothing.jacket : %v ", viper.GetString("clothing.jacket"))
	log.Printf("clothing.times : %v ", viper.GetDuration("clothing.times"))
}
