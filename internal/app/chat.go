package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/imkcat/catchat/internal/chat"
	"github.com/imkcat/catchat/internal/config"
)

func (a *App) Chat(profile config.Profile) error {
	chatModel, err := chat.NewChatModel(profile)
	if err != nil {
		return err
	}

	chatApp := tea.NewProgram(chatModel, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err = chatApp.Run()
	if err != nil {
		return err
	}
	return nil
}
