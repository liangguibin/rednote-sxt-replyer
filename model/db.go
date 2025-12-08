package model

// Message 数据库消息实体
type Message struct {
	Id          int64  `json:"id,omitempty"`
	UserId      string `json:"userId,omitempty"`
	ChatType    int    `json:"chatType,omitempty"`
	MessageType string `json:"messageType,omitempty"`
	Content     string `json:"content,omitempty"`
	CreateTime  string `json:"createTime,omitempty"`
}
