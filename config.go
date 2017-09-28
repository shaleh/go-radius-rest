package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type ServerConfig struct {
	Database struct {
		User     string
		Password string
		Host     string
		Name     string
	}
	Host string
	Port int
}

func (config *ServerConfig) BindAddress() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

func loadConfiguration(file string) (config *ServerConfig, err error) {
	configFile, err := os.Open(file)
	if err != nil {
		return
	}
	defer configFile.Close()

	config = new(ServerConfig)
	err = json.NewDecoder(configFile).Decode(config)

	return
}
