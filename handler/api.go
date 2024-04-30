package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"go-gpt-assistant/utils"

	drant_db "go-gpt-assistant/db/vector-db"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go-gpt-assistant/gpt"
	"go-gpt-assistant/tools"
	"rnd-git.valsun.cn/ebike-server/go-common/logs"
)

const uploadDir = "upload"

func upload(c *gin.Context) {
	logs.Infof("%s -- upload:  ", c.Request.RemoteAddr)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		logs.Errorf("upload file error: %v", err)
		c.JSON(http.StatusForbidden, "upload file error")
		return
	}
	defer file.Close()
	filePath := filepath.Join(uploadDir, header.Filename)
	err = utils.DealFile(file, filePath, uploadDir)
	if err != nil {
		logs.Errorf("check file error: %v", err)
		c.JSON(http.StatusForbidden, "check file error")
		return
	}
	fileContent, err := tools.LoadPdfFile(filePath)
	if err != nil {
		logs.Errorf("load file error: %v", err)
		c.JSON(http.StatusForbidden, "load file error")
		return
	}
	utils.AddFileToMongo(fileContent)
	points, err := drant_db.DocToPoints(fileContent, gpt.GptClient)
	if err != nil {
		logs.Errorf("doc to points error: %v", err)
		c.JSON(http.StatusForbidden, "doc to points error")
		return
	}
	_, err = drant_db.AddPoints("test", points)
	if err != nil {
		logs.Errorf("add points error: %v", err)
		c.JSON(http.StatusForbidden, "add points error")
		return
	}
	logs.Infof("upload success")
	// 响应客户端
	c.JSON(http.StatusOK, gin.H{"message": "upload success"})
}

func checkFile(c *gin.Context) {
	logs.Infof("%s -- checkFile:  ", c.Request.RemoteAddr)
	fileList := make([]string, 0)
	err := filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == uploadDir {
			return nil // 排除根目录
		}
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		logs.Errorf("check file error: %v", err)
		c.JSON(http.StatusForbidden, "check file error")
		return
	}
	c.JSON(http.StatusOK, fileList)
}

func chatWithFile(c *gin.Context) {
	logs.Infof("%s -- chatWithFile:  ", c.Request.RemoteAddr)
	query := c.DefaultPostForm("query", "")
	historyJson := c.DefaultPostForm("history", "")
	var history []openai.ChatCompletionMessage
	err := json.Unmarshal([]byte(historyJson), &history)
	if err != nil {
		logs.Errorf("json unmarshal error: %v", err)
		c.JSON(http.StatusForbidden, "json unmarshal error")
		return
	}
	texts, err := utils.GetDoc(query)
	if err != nil {
		logs.Errorf("get doc error: %v", err)
		c.JSON(http.StatusForbidden, "get doc error")
		return
	}
	var doc string
	for _, text := range texts {
		doc += text
	}
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	// 进行URL编码
	encodedDoc := strings.ReplaceAll(url.QueryEscape(texts[0]), "+", "%20")
	c.Writer.Header().Set("Doc", encodedDoc)
	c.Writer.(http.Flusher).Flush()
	msg := syntheticDialogue(query, doc, history)
	stream, err := gpt.GptClient.ChatWithStream(openai.GPT3Dot5Turbo, 512, msg)
	if err != nil {
		logs.Errorf("chat with stream error: %v", err)
		c.JSON(http.StatusForbidden, "chat with stream error")
		return
	}
	for {
		if stream == nil {
			continue
		}
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			logs.Errorf("stream.Recv error: %v", err)
			c.JSON(http.StatusForbidden, "stream.Recv error")
			return
		}
		if len(response.Choices) == 0 {
			continue
		}
		content := response.Choices[0].Delta.Content
		// 在 gin 中直接操作 ResponseWriter 来发送数据
		_, err = c.Writer.WriteString(content)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ServerError"})
			return
		}

		// Flush操作确保数据被立即发送
		c.Writer.Flush()
	}
}

func syntheticDialogue(query string, doc string, history []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
	message := history
	message = append(
		message, openai.ChatCompletionMessage{
			Role:    "system",
			Content: fmt.Sprintf("Now you can use this information to answer the user's question. The information provided to you is:\n\"%s\"", doc),
		},
		openai.ChatCompletionMessage{
			Role:    "user",
			Content: query,
		})
	return message
}
