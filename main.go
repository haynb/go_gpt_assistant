package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
	"io"
	"os"
	"qdrant/db"
	"qdrant/gpt"
)

func main() {
	// 定义要设置的环境变量和对应的值
	envVars := map[string]string{
		"OPENAI_API_KEY":  "sk-UhzxW8avagN8D00Q2AmHeJL1LBz6NFv9rI3USa94GiXGPa9r",
		"OPENAI_BASE_URL": "https://cfwus02.opapi.win/v1/",
		"OPENAI_MODEL":    "gpt-3.5-turbo",
	}

	// 逐个设置环境变量
	for key, value := range envVars {
		err := os.Setenv(key, value)
		if err != nil {
			panic(err)
		}
	}

	file, err := os.Open("xiaozhao.pdf")
	if err != nil {
		panic(err)
	}
	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
		return
	}
	fileSize := fileInfo.Size()
	p := documentloaders.NewPDF(file, fileSize)
	spliter := textsplitter.NewTokenSplitter()
	doc, err := p.LoadAndSplit(context.Background(), spliter)
	if err != nil {
		panic(err)
	}
	e, err := gpt.NewEm()
	if err != nil {
		panic(err)
	}
	points, err := db.DocToPoints(doc, e)
	if err != nil {
		panic(err)
	}
	err = db.DeleteCollection("test")
	if err != nil {
		panic(err)
	}
	result0, err := db.CreateCollection("test")
	if err != nil {
		panic(err)
	}
	fmt.Println("create:", string(result0))
	_, err = db.AddPoints("test", points)
	if err != nil {
		panic(err)
	}
	query := "赛维需要加班吗？"
	q_e, err := e.EmbedQuery(context.Background(), query)
	if err != nil {
		panic(err)
	}
	put_doc, err := db.Search(
		"test",
		map[string]interface{}{"exact": false, "hnsw_ef": 128},
		q_e,
		4)
	if err != nil {
		panic(err)
	}
	//fmt.Println("search:", put_doc)
	tools := gpt.NewChatGptTool()
	stream, err := tools.Ask(query, put_doc)
	if err != nil {
		panic(err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf(msg.Choices[0].Delta.Content)
	}
	fmt.Println(put_doc)
}
