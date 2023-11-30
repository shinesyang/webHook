package handle

import (
	"io/ioutil"
	"net/http"
	"webHook/sed"

	"github.com/shinesyang/common"
)

// 微信报警

type WeChatHandler struct {
}

func (e WeChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	webHookParam, _ := ioutil.ReadAll(r.Body)
	//webHookParam, err := ioutil.ReadFile("test.json")
	//if err != nil {
	//	common.Logger.Errorf("读取文件失败: %v", err)
	//	return
	//}

	common.Logger.Debugf("wechat获取webHook结果: %s", string(webHookParam))
	go sed.SedAlarmInfoTOWeChat(webHookParam) // 异步处理报警发送到微信
	w.Write([]byte("2"))
	return
}
