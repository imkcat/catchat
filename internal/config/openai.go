package config

type OpenAIModel = string

const (
	// Embeddings
	TextEmbeddingAda_002 OpenAIModel = "text-embedding-ada-002"

	// GPT-3.5
	GPT_3_5_Turbo          OpenAIModel = "gpt-3.5-turbo"
	GPT_3_5_Turbo_16k      OpenAIModel = "gpt-3.5-turbo-16k"
	GPT_3_5_Turbo_0613     OpenAIModel = "gpt-3.5-turbo-0613"
	GPT_3_5_Turbo_16k_0613 OpenAIModel = "gpt-3.5-turbo-16k-0613"
	TextDavinci_003        OpenAIModel = "text-davinci-003"
	TextDavinci_002        OpenAIModel = "text-davinci-002"
	CodeDavinci_002        OpenAIModel = "code-davinci-002"

	// GPT-4
	GPT_4          OpenAIModel = "gpt-4"
	GPT_4_0613     OpenAIModel = "gpt-4-0613"
	GPT_4_32k      OpenAIModel = "gpt-4-32k"
	GPT_4_32k_0613 OpenAIModel = "gpt-4-32k-0613"
)

// OpenAI config
type OpenAIConfig struct {
	APIKey          string      `survey:"api_key" mapstructure:"api_key" json:"api_key"`
	Model           OpenAIModel `survey:"model" mapstructure:"model" json:"model"`
	AssistantPrompt string      `survey:"assistant_prompt" mapstructure:"assistant_prompt" json:"assistant_prompt"`
}
