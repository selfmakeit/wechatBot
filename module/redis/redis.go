package redis

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"wechatbot/config"
	"wechatbot/utils"

	// l"wechatbot/module/log"
	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // url
		Password: "",
		DB:       1, // 0号数据库
	})
	_, err := RedisClient.Ping().Result()
	if err != nil {
		// l.LOG.Error("redis init success! ",zap.Error(err))
		fmt.Println("redis init err :", err)
		panic(err)
	}
	// l.LOG.Info("redis init success! ")
	fmt.Println("redis init success! ")
}
func GetContext(userId string) (string, error) {
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
	if len(queList) >= 4 || len(res) > 1000 { //上下文大于1000或者个数大于4就清除第一个
		RedisClient.ZRem(userId, queList[0].Member)
	}
	return res, nil
}
func SaveContext(q string, a string, userId string) error {
	queList := RedisClient.ZRangeWithScores(userId, 0, -1).Val()
	if len(queList) >= 4 { //上下文个数大于4就清除第一个
		RedisClient.ZRem(userId, queList[0].Member)
	}
	if len(queList) < 1 {
		RedisClient.Expire(userId, time.Second*time.Duration(config.GlobalConfig.ConversationExpire))
	}
	answer := trimQusetion(a)
	conv := redis.Z{
		Member: " " + q + " " + utils.TrimQusetionAtBegin(answer),
		Score:  float64((time.Now().Add(time.Second * time.Duration(config.GlobalConfig.QuestionExpire))).Unix()),
	}
	// fmt.Println("context saved!>>", userId,conv.Member)
	return RedisClient.ZAdd(userId, conv).Err()
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
		// fmt.Println("ascii>>>",r[i])
		if r[i] == 65311 || r[i] == 63 {
			continue
		}

		distR = append(distR, r[i])
	}
	dist = string(distR)
	return
}
func TrimQusetionAtBegin(src string) (dist string) {
	src = strings.Trim(src, "\u3000")
	if len(src) == 0 {
		return
	}

	r, _ := []rune(src), []rune{}
	if r[0] == 65311 || r[0] == 63 {
		dist, _ = utils.SubStringBetween(src, 1, len(src))
	} else {
		dist = src
	}
	return
}
