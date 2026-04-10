package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

const ConfigDir = ".config/qry"
const ConfigFile = "connections.json"

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ConfigDir, ConfigFile), nil
}

type Connection struct {
	Driver string `json:"driver"`
	DSN    string `json:"dsn"`
}

var dsnMaskRegexp = regexp.MustCompile(`(?i)(://|)([^:/@]+):([^/@]+)@`)

func (c Connection) MaskedDSN() string {
	return MaskDSN(c.Driver, c.DSN)
}

func MaskDSN(driver, dsn string) string {
	return dsnMaskRegexp.ReplaceAllString(dsn, "$1$2:******@")
}

type Config struct {
	Connections map[string]Connection `json:"connections"`
}

func LoadConfig() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Connections: make(map[string]Connection)}, nil
		}
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.Connections == nil {
		config.Connections = make(map[string]Connection)
	}

	return &config, nil
}

func SaveConfig(config *Config) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func GetConnection(name string) (*Connection, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	conn, exists := cfg.Connections[name]
	if !exists {
		return nil, fmt.Errorf("la conexión `%s` no existe", name)
	}

	return &conn, nil
}

func AddConnection(name, driver, dsn string) error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	cfg.Connections[name] = Connection{
		Driver: driver,
		DSN:    dsn,
	}

	return SaveConfig(cfg)
}

func RemoveConnection(name string) error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	if _, exists := cfg.Connections[name]; !exists {
		return fmt.Errorf("connection %q not found", name)
	}

	delete(cfg.Connections, name)
	return SaveConfig(cfg)
}

func ListConnections() (map[string]Connection, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	return cfg.Connections, nil
}
