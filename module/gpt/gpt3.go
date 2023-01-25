package gpt

import (
	"context"
	"fmt"
	"gpt3"
	"log"
	"strings"
	"wechatbot/config"
	. "wechatbot/module/redis"
	"wechatbot/utils"
)

func GetGptReply(req string, userId string) (string, error) {

	apiKey := config.GlobalConfig.GPTApiKey
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)
	var str = ""
	// client.CompletionStream(ctx, gpt3.CompletionRequest{
	// 	Prompt:    []string{req},
	// 	MaxTokens: gpt3.IntPtr(3000),
	// 	Stop:      []string{"."},
	// 	Echo:      true,
	// }, func(data *gpt3.CompletionResponse) {
	// 	fmt.Println(data.Choices[0].Text)
	// 	str += data.Choices[0].Text
	// })
	c, _ := GetContext(userId)
	// fmt.Println(">>>>>>>final request >>>>>>>>>>", utils.TrimAllSpaceAndEnter(c)+"\n"+req)
	// fmt.Println("length>>>>>>>>>:", len(c+req))
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			c +"\n" +req,
		},
		Stream		: true,
		N			: gpt3.IntPtr(1),
		Temperature : gpt3.Float32Ptr(0.3),
		MaxTokens   : gpt3.IntPtr(4096 - (4 * len(c+req))),
	}, func(resp *gpt3.CompletionResponse) {
		if len(resp.Choices)>0{
			str += resp.Choices[0].Text
		}
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// s := strings.TrimSpace(str)
	SaveContext(req, str, userId)
	// return trimFront(str), err
	return str, err
}
//去除回答前面的符号
func trimFront(src string)(dist string){
	dist = strings.TrimSpace(src)
	if len(src) == 0 {
		return ""
	}
	if utils.StartsWith(dist, "?") || utils.StartsWith(dist, "？") || utils.StartsWith(dist, "") ||utils.StartsWith(dist, ":") || utils.StartsWith(dist, "：") ||utils.StartsWith(dist, "。"){
		if len(dist)>1{
			d,_ := utils.SubString(dist, 1)
			dist = d
		}else{
			return ""
		}
	} 
	return dist
}
