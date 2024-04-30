package gpt

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"rnd-git.valsun.cn/ebike-server/go-common/logs"

	"github.com/sashabaranov/go-openai"
)

func TestGpt_ChatWithStream(t *testing.T) {
	c := NewGptClient()
	msg := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "who are you?",
		},
	}
	stream, err := c.ChatWithStream(openai.GPT3Dot5Turbo, 2048, msg)
	if err != nil {
		logs.Errorf("ChatWithStream error: %v", err)
	}
	for {
		if stream == nil {
			continue
		}
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return
		}

		if err != nil {
			logs.Errorf("stream.Recv error: %v", err)
			return
		}
		if len(response.Choices) == 0 {
			continue
		}
		fmt.Print(response.Choices[0].Delta.Content)
	}
}

func TestGpt_Embedding(t *testing.T) {
	c := NewGptClient()
	str := []string{"who are you?"}
	resp, err := c.Embedding(str)
	if err != nil {
		logs.Errorf("Embedding error: %v", err)
	}
	fmt.Println(resp)
}

func TestGetKeyWord(t *testing.T) {
	gptConfig := openai.DefaultConfig("sk-VCmJgAa7UFzmxptp55Ae3921327945E88938303f1743Aa21")
	gptConfig.BaseURL = "http://10.199.1.41:3333/v1"
	client := openai.NewClientWithConfig(gptConfig)
	c := Gpt{
		apiKey:  "sk-VCmJgAa7UFzmxptp55Ae3921327945E88938303f1743Aa21",
		baseUrl: "http://10.199.1.41:3333/v1",
		client:  client,
	}
	str := "what about the battery and the max range?"
	resp, err := c.GetKeyWord(openai.GPT3Dot5Turbo, 20, str)
	if err != nil {
		logs.Errorf("GetKeyWord error: %v", err)
	}
	fmt.Println(resp)
}
