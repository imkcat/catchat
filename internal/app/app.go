package app

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/imkcat/catchat/internal/config"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

type App struct {
	CliApp         *cli.App
	Config         *config.Config
	ConfigPath     string
	ConfigFilePath string
}

func NewApp() (*App, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := fmt.Sprintf("%s/.config/catchat", userHomeDir)
	configFilePath := fmt.Sprintf("%s/config.json", configPath)
	configClient, err := config.NewConfig(configPath, configFilePath)
	if err != nil {
		return nil, err
	}
	cliApp := cli.App{
		Name:     "CatChat",
		Usage:    "AI chat on your terminal",
		HelpName: "catchat",
	}
	appInstance := App{
		CliApp:         &cliApp,
		Config:         configClient,
		ConfigPath:     configPath,
		ConfigFilePath: configFilePath,
	}
	appInstance.InitCliApp()
	viper.OnConfigChange(func(in fsnotify.Event) {
		newConfig, err := config.ReloadConfig()
		if err != nil {
			log.Fatal(err)
		}
		appInstance.Config = newConfig
	})
	viper.WatchConfig()
	return &appInstance, nil
}

func (a *App) InitCliApp() {
	a.CliApp.Action = func(*cli.Context) error {
		err := a.MainCommand()
		if err != nil {
			return err
		}
		return nil
	}
	a.CliApp.Commands = []cli.Command{
		{
			Name:    "switch",
			Aliases: []string{"s"},
			Usage:   "Switch profile",
			Action: func(ctx *cli.Context) error {
				err := a.SwitchCommand()
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Create new profile",
			Action: func(ctx *cli.Context) error {
				err := a.CreateCommand()
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "Detele profiles",
			Action: func(ctx *cli.Context) error {
				err := a.DeleteCommand()
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List and check profile",
			Action: func(ctx *cli.Context) error {
				err := a.ListCommand()
				if err != nil {
					return err
				}
				return nil
			},
		},
	}
}

func (a *App) CheckProfiles() error {
	if len(a.Config.Profiles) == 0 {
		return errors.New("no profiles")
	}
	return nil
}
