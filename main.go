package main

import (
	"context"
	"fmt"
	"io"
	"qdrant/db"
	"qdrant/gpt"
)

func main() {
	//file, err := os.Open("666.pdf")
	//if err != nil {
	//	panic(err)
	//}
	//// 获取文件大小
	//fileInfo, err := file.Stat()
	//if err != nil {
	//	panic(err)
	//	return
	//}
	//fileSize := fileInfo.Size()
	//p := documentloaders.NewPDF(file, fileSize)
	//spliter := textsplitter.NewTokenSplitter()
	//doc, err := p.LoadAndSplit(context.Background(), spliter)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("文件已经加载完毕")
	e, err := gpt.NewEm()
	if err != nil {
		panic(err)
	}
	//points, err := db.DocToPoints(doc, e)
	//if err != nil {
	//	panic(err)
	//}
	//err = db.DeleteCollection("test")
	//if err != nil {
	//	panic(err)
	//}
	//_, err = db.CreateCollection("test")
	//if err != nil {
	//	panic(err)
	//}
	//_, err = db.AddPoints("test", points)
	//if err != nil {
	//	panic(err)
	//}
	query := "赛维周边环境怎么样？"
	q_e, err := e.EmbedQuery(context.Background(), query)
	if err != nil {
		panic(err)
	}
	put_doc, err := db.Search(
		"test",
		map[string]interface{}{"exact": false, "hnsw_ef": 256},
		q_e,
		3)
	if err != nil {
		panic(err)
	}
	tools := gpt.NewChatGptTool()
	fmt.Println("GPT思考中...")
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
}
