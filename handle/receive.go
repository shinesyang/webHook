package handle

import (
	"io/ioutil"
	"net/http"
	"webHook/sed"

	"github.com/shinesyang/common"
)

type ALLWebHookHandler struct {
}

func (n ALLWebHookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	webHookParam, _ := ioutil.ReadAll(r.Body)
	common.Logger.Debugf("receive获取webHook结果: %s", string(webHookParam))
	n.feiShuServeHTTP(webHookParam)
	n.weChatServeHTTP(webHookParam)
	w.Write([]byte("3"))
	return
}

func (n ALLWebHookHandler) feiShuServeHTTP(webHookParam []byte) {
	go sed.SedAlarmInfoTOFeishu(webHookParam) // 异步处理报警发送到飞书
	return
}

func (n ALLWebHookHandler) weChatServeHTTP(webHookParam []byte) {
	go sed.SedAlarmInfoTOWeChat(webHookParam) // 异步处理报警发送到微信
	return
}
