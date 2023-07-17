package chat

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
	"github.com/imkcat/catchat/internal/modules/config"
)

func AzureChat(config *config.AzureConfig) error {
	ctx := context.Background()
	systemRole := azopenai.ChatRoleSystem
	systemContent := "You are an AI assistant that helps people find information."
	userRole := azopenai.ChatRoleUser
	assistantRole := azopenai.ChatRoleAssistant
	messages := []azopenai.ChatMessage{
		{
			Role:    &systemRole,
			Content: &systemContent,
		},
	}

	keyCredential, err := azopenai.NewKeyCredential(config.APIKey)
	if err != nil {
		return err
	}

	client, err := azopenai.NewClientWithKeyCredential(config.APIEndpoint, keyCredential, config.DeploymentId, nil)
	if err != nil {
		return err
	}

	for {
		var message string
		err := survey.AskOne(&survey.Input{
			Message: "Input:",
		}, &message, survey.WithValidator(survey.Required))
		if err != nil {
			return err
		}
		messages = append(messages, azopenai.ChatMessage{
			Role:    &userRole,
			Content: &message,
		})
		response, err := client.GetChatCompletionsStream(ctx, azopenai.ChatCompletionsOptions{
			Messages: messages,
		}, nil)
		if err != nil {
			return err
		}

		assistantMessage := ""

		for {
			entry, err := response.ChatCompletionsStream.Read()

			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				return err
			}

			for _, choice := range entry.Choices {
				if choice.Delta.Content != nil {
					assistantMessage = fmt.Sprintf("%s%s", assistantMessage, *choice.Delta.Content)
					fmt.Print(*choice.Delta.Content)
				}
			}
		}
		fmt.Print("\n")
		if assistantMessage != "" {
			messages = append(messages, azopenai.ChatMessage{
				Role:    &assistantRole,
				Content: &assistantMessage,
			})
		}
	}
}
