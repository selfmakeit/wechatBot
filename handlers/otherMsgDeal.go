package handlers

import (
	"github.com/eatmoreapple/openwechat"
	"wechatbot/config"
	// "log"
	// "strings"
)

func dealTickleMsg(msg *openwechat.Message,isGroup bool) error {
	if msg.IsSendBySelf() {
		return nil
	}
	if !msg.IsPaiYiPai(){
		return nil
	}
	//gpt的放在最后面，这样可以在配置文件去掉gpt的前缀并且最后一个else if改成else，在没有其他触发的情况下就触发gpt
	_, _ = msg.ReplyText("请在 /price后面加上代币符号\n比如/price ETH\n" + config.GlobalConfig.MsgSuffix)
	return nil
}