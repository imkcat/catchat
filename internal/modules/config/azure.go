package config

// Azure config
type AzureConfig struct {
	APIEndpoint  string `survey:"api_endpoint"`
	APIKey       string `survey:"api_key"`
	DeploymentId string `survey:"deployment_id"`
}
