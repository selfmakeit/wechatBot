package stock

import (
	"fmt"
	"strconv"
	"strings"
	"wechatbot/config"
	"net/url"
	"github.com/gocolly/colly/v2"
)

type StockKeyStat struct {
	Name            string  `json:"stockName,omitempty"`
	Price           float32 `json:"price,omitempty"`
	PreviousClose   float32 `json:"previousClose,omitempty"`
	Change          float32 `json:"change,omitempty"`
	ChangePercent   float32 `json:"changePercent,omitempty"`
	DayRange        string  `json:"dayRange,omitempty"`
	YearRange       string  `json:"yearRange,omitempty"`
	Volume          string  `json:"volume,omitempty"`
	MarketCap       string  `json:"marketCap,omitempty"`
	PERatio         float32 `json:"peRatio,omitempty"`
	PrimaryExchange string  `json:"primaryExchange,omitempty"`
}

type StockNews struct {
	Title          string `json:"title,omitempty"`
	Source         string `json:"source,omitempty"`
	ArticleLink    string `json:"articleLink,omitempty"`
	Thumbnail_Link string `json:"thumbnailLink,omitempty"`
}

type StockClien struct {
	collector *colly.Collector
}

var SClien *StockClien

func init() {
	SClien, _ = NewClient()
}

func NewClient() (*StockClien, error) {
	collector := colly.NewCollector(
		colly.AllowedDomains("google.com", "www.google.com", "finance.google.com"),
		colly.MaxDepth(2),
		colly.Async(true),
		colly.AllowURLRevisit(),
		// colly.CacheDir("./cached_files"),
	)
	if config.GlobalConfig.UseProxy {
		collector.SetProxy("htttp://127.0.0.1:7890")
	}
	collector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})
	SClien = &StockClien{
		collector: collector,
	}
	return SClien, nil
}
func (s *StockClien) GetStockData(stock_query string) *StockKeyStat {

	url := fmt.Sprintf("https://finance.google.com/finance?q=%s", url.QueryEscape(stock_query))

	var name, dayRange, yearRange, volume, marketCap, primaryExchange string
	var price, previousClose, peRatio float32

	s.collector.Visit(url)

	s.collector.OnHTML("div.zzDege", func(element *colly.HTMLElement) {
		name = element.Text
	})

	s.collector.OnHTML("div.eYanAe > div:nth-child(2) > div", func(element *colly.HTMLElement) {
		text := strings.ReplaceAll(element.Text, ",", "")
		value, _ := strconv.ParseFloat(string([]rune(text)[1:]), 32)
		previousClose = float32(value)
	})

	s.collector.OnHTML("div.YMlKec.fxKbKc", func(element *colly.HTMLElement) {
		text := strings.ReplaceAll(element.Text, ",", "")
		value, _ := strconv.ParseFloat(string([]rune(text)[1:]), 32)
		price = float32(value)
	})

	s.collector.OnHTML("div.eYanAe > div:nth-child(2) > div.P6K39c", func(element *colly.HTMLElement) {
		dayRange = element.Text
	})

	s.collector.OnHTML("div.eYanAe > div:nth-child(3) > div.P6K39c", func(element *colly.HTMLElement) {
		yearRange = element.Text
	})

	for ctr := 4; ctr <= 12; ctr++ {

		s.collector.OnHTML(fmt.Sprintf("div.eYanAe > div:nth-child(%d)", ctr), func(element *colly.HTMLElement) {
			txt := element.ChildText("div.mfs7Fc")
			fmt.Println(">>>>",txt)
			contents := element.ChildText("div.P6K39c")
			if txt == "Avg Volume" {
				volume = contents
			} else if txt == "P/E ratio" {
				value, _ := strconv.ParseFloat(strings.ReplaceAll(contents, ",", ""), 32)
				peRatio = float32(value)
			} else if txt == "Primary exchange" {
				primaryExchange = contents
			} else if txt == "Market cap" {
				marketCap = contents
			}
		})
	}

	s.collector.Wait()

	if name == "" {
		return &StockKeyStat{}
	}

	stock := StockKeyStat{
		Name:            name,
		Price:           price,
		PreviousClose:   previousClose,
		Change:          price - previousClose,
		ChangePercent:   (((price - previousClose) / previousClose) * 100),
		DayRange:        dayRange,
		YearRange:       yearRange,
		MarketCap:       marketCap,
		Volume:          volume,
		PERatio:         peRatio,
		PrimaryExchange: primaryExchange,
	}

	return &stock
}

func (s *StockClien) GetStockNews(stock_query string) *[]StockNews {
	url := fmt.Sprintf("https://finance.google.com/finance?q=%s", url.QueryEscape(stock_query))

	var title, source, articleLink, thumbnailLink string
	allNews := make([]StockNews, 0)

	s.collector.Visit(url)
	s.collector.OnHTML(".zLrlHb.EA7tRd", func(element *colly.HTMLElement) {
		title = element.ChildText("div[class=F2KAFc]")
		title = strings.ReplaceAll(title, "\n", "") //去除回车
		source = element.ChildText("div[class=AYBNIb]")
		articleLink = element.ChildAttr("a[class=TxRU9d]", "href")
		thumbnailLink = element.ChildAttr("img.tLGtv", "src")
		if title == "" {
			return
		}

		allNews = append(allNews, StockNews{
			Title:          title,
			Source:         source,
			ArticleLink:    articleLink,
			Thumbnail_Link: thumbnailLink})
	})

	s.collector.Wait()

	return &allNews
}
