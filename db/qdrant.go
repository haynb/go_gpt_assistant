package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	em "github.com/tmc/langchaingo/embeddings/openai"
	"github.com/tmc/langchaingo/schema"
	"io"
	"net/http"
	"strings"
)

var (
	QdrantBase = "154.12.244.129"
	QdrantPort = "6333"
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

// 创建一个集合
func CreateCollection(collectionName string) ([]byte, error) {
	url := fmt.Sprintf("http://%s:%s/collections/%s", QdrantBase, QdrantPort, collectionName)
	requestBody, err := json.Marshal(map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     1536,
			"distance": "Cosine",
		},
		"on_disk_payload": true,
	})
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	result, err := io.ReadAll(response.Body)
	return result, nil
}

// 删除一个集合
func DeleteCollection(collectionName string) error {
	url := fmt.Sprintf("http://%s:%s/collections/%s", QdrantBase, QdrantPort, collectionName)
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	return nil
}

// 查询集合信息
func GetCollection(collectionName string) ([]byte, error) {
	url := fmt.Sprintf("http://%s:%s/collections/%s", QdrantBase, QdrantPort, collectionName)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	result, err := io.ReadAll(response.Body)
	return result, nil
}

// 增加向量
func AddPoints(collectionName string, points []map[string]interface{}) (string, error) {
	url := fmt.Sprintf("http://%s:%s/collections/%s/points?wait=true", QdrantBase, QdrantPort, collectionName)
	requestBody, err := json.Marshal(map[string]interface{}{
		"points": points,
	})
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	result, err := io.ReadAll(response.Body)
	if response.StatusCode != 200 {
		return string(result), fmt.Errorf("add points failed: %s, status code: %d", string(result), response.StatusCode)
	}
	return string(result), nil
}

// 搜索向量
func Search(collectionName string, params map[string]interface{}, vector []float64, limit int) (string, error) {
	url := fmt.Sprintf("http://%s:%s/collections/%s/points/search", QdrantBase, QdrantPort, collectionName)
	requestBody, err := json.Marshal(map[string]interface{}{
		"params":       params,
		"vector":       vector,
		"limit":        limit,
		"with_payload": true,
	})
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	result, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return string(result), nil
}
