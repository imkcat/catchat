package config

import (
	"os"

	"github.com/samber/lo"
	"github.com/spf13/viper"
)

// Model Provider
type Provider string

const (
	Azure  Provider = "Azure"
	OpenAI Provider = "OpenAI"
)

// Profile
type Profile struct {
	Provider Provider
	Name     string `survey:"name"`
	Azure    *AzureConfig
	OpenAI   *OpenAIConfig
}

// Config
type Config struct {
	Profiles []Profile `mapstructure:"profiles"`
}

// Create new config instance
func NewConfig(configPath, configFilePath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.SetDefault("profiles", []Profile{})
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err := os.MkdirAll(configPath, os.ModePerm)
			if err != nil {
				return nil, err
			}
			_, err = os.Create(configFilePath)
			if err != nil {
				return nil, err
			}
			err = viper.WriteConfig()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	config, err := ReloadConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

// Reload config
func ReloadConfig() (*Config, error) {
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Create new profile
func (config *Config) CreateProfile(profile Profile) error {
	viper.Set("profiles", append(config.Profiles, profile))
	err := viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

// Delete profile by id
func (config *Config) DeleteProfile(names []string) error {
	viper.Set("profiles", lo.Filter(config.Profiles, func(item Profile, index int) bool {
		return !lo.Contains(names, item.Name)
	}))
	err := viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}
