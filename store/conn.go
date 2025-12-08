package store

import (
	"database/sql"
	"github.com/gorilla/websocket"
	"net/http"
)

// HttpClient HTTP Client
var HttpClient = &http.Client{}

// WsConn WebSocket 连接
var WsConn *websocket.Conn

// DbConn 数据库连接
var DbConn *sql.DB
