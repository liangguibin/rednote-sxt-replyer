package service

import (
	"github.com/liangguibin/rednote-sxt-replyer/store"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

// InitLogger 初始化日志组件
func InitLogger() {
	logger := logrus.New()

	logger.SetOutput(&lumberjack.Logger{
		Filename:   "./log/app.log", // 日志路径及文件名
		MaxSize:    2,               // 每个文件的最大大小（MB）
		MaxBackups: 10,              // 保留最近 10 个日志文件
		MaxAge:     30,              // 日志最多保留 30 天
		Compress:   false,           // 是否压缩旧日志
	})

	logger.SetFormatter(&logrus.TextFormatter{ // 设置日志格式
		DisableColors: true,
		FullTimestamp: true,
	})

	store.Logger = logger
}
