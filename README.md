# 微信机器人说明

**功能：**

* 接入2chatgpt

* 虚拟币实时价格和信息查询

* 虚拟币watch list价格定时推送

* 股票价格查询

* 股票热门新闻抓取

* web3热门新闻抓取推送

**说明**

1. **关于chatgpt： **接入了gpt-3,利用redis实现了对话上下文缓存，在使用的时候可根据自己的需求更改module/redis/redis.go里的内容。对于gpt3使用的是

2. **关于虚拟币： **如果有coingecko的apikey可以在配置文件中填写，没有也可以用，只是频率有一定限制。

3. **日志：** 日志记录使用了zap，并做了日志文件分割，对于微信信息部分的日志因为是使用的openwechat(https://github.com/eatmoreapple/openwechat) ，没有去改内部的日志处理，需要的可自行定制。

使用截图：

![image-20230126003032724](https://raw.githubusercontent.com/selfmakeit/resource/main/image-20230126003032724.png)

![image-20230126003058183](https://raw.githubusercontent.com/selfmakeit/resource/main/image-20230126003058183.png)



<img src="https://raw.githubusercontent.com/selfmakeit/resource/main/image-20230210002700406.png" alt="image-20230210002700406" style="zoom:67%;" />

![image-20230210002733091](https://raw.githubusercontent.com/selfmakeit/resource/main/image-20230210002733091.png)

![image-20230211230331742](https://raw.githubusercontent.com/selfmakeit/resource/main/image-20230211230331742.png)

# 定制建议：
1. 日志：如果openwechat内部部分的日志保持一致，需要下载openwechat代码到本地进行更改。
2. 如果要扩展其他功能，可自己在config.go文件中自定义消息触发前缀，然后在消息处理的地方(群消息处理和私聊处理)添加消息类型判断。
3. 后续在消息推送方面扩展开发的时候要限制一下消息推送频率，不然容易被封号。
