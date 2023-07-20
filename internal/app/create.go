package app

import (
	"errors"
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/imkcat/catchat/internal/config"
	"github.com/imkcat/catchat/internal/essentials"
	"github.com/samber/lo"
)

func (a *App) CreateCommand() error {
	var newProfileProviderString string
	err := survey.AskOne(&survey.Select{
		Message: "Select provider:",
		Options: lo.Map([]config.Provider{config.OpenAI, config.Azure}, func(item config.Provider, index int) string {
			return string(item)
		}),
	}, &newProfileProviderString)
	if err != nil {
		return err
	}
	newProfileProvider := config.Provider(newProfileProviderString)

	newProfile := config.Profile{}
	switch newProfileProvider {
	// OpenAI
	case config.OpenAI:
		var newOpenAIConfig config.OpenAIConfig
		err := survey.Ask([]*survey.Question{
			{
				Name: "api_key",
				Prompt: &survey.Input{
					Message: "OpenAI API key:",
				},
				Validate: survey.Required,
			},
			{
				Name: "model",
				Prompt: &survey.Select{
					Message: "Model:",
					Options: lo.Map([]config.OpenAIModel{
						config.GPT_3_5_Turbo_0613,
						config.GPT_3_5_Turbo_16k_0613,
						config.GPT_3_5_Turbo,
						config.GPT_3_5_Turbo_16k,
						config.GPT_4,
						config.GPT_4_0613,
						config.GPT_4_32k,
						config.GPT_4_32k_0613,
						config.TextEmbeddingAda_002,
						config.TextDavinci_003,
						config.TextDavinci_002,
						config.CodeDavinci_002,
					}, func(item config.OpenAIModel, index int) string {
						return string(item)
					}),
				},
				Validate: survey.Required,
			},
		}, &newOpenAIConfig)
		if err != nil {
			return err
		}
		var assistantPrompt string
		err = survey.AskOne(&survey.Input{
			Message: "Assistant Prompt:",
			Help:    fmt.Sprintf("The prompt that help set the behavior of the assistant.(Default: %s)", essentials.DefaultAssistantPrompt),
			Default: essentials.DefaultAssistantPrompt,
		}, &assistantPrompt)
		if err != nil {
			return err
		}
		if assistantPrompt == "" {
			newOpenAIConfig.AssistantPrompt = essentials.DefaultAssistantPrompt
		}
		newOpenAIConfig.AssistantPrompt = assistantPrompt
		newProfile.OpenAI = &newOpenAIConfig
		newProfile.Provider = config.OpenAI
	// Azure
	case config.Azure:
		var newAzureConfig config.AzureConfig
		err := survey.Ask([]*survey.Question{
			{
				Name: "api_endpoint",
				Prompt: &survey.Input{
					Message: "API Enpoint:",
					Help:    "Azure OpenAI service endpoint, for example: https://{your-resource-name}.openai.azure.com",
				},
				Validate: survey.Required,
			},
			{
				Name: "api_key",
				Prompt: &survey.Input{
					Message: "Azure OpenAI service key:",
				},
				Validate: survey.Required,
			},
			{
				Name: "deployment_id",
				Prompt: &survey.Input{
					Message: "Deployment Id:",
					Help:    "The deployment Id of the model.",
				},
				Validate: survey.Required,
			},
		}, &newAzureConfig)
		if err != nil {
			return err
		}
		var assistantPrompt string
		err = survey.AskOne(&survey.Input{
			Message: "Assistant Prompt:",
			Help:    fmt.Sprintf("The prompt that help set the behavior of the assistant.(Default: %s)", essentials.DefaultAssistantPrompt),
			Default: essentials.DefaultAssistantPrompt,
		}, &assistantPrompt)
		if err != nil {
			return err
		}
		if assistantPrompt == "" {
			newAzureConfig.AssistantPrompt = essentials.DefaultAssistantPrompt
		}
		newAzureConfig.AssistantPrompt = assistantPrompt
		newProfile.Azure = &newAzureConfig
		newProfile.Provider = config.Azure
	}
	var newProfileName string
	err = survey.AskOne(&survey.Input{
		Message: "Profile Name:",
	}, &newProfileName, survey.WithValidator(func(ans interface{}) error {
		for _, v := range a.Config.Profiles {
			if v.Name == ans {
				return errors.New("profile name conflicted")
			}
		}
		if len(ans.(string)) == 0 {
			return errors.New("profile name is empty")
		}
		return nil
	}))
	if err != nil {
		return err
	}
	newProfile.Name = newProfileName
	err = a.Config.CreateProfile(newProfile)
	if err != nil {
		return err
	}
	return nil
}
