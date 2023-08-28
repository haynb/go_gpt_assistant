package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
	"io"
	"kefu/db"
	"kefu/gpt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var (
	uploadDir = "upload"
	// 一个锁
	lock sync.Mutex
)

func writeStringJson(w http.ResponseWriter, statusCode int, message string) error {
	jsStr, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(jsStr)
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return err
	}
	return nil
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	type errorResponse struct {
		Message string `json:"message"`
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errRes := errorResponse{
		Message: message,
	}
	jsonData, _ := json.Marshal(errRes)
	w.Write(jsonData)
}
func upload(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Printf("%s --upload:  ", r.RemoteAddr)
	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSONError(w, http.StatusForbidden, "ServerError")
		fmt.Println(err)
		return
	}
	defer file.Close()
	var buf [512]byte
	n, err := io.ReadFull(file, buf[:])
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		fmt.Println(err)
		return
	}

	contentType := http.DetectContentType(buf[:n])
	if contentType != "application/pdf" {
		writeJSONError(w, http.StatusBadRequest, "Only PDF files are allowed")
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return
	}
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return
	}
	filePath := filepath.Join(uploadDir, header.Filename)
	fmt.Println(filePath)
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return
	}
	lock.Lock()
	err = loadAndMakeFile(filePath)
	lock.Unlock()
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "pdf broken,change a RIGHT pdf file")
		removeErr := os.Remove(filePath)
		if removeErr != nil {
			fmt.Println(removeErr)
		}
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return
	}
	err = writeStringJson(w, http.StatusOK, "upload success")
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return
	}
}
func loadAndMakeFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
		return err
	}
	fileSize := fileInfo.Size()
	p := documentloaders.NewPDF(file, fileSize)
	spliter := textsplitter.NewTokenSplitter()
	doc, err := p.LoadAndSplit(context.Background(), spliter)
	if err != nil {
		return err
	}
	fmt.Println("文件已经加载完毕")
	e := gpt.NewEm()
	points, err := db.DocToPoints(doc, e)
	if err != nil {
		return err
	}
	_, err = db.AddPoints("test", points)
	if err != nil {
		return err
	}
	return nil
}
func checkFile(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println(r.RemoteAddr, "---:checkFile")
	file_list := make([]string, 0)
	err := filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if path == uploadDir {
			return nil
		}
		file_list = append(file_list, path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return
	}
	w.WriteHeader(http.StatusOK) // 返回状态码 200
	jsonData, err := json.Marshal(file_list)
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
	return
}
func chatWithFile(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := r.FormValue("query")
	historyJson := r.FormValue("history")
	var history []gpt.Gpt3Dot5Message
	err := json.Unmarshal([]byte(historyJson), &history)
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return
	}
	tools := gpt.NewChatGptTool()
	keyWords, err := tools.GetKeyWord(query)
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return
	}
	e := gpt.NewEm()
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return
	}
	q_e, err := e.EmbedQuery(keyWords)
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return
	}
	put_doc, err := db.Search(
		"test",
		map[string]interface{}{"exact": false, "hnsw_ef": 256},
		q_e,
		3)
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return
	}
	stream, err := tools.Ask(query, put_doc, history)
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	fmt.Println(r.RemoteAddr, "----ASK:   ", query, "---keyWords:   ", keyWords)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			writeJSONError(w, http.StatusInternalServerError, "ServerError")
			return
		}
		content := msg.Choices[0].Delta.Content
		fmt.Fprintf(w, "%s", content)
		flusher.Flush()
	}
	putDocJSON, err := json.Marshal(put_doc)
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return
	}
	w.Write(putDocJSON)
}

func message(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println(r.RemoteAddr, "----GetMessage")
	msg := []string{
		"本系统为智能客服测试版，不代表最终质量~",
		"基本功能：上传pdf文件，系统会使用人工智能ai自动提取关键词，然后根据关键词智能搜索文档的相关信息，根据信息再使用人工智能回答相关问题的答案~",
		"强烈建议上传文字较多的pdf文件，否则关键词提取效果不好，搜索结果也可能不好~",
		"暂时关闭删除文件功能，防止误删~",
		"可以点击ai回答上方的按钮显示后台查询到的相关文档，鼠标移上去就会出现哦，可能有乱码很正常，人工智能在回答的时候会自动处理。pdf余额符合规范，乱码越少~",
		"上传文件的时候会在后台进行处理，处理时间可能较长，请尽量避免同时上传文件哦，处理完毕后会自动显示在下方的文件列表中，后续应该会更新改进的哦，请耐心等待~",
		"欢迎大家体验评测，前端反馈请企微联系刘宽同学，后端反馈请企微联系何安阳同学，谢谢~",
	}
	msgJson, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		writeJSONError(w, http.StatusInternalServerError, "ServerError")
		return
	}
	w.Write(msgJson)
}
func allowCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	router := httprouter.New()
	router.POST("/upload", upload)
	router.GET("/check_file", checkFile)
	router.POST("/chat_with_file", chatWithFile)
	router.GET("/message", message)
	handler := allowCORS(router)
	fmt.Println("服务已经启动")
	err := http.ListenAndServe(":9964", handler)
	if err != nil {
		panic(err)
	}
}
