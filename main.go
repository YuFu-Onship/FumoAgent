package main

import (
	_ "embed"
	"myapp/src/Ctrl"
	"myapp/src/config"
	"myapp/src/model"
	"os"
	"runtime/debug"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// 初始化
	debug.SetGCPercent(40)
	config.Error_Init()

	trunk := model.New_Trunk()
	trunk.RootPath = exePath

	app := Ctrl.New_AppCtrl(trunk)
	app.Run()
}
