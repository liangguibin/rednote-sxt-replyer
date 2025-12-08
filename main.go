package main

import (
	"fmt"
	"github.com/liangguibin/rednote-sxt-replyer/cmd/version"
	"github.com/liangguibin/rednote-sxt-replyer/service"
	"github.com/spf13/cobra"
	"os"
)

// 入口方法
func main() {
	// 根命令
	var rootCmd = &cobra.Command{
		Use:   "replyer",
		Short: "rednote sxt replyer",
		Long:  `rednote sxt automatic reply toolkit`,
		Run: func(cmd *cobra.Command, args []string) {
			// 初始化
			accessToken, _ := cmd.Flags().GetString("cookie")
			service.InitCookie(accessToken)
			// 初始化日志组件
			service.InitLogger()
			// 获取私信通用户数据
			service.GetUserInfo()
			// 获取 Flow 用户数据
			service.GetFlowUserInfo()
			// 初始化数据库
			service.InitDb()
			defer service.CloseDb()
			// 初始化 WebSocket
			service.InitWebSocket()
			defer service.CloseWebSocket()
			// 启用协程监听消息
			service.OnMessage()
			// 发送认证消息
			service.Authentication()
			// 启动心跳机制 - 客户端
			heartbeatTicker := service.StartHeartbeat()
			defer heartbeatTicker.Stop()
			// 阻塞主线程
			select {}
		},
	}
	// 添加 flag
	rootCmd.Flags().StringP("cookie", "c", "", "set cookie")
	// 添加子命令
	rootCmd.AddCommand(version.RootCmd)
	// 执行
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
