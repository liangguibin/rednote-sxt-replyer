package service

import (
	"database/sql"
	"fmt"
	"github.com/liangguibin/rednote-sxt-replyer/model"
	"github.com/liangguibin/rednote-sxt-replyer/store"
	"github.com/liangguibin/rednote-sxt-replyer/util"
	"log"
	_ "modernc.org/sqlite"
)

// InitDb 初始化数据库
func InitDb() {
	// 打开数据库连接池 如果数据库文件不存在则创建
	db, err := sql.Open("sqlite", "./database/replyer.db?_journal_mode=WAL&_busy_timeout=-1&_synchronous=NORMAL&_cache_size=1000")
	if err != nil {
		log.Fatal(util.GetTime(), " 初始化数据库失败: ", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(5)

	fmt.Println(util.GetTime(), " 数据库连接成功")
	// 创建表 如果存在则什么也不做
	// chat_type 1 客户消息 2 回复消息
	query := `
		CREATE TABLE IF NOT EXISTS message (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			chat_type  INTEGER NOT NULL,
			message_type TEXT NOT NULL,
			content TEXT NOT NULL,
			create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(util.GetTime(), " 表结构创建失败: ", err)
	}
	fmt.Println(util.GetTime(), " 表结构创建成功")
	// 保持连接
	store.DbConn = db
}

// CloseDb 关闭数据库
func CloseDb() {
	if store.DbConn != nil {
		err := store.DbConn.Close()
		if err != nil {
			fmt.Println(util.GetTime(), " 数据库连接关闭失败: ", err)
		} else {
			fmt.Println(util.GetTime(), " 数据库连接已关闭")
		}
	}
}

// Insert 插入消息记录
func Insert(userId string, chatType int, messageType string, content string) (int64, error) {
	res, err := store.DbConn.Exec(`
		INSERT INTO message (user_id, chat_type, message_type, content) VALUES (?, ?, ?, ?)
	`, userId, chatType, messageType, content)
	if err != nil {
		fmt.Println(util.GetTime(), " 插入消息记录失败: ", err)
		return 0, err
	}

	return res.LastInsertId()
}

// Read 查询消息记录
func Read(userId string) ([]model.Message, error) {
	rows, err := store.DbConn.Query(`
		SELECT id, user_id, chat_type, message_type, content, datetime(create_time, 'localtime') AS create_time FROM message WHERE user_id = ? ORDER BY create_time
	`, userId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(util.GetTime(), " 消息记录查询失败: ", err)
		}
	}(rows)

	var messages []model.Message
	for rows.Next() {
		var message model.Message
		err = rows.Scan(&message.Id, &message.UserId, &message.ChatType, &message.MessageType, &message.Content, &message.CreateTime)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
