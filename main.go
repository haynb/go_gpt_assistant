package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/documentloaders"
	em "github.com/tmc/langchaingo/embeddings/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"os"
	"strings"
)

func docToPo(docs []schema.Document) ([]string, error) {
	texts := make([]string, len(docs))
	// Iterate over the docs slice
	for i, doc := range docs {
		// Combine the PageContent and Metadata into a single string
		metadataStrs := make([]string, 0, len(doc.Metadata))
		for k, v := range doc.Metadata {
			metadataStrs = append(metadataStrs, fmt.Sprintf("%s: %v", k, v))
		}
		texts[i] = fmt.Sprintf("%s\nMetadata:\n%s", doc.PageContent, strings.Join(metadataStrs, "\n"))
	}

	return texts, nil
}
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
	//llm, err := openai.NewChat()
	//if err != nil {
	//	panic(err)
	//}
	//stuffQAChain := chains.LoadStuffQA(llm)
	file, err := os.Open("xiaozhao.pdf")
	if err != nil {
		panic(err)
	}
	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fileSize := fileInfo.Size()
	p := documentloaders.NewPDF(file, fileSize)
	spliter := textsplitter.NewTokenSplitter()
	doc, err := p.LoadAndSplit(context.Background(), spliter)
	if err != nil {
		panic(err)
	}

	e, err := em.NewOpenAI()
	if err != nil {
		panic(err)
	}
	docs, err := docToEm(doc)
	if err != nil {
		panic(err)
	}
	put_doc, err := e.EmbedDocuments(context.Background(), docs)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(put_doc), put_doc[0])
	//answer, err := chains.Call(context.Background(), stuffQAChain, map[string]any{
	//	"input_documents": put_doc,
	//	"question":        "什么是校招生?",
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Print(answer)
}
