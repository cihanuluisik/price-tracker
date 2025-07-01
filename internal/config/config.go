package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	WebSocket   WebSocketConfig `yaml:"websocket"`
	Server      ServerConfig    `yaml:"server"`
	Database    DatabaseConfig  `yaml:"database"`
	OHLCPeriods []int           `yaml:"ohlc_periods"`
}

type WebSocketConfig struct {
	URL     string   `yaml:"url"`
	Symbols []string `yaml:"symbols"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
