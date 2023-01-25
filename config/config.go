package config

import (
	"encoding/json"
	"fmt"
	io "io/ioutil"
	"strings"
	"sync"

	"github.com/superoo7/go-gecko/v3/types"
	"go.uber.org/zap/zapcore"
	// "net/http"
	// "net/url"
	// "time"
	// "log"
	// coingecko "github.com/superoo7/go-gecko/v3"
)

// Configuration 项目配置
type Configuration struct {
	// gtp apikey(token)
	GPTApiKey    string `json:"GPTApiKey"`
	GeckoApiKey  string `json:"GeckoApiKey"`
	// BearerToken  string `json:"BearerToken"`
	// SessionToken string `json:"SessionToken"` //对应"__Secure-next-auth.session-token"，To avoid needing to refresh bearer token every hour,
	//触发前缀
	GPTPrefix       string `json:"GPTPrefix"`
	CoinMsgPrefix   string `json:"CoinMsgPrefix"`
	CoinPricePrefix string `json:"CoinPricePrefix"`
	GasPrefix       string `json:"GasPrefix"`
	//用户过滤设置
	SkipSelf     bool   `json:"SkipSelf"`     //过滤自己发的消息，自己的消息不触发
	OnlySelf     bool   `json:"OnlySelf"`     //过滤别人发的消息，只触发自己发的消息
	OnlySomebody bool   `json:"OnlySomebody"` //只能指定某个人发的消息才出发
	SomebodyID   string `json:"SomebodyID"`   //指定的那个人的id 2285462961

	//redis
	QuestionExpire     int `json:"QuestionExpire"`     //上下文中一个问答的持续时长，redis的存储时长
	ConversationExpire int `json:"ConversationExpire"` //上下文中一个会话的持续时长，redis的存储时长
	//日志
	LogFormat     string `json:"LogFormat"`     //日志格式
	LogDirectory  string `json:"LogDirectory"`  //日志存放文件夹
	ShowInConsole bool   `json:"ShowInConsole"` //是否在控制台显示日志
	LogKeepDays   int    `json:"LogKeepDays"`   //日志存留时间
	LogLevel      string `json:"LogLevel"`      //日志级别
	ShowLine      bool   `json:"ShowLine"`      //显示行

	// 自动通过好友
	AutoPass      bool `json:"AutoPass"`
	UseProxy      bool `json:"UseProxy"`
	GeckoUseProxy bool `json:"GeckoUseProxy"` //gecko在本地需要代理，不管开没开代理软件
	//消息
	MsgSuffix	  string `json:"MsgSuffix"`  //消息后缀
}

var (
	GlobalConfig *Configuration
	CoinList     *types.CoinList
)
var file_locker sync.Mutex

func init() {
	loadConfig()
	// initGeckoClient()
	// initCoinList()
}

// LoadConfig 加载配置
func loadConfig() *Configuration {

	GlobalConfig = &Configuration{}
	file_locker.Lock()
	data, err := io.ReadFile("config/config.json") //read config file
	file_locker.Unlock()
	if err != nil {
		fmt.Println("read cinfig file error")
		panic(err)
	}
	datajson := []byte(data)
	err = json.Unmarshal(datajson, &GlobalConfig)
	if err != nil {
		fmt.Println("unmarshal json file error")
		panic(err)
	}
	// fmt.Println("== BearerToken ==>", config.BearerToken)
	// fmt.Println("== SessionToken ==>", config.SessionToken)
	// fmt.Println("== ApiKey ==>", config.ApiKey)
	return GlobalConfig
}

// TransportLevel 根据字符串转化为 zapcore.Level
func GetLogLevel() zapcore.Level {
	l := strings.ToLower(GlobalConfig.LogLevel)
	switch l {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.WarnLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

// func initGeckoClient(){
// 	var httpClient *http.Client
// 	if GlobalConfig.UseProxy {
// 		u, _ := url.Parse("htttp://127.0.0.1:7890")
// 		t := &http.Transport{
// 			MaxIdleConns:    10,
// 			MaxConnsPerHost: 10,
// 			IdleConnTimeout: time.Duration(10) * time.Second,
// 			//Proxy: http.ProxyURL(url),
// 			Proxy: http.ProxyURL(u),
// 		}
// 		httpClient = &http.Client{
// 			Transport: t,
// 			Timeout:   time.Duration(10) * time.Second,
// 		}
// 	} else {
// 		httpClient = &http.Client{
// 			Timeout: time.Second * 100,
// 		}
// 	}
// 	CoinGeckoClient = coingecko.NewClient(httpClient)
// }
// func initCoinList(){
// 	CoinList, err := CoinGeckoClient.CoinsList()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Coin List inited !Available coins:", len(*CoinList))
// }
