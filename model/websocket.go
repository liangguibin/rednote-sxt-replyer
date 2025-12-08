package model

// WebSocketAdditionalInfo WebSocket Additional Info
type WebSocketAdditionalInfo struct {
	SellerId string `json:"sellerId,omitempty"`
	UserId   string `json:"userId,omitempty"`
}

// WebSocketUserAgent WebSocket User Agent
type WebSocketUserAgent struct {
	AppName    string `json:"appName,omitempty"`
	AppVersion string `json:"appVersion,omitempty"`
}

// WebSocketInfo WebSocket 客户端信息
type WebSocketInfo struct {
	AdditionalInfo WebSocketAdditionalInfo `json:"additionalInfo,omitempty"`
	UserAgent      WebSocketUserAgent      `json:"userAgent,omitempty"`
}

// WebSocketMessage WebSocket 消息 - 仅包含常用字段
type WebSocketMessage struct {
	AppId     string        `json:"appId,omitempty"`
	Token     string        `json:"token,omitempty"`
	Type      int           `json:"type,omitempty"`
	Seq       int           `json:"seq,omitempty"`
	Ack       int           `json:"ack,omitempty"`
	Topic     string        `json:"topic,omitempty"`
	SecureKey string        `json:"secureKey,omitempty"`
	Encrypt   bool          `json:"encrypt,omitempty"`
	Info      WebSocketInfo `json:"data,omitempty"`
}
