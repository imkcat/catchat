package app

import (
	survey "github.com/AlecAivazis/survey/v2"
	"github.com/imkcat/catchat/internal/modules/chat"
	"github.com/imkcat/catchat/internal/modules/config"
	"github.com/samber/lo"
)

func (a *App) MainCommand() error {
	if len(a.Config.Profiles) == 0 {
		createNewProfile := true
		err := survey.AskOne(&survey.Confirm{
			Message: "There is no profile. Would you like to create a new one now?",
			Default: true,
		}, &createNewProfile)
		if err != nil {
			return err
		}
		if createNewProfile {
			err := a.CreateCommand()
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	}
	var selectedProfileIndex int
	err := survey.AskOne(&survey.Select{
		Message: "Select profile:",
		Options: lo.Map(a.Config.Profiles, func(item config.Profile, index int) string {
			return item.Name
		}),
	}, &selectedProfileIndex, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}
	selectedProfile := a.Config.Profiles[selectedProfileIndex]
	switch selectedProfile.Provider {
	case config.Azure:
		err = chat.AzureChat(selectedProfile.Azure)
		if err != nil {
			return err
		}
	}
	return nil
}
