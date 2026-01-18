package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	OrgFiles    []string      `mapstructure:"org_files"`
	DefaultFile string        `mapstructure:"default_file"`
	Capture     CaptureConfig `mapstructure:"capture"`
}

type CaptureConfig struct {
	DefaultFile string `mapstructure:"default_file"`
	Format      string `mapstructure:"format"`
	Prepend     bool   `mapstructure:"prepend"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
