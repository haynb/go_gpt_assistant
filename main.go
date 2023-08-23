package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/documentloaders"
	em "github.com/tmc/langchaingo/embeddings/openai"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"os"
	"qdrant/db"
	"strings"
)

func DocToPoints(docs []schema.Document, e em.OpenAI) ([]map[string]interface{}, error) {
	points := make([]map[string]interface{}, len(docs))
	for i, doc := range docs {
		metadataStrs := make([]string, 0, len(doc.Metadata))
		for k, v := range doc.Metadata {
			metadataStrs = append(metadataStrs, fmt.Sprintf("%s: %v", k, v))
		}
		fullText := fmt.Sprintf("%s\nMetadata:\n%s", doc.PageContent, strings.Join(metadataStrs, "\n"))
		// 获取向量
		embeddingResponse, err := e.EmbedQuery(context.Background(), fullText)
		if err != nil {
			return nil, fmt.Errorf("Failed to get embedding for document %d: %v", i, err)
		}
		// 芜湖起飞
		points[i] = map[string]interface{}{
			"id": i + 1,
			"payload": map[string]interface{}{
				"text": fullText,
			},
			"vectors": embeddingResponse,
		}
	}
	return points, nil
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
	llm, err := openai.NewChat()
	if err != nil {
		panic(err)
	}
	stuffQAChain := chains.LoadStuffQA(llm)
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
	e, err := em.NewOpenAI()
	if err != nil {
		panic(err)
	}
	points, err := DocToPoints(doc, e)
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
	result1, err := db.AddPoints("test", points)
	if err != nil {
		panic(err)
	}
	fmt.Println("add:", result1)
	query := "什么是校招生？"
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
	fmt.Println("search:", put_doc)
	answer, err := chains.Call(context.Background(), stuffQAChain, map[string]any{
		"input_documents": put_doc,
		"question":        query,
	})
	if err != nil {
		panic(err)
	}
	fmt.Print(answer)
}
