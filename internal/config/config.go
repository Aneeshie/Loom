package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AgentConfig struct {
	Server struct {
		Addr string `yaml:"addr"`
	} `yaml:"server"`

	Service struct {
		Name string `yaml:"name"`
		Host string `yaml:"host"`
	} `yaml:"service"`

	Source struct {
		Path string `yaml:"path"`
	} `yaml:"source"`

	Query struct {
		Limit int64 `yaml:"limit"`
	} `yaml:"query"`
}

type ServerConfig struct {
	GRPC struct {
		Addr string `yaml:"addr"`
	} `yaml:"grpc"`

	Database struct {
		URL string `yaml:"url"`
	} `yaml:"database"`
}

func LoadAgentConfig(path string) (*AgentConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg AgentConfig

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func LoadServerConfig(path string) (*ServerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ServerConfig

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
