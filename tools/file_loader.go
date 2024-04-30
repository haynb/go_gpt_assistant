package tools

import (
	"context"
	"fmt"
	"os"

	"rnd-git.valsun.cn/ebike-server/go-common/logs"

	"github.com/tmc/langchaingo/schema"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
)

func LoadPdfFile(filePath string) ([]schema.Document, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()
	p := documentloaders.NewPDF(file, fileSize)
	spliter := textsplitter.NewTokenSplitter()
	doc, err := p.LoadAndSplit(context.Background(), spliter)
	if err != nil {
		return nil, err
	}
	logs.Infof("文件已经加载完毕")
	return doc, nil
}
