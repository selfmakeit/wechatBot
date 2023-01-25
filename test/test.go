package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"wechatbot/config"
	"wechatbot/module/coin"
	. "wechatbot/module/redis"
	"wechatbot/utils"

	"github.com/go-redis/redis"

	"gpt3"
)

func main() {

	getChieseEnter()
	// id := "23323"
	// testgpt("java写一个快速排序算法", id)
	// time.Sleep(time.Second * 3)
	// testgpt("能够再优化使其简短一点吗", id)
	// time.Sleep(time.Second * 3)
	// testgpt("爷爷呢", id)
	// time.Sleep(time.Second * 20)
	// test2("美国的首都是哪里"+""+as+"那里风土人情怎么样")
	// testGroup("写一首冬天的诗")
	// testGecko()
	// if r.RedisClient == nil {
	// 	fmt.Println("NILLLLLLLLLLL")
	// }
	// redis.RedisClient.Set("aa", "aavlaue", time.Second*5)
	// fmt.Println(redis.RedisClient.Get("aa").Result())
	// time.Sleep(time.Second * 6)
	// v,err := r.RedisClient.Get("aa").Result()
	// if err ==redis.Nil{
	// 	fmt.Println("errr")
	// 	fmt.Println(err)
	// }
	// fmt.Println("value-->",v)

}
func testgg() {
	fmt.Println("dur>>>>", config.GlobalConfig.ConversationExpire)

}
func GetContext2(userId string) (string, error) {
	queList, err := RedisClient.ZRangeWithScores(userId, 0, -1).Result()
	if err != err {
		fmt.Println(err)
		return "", err
	}
	if len(queList) < 1 {
		return "", nil
	}
	var res = ""
	for _, v := range queList {
		if int64(v.Score) < time.Now().Unix() {
			err := RedisClient.ZRemRangeByScore(userId, "-inf", "("+strconv.FormatInt(time.Now().Unix(), 10)).Err()
			if err != nil {
				return "", err
			}
		} else {
			res += fmt.Sprintf("%v", v.Member)
		}
	}
	return res, nil
}
func SaveContext2(q string, a string, userId string) error {

	queList0 := RedisClient.ZRangeWithScores(userId, 0, -1).Val()
	if len(queList0) < 1 {
		RedisClient.Expire(userId, time.Second*time.Duration(config.GlobalConfig.ConversationExpire))
	}
	anwser := trimQusetion(a)
	conv := redis.Z{
		// Member: "问题:"+q+ "回答:"+TrimAllSpaceAndEnter(anwser),
		Member: " " + q + " " + anwser,
		Score:  float64((time.Now().Add(time.Second * time.Duration(config.GlobalConfig.QuestionExpire))).Unix()),
	}
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>context saved!>>", userId, conv.Member)
	return RedisClient.ZAdd(userId, conv).Err()
}
func test1() {

}
func testgpt(pr string, userId string) string {
	// config := config.LoadConfig()
	apiKey := config.GlobalConfig.GPTApiKey
	if apiKey == "" {
		fmt.Println("Missing API KEY")
	}
	c, _ := GetContext(userId)
	println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>get  context>>>>>", c)
	ctx := context.Background()
	client := gpt3.NewClient(apiKey)
	// client.CompletionStream(ctx, gpt3.CompletionRequest{
	// 	Prompt:    []string{req},
	// 	MaxTokens: gpt3.IntPtr(3000),
	// 	Stop:      []string{"."},
	// 	Echo:      true,
	// }, func(data *gpt3.CompletionResponse) {
	// 	fmt.Println(data.Choices[0].Text)
	// 	str += data.Choices[0].Text
	// })
	var str = ""
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>最终请求>>>>", c+pr)
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			c + pr,
		},
		MaxTokens:   gpt3.IntPtr(4096-(3*len(c+pr))),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		str += resp.Choices[0].Text
		// fmt.Print(resp.Choices[0].Text)
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("回答1>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println(str)
	fmt.Println("回答2>>>>>>>>>>>>>>>>>>>>>>>>>>")
	SaveContext(pr, str, userId)
	return str
}
func testGroup(pr string) {
	fmt.Println("===================================TEST-group===================================")
	msg := " @CONTAC /gpt" + pr
	replaceGpt := "@" + "CONTAC" + config.GlobalConfig.GPTPrefix
	replaceCoinPrice := "@" + "CONTAC" + config.GlobalConfig.CoinPricePrefix
	replaceCoinMsg := "@" + "CONTAC" + config.GlobalConfig.CoinMsgPrefix
	replaceGas := "@" + "CONTAC" + config.GlobalConfig.GasPrefix
	content := utils.TrimAllSpace(msg)

	textContent := strings.ReplaceAll(content, replaceGpt, "")
	textContent = strings.ReplaceAll(textContent, replaceCoinPrice, "")
	textContent = strings.ReplaceAll(textContent, replaceCoinMsg, "")
	textContent = strings.ReplaceAll(textContent, replaceGas, "")
	fmt.Println(textContent)
}

func testGecko() {
	client := coin.GekClient
	// d, _ := client.GetPriceBySymbol("apecoin", "usd")
	// d1, _ := client.GetCoinPriceById("apecoin", "cny")
	// d2, _ := client.GetChainTokenPrice("ethereum", "cny")
	// d3, _ := client.GetChainTokenPrice("ethereum", "usd")
	d4, _ := client.SimpleCoinMarket("usd", "ethereum")
	// fmt.Println("price->", d)
	// fmt.Println("price->", d1)
	// fmt.Println("price->", d2)
	// fmt.Println("price->", d3)
	fmt.Println("price->", d4.Low24)
	var msg = fmt.Sprintf(`
	代币名称 : %v
	代币标识 : %v
	当前价格 : %v
	总 市 值 : %v
	市值排名 : %v
	完全稀释估值 : %v
	总交易额 : %v
	流 通 量 : %v
	总供应量 : %v
	24h最高价 : %v
	24h最低价 : %v
	24h价格变化 : %v
	24h变化百分比 : %v%%
	24h市值变化 : %v
	24h市值变化百分比 : %v%%
	历史最高价 : %v
	历史最高价时间 : %v
	距离最高价涨跌 : %v%%
	历史最低价 : %v
	历史最低价时间 : %v
	`, d4.Name, d4.Symbol, d4.CurrentPrice, d4.MarketCap, d4.MarketCapRank, d4.FDV,
		d4.TotalVolume, d4.CirculatingSupply, d4.TotalSupply, d4.High24, d4.Low24, d4.PriceChange24h, d4.PriceChangePercentage24h,
		d4.MarketCapChange24h, d4.MarketCapChangePercentage24h, d4.ATH, d4.ATHDate, d4.ATHChangePercentage, d4.ATL, d4.ATLDate)
	fmt.Println(msg)
}

func getChieseEnter() {
	a := `？
	`
	b := utils.TrimQusetionAtBegin(a)
	fmt.Println(b)
}
func TrimAllSpaceAndEnter(src string) (dist string) {
	src = strings.Trim(src, "\u3000")
	if len(src) == 0 {
		return
	}

	r, distR := []rune(src), []rune{}
	for i := 0; i < len(r); i++ {
		/*
			10 \n
			32 空格
			8197 中文空格
		*/
		// fmt.Println("ascii>>>",r[i])
		if r[i] == 8197 || r[i] == 32 || r[i] == 10 || r[i] == 9 || r[i] == 9 {
			continue
		}

		distR = append(distR, r[i])
	}
	dist = string(distR)
	return
}
func trimQusetion(src string) (dist string) {
	src = strings.Trim(src, "\u3000")
	if len(src) == 0 {
		return
	}

	r, distR := []rune(src), []rune{}
	for i := 0; i < len(r); i++ {
		/*
			10 \n
			32 空格
			8197 中文空格
		*/
		fmt.Println("ascii>>>",r[i])
		if r[i] == 65311 || r[i] == 63 {
			continue
		}

		distR = append(distR, r[i])
	}
	dist = string(distR)
	return
}
