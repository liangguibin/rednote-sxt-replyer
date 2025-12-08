package service

import (
	"context"
	"fmt"
	"github.com/liangguibin/rednote-sxt-replyer/constant"
	"github.com/liangguibin/rednote-sxt-replyer/model"
	"github.com/liangguibin/rednote-sxt-replyer/prompt"
	"github.com/liangguibin/rednote-sxt-replyer/util"
	"github.com/sashabaranov/go-openai"
)

// Chat AI 对话
func Chat(userId string, history []model.Message) {
	go func() {
		// 聊天记录
		messages := make([]openai.ChatCompletionMessage, 0)
		systemMessage := openai.ChatCompletionMessage{
			Role:    "system",
			Content: prompt.SystemPrompt,
		}
		messages = append(messages, systemMessage)
		for _, value := range history {
			var role string
			if value.ChatType == 1 {
				role = "user"
			} else {
				role = "assistant"
			}
			message := openai.ChatCompletionMessage{
				Role:    role,
				Content: value.Content,
			}
			messages = append(messages, message)
		}
		// AI 配置
		config := openai.DefaultConfig(constant.ApiKey)
		config.BaseURL = constant.BaseUrl

		client := openai.NewClientWithConfig(config)

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    constant.Model,
				Messages: messages,
			},
		)

		if err != nil {
			fmt.Println(util.GetTime(), " AI 调用失败: ", err)
			reply := "感谢您的耐心等待，我们会尽快给您答复。"
			SendHttpMessage(reply, "TEXT", userId)
		} else {
			reply := resp.Choices[0].Message.Content
			SendHttpMessage(reply, "TEXT", userId)
		}
	}()
}
