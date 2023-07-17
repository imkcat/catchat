package app

import (
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/imkcat/catchat/internal/modules/config"
	"github.com/samber/lo"
)

func (a *App) ListCommand() error {
	if err := a.CheckProfiles(); err != nil {
		return err
	}
	var selectProfileIndex int
	err := survey.AskOne(&survey.Select{
		Message: "All profiles, select profile to check the configuration:",
		Options: lo.Map(a.Config.Profiles, func(item config.Profile, index int) string {
			return item.Name
		}),
	}, &selectProfileIndex, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}
	selectedProfile := a.Config.Profiles[selectProfileIndex]
	fmt.Printf("Name: %s\n", selectedProfile.Name)
	fmt.Printf("Provider: %s\n", selectedProfile.Provider)
	switch selectedProfile.Provider {
	case config.Azure:
		fmt.Printf("API Endpoint: %s\n", selectedProfile.Azure.APIEndpoint)
		fmt.Printf("API Key: %s\n", selectedProfile.Azure.APIKey)
		fmt.Printf("Deployment Id: %s\n", selectedProfile.Azure.DeploymentId)
	}
	return nil
}
