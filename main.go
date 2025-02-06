package main

import (
	"second_hand_mall/config"
	"second_hand_mall/internal/core"
	"second_hand_mall/internal/core/initialize"
	"second_hand_mall/internal/core/initialize/db"
	"second_hand_mall/internal/global"
	"sync"
)

func main() {
	// 初始化配置文件
	global.GVAL_CONFIG = config.LoadLocalConfig("./config.yaml")

	// 初始化日志
	global.GVAL_LOG = initialize.InitZap()

	// 初始化数据库
	if err := db.InitEngine(global.GVAL_CONFIG.Db); err != nil {
		global.GVAL_LOG.DPanic("初始化数据库失败:" + err.Error())
		return
	}

	global.GVAL_LOG.Debug("测试日志")

	var wg sync.WaitGroup
	wg.Add(2)
	// 启动http服务
	go core.RunHttpServer(&wg)
	// 启动Rpc服务
	go core.RunRpcServer(&wg)
	wg.Wait()
}
