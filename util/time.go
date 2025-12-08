package util

import (
	"strconv"
	"time"
)

// GetTimestamp 获取时间戳 - 13 位
func GetTimestamp() string {
	now := time.Now()
	timestamp := now.UnixNano() / int64(time.Millisecond)
	return strconv.FormatInt(timestamp, 10)
}

// GetTime 获取格式化时间
func GetTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
