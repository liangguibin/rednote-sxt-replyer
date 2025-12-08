package service

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/liangguibin/rednote-sxt-replyer/constant"
	"github.com/liangguibin/rednote-sxt-replyer/model"
	"github.com/liangguibin/rednote-sxt-replyer/store"
	"github.com/liangguibin/rednote-sxt-replyer/util"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// CreateHttpHeaders 创建 http 请求头
func CreateHttpHeaders(req *http.Request) {
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-encoding", "gzip, deflate, br, zstd")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("no-cache", "212")
	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("cookie", store.Cookie)
	req.Header.Set("origin", "https://sxt.xiaohongshu.com")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://sxt.xiaohongshu.com/im/multiCustomerService?uba_pre=115.session_manage..1754618064773&uba_ppre=115.multi_customer_service..1754618048985&uba_index=5")
	req.Header.Set("sec-ch-ua", "\"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"138\", \"Google Chrome\";v=\"138\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	req.Header.Set("x-subsystem", "sxt")
}

// InitCookie 初始化 Cookie
func InitCookie(accessToken string) {
	if accessToken == "" {
		// 读取命令行输入
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("请输入 Access Token: ")
		accessToken, _ = reader.ReadString('\n')
		accessToken = strings.TrimSpace(accessToken)
	}
	store.Cookie = "access-token-sxt.xiaohongshu.com=" + accessToken + ";"
}

// GetUserInfo 获取私信用用户数据
func GetUserInfo() {
	req, _ := http.NewRequest("GET", constant.UserInfoUrl, nil)

	CreateHttpHeaders(req)

	resp, err := store.HttpClient.Do(req)
	if err != nil {
		log.Fatal(util.GetTime(), " 获取私信通用户数据失败: ", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode == 401 {
		SendNotice("登录状态失效")
		log.Fatal(util.GetTime(), " 登录状态失效")
	}

	gzipReader, _ := gzip.NewReader(resp.Body)
	defer func(gzipReader *gzip.Reader) {
		_ = gzipReader.Close()
	}(gzipReader)

	body, _ := io.ReadAll(gzipReader)

	var response model.UserInfo
	_ = json.Unmarshal(body, &response)

	store.AccountNo = response.Data.AccountNo
	store.BUserId = response.Data.BUserId
	store.CUserId = response.Data.CUserId
}

// GetFlowUserInfo 获取 Flow 用户信息
func GetFlowUserInfo() {
	fullUrl := constant.FlowUserInfoUrl + "?account_no=" + store.AccountNo + "&contact_way=octopus"
	req, _ := http.NewRequest("GET", fullUrl, nil)

	CreateHttpHeaders(req)

	resp, err := store.HttpClient.Do(req)
	if err != nil {
		log.Fatal(util.GetTime(), " 获取 Flow 用户数据失败: ", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode == 401 {
		SendNotice("登录状态失效")
		log.Fatal(util.GetTime(), " 登录状态失效")
	}

	body, _ := io.ReadAll(resp.Body)

	var response model.FlowUserInfo
	_ = json.Unmarshal(body, &response)

	store.CsProviderId = response.Data.FlowUser.CsProviderId
}

// GetChatInfo 获取会话历史记录 - 首次进线使用
func GetChatInfo(userId string) []map[string]interface{} {
	// 调用接口
	fullUrl := constant.ChatInfoUrl + "?grantor_user_id=&customer_user_id=" + userId + "&limit=10"
	req, _ := http.NewRequest("GET", fullUrl, nil)
	CreateHttpHeaders(req)
	resp, err := store.HttpClient.Do(req)
	if err != nil {
		fmt.Println(util.GetTime(), " 获取会话历史失败: ", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	// 获取响应数据
	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return make([]map[string]interface{}, 0)
	}
	defer func(gzipReader *gzip.Reader) {
		_ = gzipReader.Close()
	}(gzipReader)

	body, _ := io.ReadAll(gzipReader)
	var response map[string]interface{}
	_ = json.Unmarshal(body, &response)

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		fmt.Println(fmt.Println(util.GetTime(), " 会话历史接口返回异常: ", string(body)))
		return make([]map[string]interface{}, 0)
	}
	messages, ok := data["messages"].([]interface{})
	if !ok {
		fmt.Println(fmt.Println(util.GetTime(), " 会话历史接口返回异常: ", string(body)))
		return make([]map[string]interface{}, 0)
	}
	history := make([]map[string]interface{}, 0)
	for _, value := range messages {
		formatVal := value.(map[string]interface{})
		history = append(history, formatVal)
	}
	return history
}

// SendHttpMessage 发送 http 消息
func SendHttpMessage(message string, messageType string, receiverId string) {
	params := map[string]interface{}{
		"c_user_id":       store.CUserId,
		"content":         message,
		"message_type":    messageType,
		"platform":        1,
		"receiver_id":     receiverId,
		"sender_porch_id": store.BUserId,
		"uuid":            "1754631353739-39893284",
	}
	// json 编码
	jsonData, _ := json.Marshal(params)
	// 请求
	req, _ := http.NewRequest("POST", constant.MessageUrl, bytes.NewBuffer(jsonData))
	CreateHttpHeaders(req)
	// 执行请求
	resp, err := store.HttpClient.Do(req)
	if err != nil {
		fmt.Println(util.GetTime(), " 发送消息失败: ", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode == 401 {
		SendNotice("登录状态失效")
		log.Fatal(util.GetTime(), " 登录状态失效")
		return
	}
}
