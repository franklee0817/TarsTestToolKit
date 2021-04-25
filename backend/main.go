package main

import (
	"fmt"
	"os"

	"github.com/TarsCloud/TarsGo/tars"

	"github.com/franklee0817/t3k/backend/controllers"
	"github.com/franklee0817/t3k/backend/impl"
	"github.com/franklee0817/t3k/backend/tars-protocol/apitars"

	_ "github.com/franklee0817/t3k/backend/config"
)

func main() {
	serveHttp()
}

func serveHttp() {
	mux := &tars.TarsHttpMux{}
	// New servant
	mux.HandleFunc("/testFunc", controllers.HandleDoFuncTest)        // 执行功能测试
	mux.HandleFunc("/testPerf", controllers.HandleDoPerfTest)        // 执行性能测试
	mux.HandleFunc("/detail", controllers.HandleGetTestDetail)       // 获取测试详情
	mux.HandleFunc("/histories", controllers.HandleGetTestHistories) // 获取历史测试列表

	cfg := tars.GetServerConfig()
	tars.AddHttpServant(mux, cfg.App+"."+cfg.Server+".apiObj") //Register http server
	tars.Run()
}

func serveTars() {
	// Get server config
	cfg := tars.GetServerConfig()

	// New servant imp
	imp := new(impl.APIImpl)
	err := imp.Init()
	if err != nil {
		fmt.Printf("apiImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	// New servant
	app := new(apitars.Api)
	// Register Servant
	app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".apiObj")

	// Run application
	tars.Run()
}
