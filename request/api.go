package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"webHook/parse"

	"github.com/shinesyang/common"
)

// 请求飞书 api webhook接口

func RequestFeishu(WarnDataAll []map[string]string) {
	contentList := [][]map[string]string{}

	count := 0
	for _, WarnData := range WarnDataAll {
		count++
		contents := []map[string]string{}
		for key, value := range WarnData {
			tag := "text"
			textMap := map[string]string{}
			if key == "告警面板" || key == "恢复面板" {
				tag = "a"
				textMap["text"] = key + ":" + value + "\n"
				textMap["tag"] = tag
				textMap["href"] = value
			} else {
				textMap["text"] = key + ":" + value + "\n"
				textMap["tag"] = tag
			}
			contents = append(contents, textMap)
		}
		// 存在多个告警，则空多行
		contents = append(contents, map[string]string{"tag": "at", "user_id": "all", "user_name": "所有人"})
		if count < len(WarnDataAll) {
			contents = append(contents, map[string]string{"tag": "text", "text": "\n\n\n"})
		}
		contentList = append(contentList, contents)
	}

	//contentMarshal, _ := json.Marshal(contentList)
	param := map[string]interface{}{
		"msg_type": "post",
		"content": map[string]interface{}{
			"post": map[string]interface{}{
				"zh_cn": map[string]interface{}{
					"title":   "PrometheusAlert告警消息",
					"content": contentList,
				},
			},
		},
	}

	paramMarshal, err := json.Marshal(param)
	common.Logger.Debugf("飞书参数: %s", string(paramMarshal))
	if err != nil {
		common.Logger.Errorf("param格式错误: %v", err)
		return
	}

	reader := bytes.NewReader(paramMarshal)
	request, err := http.NewRequest("POST", parse.CONF.FeiShuApiUrl, reader)
	if err != nil {
		common.Logger.Errorf("创建飞书 request实例失败: %v", err)
	}
	request.Header.Add(
		"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	)
	request.Header.Add(
		"Content-Type", "application/json",
	)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		common.Logger.Errorf("wechat http请求失败: %v", err)
		return
	}

	defer response.Body.Close()

	// 请求消息返回
	readAll, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		common.Logger.Infof("请求飞书机器人接口失败: %s,状态码: %d", string(readAll), response.StatusCode)
		return
	}

	resMap := map[string]interface{}{}
	_ = json.Unmarshal(readAll, &resMap)

	code := resMap["code"]
	codeInt := code.(float64)
	if codeInt != 0 {
		msg := resMap["msg"]
		common.Logger.Errorf("请求飞书机器人接口失败: %v", msg)
	} else {
		common.Logger.Infof("请求飞书机器人接口成功")
	}

}

// 请求微信 api

func RequestWeChat(WarnDataAll []map[string]string) {
	token, err := GetToken()
	if err != nil {
		return
	}

	MsgsendUrl := fmt.Sprintf("%s%v", parse.CONF.WeChatMsgUrl, token)

	MessageList := make([]string, 0, len(WarnDataAll))
	count := 0
	for _, WarnData := range WarnDataAll {
		count++
		for key, value := range WarnData {
			MessageList = append(MessageList, key+":"+value)
		}
		if count < len(WarnDataAll) {
			MessageList = append(MessageList, "\n\n", ">>>>>>>>>>>分割线<<<<<<<<<<<", "\n\n")
		}
	}

	Message := strings.Join(MessageList, "\n")

	tousers := "" // 默认告警接收人

	// 读取配置文件告警接收人
	users := parse.CONF.Users
	if users != "" {
		tousers = users
	}

	Params := map[string]interface{}{
		"touser":  tousers,
		"msgtype": "text",
		"agentid": 1000002,
		"text": map[string]string{
			"content": Message,
		},
		"safe": 0,
	}

	marshal, err := json.Marshal(Params)
	common.Logger.Debugf("微信参数: %s", string(marshal))
	if err != nil {
		common.Logger.Errorf("解析WeChat请求接口参数失败: %v", err)
		return
	}
	reader := bytes.NewReader(marshal)
	request, err := http.NewRequest("POST", MsgsendUrl, reader)
	if err != nil {
		common.Logger.Errorf("创建wechat 接口请求 request实例失败: %v", err)
		return
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		common.Logger.Errorf("wechat 接口 http请求失败: %v", err)
		return
	}

	defer response.Body.Close()

	readAll, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		common.Logger.Info("请求获取wechat接口失败: %s,状态码: %d", string(readAll), response.StatusCode)
		return
	}

	resMap := map[string]interface{}{}

	_ = json.Unmarshal(readAll, &resMap)

	errmsg := resMap["errmsg"].(string)
	errcode := resMap["errcode"].(float64)
	if errcode == float64(0) && errmsg == "ok" {
		common.Logger.Infof("请求wechat接口成功")
	} else {
		common.Logger.Errorf("请求wechat接口失败: %s", string(readAll))
	}
}

//获取token
func GetToken() (interface{}, error) {
	GetTokenUrl := parse.CONF.WeChatTokenIrl + parse.CONF.CorPid + "&corpsecret=" + parse.CONF.AppSecret

	request, err := http.NewRequest("GET", GetTokenUrl, nil)
	if err != nil {
		common.Logger.Errorf("创建wechat request实例失败: %v", err)
		return nil, err
	}
	request.Header.Add(
		"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	)
	client := http.Client{}

	response, err := client.Do(request)

	if err != nil {
		common.Logger.Errorf("wechat 获取token http请求失败: %v", err)
		return nil, err
	}

	defer response.Body.Close()

	readAll, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		msg := fmt.Sprintf("请求获取token接口失败: %s,状态码: %d", string(readAll), response.StatusCode)
		common.Logger.Info(msg)
		return nil, errors.New(msg)
	}

	resMap := map[string]interface{}{}
	_ = json.Unmarshal(readAll, &resMap)

	accessToken, ok := resMap["access_token"]
	if !ok {
		marshal, _ := json.Marshal(resMap)
		msg := fmt.Sprintf("获取WeChat token失败,返回结果为: %s", string(marshal))
		common.Logger.Info(msg)
		return nil, errors.New(msg)
	}

	return accessToken, nil
}
