package flash

import (
	"time"
	"fmt"
	// "strconv"
	"strings"
	"wechatbot/config"
	"github.com/gocolly/colly/v2"
)


type News struct {
	Title          string `json:"title,omitempty"`
	Summary          string `json:"summary,omitempty"`
}

type FlashNews struct{
	News *[]News `json:"news,omitempty"`
	Date string `json:"date,omitempty"`
}
type FlashClien struct {
	collector *colly.Collector
}

var FClien *FlashClien

func init() {
	FClien, _ = NewFlashClient()
}

func NewFlashClient() (*FlashClien, error) {
	collector := colly.NewCollector(
		colly.AllowedDomains("theblockbeats.info", "www.theblockbeats.info"),
		colly.MaxDepth(2),
		colly.Async(true),
		colly.AllowURLRevisit(),
		// colly.CacheDir("./cached_files"),
	)
	if config.GlobalConfig.UseProxy {
		collector.SetProxy("htttp://127.0.0.1:7890")
	}
	collector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})
	FClien = &FlashClien{
		collector: collector,
	}
	return FClien, nil
}

func (s *FlashClien) GetFlashNews() *FlashNews {
	url := fmt.Sprintf("https://www.theblockbeats.info/newsflash")

	var title, summary string
	allNews := make([]News, 0)

	s.collector.Visit(url)
	
	s.collector.OnHTML("div.news-flash-wrapper", func(element *colly.HTMLElement) {
		title = element.ChildText("a.news-flash-item-title.news-flash-ios")
		title = strings.ReplaceAll(title, "\n", "") //去除回车
		// fmt.Println("aaaa--->",title)
		summary = element.ChildText("div.news-flash-item-content")
		// fmt.Println("bbb--->",summary)
		if title == "" {
			return
		}
		allNews = append(allNews, News{
			Title:          title,
			Summary:		summary,
			})
	})
	flashNews :=&FlashNews{
		News:&allNews,
		Date:time.Now().Format("2006-01-02"),
	}
	s.collector.Wait()

	return flashNews
}
