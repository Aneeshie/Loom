package config

import (
	"os"
	"path/filepath"

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

func resolvePath(path string) string {
	// 1. If path exists directly, return it.
	if _, err := os.Stat(path); err == nil {
		return path
	}

	filename := filepath.Base(path)

	// 2. Try looking in ./configs/filename
	cfgPath := filepath.Join("configs", filename)
	if _, err := os.Stat(cfgPath); err == nil {
		return cfgPath
	}

	// 3. Try looking in the current directory ./filename
	if _, err := os.Stat(filename); err == nil {
		return filename
	}

	// 4. Try resolving relative to the executable's directory
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		// Check exeDir + "/configs/" + filename
		checkPath := filepath.Join(exeDir, "configs", filename)
		if _, err := os.Stat(checkPath); err == nil {
			return checkPath
		}
		// Check exeDir + "/../configs/" + filename
		checkPath = filepath.Join(exeDir, "..", "configs", filename)
		if _, err := os.Stat(checkPath); err == nil {
			return checkPath
		}
	}

	return path
}

func LoadAgentConfig(path string) (*AgentConfig, error) {
	resolved := resolvePath(path)
	data, err := os.ReadFile(resolved)
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
	resolved := resolvePath(path)
	data, err := os.ReadFile(resolved)
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
