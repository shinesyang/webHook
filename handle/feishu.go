package handle

import (
	"io/ioutil"
	"net/http"
	"webHook/sed"

	"github.com/shinesyang/common"
)

// 飞书报警

type FeishuHandler struct {
}

func (f FeishuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	webHookParam, _ := ioutil.ReadAll(r.Body)
	//webHookParam, err := ioutil.ReadFile("test.json")
	//if err != nil {
	//	common.Logger.Errorf("读取文件失败: %v", err)
	//	return
	//}

	common.Logger.Debugf("飞书获取webHook结果: %s", string(webHookParam))
	go sed.SedAlarmInfoTOFeishu(webHookParam) // 异步处理报警发送到飞书
	w.Write([]byte("1"))
	return
}
