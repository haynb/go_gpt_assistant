package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	QdrantBase = "154.12.244.129"
	QdrantPort = "6333"
)

// 创建一个集合
func CreateCollection(collectionName string) error {
	url := fmt.Sprintf("http://%s:%s/collections", QdrantBase, QdrantPort)
	requestBody, err := json.Marshal(map[string]interface{}{
		"name": collectionName,
		"vectors": map[string]interface{}{
			"sizes": 1536,
			"dists": "cosine",
		},
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
	return nil
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
	url := fmt.Sprintf("http://%s:%s/collections/%s", QdrantBase, QdrantPort, collectionName)
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
func Search(collectionName string, query map[string]interface{}, vector []float64, limit int) ([]byte, error) {
	url := fmt.Sprintf("http://%s:%s/collections/%s/search", QdrantBase, QdrantPort, collectionName)
	requestBody, err := json.Marshal(map[string]interface{}{
		"params":       query,
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
	return result, nil
}
