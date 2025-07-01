package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Symbols     []string `yaml:"symbols"`
	Server      Server   `yaml:"server"`
	Database    Database `yaml:"database"`
	OHLCPeriods []string `yaml:"ohlc_periods"`
}

type Server struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type Database struct {
	Path string `yaml:"path"`
}

// PeriodInfo represents a parsed OHLC period
type PeriodInfo struct {
	Value    int    // Value in seconds
	Original string // Original string (e.g., "10s", "1m")
	Unit     string // Unit (s for seconds, m for minutes)
}

// ParsePeriod parses a period string like "10s" or "1m" and returns seconds
func ParsePeriod(periodStr string) (int, error) {
	if len(periodStr) < 2 {
		return 0, fmt.Errorf("invalid period format: %s", periodStr)
	}

	unit := periodStr[len(periodStr)-1:]
	valueStr := periodStr[:len(periodStr)-1]

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid period value: %s", periodStr)
	}

	switch unit {
	case "s":
		return value, nil
	case "m":
		return value * 60, nil
	default:
		return 0, fmt.Errorf("unsupported period unit: %s", unit)
	}
}

// GetPeriodInfo parses a period string and returns PeriodInfo
func GetPeriodInfo(periodStr string) (*PeriodInfo, error) {
	value, err := ParsePeriod(periodStr)
	if err != nil {
		return nil, err
	}

	unit := periodStr[len(periodStr)-1:]
	return &PeriodInfo{
		Value:    value,
		Original: periodStr,
		Unit:     unit,
	}, nil
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
