package main

import (
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"server/conf"
	"server/game"
	"server/gate"
	"server/login"
)

func main() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath

	//go install server 构建
	//这里按顺序注册了 game、gate、login 三个模块。每个模块都需要实现接口：
	// Leaf 首先会在同一个 goroutine 中按模块注册顺序执行模块的 OnInit 方法，等到所有模块 OnInit 方法执行完成后
	// 则为每一个模块启动一个 goroutine 并执行模块的 Run 方法。最后，游戏服务器关闭时（Ctrl + C 关闭游戏服务器）
	// 将按模块注册相反顺序在同一个 goroutine 中执行模块的 OnDestroy 方法。
	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
