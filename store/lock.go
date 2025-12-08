package store

import "sync"

// WsWriteMutex WebSocket 写锁
var WsWriteMutex sync.Mutex
