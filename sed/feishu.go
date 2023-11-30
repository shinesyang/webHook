package sed

import (
	"time"
	"webHook/alarm"
	"webHook/global"
	"webHook/request"

	"github.com/shinesyang/common"
)

func SedAlarmInfoTOFeishu(webHookParam []byte) {
	message, alarmNumber, err := alarm.CleanMessage(webHookParam)
	if err != nil {
		return
	}

	common.Logger.Debugf("告警信息: %v", message)

	common.Logger.Info("发送飞书告警")

	/*默认第一次告警是5分钟,之后每次告警加2分钟*/
	defaultTime := 5
	for i := 1; i <= alarmNumber; i++ {
		requestMessage := make([]map[string]string, 0, len(message))
		// 不存在则表示告警已经恢复了
		for alarmItemMd5, msg := range message {
			_, ok := global.AlarmMap.Load(alarmItemMd5)
			if !ok {
				status := msg["恢复状态"]
				if status == "resolved" {
					requestMessage = append(requestMessage, msg)
					delete(message, alarmItemMd5)
				} else {
					common.Logger.Infof("%s对应的告警已经恢复", msg["告警信息"])
				}
			} else {
				requestMessage = append(requestMessage, msg)
			}
		}

		// message全部清空则取消发送
		if len(requestMessage) <= 0 {
			break
		}

		//marshal, _ := json.Marshal(requestMessage)
		//common.Logger.Info(string(marshal))
		request.RequestFeishu(requestMessage)
		if i == alarmNumber {
			break
		}

		<-time.After(time.Minute * time.Duration(defaultTime+(i-1)*2))
	}
	return

}
