package chat

import (
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
	"github.com/imkcat/catchat/internal/config"
)

func NewAzureChatClient(config config.AzureConfig) (*azopenai.Client, error) {
	keyCredential, err := azopenai.NewKeyCredential(config.APIKey)
	if err != nil {
		return nil, err
	}

	client, err := azopenai.NewClientWithKeyCredential(config.APIEndpoint, keyCredential, config.DeploymentId, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewOpenAIChatClient(config config.OpenAIConfig) (*azopenai.Client, error) {
	keyCredential, err := azopenai.NewKeyCredential(config.APIKey)
	if err != nil {
		return nil, err
	}

	client, err := azopenai.NewClientForOpenAI("https://api.openai.com/v1", keyCredential, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
