package service

import (
	"encoding/json"
	"github.com/liangguibin/rednote-sxt-replyer/constant"
	"github.com/liangguibin/rednote-sxt-replyer/store"
	"strings"
)

// HandleMessage 处理消息
func HandleMessage(msg []byte) {
	// 解析消息
	var message map[string]interface{}
	_ = json.Unmarshal(msg, &message)
	// 数据字段
	data, ok := message["data"].(map[string]interface{})
	if !ok {
		return
	}
	// 推送类型
	pushType, ok := data["type"].(string)
	if !ok {
		return
	}
	// 只处理 PUSH_SIXINTONG_CHAT 和 PUSH_SIXINTONG_MSG 消息
	if pushType != "PUSH_SIXINTONG_CHAT" && pushType != "PUSH_SIXINTONG_MSG" {
		return
	}
	// payload
	payload, ok := data["payload"].(map[string]interface{})
	if !ok {
		return
	}
	// sixin_chat
	siXinChat, ok := payload["sixin_chat"].(map[string]interface{})
	if !ok {
		return
	}
	// 用户 ID、用户昵称、消息类型、消息内容、发送人 ID
	var visitorId, nickname, messageType, content, senderId string
	// 获取数据
	nickname = siXinChat["nickname"].(string)
	if pushType == "PUSH_SIXINTONG_CHAT" {
		// 首次进线
		// session
		session, ok := payload["session"].(map[string]interface{})
		if ok {
			state := session["state"].(string)
			if state != "PROCESSING" {
				// 关闭会话 等动作 不需要回复
				return
			}
		}

		visitorId = siXinChat["user_id"].(string)
		remoteHistory := GetChatInfo(visitorId)
		if len(remoteHistory) == 0 {
			return
		}
		lastMessage := remoteHistory[0]
		siXinMessage, ok := lastMessage["sixin_message"].(map[string]interface{})
		if !ok {
			return
		}
		messageType = siXinMessage["message_type"].(string)
		content = siXinMessage["content"].(string)
		senderId = siXinMessage["sender_id"].(string)
	} else {
		// 非首次进线
		siXinMessage, ok := payload["sixin_message"].(map[string]interface{})
		if !ok {
			return
		}
		visitorId = payload["visitor_id"].(string)
		messageType = siXinMessage["message_type"].(string)
		content = siXinMessage["content"].(string)
		senderId = siXinMessage["sender_id"].(string)
	}
	// 处理消息
	if senderId == store.CUserId {
		_, _ = Insert(visitorId, constant.Two, messageType, content)
		store.Logger.Info("发送对话消息: ", visitorId, " ", nickname, " ", messageType, " ", content)
	} else {
		// 系统消息 不需要回复 - 温馨提示、关注前限制、点击卡片等
		// if messageType == "HINT" {
		// 	return
		// }
		// 系统消息 不需要回复 - 对方通过某某笔记进入等
		// if messageType == "RICH_HINT" {
		// 	return
		// }

		// HINT 类型消息 只回复 - 点击卡片
		if messageType == "HINT" {
			hintMessage := make(map[string]interface{})
			_ = json.Unmarshal([]byte(content), &hintMessage)
			content = hintMessage["content"].(string)

			if content != "对方已点击你的企业微信联系卡" {
				return
			}
		}

		// RICH_HINT 类型消息 全部回复
		if messageType == "RICH_HINT" {
			// RICH_HINT 消息会发送两遍 只回复 PUSH_SIXINTONG_CHAT
			if pushType == "PUSH_SIXINTONG_MSG" {
				return
			}
			richHintMessage := make(map[string]interface{})
			_ = json.Unmarshal([]byte(content), &richHintMessage)

			richHintContent, ok := richHintMessage["content"].(string)
			if !ok {
				return
			}
			replaceLink, ok := richHintMessage["replaceLink"].([]interface{})
			if !ok {
				return
			}
			replaceLinkIns, ok := replaceLink[0].(map[string]interface{})
			if !ok {
				return
			}
			replaceLinkName, ok := replaceLinkIns["name"].(string)
			if !ok {
				return
			}
			content = strings.Replace(richHintContent, "/@/#/", replaceLinkName, -1)
		}

		_, _ = Insert(visitorId, constant.One, messageType, content)
		store.Logger.Info("收到对话消息: ", visitorId, " ", nickname, " ", messageType, " ", content)
		// 自动回复
		history, _ := Read(visitorId)
		Chat(visitorId, history)
	}
}
