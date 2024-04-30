package gpt

import (
	"context"
	"strconv"

	"go-gpt-assistant/config"

	"github.com/sashabaranov/go-openai"
)

type Gpt struct {
	apiKey  string
	baseUrl string
	client  *openai.Client
}

var GptClient *Gpt

func InitGpt() {
	GptClient = NewGptClient()
}

func NewGptClient() *Gpt {
	baseUrl := config.GetAppConf().GptUrl
	apiKey := config.GetAppConf().GptApiKey
	gptConfig := openai.DefaultConfig(apiKey)
	gptConfig.BaseURL = baseUrl
	c := openai.NewClientWithConfig(gptConfig)
	return &Gpt{
		apiKey:  apiKey,
		baseUrl: baseUrl,
		client:  c,
	}
}

func (gc *Gpt) ChatWithStream(model string, maxTokens int, messages []openai.ChatCompletionMessage) (*openai.ChatCompletionStream, error) {
	req := openai.ChatCompletionRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Messages:  messages,
		Stream:    true,
	}
	ctx := context.Background()
	stream, err := gc.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func (gc *Gpt) Embedding(str []string) ([]float32, error) {
	queryReq := openai.EmbeddingRequest{
		Input: str,
		Model: openai.AdaEmbeddingV2,
	}
	ctx := context.Background()
	queryResponse, err := gc.client.CreateEmbeddings(ctx, queryReq)
	if err != nil {
		return nil, err
	}
	return queryResponse.Data[0].Embedding, nil
}

func (gc *Gpt) GetKeyWord(model string, maxTokens int, query string) (string, error) {
	message := make([]openai.ChatCompletionMessage, 0)
	message = append(message, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleSystem,
		Content: "You are a keyword extraction robot. Users will input a paragraph of text, and you need to extract " +
			"the keywords from it. Then output them in this format: \"keyword1,keyword2,keyword3...\". You should use the " +
			"same language as the user's input text.",
	},
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: query,
		},
	)
	req := openai.ChatCompletionRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Messages:  message,
		Stream:    false,
	}
	ctx := context.Background()
	queryResponse, err := gc.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	return queryResponse.Choices[0].Message.Content, nil
}

func (gc *Gpt) ChooseDoc(model string, maxTokens int, query string, texts []string) (string, error) {
	message := make([]openai.ChatCompletionMessage, 0)
	var text string
	for i, t := range texts {
		text += "The information No." + strconv.Itoa(i) + "\n[" + t + "]" + "\n"
	}
	message = append(message, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleSystem,
		Content: "You are a document selection robot. From the several pieces of information the user inputs to you, " +
			"you need to select the information most relevant to the query after." +
			" Then output in this format: If you think the first document is most relevant, output \"1\". If it's the " +
			"second one, output \"2\". And the query is : " + query,
	},
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: text,
		},
	)
	req := openai.ChatCompletionRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Messages:  message,
		Stream:    false,
	}
	ctx := context.Background()
	queryResponse, err := gc.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	return queryResponse.Choices[0].Message.Content, nil
}
