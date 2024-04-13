package configuration

import (
	"encoding/json"
	"os"
	"tiger-tracker-api/configuration/models"
)

type ConfigLoader interface {
	LoadConfig(fileName string) (*ConfigData, error)
}

type configLoader struct {
}

func NewConfigLoader() ConfigLoader {
	return configLoader{}
}

func (c configLoader) LoadConfig(fileName string) (*ConfigData, error) {
	configuration := ConfigData{}
	file, err := os.Open(fileName)
	if err != nil {
		return &configuration, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	return &configuration, err
}

type ConfigData struct {
	Environment      string                  `json:"environment"`
	ListenPort       string                  `json:"listenPort"`
	DbConnectionPool models.DbConnectionPool `json:"dbConnectionPool"`
}
