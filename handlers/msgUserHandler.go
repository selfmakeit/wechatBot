package handlers

import (
	"fmt"
	"strings"
"time"
	"wechatbot/config"
	"wechatbot/module/coin"
	"wechatbot/module/gpt"
	"wechatbot/module/stock"

	// ."wechatbot/module/redis"
	l "wechatbot/module/log"
	// "go.uber.org/zap"
	"wechatbot/utils"

	"github.com/eatmoreapple/openwechat"
)

var _ MessageHandlerInterface = (*UserMessageHandler)(nil)

// UserMessageHandler 私聊消息处理
type UserMessageHandler struct {
}

// handler 分类 处理消息
func (g *UserMessageHandler) handler(msg *openwechat.Message) error {
	if msg.IsSendBySelf() {
		return nil
	} else {
		if msg.IsText() {
			message := utils.TrimAllSpace(msg.Content)
			if utils.ContainsIgnoreCase(message, config.GlobalConfig.CoinMsgPrefix) {
				//前缀只能在消息最前面
				s, err := utils.SubStringBetween(message, 0, len(config.GlobalConfig.CoinMsgPrefix))
				if err == nil && utils.ContainsIgnoreCase(s, config.GlobalConfig.CoinMsgPrefix) {
					return g.ReplyText(msg, coinMsgType)
				}
			} else if utils.ContainsIgnoreCase(message, config.GlobalConfig.CoinPricePrefix) {
				s, err := utils.SubStringBetween(message, 0, len(config.GlobalConfig.CoinPricePrefix))
				if err == nil && utils.ContainsIgnoreCase(s, config.GlobalConfig.CoinPricePrefix) {
					return g.ReplyText(msg, coinType)
				}
			} else if utils.ContainsIgnoreCase(message, config.GlobalConfig.GPTPrefix) {
				s, err := utils.SubStringBetween(message, 0, len(config.GlobalConfig.GPTPrefix))
				if err == nil && utils.ContainsIgnoreCase(s, config.GlobalConfig.GPTPrefix) {
					return g.ReplyText(msg, gptType)
				}
			}else if utils.ContainsIgnoreCase(message, config.GlobalConfig.StockNewsPrefix) {
				s, err := utils.SubStringBetween(message, 0, len(config.GlobalConfig.StockNewsPrefix))
				if err == nil && utils.ContainsIgnoreCase(s, config.GlobalConfig.StockNewsPrefix) {
					return g.ReplyText(msg, stockNewsType)
				}
			}else if utils.ContainsIgnoreCase(message, config.GlobalConfig.StockPrefix) {
				s, err := utils.SubStringBetween(message, 0, len(config.GlobalConfig.StockPrefix))
				if err == nil && utils.ContainsIgnoreCase(s, config.GlobalConfig.StockPrefix) {
					return g.ReplyText(msg, stockType)
				}
			}else {
				return g.ReplyText(msg, gptType)
			}
		} else if msg.IsPaiYiPai() {
			msg.ReplyText("咋啦")
		} else if msg.IsPicture() {

		} else if msg.IsVoice() {

		} else if msg.IsVideo() {

		}
	}

	return nil
}

// NewUserMessageHandler 创建私聊处理器
func NewUserMessageHandler() MessageHandlerInterface {
	return &UserMessageHandler{}
}

// ReplyText 发送文本消息
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message, msgType int16) error {
	// 接收私聊消息
	if msg.IsSendBySelf() {
		return nil
	}
	sender, err := msg.Sender()
	if err != nil {
		return err
	}
	// l.LOG.Info("Received Private Message From (User " + sender.NickName + " ID:" + sender.ID() + ") =>Text Msg : " + msg.Content)
	//去除两头空格
	//去除回车
	requestText := getPrvReqText(msg, msgType)
	if requestText == "" {
		switch msgType {
		case gptType:
			_, _ = msg.ReplyText("请在 /gpt后面加上你想提的问题" + config.GlobalConfig.MsgSuffix)
		case coinType:
			_, _ = msg.ReplyText("请在 /price后面加上代币符号\n比如/price ETH\n" + config.GlobalConfig.MsgSuffix)
		case coinMsgType:
			_, _ = msg.ReplyText("请在 /coin后面加上代币符号\n比如/coin ETH\n" + config.GlobalConfig.MsgSuffix)
		case stockNewsType:
			_, _ = msg.ReplyText("请在 /news后面加上代币符号\n比如/news apple\n" + config.GlobalConfig.MsgSuffix)
		case stockType:
			_, _ = msg.ReplyText("请在 /stock后面加上代币符号\n比如/stock tesla\n" + config.GlobalConfig.MsgSuffix)
		}
		return nil
	}
	var reply = ""
	switch msgType {
	case gptType:
		//获取上下文
		l.LOG.Sugar().Infof("personal [ %v ] ask [ %v ] to gpt", sender.NickName, requestText)
		reply, err = gpt.GetGptReply(requestText, sender.ID())
		if err != nil {
			_, _ = msg.ReplyText("gpt网络出问题了,请一会儿再试")
			return fmt.Errorf("gpt request error: %v", err)
		}
	case coinType:
		requestText = utils.TrimAllSpace(requestText)
		coinmsg, err := coin.GekClient.GetCoinMsgBySymbol(requestText)
		if err != nil {
			return fmt.Errorf("get coin error: %v", err)
		}
		l.LOG.Sugar().Infof("personal [ %v ] search coin price[symbol: %v name: %v]", sender.NickName, coinmsg.Symbol, coinmsg.Name)
		price, err := coin.GekClient.GetCoinPriceById(coinmsg.ID, "usd")
		requestText = requestText + "(" + coinmsg.Name + ")"
		if err != nil {
			return fmt.Errorf("get coin price error: %v(id:%v)", err, coinmsg.ID)
		}
		if price.IsZero() {
			reply = ""
		} else {
			reply = price.String()
		}

	case coinMsgType:
		requestText = utils.TrimAllSpace(requestText)
		coinmsg, err := coin.GekClient.GetCoinMsgBySymbol(requestText)
		if err != nil {
			return fmt.Errorf("get coin error: %v", err)
		}
		l.LOG.Sugar().Infof("personal [ %v ] search coin message[symbol: %v name: %v]", sender.NickName, coinmsg.Symbol, coinmsg.Name)
		msg, err := coin.GekClient.SimpleCoinMarketById("usd", coinmsg.ID)
		if err != nil {
			return fmt.Errorf("get coin price error: %v", err)
		}
		reply = formatCoinMsg(msg)
	case stockNewsType:
		requestText = utils.TrimAllSpace(requestText)
		fmt.Println(requestText)
		for i :=0;i<5;i++{//经常出问题，请求多次
			res := stock.SClien.GetStockNews(requestText)
			reply = formatStockNewsMsg(res)
			if reply !=""{
				break
			}else{
				time.Sleep(time.Microsecond*200)
			}
		}
	case stockType:
		requestText = utils.TrimAllSpace(requestText)
		res := stock.SClien.GetStockData(requestText)
		reply = formatStockMsg(res)
	}

	// reply, err := gpt.Completions(requestText)
	if reply == "" {
		_, err = msg.ReplyText("我也不知道[尴尬]" + config.GlobalConfig.MsgSuffix)
		if err != nil {
			return fmt.Errorf("response error error: %v", err)
		}
		return fmt.Errorf("没有获取到消息: %v", err)
	}
	reply = strings.TrimSpace(reply)
	reply = createPrvReplyContent(requestText, reply, msgType)
	// 回复用户
	_, err = msg.ReplyText(reply)
	if err != nil {
		return fmt.Errorf("response error error: %v", err)
	}
	return nil
}
func createPrvReplyContent(request string, content string, msgType int16) string {
	var re = ""
	switch msgType {
	case gptType:
		re = /* "\"" + request + "\"" + "\n---------------------------\n" + */ content + config.GlobalConfig.MsgSuffix
	case coinType:
		re = request + "当前价格:\n\n$" + content + config.GlobalConfig.MsgSuffix

	case coinMsgType:
		re = "代币" + request + "信息(单位：美元):\n\n" + content + config.GlobalConfig.MsgSuffix
	case stockType:
		re = "股票" + request + "信息:\n\n" + content + config.GlobalConfig.MsgSuffix
	case stockNewsType:
		re = request + "热门新闻:\n\n" + content + config.GlobalConfig.MsgSuffix
	}
	return re
}
func getPrvReqText(msg *openwechat.Message, msgyType int16) string {

	var textContent = ""
	switch msgyType {
	case gptType:
		replaceGpt := config.GlobalConfig.GPTPrefix
		//这里把文本中间的那些空格也去掉的话在英文环境下无法阅读，英文是利用空格来间隔单词,所以这里要把gpt的单独出来另外提取
		// content := utils.TrimAllSpace(msg.Content)
		content := utils.TrimAllChineseSpace(msg.Content)
		content = strings.Trim(content, "\n") //去除回车

		textContent = strings.ReplaceAll(content, replaceGpt, "")
	default:
		content := utils.TrimAllSpace(msg.Content)
		//这种方式不区分
		// replaceCoinPrice, _ := utils.SubStringBetween(content, 0, len(config.GlobalConfig.CoinPricePrefix))
		// replaceCoinMsg, _ := utils.SubStringBetween(content, 0, len(config.GlobalConfig.CoinMsgPrefix))
		// replaceGas, _ := utils.SubStringBetween(content, 0, len(config.GlobalConfig.GasPrefix))
		// replaceNews, _ := utils.SubStringBetween(content, 0, len(config.GlobalConfig.StockNewsPrefix))
		// replaceStock, _ := utils.SubStringBetween(content, 0, len(config.GlobalConfig.StockPrefix))

		// content :=strings.TrimSpace(msg.Content)
		content = strings.Trim(content, "\n") //去除回车
		textContent = strings.ReplaceAll(content, config.GlobalConfig.CoinPricePrefix, "")
		textContent = strings.ReplaceAll(textContent, config.GlobalConfig.CoinMsgPrefix, "")
		textContent = strings.ReplaceAll(textContent, config.GlobalConfig.GasPrefix, "")
		textContent = strings.ReplaceAll(textContent, config.GlobalConfig.StockNewsPrefix, "")
		textContent = strings.ReplaceAll(textContent, config.GlobalConfig.StockPrefix, "")
	}
	return textContent
}
