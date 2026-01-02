package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	OrgFiles    []string `mapstructure:"org_files"`
	DefaultFile string   `mapstructure:"default_file"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
