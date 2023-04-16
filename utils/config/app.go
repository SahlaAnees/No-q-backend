package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type App struct {
	Port int    `yaml:"service-port"`
	Host string `yaml:"service-host"`
}

func (app *App) Parse() error {

	yamlFile, err := os.ReadFile("configurations/app.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, app)
	if err != nil {
		return err
	}

	return nil
}
