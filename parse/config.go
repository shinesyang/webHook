package parse

import (
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/shinesyang/common"
)

// 配置文件参数解析

type Config struct {
	Username       string `json:"username" toml:"username"`
	Password       string `json:"password" toml:"password"`
	Users          string `json:"users" toml:"users"`
	CorPid         string `json:"corPid" toml:"corPid"`
	AppSecret      string `json:"appSecret" toml:"appSecret"`
	WeChatMsgUrl   string `json:"weChatMsgUrl" toml:"weChatMsgUrl"`
	WeChatTokenIrl string `json:"weChatTokenIrl" toml:"weChatTokenIrl"`
	FeiShuApiUrl   string `json:"feiShuApiUrl" toml:"feiShuApiUrl"`
	Panel          string `json:"panel" toml:"panel"`
}

var CONF Config

func init() {
	dir, _ := os.Getwd()

	file := path.Join(dir, "config", "config")

	_, err := toml.DecodeFile(file, &CONF)
	if err != nil {
		common.Logger.Warnf("解析config配置文件失败: %v", err)
	}
}
