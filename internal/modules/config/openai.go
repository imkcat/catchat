package config

// OpenAI config
type OpenAIConfig struct {
	Profile
	Organization string
	APIKey       string
	ModelId      string
}
