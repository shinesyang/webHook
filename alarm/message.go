package alarm

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"webHook/global"
	"webHook/parse"

	"github.com/shinesyang/common"
)

// 处理grafana报警信息

func CleanMessage(webHookParam []byte) (map[string]map[string]string, int, error) {
	dataMap := make(map[string]interface{}, 100)
	err := json.Unmarshal(webHookParam, &dataMap)
	if err != nil {
		common.Logger.Errorf("解析参数文件失败: %v", err)
		return nil, 0, err
	}

	alerts := dataMap["alerts"]

	alertsList := make([]parse.Alert, 0, len(alerts.([]interface{}))) // alertsList的长度等于断言成切片的alerts长度

	// alerts转成[]byte再转成[]parse.Alert
	marshal, _ := json.Marshal(alerts)
	err = json.Unmarshal(marshal, &alertsList)
	if err != nil {
		common.Logger.Errorf("解析alerts失败: %v", err)
		return nil, 0, err
	}

	WarnDataAll := make(map[string]map[string]string, len(alertsList))
	alarmNumber := 1 // 默认告警次数
	for _, alert := range alertsList {
		var service string
		hostName := alert.Labels.HostName
		if hostName == "" {
			service = alert.Labels.Instance
		} else {
			service = hostName
		}

		/*不同的主题提示以及根据不同的告警级别触发告警次数*/
		level := alert.Annotations.Level
		if level == "" {
			level = "Default"
		}
		var theme string
		if level == "Average" {
			theme = alert.Labels.AlertName + "已经无法提供服务告警,请立即处理，请立即处理"
			alarmNumber = 5
		} else if level == "Severity" {
			theme = alert.Labels.AlertName + "严重告警,请尽快处理"
			alarmNumber = 3
		} else if level == "Warning" {
			theme = alert.Labels.AlertName + "已经超出预定阀值,请及时关注处理"
		} else {
			theme = alert.Labels.AlertName + "告警提示,请注意"
		}

		// 恢复告警不管什么级别都是告警一次
		if alert.Status == "resolved" {
			alarmNumber = 1
		}

		// 当前告警时的值
		values := alert.Values
		valueString := alert.ValueString
		nowWarningValue := CleanValues(values, valueString)
		var WarnData map[string]string

		// 处理告警时间
		startsAt := alert.StartsAt
		timeSplit1 := strings.Split(startsAt, "T")
		timeJoin := strings.Join(timeSplit1, " ")
		timeSplit2 := strings.Split(timeJoin, "+")
		var dateTime string
		if len(timeSplit2) == 2 {
			dateTime = timeSplit2[0]
		} else {
			dateTime = time.Now().Format("2006-01-02 15:04:05")
		}

		// 区分不同的报警信息
		if alert.Status == "resolved" {
			WarnData = map[string]string{
				"恢复主题": alert.Labels.AlertName + "当前已经恢复,请注意查看",
				"恢复类型": alert.Labels.AlertName,
				"恢复级别": level,
				"恢复主机": service,
				"恢复信息": service + " >>>> " + alert.Annotations.Description,
				"恢复时间": dateTime,
				"恢复状态": alert.Status,
				"恢复说明": "阀值已经恢复,请注意查看",
				//"恢复面板":  alert.PanelURL,
				//"当前值": nowWarningValue,.
			}
		} else {
			WarnData = map[string]string{
				"告警主题": theme,
				"告警类型": alert.Labels.AlertName,
				"告警级别": level,
				"告警主机": service,
				"告警信息": service + " >>>> " + alert.Annotations.Description,
				"告警时间": dateTime,
				"告警状态": alert.Status,
				//"告警面板":  alert.PanelURL,
				//"当前值": nowWarningValue,
			}
		}

		panelURL := alert.PanelURL
		if panelURL != "" {
			// 替换掉panelURL的前缀url
			panelURLSplit := strings.Split(panelURL, "//")
			uRLTwo := panelURLSplit[1]
			uRLTwoSplitAfterN := strings.SplitN(uRLTwo, "/", 2)
			router := uRLTwoSplitAfterN[1]
			newUrl := parse.CONF.Panel + router
			// 根据不同的alert组装不同的panelURL
			job := alert.Labels.Job
			project := alert.Labels.Project
			instance := alert.Labels.Instance
			params := make([]string, 0, 3)
			if job != "" {
				params = append(params, fmt.Sprintf("var-job=%s", job))
			}
			if project != "" {
				params = append(params, fmt.Sprintf("var-project=%s", project))
			}
			if instance != "" {
				params = append(params, fmt.Sprintf("var-instance=%s", instance))
			}

			conParams := strings.Join(params, "&")
			common.Logger.Debugf("对应可以访问的panelURL: %s", newUrl+"&"+conParams)

			if alert.Status == "resolved" {
				WarnData["恢复面板"] = newUrl + "&" + conParams
			} else {
				WarnData["告警面板"] = newUrl + "&" + conParams
			}

		}
		if nowWarningValue != "" {
			WarnData["当前值"] = nowWarningValue
		}

		project := alert.Labels.Project
		if project != "" && alert.Status == "resolved" {
			WarnData["恢复项目"] = project
		} else if project != "" && alert.Status != "resolved" {
			WarnData["告警项目"] = project
		}

		/* 根据alert.Labels.Instance/alert.Labels.AlertName/alert.Labels.Project添加一个告警map
		这个告警map，在触发告警次数时先判断map是否存在值(该告警是否恢复，恢复则删除map key),不存在则告警恢复
		*/
		alarmItemMd5 := CreateMd5(alert.Labels.Instance, alert.Labels.AlertName, alert.Labels.Project)
		if alert.Status == "resolved" {
			_, ok := global.AlarmMap.Load(alarmItemMd5)
			if ok {
				global.AlarmMap.Delete(alarmItemMd5)
			}
		} else {
			global.AlarmMap.Store(alarmItemMd5, alarmNumber)
		}

		WarnDataAll[alarmItemMd5] = WarnData
	}

	return WarnDataAll, alarmNumber, nil
}

// values取值处理

func CleanValues(values interface{}, valueString string) string {
	// 新方式在grafana报警的alertrules ：Set a query and alert condition 设置固定的 Add expression 的值为: (lastValue)
	nowValue := "0"
	if values != nil {
		valuesMap := values.(map[string]interface{})
		lastValue, ok := valuesMap["lastValue"]
		if ok {
			nowValue = fmt.Sprintf("%.2f", lastValue)
			return nowValue
		}
	} else {
		valueStringList := strings.Split(valueString, "],")
		for _, value := range valueStringList {
			if strings.Contains(value, "lastValue") {
				valueSplit := strings.Split(value, "value=")
				if len(valueSplit) == 2 {
					i := strings.TrimSpace(valueSplit[1])
					float, err := strconv.ParseFloat(i, 64)
					if err == nil {
						nowValue = fmt.Sprintf("%.2f", float)
						return nowValue
					}
				}
			}
		}

	}

	return nowValue
}

// 根据alert.Labels.Instance/alert.Labels.AlertName/alert.Labels.Project 生成MD5值

func CreateMd5(instance, alertName, project string) string {
	var data []byte
	if project != "" {
		data = []byte(instance + alertName + project)
	} else {
		data = []byte(instance + alertName)
	}
	sum := md5.Sum(data)
	encodeToString := hex.EncodeToString(sum[:])
	return encodeToString
}
