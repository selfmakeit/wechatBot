package handlers

import (
	"time"
	"fmt"
	// "log"
	"strings"
	"wechatbot/config"
	"wechatbot/module/coin"
	"wechatbot/module/gpt"
	l "wechatbot/module/log"
	"wechatbot/module/stock"
	"wechatbot/utils"

	"github.com/eatmoreapple/openwechat"
)

var _ MessageHandlerInterface = (*GroupMessageHandler)(nil)

// GroupMessageHandler 群消息处理
type GroupMessageHandler struct {
}

// handle 处理消息
func (g *GroupMessageHandler) handler(msg *openwechat.Message) error {
	if msg.IsAt() {
		if msg.IsSendBySelf() {
			return nil
		}
		if msg.IsText() {
			//去除@我的字符串
			textContent := removeAtString(msg)
			textContent = utils.TrimAllSpace(textContent)
			if utils.ContainsIgnoreCase(textContent, config.GlobalConfig.CoinMsgPrefix) {
				s, err := utils.SubStringBetween(textContent, 0, len(config.GlobalConfig.CoinMsgPrefix)) //前缀必须在前面
				if err == nil && utils.ContainsIgnoreCase(s, config.GlobalConfig.CoinMsgPrefix) {
					return g.ReplyText(msg, coinMsgType)
				}
				return g.ReplyText(msg, coinMsgType)
			} else if utils.ContainsIgnoreCase(textContent, config.GlobalConfig.CoinPricePrefix) {
				s, err := utils.SubStringBetween(textContent, 0, len(config.GlobalConfig.CoinPricePrefix))
				if err == nil && utils.ContainsIgnoreCase(s, config.GlobalConfig.CoinPricePrefix) {
					return g.ReplyText(msg, coinType)
				}
				// return g.ReplyText(msg, coinType)
			} else if utils.ContainsIgnoreCase(textContent, config.GlobalConfig.GPTPrefix) {
				s, err := utils.SubStringBetween(textContent, 0, len(config.GlobalConfig.GPTPrefix))
				if err == nil && utils.ContainsIgnoreCase(s, config.GlobalConfig.GPTPrefix) {
					return g.ReplyText(msg, gptType)
				}
			} else if utils.ContainsIgnoreCase(textContent, config.GlobalConfig.StockNewsPrefix) {
				s, err := utils.SubStringBetween(textContent, 0, len(config.GlobalConfig.StockNewsPrefix))
				if err == nil && utils.ContainsIgnoreCase(s, config.GlobalConfig.StockNewsPrefix) {
					return g.ReplyText(msg, stockNewsType)
				}
			} else if utils.ContainsIgnoreCase(textContent, config.GlobalConfig.StockPrefix) {
				s, err := utils.SubStringBetween(textContent, 0, len(config.GlobalConfig.StockPrefix))
				if err == nil && utils.ContainsIgnoreCase(s, config.GlobalConfig.StockPrefix) {
					return g.ReplyText(msg, stockType)
				}
			} else {
				return g.ReplyText(msg, gptType)
			}
		}
	} else {
		if msg.IsPaiYiPai() {
			// _, err := msg.ReplyText("咋啦")
			// return err
		}
	}
	return nil
}

// NewGroupMessageHandler 创建群消息处理器
func NewGroupMessageHandler() MessageHandlerInterface {
	return &GroupMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *GroupMessageHandler) ReplyText(msg *openwechat.Message, msgType int16) error {
	// 接收群消息
	sender, err := msg.Sender()
	if err != nil {
		return fmt.Errorf("get sender in group error :%v", err)
	}
	group := openwechat.Group{sender}
	// l.LOG.Info("Received Group Message From (User " + group.NickName + " ID:" + group.ID() + ") =>Text Msg : " + removeAtString(msg))
	// 获取@我的用户
	groupSender, err := msg.SenderInGroup()
	if err != nil {
		return fmt.Errorf("get sender in group error :%v", err)
	}
	req := getGroupReqText(msg, msgType)
	if req == "" {
		_, err = msg.ReplyText("@" + groupSender.NickName + "\n  " + "请加上你想提的问题:\n\n目前支持:\n1. /gpt + 随意内容交互, 比如：\n    /gpt 请用Go写一段代码\n2. /price +代币查询价格,比如：\n   /price ETH\n3. /coin +代币查询币信息,比如:\n   /coin ETH\n" + config.GlobalConfig.MsgSuffix)
		return err
	}
	var reply = ""
	switch msgType {
	case gptType:
		l.LOG.Sugar().Infof("group:[ %v ] in group [ %v ] ask [ %v ] to gpt", groupSender.NickName, group.NickName, removeAtString(msg))
		reply, err = gpt.GetGptReply(req, sender.ID())
		if err != nil {
			_, err = msg.ReplyText("@" + groupSender.NickName + "\n  " + "网络出问题了,请一会儿再试")
			if err != nil {
				return fmt.Errorf("response group error: %v ", err)
			}
			return fmt.Errorf("get gpt reply error :%v", err)
		}
	case coinType:
		coinmsg, err := coin.GekClient.GetCoinMsgBySymbol(req)
		if err != nil {
			return fmt.Errorf("get coin error: %v", err)
		}
		l.LOG.Sugar().Infof("group:[ %v ] in group [ %v ]  search coin price[symbol: %v name: %v]", groupSender.NickName, group.NickName, coinmsg.Symbol, coinmsg.Name)
		price, err := coin.GekClient.GetCoinPriceById(coinmsg.ID, "usd")
		if err != nil {
			return fmt.Errorf("get coin price error: %v", err)
		}
		if price.IsZero() {
			reply = ""
		} else {
			reply = price.String()
		}
	case coinMsgType:
		coinmsg, err := coin.GekClient.GetCoinMsgBySymbol(req)
		if err != nil {
			return fmt.Errorf("get coin error: %v", err)
		}
		l.LOG.Sugar().Infof("group:[ %v ] in group [ %v ]  search coin message[symbol: %v name: %v]", groupSender.NickName, group.NickName, coinmsg.Symbol, coinmsg.Name)
		msg, err := coin.GekClient.SimpleCoinMarketById("usd", coinmsg.ID)
		if err != nil {
			return fmt.Errorf("get coin message error: %v", err)
		}
		reply = formatCoinMsg(msg)
	case stockNewsType:
		for i := 0; i < 5; i++ { //经常出问题，请求多次
			res := stock.SClien.GetStockNews(req)
			reply = formatStockNewsMsg(res)
			if reply != "" {
				break
			}else{
				time.Sleep(time.Microsecond*200)
			}
		}
	case stockType:
		res := stock.SClien.GetStockData(req)
		reply = formatStockMsg(res)
	}
	if reply == "" {
		_, _ = msg.ReplyText("@" + groupSender.NickName + "我也不知道[尴尬]" + config.GlobalConfig.MsgSuffix)
		return fmt.Errorf("got a nil reply")
	}

	// 回复@我的用户
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	replyText := createGroupReplyContent(groupSender.NickName, req, reply, msgType)
	if replyText == "" {
		_, err = msg.ReplyText("@" + groupSender.NickName + "我也不知道[尴尬]" + config.GlobalConfig.MsgSuffix)
		if err != nil {
			return fmt.Errorf("response group error: %v", err)
		}
	}
	_, err = msg.ReplyText(replyText)
	if err != nil {
		return fmt.Errorf("response group error: %v", err)
	}
	return nil
}

/* 去除@我的字符串 */
func removeAtString(msg *openwechat.Message) string {
	sender, err := msg.Sender()
	if err != nil {
		return ""
	}
	replaceText := "@" + sender.Self().NickName
	// content := utils.TrimAllSpace(msg.Content)
	content := strings.TrimSpace(msg.Content)
	textContent := strings.ReplaceAll(content, replaceText, "")
	return textContent
}
func createGroupReplyContent(to string, request string, content string, msgType int16) string {
	var reply = ""
	switch msgType {
	case gptType:
		reply = "@" + to + " \"" + request + "\":\n" /* + "\n---------------------------\n" */ + content + "\n" + config.GlobalConfig.MsgSuffix
	case coinType:
		reply = "@" + to + " >> " + request + "当前价格:\n\n$" + content + config.GlobalConfig.MsgSuffix
	case coinMsgType:
		reply = "@" + to + ">> 代币" + request + "信息(单位：美元):\n\n" + content + config.GlobalConfig.MsgSuffix
	case stockType:
		reply = "@" + to + ">> 股票" + request + "信息:\n\n" + content + config.GlobalConfig.MsgSuffix
	case stockNewsType:
		reply = "@" + to + ">> " + request + "热门新闻:\n\n" + content + config.GlobalConfig.MsgSuffix
	}
	return reply
}

/* 获取问题：去除@我和前缀 */
func getGroupReqText(msg *openwechat.Message, msgType int16) string {

	var textContent = ""
	sender, err := msg.Sender()
	if err != nil {
		return ""
	}
	switch msgType {
	case gptType:
		replaceGpt := "@" + sender.Self().NickName + config.GlobalConfig.GPTPrefix
		replaceGpt2 := "@" + sender.Self().NickName
		//这里把文本中间的那些空格也去掉的话在英文环境下无法阅读，英文是利用空格来间隔单词,所以这里要把gpt的单独出来另外提取
		// content := utils.TrimAllSpace(msg.Content)
		content := utils.TrimAllChineseSpace(msg.Content)
		content = strings.Trim(content, "\n") //去除回车

		textContent = strings.ReplaceAll(content, replaceGpt, "")
		textContent = strings.ReplaceAll(textContent, replaceGpt2, "")
	default:
		replaceCoinPrice := "@" + sender.Self().NickName + config.GlobalConfig.CoinPricePrefix
		replaceCoinMsg := "@" + sender.Self().NickName + config.GlobalConfig.CoinMsgPrefix
		replaceGas := "@" + sender.Self().NickName + config.GlobalConfig.GasPrefix
		replaceStock := "@" + sender.Self().NickName + config.GlobalConfig.StockPrefix
		replaceStockNews := "@" + sender.Self().NickName + config.GlobalConfig.StockNewsPrefix

		content := utils.TrimAllSpace(msg.Content)
		// content :=strings.TrimSpace(msg.Content)
		content = strings.Trim(content, "\n") //去除回车

		textContent = strings.ReplaceAll(content, replaceCoinPrice, "")
		textContent = strings.ReplaceAll(textContent, replaceCoinMsg, "")
		textContent = strings.ReplaceAll(textContent, replaceGas, "")
		textContent = strings.ReplaceAll(textContent, replaceStock, "")
		textContent = strings.ReplaceAll(textContent, replaceStockNews, "")
	}
	return textContent
}
