package config

import (
	"fmt"
	"os"
	"strings"

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
	Provider Provider      `mapstructure:"provider" json:"provider"`
	Name     string        `survey:"name" mapstructure:"name" json:"name"`
	Azure    *AzureConfig  `mapstructure:"azure,omitempty" json:"azure"`
	OpenAI   *OpenAIConfig `mapstructure:"open_ai,omitempty" json:"open_ai"`
}

// Config
type Config struct {
	Profiles []Profile `mapstructure:"profiles" json:"profiles"`
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
	fmt.Printf("Profile %s has been created! Please re-run catchat command and select it to start chatting!\n", profile.Name)
	return nil
}

// Switch profile
func (config *Config) SwitchProfile(profile Profile) error {
	viper.Set("current_profile_name", profile.Name)
	err := viper.WriteConfig()
	if err != nil {
		return err
	}
	fmt.Printf("Switched to profile: %s", profile.Name)
	return nil
}

// Check profile
func CheckProfile(profile Profile) {
	fmt.Printf("%s\n", ProfileInformation(profile))
}

func ProfileInformation(profile Profile) string {
	informations := make([]string, 0)
	informations = append(informations, fmt.Sprintf("Name: %s", profile.Name))
	informations = append(informations, fmt.Sprintf("Provider: %s", profile.Provider))
	switch profile.Provider {
	case Azure:
		informations = append(informations, fmt.Sprintf("API Endpoint: %s", profile.Azure.APIEndpoint))
		informations = append(informations, fmt.Sprintf("API Key: %s", profile.Azure.APIKey))
		informations = append(informations, fmt.Sprintf("Deployment Id: %s", profile.Azure.DeploymentId))
		informations = append(informations, fmt.Sprintf("Assistant Prompt: %s", profile.Azure.AssistantPrompt))
	case OpenAI:
		informations = append(informations, fmt.Sprintf("API Key: %s", profile.OpenAI.APIKey))
		informations = append(informations, fmt.Sprintf("Assistant Prompt: %s", profile.OpenAI.AssistantPrompt))
	}
	return strings.Join(informations, "\n")
}

func ProfileAssistantPrompt(profile Profile) string {
	switch profile.Provider {
	case Azure:
		return profile.Azure.AssistantPrompt
	case OpenAI:
		return profile.OpenAI.AssistantPrompt
	}
	return ""
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
	fmt.Printf("%d profiles deleted\n", len(names))
	return nil
}
