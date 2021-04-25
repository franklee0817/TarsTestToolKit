package main

import (
	"fmt"
	"os"

	"github.com/TarsCloud/TarsGo/tars"

	"github.com/franklee0817/golang/TestUnits"
)

func main() {
	// Get server config
	cfg := tars.GetServerConfig()

	// New servant imp
	imp := new(GolangImp)
	err := imp.Init()
	if err != nil {
		fmt.Printf("testImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	// New servant
	app := new(TestUnits.Golang)
	// Register Servant
	app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".testObj")

	// Run application
	tars.Run()
}
