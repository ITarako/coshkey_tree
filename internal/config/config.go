package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var cfg *Config

// GetConfigInstance returns service config
func GetConfigInstance() Config {
	if cfg != nil {
		return *cfg
	}

	return Config{}
}

// Database - contains all parameters database connection.
type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SslMode  string `yaml:"sslmode"`
}

// Rest - contains parameter rest json connection.
type Rest struct {
	Port         int    `yaml:"port"`
	Host         string `yaml:"host"`
	WriteTimeout int    `yaml:"writeTimeout"`
	ReadTimeout  int    `yaml:"readTimeout"`
}

// Project - contains all parameters project information.
type Project struct {
	Debug        bool   `yaml:"debug"`
	Name         string `yaml:"name"`
	Environment  string `yaml:"environment"`
	CoshkeyUrl   string `yaml:"coshkeyUrl"`
	CoshkeyToken string `yaml:"coshkeyToken"`
}

// Config - contains all configuration parameters in config package.
type Config struct {
	Project  Project  `yaml:"project"`
	Rest     Rest     `yaml:"rest"`
	Database Database `yaml:"database"`
}

// ReadConfigYML - read configurations from file and init instance Config.
func ReadConfigYML(filePath string) error {
	if cfg != nil {
		return nil
	}

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		return err
	}

	return nil
}
