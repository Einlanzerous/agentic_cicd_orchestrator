package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Profiles map[string]Profile `mapstructure:"profiles"`
}

type Profile struct {
	Planner  ProviderConfig `mapstructure:"planner"`
	Coder    ProviderConfig `mapstructure:"coder"`
	Executor ProviderConfig `mapstructure:"executor"`
}

type ProviderConfig struct {
	Type   string            `mapstructure:"type"`
	Model  string            `mapstructure:"model"`
	Params map[string]string `mapstructure:"params"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// Set defaults
	viper.SetDefault("profiles", map[string]Profile{})

	// Enable env var override (e.g. LOCALSPRITE_PROFILES_WORK_CODER_MODEL)
	viper.SetEnvPrefix("localsprite")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
