package config

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

type Config struct {
	ConsumerServer ServerConnection
	Pinger         PingerParams
	Docker         ServerConnection
	Logger         zap.Config
}

type PingerParams struct {
	PingPeriodSec     int `koanf:"ping-period"`
	RequestTimeoutSec int `koanf:"request-timeout"`
	Workers           int `koanf:"worker-pool"`
}

type ServerConnection struct {
	Host       string `koanf:"host"`
	MaxRetries int    `koanf:"max-retries"`
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
		loggerConf   zap.Config
		pingerConf   PingerParams
		dockerConf   ServerConnection
		consumerConf ServerConnection
	)
	if err := k.Unmarshal("application.docker", &dockerConf); err != nil {
		return zero, fmt.Errorf("failed to unmarshal db configuration: %w", err)
	}
	if err := k.Unmarshal("application.pinger", &pingerConf); err != nil {
		return zero, fmt.Errorf("failed to unmarshal db configuration: %w", err)
	}
	if err := k.Unmarshal("application.consumer-server", &consumerConf); err != nil {
		return zero, fmt.Errorf("failed to unmarshal server configuration: %w", err)
	}
	if err := k.UnmarshalWithConf("application.logger", &loggerConf, koanf.UnmarshalConf{Tag: "yaml"}); err != nil {
		return zero, fmt.Errorf("failed to unmarshal logger configuration: %w", err)
	}

	return Config{
		ConsumerServer: consumerConf,
		Pinger:         pingerConf,
		Logger:         loggerConf,
		Docker:         dockerConf,
	}, nil
}
