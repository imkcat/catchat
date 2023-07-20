package config

// Azure config
type AzureConfig struct {
	APIEndpoint     string `survey:"api_endpoint" mapstructure:"api_endpoint" json:"api_endpoint"`
	APIKey          string `survey:"api_key" mapstructure:"api_key" json:"api_key"`
	DeploymentId    string `survey:"deployment_id" mapstructure:"deployment_id" json:"deployment_id"`
	AssistantPrompt string `survey:"assistant_prompt" mapstructure:"assistant_prompt" json:"assistant_prompt"`
}
