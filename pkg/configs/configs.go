package configs

import (
	"encoding/json"
	"fmt"
	"os"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ServerConfig struct {
	Port int `json:"port"`
}

type Config struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to open configuration file: %v", err)
	}
	defer file.Close()

	cfg := &Config{}
	err = json.NewDecoder(file).Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse configuration file: %v", err)
	}

	return cfg, nil
}
