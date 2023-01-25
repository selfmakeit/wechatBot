package main

import (
	_ "wechatbot/config"
	"wechatbot/handlers"
	_ "wechatbot/module/coin"
	l"wechatbot/module/log"
	"go.uber.org/zap"
	"github.com/eatmoreapple/openwechat"
)

func Run() {
	l.InitZap()
	//bot := openwechat.DefaultBot()
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式
	// 注册消息处理函数
	bot.MessageHandler = handlers.Handler
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")
	// 执行热登录
	err := bot.HotLogin(reloadStorage)
	if err != nil {
		if err = bot.Login(); err != nil {
			l.LOG.Error("登录错误!", zap.Error(err))
			return
		}
	}
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
func main() {
	Run()
	// defer RedisClient.Close()
}
