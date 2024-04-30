package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	drant_db "go-gpt-assistant/db/vector-db"
	"go-gpt-assistant/gpt"

	"github.com/tmc/langchaingo/schema"
	mongo_db "go-gpt-assistant/db/mongo-db"

	"rnd-git.valsun.cn/ebike-server/go-common/logs"
)

func DealFile(file multipart.File, filePath string, uploadDir string) error {
	var buf [512]byte
	n, err := io.ReadFull(file, buf[:])
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return fmt.Errorf("read file error: %v", err)
	}
	// 判断上传的文件类型
	contentType := http.DetectContentType(buf[:n])
	if contentType != "application/pdf" {
		return fmt.Errorf("file type error: %v", contentType)
	}
	// 重新定位文件指针到起始位置
	file.Seek(0, io.SeekStart)
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create dir error: %v", err)
	}
	logs.Infof("file path: %s", filePath)
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file error: %v", err)
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		return fmt.Errorf("copy file error: %v", err)
	}
	return nil
}

func AddFileToMongo(file []schema.Document) {
	for _, doc := range file {
		mongoDoc := mongo_db.Document{
			Content:  doc.PageContent,
			MetaData: doc.Metadata,
		}
		err := mongo_db.InsertOne(mongoDoc)
		if err != nil {
			logs.Errorf("MongoDB insert one error: %v", err)
		}
	}
}

func GetDoc(query string) ([]string, error) {
	var doc string
	docs := make([]string, 0)
	keyWords, err := gpt.GptClient.GetKeyWord(openai.GPT3Dot5Turbo, 20, query)
	if err != nil {
		logs.Warnf("get key words error: %v", err)
		doc, err = EmAndSearchDoc(query)
	} else {
		doc, err = EmAndSearchDoc(keyWords)
		logs.Infof("query: %s ,key words: %s", query, keyWords)
	}
	docs = append(docs, doc)
	if err != nil {
		return docs, err
	}
	keyWord := strings.Split(keyWords, ",")
	for _, word := range keyWord {
		text, err := mongo_db.Searchtext(word)
		if err != nil {
			logs.Errorf("Mongo search text error: %v", err)
			return docs, nil
		}
		for _, t := range text {
			docs = append(docs, t.Content)
		}
	}
	return docs, nil
}

func EmAndSearchDoc(query string) (string, error) {
	q_e, err := gpt.GptClient.Embedding([]string{query})
	if err != nil {
		return "", fmt.Errorf("embedding error: %v", err)
	}
	putDoc, err := drant_db.Search(
		"test",
		map[string]interface{}{"exact": false, "hnsw_ef": 256},
		q_e,
		2,
	)
	if err != nil {
		return "", fmt.Errorf("search error: %v", err)
	}
	return putDoc, nil
}
