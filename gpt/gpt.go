package gpt

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	em "github.com/tmc/langchaingo/embeddings/openai"
	"os"
)

type ChatGptTool struct {
	Secret string
	Url    string
	Client *openai.Client
}
type Gpt3Dot5Message openai.ChatCompletionMessage

func init() {
	// 定义要设置的环境变量和对应的值
	envVars := map[string]string{
		"OPENAI_API_KEY":  "sk-UhzxW8avagN8D00Q2AmHeJL1LBz6NFv9rI3USa94GiXGPa9r",
		"OPENAI_BASE_URL": "https://cfwus02.opapi.win/v1",
		"OPENAI_MODEL":    "gpt-3.5-turbo",
	}

	// 逐个设置环境变量
	for key, value := range envVars {
		err := os.Setenv(key, value)
		if err != nil {
			panic(err)
		}
	}
}

func NewEm() (em.OpenAI, error) {
	e, err := em.NewOpenAI()
	if err != nil {
		panic(err)
	}
	return e, nil
}
func NewChatGptTool() *ChatGptTool {
	key := os.Getenv("OPENAI_API_KEY")
	url := os.Getenv("OPENAI_BASE_URL")
	config := openai.DefaultConfig(key)
	config.BaseURL = url
	client := openai.NewClientWithConfig(config)
	return &ChatGptTool{
		Secret: key,
		Client: client,
		Url:    url,
	}
}
func (this *ChatGptTool) chatGPT3Dot5TurboStream(messages []Gpt3Dot5Message) (*openai.ChatCompletionStream, error) {
	c := this.Client
	ctx := context.Background()
	reqMessages := make([]openai.ChatCompletionMessage, 0)
	for _, row := range messages {
		reqMessage := openai.ChatCompletionMessage{
			Role:    row.Role,
			Content: row.Content,
			Name:    row.Name,
		}
		reqMessages = append(reqMessages, reqMessage)
	}
	req := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: reqMessages,
		Stream:   true,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return stream, err
	}
	return stream, nil
}
func (this *ChatGptTool) Ask(query string, inputdoc string) (*openai.ChatCompletionStream, error) {
	message := make([]Gpt3Dot5Message, 0)
	message = append(message, Gpt3Dot5Message{
		Role:    "system",
		Content: "你是一个问答机器人，请严格根据提供的信息回答问题并详细解释。\n忽略与问题无关的异常搜索结果。\n对于与信息无关的问题或者不理解的问题,有错误的答案等，你应拒绝并告知用户“未查询到相关信息，请提供详细的问题信息。”\n避免引用任何当前或过去的政治人物或事件，以及可能引起争议或分裂的历史人物或事件。",
	})
	message = append(message, Gpt3Dot5Message{
		Role:    "user",
		Content: fmt.Sprintf("问题是：\"%s\"\n给你提供的文件如下:\"%s\"", query, inputdoc),
	})
	stream, err := this.chatGPT3Dot5TurboStream(message)
	if err != nil {
		panic(err)
	}
	return stream, nil
}
