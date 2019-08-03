package main

import (
	"flag"
	"fmt"
	"github.com/sishen007/gocrontab/master"
	"runtime"
	"time"
)

var (
	confFile string // 配置文件路径
)

// 解析命令行参数
func initArgs() {
	flag.StringVar(&confFile, "config", "D:/data/Go/uxinFort/pro1/src/github.com/sishen007/gocrontab/master/main/master.json", "指定master.json")
	flag.Parse()
}

// 初始化线程数量
func initEnv() {
	// 设置线程数
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	var (
		err error
	)
	// 初始化命令行参数
	initArgs()

	// 初始化线程
	initEnv()

	// 初始化配置
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	// 启动服务发现管理器
	if err = master.InitWorkerMgr(); err != nil {
		goto ERR
	}

	// 启动日志管理器
	if err = master.InitLogMgr(); err != nil {
		goto ERR
	}

	// 启动任务管理器
	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	// 启动Api Http服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}
	// 正常退出
	for {
		time.Sleep(1 * time.Second)
	}
	return

ERR:
	fmt.Println(err)
}
