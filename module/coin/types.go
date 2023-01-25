package coin

import (
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type CoinPrice struct {
	USD float32 `json:"usd"`
	CNY float32 `json:"cny"`
}
type SimpleSinglePrice struct {
	ID          string
	Currency    string
	MarketPrice decimal.Decimal
}
type Gecko struct {
	Client       *http.Client
	CoinList     []Coin
	ListUpdateAt time.Time
	ApiKey       string
}
type Coin struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}
type CoinMsg struct {
}
type CoinsMarket []CoinsMarketItem
type CoinsMarketItem struct {
	Coin
	Image                               string          		`json:"image"`
	CurrentPrice                        decimal.Decimal         `json:"current_price"`
	MarketCap                           decimal.Decimal 		`json:"market_cap"`
	MarketCapRank                       int16           		`json:"market_cap_rank"`
	TotalVolume                         decimal.Decimal 		`json:"total_volume"`
	High24                              decimal.Decimal         `json:"high_24h"`
	Low24                               decimal.Decimal         `json:"low_24h"`
	PriceChange24h                      decimal.Decimal         `json:"price_change_24h"`
	PriceChangePercentage24h            decimal.Decimal         `json:"price_change_percentage_24h"`
	MarketCapChange24h                  decimal.Decimal 		`json:"market_cap_change_24h"`
	MarketCapChangePercentage24h        decimal.Decimal         `json:"market_cap_change_percentage_24h"`
	CirculatingSupply                   decimal.Decimal 		`json:"circulating_supply"`
	TotalSupply                         decimal.Decimal 		`json:"total_supply"`
	FDV                                 decimal.Decimal 		`json:"fully_diluted_valuation"`
	ATH                                 decimal.Decimal         `json:"ath"`
	ATL                                 decimal.Decimal         `json:"atl"`
	ATHChangePercentage                 decimal.Decimal         `json:"ath_change_percentage"`
	ATHDate                             string          		`json:"ath_date"`
	ATLDate                             string          		`json:"atl_date"`
	ROI                                 *ROIItem        		`json:"roi"`
	LastUpdated                         string          		`json:"last_updated"`
	SparklineIn7d                       *SparklineItem  		`json:"sparkline_in_7d"`
	PriceChangePercentage1hInCurrency   *float64        		`json:"price_change_percentage_1h_in_currency"`
	PriceChangePercentage24hInCurrency  *float64        		`json:"price_change_percentage_24h_in_currency"`
	PriceChangePercentage7dInCurrency   *float64        		`json:"price_change_percentage_7d_in_currency"`
	PriceChangePercentage14dInCurrency  *float64        		`json:"price_change_percentage_14d_in_currency"`
	PriceChangePercentage30dInCurrency  *float64        		`json:"price_change_percentage_30d_in_currency"`
	PriceChangePercentage200dInCurrency *float64        		`json:"price_change_percentage_200d_in_currency"`
	PriceChangePercentage1yInCurrency   *float64        		`json:"price_change_percentage_1y_in_currency"`
}
type ROIItem struct {
	Times      float64 `json:"times"`
	Currency   string  `json:"currency"`
	Percentage float64 `json:"percentage"`
}

type SparklineItem struct {
	Price []float64 `json:"price"`
}

// OrderType in CoinGecko
type OrderType struct {
	MarketCapDesc string
	MarketCapAsc  string
	GeckoDesc     string
	GeckoAsc      string
	VolumeAsc     string
	VolumeDesc    string
}

// OrderTypeObject for certain order
var OrderTypeObject = &OrderType{
	MarketCapDesc: "market_cap_desc",
	MarketCapAsc:  "market_cap_asc",
	GeckoDesc:     "gecko_desc",
	GeckoAsc:      "gecko_asc",
	VolumeAsc:     "volume_asc",
	VolumeDesc:    "volume_desc",
}
