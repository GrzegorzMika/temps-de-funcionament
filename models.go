package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Services []Services `yaml:"services"`
}
type Services struct {
	Name          string `yaml:"name"`
	URL           string `yaml:"url"`
	AllowedCodes  []int  `yaml:"allowed_codes"`
	PeriodSeconds int    `yaml:"periodSeconds"`
	MaxFailures   int    `yaml:"maxFailures"`
	SlackChannel  string `yaml:"slackChannel"`
}

func GetConfiguration(path string) (*Configuration, error) {
	c := &Configuration{}
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling yaml: %v", err)
	}

	return c, nil
}
