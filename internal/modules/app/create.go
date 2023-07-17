package app

import (
	"errors"
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/imkcat/catchat/internal/modules/config"
	"github.com/samber/lo"
)

func (a *App) CreateCommand() error {
	var newProfileProviderString string
	err := survey.AskOne(&survey.Select{
		Message: "Select provider:",
		Options: lo.Map([]config.Provider{config.Azure, config.OpenAI}, func(item config.Provider, index int) string {
			return string(item)
		}),
	}, &newProfileProviderString)
	if err != nil {
		return err
	}
	newProfileProvider := config.Provider(newProfileProviderString)

	newProfile := config.Profile{}
	switch newProfileProvider {
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
					Help:    "The deployment Id of the model",
				},
				Validate: survey.Required,
			},
		}, &newAzureConfig)
		if err != nil {
			return err
		}
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
	fmt.Printf("Profile %s has been created! Please re-run catchat command and select it to start chatting!\n", newProfileName)
	return nil
}
