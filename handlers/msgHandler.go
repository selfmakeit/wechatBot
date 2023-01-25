package handlers

import (
	"fmt"
	"strings"
	"wechatbot/config"
	"wechatbot/module/coin"
	l "wechatbot/module/log"
	"go.uber.org/zap"
	"wechatbot/utils"

	"github.com/eatmoreapple/openwechat"
)

// MessageHandlerInterface 消息处理接口
type MessageHandlerInterface interface {
	handler(*openwechat.Message) error
	ReplyText(*openwechat.Message, int16) error
}

type HandlerType string

var (
	gptType     int16 = 1
	coinType    int16 = 2
	coinMsgType int16 = 3
)

const (
	GroupHandler = "group"
	UserHandler  = "user"
)

// handlers 所有消息类型类型的处理器
var handlers map[HandlerType]MessageHandlerInterface

func init() {
	handlers = make(map[HandlerType]MessageHandlerInterface)
	handlers[GroupHandler] = NewGroupMessageHandler()
	handlers[UserHandler] = NewUserMessageHandler()
}

// Handler 全局处理入口
func Handler(msg *openwechat.Message) {
	// log.Printf("hadler Received msg : %v", msg.Content)
	// 处理群消息
	if msg.IsSendByGroup() {
		err := handlers[GroupHandler].handler(msg)
		if err!=nil{
			l.LOG.Error("", zap.Error(err))
			msg.ReplyText("发生错误了：\n"+err.Error())
		}
		return
	}
	// 好友申请
	if msg.IsFriendAdd() {
		if config.GlobalConfig.AutoPass {
			_, err := msg.Agree("你好我是基于chatGPT引擎开发的微信机器人，你可以向我提问任何问题。")
			if err != nil {
				l.LOG.Error("add friend agree error", zap.Error(err))
				return
			}
		}
	}

	// 私聊
	err := handlers[UserHandler].handler(msg)
	if err!=nil{
		l.LOG.Error("", zap.Error(err))
		msg.ReplyText("发生错误了：\n"+err.Error())
	}
}
func dealErr(err error, msg string) {

}
func formatCoinMsg(coin *coin.CoinsMarketItem) string {
	if coin == nil {
		return ""
	}
	var mcd, _ = utils.SubStringBetween(coin.ATHDate, 0, 10)
	var ald, _ = utils.SubStringBetween(coin.ATLDate, 0, 10)
	var msg = fmt.Sprintf(
		`
代币名称 : %v
代币标识 : %v
当前价格 : %v
总 市 值 : %v
市值排名 : %v
完全稀释估值 : %v
总交易额 : %v
流 通 量 : %v
总供应量(枚) : %v
24小时最高价 : %v
24小时最低价 : %v
24小时价格变化 : %v
24小时变化百分比 : %v%%
24小时市值变化 : %v
24小时市值变化百分比 : %v%%
历史最高价 : %v
历史最高价时间 : %v
距离最高价涨跌 : %v%%
历史最低价 : %v
历史最低价时间 : %v`,
		coin.Name,
		strings.ToUpper(coin.Symbol),
		coin.CurrentPrice,
		coin.MarketCap,
		coin.MarketCapRank,
		coin.FDV,
		coin.TotalVolume,
		coin.CirculatingSupply.Round(2),
		coin.TotalSupply.Round(2),
		coin.High24, coin.Low24,
		coin.PriceChange24h.Round(2),
		coin.PriceChangePercentage24h.Round(2),
		coin.MarketCapChange24h,
		coin.MarketCapChangePercentage24h.Round(2),
		coin.ATH,
		mcd,
		coin.ATHChangePercentage.Round(2),
		coin.ATL,
		ald)

	return msg
}
