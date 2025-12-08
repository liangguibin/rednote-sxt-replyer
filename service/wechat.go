package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/liangguibin/rednote-sxt-replyer/constant"
	"github.com/liangguibin/rednote-sxt-replyer/store"
	"github.com/liangguibin/rednote-sxt-replyer/util"
	"io"
	"net/http"
)

// SendNotice 发送企业微信通知
func SendNotice(msg string) {
	if constant.WechatNoticeKey != "" {
		// 参数
		atList := make([]string, 0)
		atList = append(atList, "@all")
		text := map[string]interface{}{
			"content":        msg,
			"mentioned_list": atList,
		}
		params := map[string]interface{}{
			"msgtype": "text",
			"text":    text,
		}
		// URL
		fullUrl := constant.WechatNoticeUrl + constant.WechatNoticeKey
		// json 编码
		jsonData, _ := json.Marshal(params)
		// 请求
		req, _ := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		// 执行请求
		resp, err := store.HttpClient.Do(req)
		if err != nil {
			fmt.Println(util.GetTime(), " 企业微信通知发送失败: ", err)
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
	}
}
