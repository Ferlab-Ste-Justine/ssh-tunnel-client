package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type TunnelConfigBinding struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
}

type TunnelConfig struct {
	HostMd5FingerPrint string                `json:"host_md5_fingerprint"`
	HostUrl            string                `json:"host_url"`
	HostUser           string                `json:"host_user"`
	AuthMethod         string                `json:"auth_method"`
	Bindings           []TunnelConfigBinding `json:"bindings"`
}

func getTunnelConfig(fallbackVal string) (TunnelConfig, error) {
	_, err := os.Stat("tunnel_config.json")
	configFileExists := err == nil
    
	if !configFileExists && fallbackVal == "" {
		return TunnelConfig{}, errors.New("tunnel_config.json is not present in running directory nor embedded in binary")
	}

	var config TunnelConfig

	var buffer []byte

	if configFileExists {
		var err error
		buffer, err = ioutil.ReadFile("tunnel_config.json")
		if err != nil {
			return config, err
		}
	} else {
		buffer = []byte(fallbackVal)
	}


	if err := json.Unmarshal(buffer, &config); err != nil {
		return config, err
	}

	return config, nil

}