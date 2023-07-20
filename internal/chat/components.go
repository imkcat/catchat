package chat

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/imkcat/catchat/internal/config"
	"github.com/imkcat/catchat/internal/essentials"
	"github.com/samber/lo"
)

func (m ChatModel) Divider(topMargin bool, bottomMargin bool) string {
	topMarginString := ""
	bottomMarginString := ""
	if topMargin {
		topMarginString = "\n"
	}
	if bottomMargin {
		bottomMarginString = "\n"
	}
	return fmt.Sprintf("%s%s%s", topMarginString, strings.Repeat("-", m.MainView.Width), bottomMarginString)
}

func (m ChatModel) Message(message ChatMessage) string {
	messageMarkdown, err := m.MarkdownRender.Render(message.Message)
	if err != nil {
		messageMarkdown = ""
	}
	name := ""
	switch message.Role {
	case essentials.Assistant:
		name = lipgloss.
			NewStyle().
			Background(lipgloss.Color(essentials.ColorRoleMap[message.Role])).
			Foreground(lipgloss.Color("#ffffff")).
			PaddingLeft(1).
			PaddingRight(1).
			Render("Assistant")
	case essentials.User:
		name = lipgloss.
			NewStyle().
			Background(lipgloss.Color(essentials.ColorRoleMap[message.Role])).
			Foreground(lipgloss.Color("#ffffff")).
			PaddingLeft(1).
			PaddingRight(1).
			Render("You")
	}
	return fmt.Sprintf("%s - %s\n%s",
		name,
		message.CreatedAt.Format("2006-01-02 15:04:05"),
		messageMarkdown,
	)
}

func (m ChatModel) MessagesHistory() string {
	if len(m.Messages) == 0 {
		return "There is no messages yet"
	}
	messageContents := lo.Map(m.Messages, func(item ChatMessage, index int) string {
		return m.Message(item)
	})
	if m.Streaming {
		messageContents = append(messageContents, m.Message(ChatMessage{
			Role:      essentials.Assistant,
			CreatedAt: m.StreamingStartedAt,
			Message:   m.StreamingMessage,
		}))
	}
	return strings.Join(messageContents, "\n")
}

type HelpKeyMap struct {
	ScrollUp   key.Binding
	ScrollDown key.Binding
	Send       key.Binding
	Quit       key.Binding
}

func (k HelpKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Send, k.ScrollUp, k.ScrollDown, k.Quit}
}

func (k HelpKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Send, k.ScrollUp, k.ScrollDown, k.Quit}, // first column
	}
}

func (m ChatModel) Help() string {
	return m.HelpView.View(HelpKeyMap{
		ScrollUp:   key.NewBinding(key.WithKeys("pgup"), key.WithHelp("pgup", "Scroll Up")),
		ScrollDown: key.NewBinding(key.WithKeys("pgdown"), key.WithHelp("pgdown", "Scroll Down")),
		Send:       key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "Send Message")),
		Quit:       key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "Quit")),
	})
}

func (m ChatModel) MainViewContent() string {
	contents := make([]string, 0)
	contents = append(contents,
		lipgloss.
			NewStyle().
			Bold(true).
			Background(lipgloss.Color(essentials.ColorMain)).
			Foreground(lipgloss.Color("#ffffff")).
			PaddingLeft(1).
			PaddingRight(1).
			Render("catchat"),
	)
	contents = append(contents,
		fmt.Sprintf("Version: %s",
			lipgloss.
				NewStyle().
				Bold(true).
				Render(essentials.Version),
		),
	)
	contents = append(contents,
		fmt.Sprintf("Current Profile: %s",
			lipgloss.
				NewStyle().
				Bold(true).
				Render(m.Profile.Name),
		),
	)
	contents = append(contents,
		fmt.Sprintf("Assistant Prompt: %s",
			lipgloss.
				NewStyle().
				Bold(true).
				Render(config.ProfileAssistantPrompt(m.Profile)),
		),
	)
	contents = append(contents, m.Divider(false, true))
	contents = append(contents, m.MessagesHistory())
	contents = append(contents, m.Divider(true, false))
	contents = append(contents, m.InputTextarea.View())
	contents = append(contents, m.Divider(false, false))
	contents = append(contents, m.Help())
	return lipgloss.NewStyle().Width(m.MainView.Width).Render(strings.Join(contents, "\n"))
}
