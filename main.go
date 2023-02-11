package main

import (
	_ "wechatbot/config"
	"wechatbot/handlers"
	"context"
	"strings"
	"strconv"
    "fmt"
    "os"
    "time"
    "os/signal"
	. "wechatbot/module/coin"
	. "wechatbot/module/flash"
	l "wechatbot/module/log"
// "github.com/robfig/cron"
	"github.com/eatmoreapple/openwechat"
	"go.uber.org/zap"
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
	usr,_ := bot.GetCurrentUser()
	groups,er := usr.Groups(true)
	if er !=nil{
		fmt.Println("ssssss--->",er)
	}
	grs := groups.SearchByNickName(1,"核心猛冲群")
	ticker := time.NewTicker(time.Minute)
	pushMsg(grs,ticker)
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
func pushMsg(groups openwechat.Groups,ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			h, m, _ := time.Now().Clock()
			//推送币价
			if (h ==9 || h ==22)&& (m==1){
				for i :=0;i<3;i++{
					res,err :=GekClient.CoinPushToWechat()
					if err==nil{
						groups.First().SendText(res)
						break
					}
				}
			}
			//推送资讯
			if(h ==8 || h ==12|| h ==14|| h ==18|| h ==20|| h ==22)&& m==1 {
				res := FClien.GetFlashNews()
				s := formatNews(res)
				groups.First().SendText(s)
			}

		}
	}
    
}
func formatNews(news *FlashNews)string{
	if news ==nil{
		return ""
	}
	var msg =news.Date+"\n"
	for i,v := range *news.News{
		ss := strings.ReplaceAll(v.Summary, "[原文链接]", "")
		ss = strings.ReplaceAll(ss, "\n", "")
		msg += "\n\n"+strconv.Itoa(i+1)+". "+v.Title+":\n"+"  "+ss
	}
	return msg

}

func pushMsg2() {
	ticker := time.NewTicker(time.Minute)
    done := make(chan bool)
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()
    go func() {
        for {
            select {
            case <-done:
                return
            case <-ticker.C:
                _, _, s := time.Now().Clock()
                // h, m, s := time.Now().Clock()
				if s ==30{
					fmt.Printf("Doing the job")
				}
                /* if m == 0 && (  h == 9 || h == 15 ) {
                    fmt.Printf("Doing the job")
                } */
            }
        }
    }()

    <-ctx.Done()
    stop()
    done <- true
}
func main() {
	// pushMsg()
	Run()
	// defer RedisClient.Close()
}
