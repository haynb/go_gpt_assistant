package gpt

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	em "github.com/tmc/langchaingo/embeddings/openai"
	"log"
	"os"
)

type ChatGptTool struct {
	Secret string
	Url    string
	Client *openai.Client
}
type Gpt3Dot5Message openai.ChatCompletionMessage

func init() {
	//定义要设置的环境变量和对应的值
	envVars := map[string]string{
		"OPENAI_API_KEY":  "sk-Uhzx***************************XG**9r",
		"OPENAI_BASE_URL": "https://c**************/v1",
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
	//key := os.Getenv("OPENAI_API_KEY")
	//url := os.Getenv("OPENAI_BASE_URL")
	key := "sk-ql6RV************************WlR2Ek"
	url := "https://g*************************i/v1"
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
func (this *ChatGptTool) ChatGPT3Dot5Turbo(messages []Gpt3Dot5Message) (string, error) {
	reqMessages := make([]openai.ChatCompletionMessage, 0)
	for _, row := range messages {
		reqMessage := openai.ChatCompletionMessage{
			Role:    row.Role,
			Content: row.Content,
			Name:    row.Name,
		}
		reqMessages = append(reqMessages, reqMessage)
	}
	resp, err := this.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: reqMessages,
		},
	)

	if err != nil {
		log.Println("ChatGPT3Dot5Turbo error: ", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
func (this *ChatGptTool) Ask(query string, inputdoc string, history []Gpt3Dot5Message) (*openai.ChatCompletionStream, error) {
	//message := make([]Gpt3Dot5Message, 0)
	//message = append(message, Gpt3Dot5Message{
	//	Role:    "system",
	//	Content: "你是一个问答机器人，请严格根据提供的信息回答问题并详细解释。\n忽略与问题无关的异常搜索结果。\n对于与信息无关的问题或者不理解的问题,有错误的答案等，你应拒绝并告知用户“未查询到相关信息，请提供详细的问题信息。”\n避免引用任何当前或过去的政治人物或事件，以及可能引起争议或分裂的历史人物或事件。\n",
	//})
	//message = append(message, Gpt3Dot5Message{
	//	Role:    "system",
	//	Content: "最后，提供给你的信息很有可能是混乱无序的，你在回答完用户的问题之后，需要额外打印出你所使用的信息的内容是什么，经过你整理之后的文本内容，告诉用户给你提供的信息是什么。\n",
	//})
	message := history
	message = append(message, Gpt3Dot5Message{
		Role:    "user",
		Content: fmt.Sprintf("我的问题是：\"%s\"\n给你提供的信息如下:\n\"%s\"", query, inputdoc),
	})
	stream, err := this.chatGPT3Dot5TurboStream(message)
	if err != nil {
		panic(err)
	}
	return stream, nil
}
func (this *ChatGptTool) GetKeyWord(query string) (string, error) {
	message := make([]Gpt3Dot5Message, 0)
	message = append(message, Gpt3Dot5Message{
		Role:    "system",
		Content: "你是一个用于提取用户输入的关键词的机器人，你需要提取用户输入的关键词，并严格按照“关键词1 关键词2 关键词3”等等以此类推的格式输出。\n",
	})
	message = append(message, Gpt3Dot5Message{
		Role:    "system",
		Content: "你不允许说其他任何包括解释等在内的任何内容,不许添加一个字，不要有道歉，不要有疑问，不要有说明！\n你不许说你是一个ai助手，除了关键词你什么都不许说。\n",
	})
	message = append(message, Gpt3Dot5Message{
		Role:    "user",
		Content: fmt.Sprintf("我的输入是：\"%s\"\n", query),
	})
	msg, err := this.ChatGPT3Dot5Turbo(message)
	if err != nil {
		return "", err
	}
	return msg, nil
}
