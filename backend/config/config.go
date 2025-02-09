package config

import (
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"go.uber.org/zap"
)

type Config struct {
	Database DatabaseConnection
	Logger   zap.Config
	Server   ServerParams
}

type ServerParams struct {
	Hostname           string `koanf:"host"`
	Port               int    `koanf:"port"`
	IdleTimeoutSec     int    `koanf:"idle-timeout"`
	ShutdownTimeoutSec int    `koanf:"shutdown-timeout"`
}

type DatabaseConnection struct {
	Hostname   string    `koanf:"host"`
	Port       int       `koanf:"port"`
	Database   string    `koanf:"name"`
	User       BasicAuth `koanf:"user"`
	TLSEnabled bool      `koanf:"tls-enabled"`
}

type BasicAuth struct {
	Login    string `koanf:"login"`
	Password string `koanf:"password"`
}

func Load() (Config, error) {
	var (
		zero Config
		k    = koanf.New(".")
	)

	// Load main config
	if err := k.Load(file.Provider("application.yml"), yaml.Parser()); err != nil {
		return zero, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Merge with local config
	_ = k.Load(file.Provider("local.yml"), yaml.Parser())

	var (
		loggerConfig zap.Config
		dbConfig     DatabaseConnection
		serverConfig ServerParams
	)
	if err := k.Unmarshal("application.db", &dbConfig); err != nil {
		return zero, fmt.Errorf("failed to unmarshal db configuration: %w", err)
	}
	if err := k.Unmarshal("application.server", &serverConfig); err != nil {
		return zero, fmt.Errorf("failed to unmarshal server configuration: %w", err)
	}
	if err := k.UnmarshalWithConf("application.logger", &loggerConfig, koanf.UnmarshalConf{Tag: "yaml"}); err != nil {
		return zero, fmt.Errorf("failed to unmarshal logger configuration: %w", err)
	}

	return Config{
		Database: dbConfig,
		Logger:   loggerConfig,
		Server:   serverConfig,
	}, nil
}
