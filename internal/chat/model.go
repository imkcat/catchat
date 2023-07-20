package chat

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/imkcat/catchat/internal/config"
	"github.com/imkcat/catchat/internal/essentials"
)

type ChatMessage struct {
	Role      essentials.ChatRole
	CreatedAt time.Time
	Message   string
}

type ActionStatus = string

const (
	SentMessage    ActionStatus = "sendMessage"
	StreamStarted  ActionStatus = "streamStarted"
	Streaming      ActionStatus = "streaming"
	StreamFinished ActionStatus = "streamFinished"
	Error          ActionStatus = "error"
)

type ChatAction struct {
	status  ActionStatus
	message string
	err     error
}

type ChatModel struct {
	Profile            config.Profile
	MainView           viewport.Model
	InputTextarea      textarea.Model
	HelpView           help.Model
	MarkdownRender     glamour.TermRenderer
	Messages           []ChatMessage
	ActionSub          chan ChatAction
	Streaming          bool
	StreamingMessage   string
	StreamingStartedAt time.Time
}

func NewChatModel(profile config.Profile) (*ChatModel, error) {
	inputTextarea := textarea.New()
	inputTextarea.Placeholder = "Ask me anything..."
	inputTextarea.Prompt = ""
	inputTextarea.SetHeight(10)
	inputTextarea.Focus()

	helpView := help.New()

	mainView := viewport.New(100, 100)

	markdownRender, _ := glamour.NewTermRenderer(glamour.WithAutoStyle())

	chatModel := ChatModel{
		Profile:        profile,
		InputTextarea:  inputTextarea,
		MainView:       mainView,
		HelpView:       helpView,
		MarkdownRender: *markdownRender,
		ActionSub:      make(chan ChatAction),
		Streaming:      false,
	}

	return &chatModel, nil
}

func (m ChatModel) Init() tea.Cmd {
	return WaitForChatAction(m.ActionSub)
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.MainView.Height = msg.Height
		m.MainView.Width = msg.Width
		m.InputTextarea.SetWidth(msg.Width)
		m.MainView.SetContent(m.MainViewContent())
	case tea.MouseMsg:
		var mainViewCmd tea.Cmd
		m.MainView, mainViewCmd = m.MainView.Update(msg)
		cmds = append(cmds, mainViewCmd)
		m.MainView.SetContent(m.MainViewContent())
	case tea.KeyMsg:
		var inputTextareaCmd tea.Cmd
		if msg.String() == "ctrl+s" {
			if len(m.InputTextarea.Value()) != 0 {
				newMessage := m.InputTextarea.Value()
				return m, tea.Batch(SendMessage(m.ActionSub, newMessage))
			}
		} else {
			switch msg.Type {
			case tea.KeyCtrlC:
				return m, tea.Quit
			case tea.KeyPgUp:
				m.MainView.ViewUp()
			case tea.KeyPgDown:
				m.MainView.ViewDown()
			}
			m.InputTextarea, inputTextareaCmd = m.InputTextarea.Update(msg)
			cmds = append(cmds, inputTextareaCmd)
		}
		m.MainView.SetContent(m.MainViewContent())
	case ChatAction:
		switch msg.status {
		case SentMessage:
			m.Messages = append(m.Messages, ChatMessage{
				Role:      essentials.User,
				CreatedAt: time.Now(),
				Message:   msg.message,
			})
			m.InputTextarea.Reset()
			m.InputTextarea.Blur()
			cmds = append(cmds, ChatCompletionsStream(m.ActionSub, m.Profile, m.Messages))
		case StreamStarted:
			m.StreamingStartedAt = time.Now()
			m.StreamingMessage = ""
			m.Streaming = true
		case Streaming:
			m.StreamingMessage = msg.message
			m.MainView.GotoBottom()
		case StreamFinished:
			m.Streaming = false
			m.StreamingMessage = ""
			m.Messages = append(m.Messages, ChatMessage{
				Role:      essentials.Assistant,
				CreatedAt: m.StreamingStartedAt,
				Message:   msg.message,
			})
			m.InputTextarea.Focus()
		}
		cmds = append(cmds, WaitForChatAction(m.ActionSub))
		m.MainView.SetContent(m.MainViewContent())
	}
	return m, tea.Batch(cmds...)
}

func (m ChatModel) View() string {
	return m.MainView.View()
}

func WaitForChatAction(actionSub chan ChatAction) tea.Cmd {
	return func() tea.Msg {
		return ChatAction(<-actionSub)
	}
}

func SendMessage(actionSub chan ChatAction, newMessage string) tea.Cmd {
	return func() tea.Msg {
		actionSub <- ChatAction{
			status:  SentMessage,
			message: newMessage,
		}
		return nil
	}
}

func ChatCompletionsStream(actionSub chan ChatAction, profile config.Profile, messages []ChatMessage) tea.Cmd {
	return func() tea.Msg {
		switch profile.Provider {
		case config.Azure, config.OpenAI:
			var chatClient *azopenai.Client
			var err error
			var model *string
			if profile.Provider == config.Azure {
				chatClient, err = NewAzureChatClient(*profile.Azure)
			}
			if profile.Provider == config.OpenAI {
				chatClient, err = NewOpenAIChatClient(*profile.OpenAI)
				model = &profile.OpenAI.Model
			}
			if err != nil {
				actionSub <- ChatAction{
					status: Error,
					err:    err,
				}
				return nil
			}
			systemRole := azopenai.ChatRoleSystem
			assistantRole := azopenai.ChatRoleAssistant
			userRole := azopenai.ChatRoleUser
			assistantPrompt := config.ProfileAssistantPrompt(profile)
			newMessages := make([]azopenai.ChatMessage, 0)
			newMessages = append(newMessages, azopenai.ChatMessage{
				Role:    &systemRole,
				Content: &assistantPrompt,
			})
			for _, v := range messages {
				newMessage := v.Message
				switch v.Role {
				case essentials.Assistant:
					newMessages = append(newMessages, azopenai.ChatMessage{
						Role:    &assistantRole,
						Content: &newMessage,
					})
				case essentials.User:
					newMessages = append(newMessages, azopenai.ChatMessage{
						Role:    &userRole,
						Content: &newMessage,
					})
				}
			}
			response, err := chatClient.GetChatCompletionsStream(context.Background(), azopenai.ChatCompletionsOptions{
				Messages: newMessages,
				Model:    model,
			}, nil)
			if err != nil {
				actionSub <- ChatAction{
					status: Error,
					err:    err,
				}
				return nil
			}

			actionSub <- ChatAction{
				status: StreamStarted,
			}

			streamingMessage := ""

			for {
				entry, err := response.ChatCompletionsStream.Read()

				if errors.Is(err, io.EOF) {
					break
				}

				if err != nil {
					actionSub <- ChatAction{
						status: Error,
						err:    err,
					}
				}

				for _, choice := range entry.Choices {
					if choice.Delta.Content != nil {
						streamingMessage = fmt.Sprintf("%s%s", streamingMessage, *choice.Delta.Content)
						actionSub <- ChatAction{
							status:  Streaming,
							message: streamingMessage,
						}
					}
				}
			}
			actionSub <- ChatAction{
				status:  StreamFinished,
				message: streamingMessage,
			}
		}
		return nil
	}
}
