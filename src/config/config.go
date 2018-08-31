package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type ConfigStruct struct {
	APIServer struct {
		ListenAddress string `json:"listenAddress"`
	} `json:"apiServer"`
	Mower struct {
		Name           string `json:"name"`
		CameraDeviceID int    `json:"cameraDeviceID"`
	} `json:mower`
}

var (
	ConfigFile string
	Config     *ConfigStruct
)

func LoadConfig(file string) (*ConfigStruct, error) {
	var cfg ConfigStruct

	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&cfg)

	Config = &cfg

	return &cfg, err
}
