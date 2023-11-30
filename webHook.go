package main

import (
	"net/http"
	"webHook/auth"
	"webHook/handle"
)

// webHook使用Api,用于调用grafana发送报警

func main() {
	http.Handle("/webhook/receive", auth.BaseAuth(handle.ALLWebHookHandler{}))
	http.Handle("/webhook/feishu", auth.BaseAuth(handle.FeishuHandler{}))
	http.Handle("/webhook/wechat", auth.BaseAuth(handle.WeChatHandler{}))

	// 启动http服务
	http.ListenAndServe(":8066", nil)
}
