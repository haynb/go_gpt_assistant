package main

import (
	"flag"
	"sync"

	mongo_db "go-gpt-assistant/db/mongo-db"

	drant_db "go-gpt-assistant/db/vector-db"

	"go-gpt-assistant/handler"

	"go-gpt-assistant/gpt"

	"go-gpt-assistant/config"

	"rnd-git.valsun.cn/ebike-server/go-common/logs"
	comSvc "rnd-git.valsun.cn/ebike-server/go-common/server"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "./conf/app.yaml", "service configure file")
	flag.Parse()

	// do system initial
	if err := config.LoadFromFile(configFile); err != nil {
		panic(err)
	}

	// get configure
	cfg := config.GetAppConf()
	if err := cfg.PrepareEnv(); err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	// with grageful exit
	go comSvc.GracefulExit(wg)

	// start the Gpt system
	gpt.InitGpt()
	drant_db.InitQdrant()
	mongo_db.InitMongo()
	//// start grpc server
	//go server.StartGrpcServer(wg, cfg.GrpcAddress, handler.AddGrpcService)
	//
	// start the http server
	go comSvc.StartGinServer(wg, cfg.ServerAddress, handler.AddRouter, cfg.EnableDebug)

	// wait all exit
	wg.Wait()
	logs.Infof("main exit")
}
