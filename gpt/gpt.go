package gpt

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
)

type ChatGptTool struct {
	Secret string
	Url    string
	Client *openai.Client
}
type OpenAIEmbedder struct {
	Client *openai.Client
}
type Gpt3Dot5Message openai.ChatCompletionMessage

func NewEm() *OpenAIEmbedder {
	config := openai.DefaultAzureConfig("4beb6d9440bb4e93ab56e39f230c9f45", "https://heanyang.openai.azure.com/")
	config.APIVersion = "2023-05-15" // optional update to latest API version
	client := openai.NewClientWithConfig(config)
	return &OpenAIEmbedder{Client: client}
}
func (embedder *OpenAIEmbedder) EmbedQuery(input string) ([]float64, error) {
	resp, err := embedder.Client.CreateEmbeddings(
		context.Background(),
		openai.EmbeddingRequest{
			Input: []string{input},
			Model: openai.AdaEmbeddingV2,
		})
	if err != nil {
		log.Printf("CreateEmbeddings error: %v\n", err)
		return nil, err
	}
	vectors := resp.Data[0].Embedding
	vectors64 := make([]float64, len(vectors))
	for i, v := range vectors {
		vectors64[i] = float64(v)
	}

	return vectors64, nil
}
func NewChatGptTool() *ChatGptTool {
	key := "4beb6d9440bb4e93ab56e39f230c9f45"
	url := "https://heanyang.openai.azure.com/"
	config := openai.DefaultAzureConfig(key, url)
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
		Content: "You are a robot for extracting keywords entered by users. You need to extract the keywords entered by users and output them in strict accordance with the format of \"Keyword 1 Keyword 2 Keyword 3\" and so on.\n",
	})
	message = append(message, Gpt3Dot5Message{
		Role:    "system",
		Content: "You are not allowed to say anything else including explanations, etc., not to add a word, not to apologize\nYou are not allowed to say that you are an AI assistant, and you are not allowed to say anything except keywords.\n",
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
