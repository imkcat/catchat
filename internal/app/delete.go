package app

import (
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/imkcat/catchat/internal/config"
	"github.com/samber/lo"
)

func (a *App) DeleteCommand() error {
	if err := a.CheckProfiles(); err != nil {
		return err
	}
	profileIndexes := make([]int, 0)
	err := survey.AskOne(&survey.MultiSelect{
		Message: "Select providers to delete:",
		Options: lo.Map(a.Config.Profiles, func(item config.Profile, index int) string {
			return fmt.Sprintf("%s - %s", item.Name, item.Provider)
		}),
	}, &profileIndexes)
	if err != nil {
		return err
	}
	if len(profileIndexes) == 0 {
		fmt.Println("Nothing deleted")
		return nil
	}
	deleteProfileNames := make([]string, 0)
	for _, v := range profileIndexes {
		deleteProfileNames = append(deleteProfileNames, a.Config.Profiles[v].Name)
	}
	err = a.Config.DeleteProfile(deleteProfileNames)
	if err != nil {
		return err
	}
	return nil
}
