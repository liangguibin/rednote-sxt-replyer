## 小红书-私信通-AI自动回复工具

### 技术栈

- modernc.org/sqlite
- github.com/spf13/cobra
- github.com/gorilla/websocket
- github.com/sashabaranov/go-openai
- github.com/natefinch/lumberjack
- github.com/sirupsen/logrus

### 线程模型

- 主线程
  - 心跳机制线程
  - 消息监听线程
    - AI交互线程 1
    - AI交互线程 2
    - AI交互线程 3
    - AI交互线程 N

### 修改配置

- constant/ai.go - AI大模型相关配置
- constant/wechat.go - 企业微信群机器人配置（可选，用于发送登录状态过期提醒）

### 构建

- Windows

`GOOS=windows GOARCH=amd64 go build -o replyer.exe main.go`

- macOS

`GOOS=darwin GOARCH=amd64 go build -o replyer main.go`

- Linux

`GOOS=linux GOARCH=amd64 go build -o replyer main.go`

### 使用

- 命令行环境-直接运行

`./replyer`

- 命令行环境-指定 AccessToken

`./replyer -c "<AccessToken>"`

### AccessToken 获取

- 登录私信通网页版 https://sxt.xiaohongshu.com/im
- 打开 开发者工具
- 在 Application - Cookies 中找到 access-token* 开头的项目
- 复制其值得到 Access Token

### 注意
- 部署在 Windows 系统时，需要彻底关闭或卸载 Microsoft Defender 和其他杀毒软件

### 声明

- 本项目仅出于学习或研究为目的，请勿用作非法用途。任何法律风险及责任，需使用者自行承担，与本项目无关。

