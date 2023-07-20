package app

import (
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/imkcat/catchat/internal/config"
	"github.com/samber/lo"
)

func (a *App) SwitchCommand() error {
	if err := a.CheckProfiles(); err != nil {
		return err
	}
	var switchedProfileIndex int
	err := survey.AskOne(&survey.Select{
		Message: "Switch profile to:",
		Options: lo.Map(a.Config.Profiles, func(item config.Profile, index int) string {
			return fmt.Sprintf("%s - %s", item.Name, item.Provider)
		}),
	}, &switchedProfileIndex)
	if err != nil {
		return err
	}
	switchedProfile := a.Config.Profiles[switchedProfileIndex]
	err = a.Config.SwitchProfile(switchedProfile)
	if err != nil {
		return err
	}
	return nil
}
