package coin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"wechatbot/common"
	. "wechatbot/config"
	"wechatbot/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcom "github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

var baseURL = "https://api.coingecko.com/api/v3"
var GekClient *Gecko

func init() {
	GekClient, _ = NewClient(GlobalConfig.GeckoApiKey)
}

func NewClient(apiKey string) (*Gecko, error) {
	var httpClient *http.Client
	if GlobalConfig.GeckoUseProxy {
		url, _ := url.Parse("htttp://127.0.0.1:7890")
		t := &http.Transport{
			MaxIdleConns:    10,
			MaxConnsPerHost: 10,
			IdleConnTimeout: time.Duration(10) * time.Second,
			Proxy:           http.ProxyURL(url),
		}
		httpClient = &http.Client{
			Transport: t,
			Timeout:   time.Duration(10) * time.Second,
		}
	} else {
		httpClient = &http.Client{
			Timeout: time.Second * 100,
		}
	}
	var key string
	if apiKey != "" {
		key = apiKey
	}
	var gecko = &Gecko{
		ApiKey:       key,
		CoinList:     nil,
		Client:       httpClient,
		ListUpdateAt: time.Now(),
	}
	gecko.UpdateCoinList()
	return gecko, nil
}
func (g *Gecko) MakeReq(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := doReq(req, g.Client)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func doReq(req *http.Request, client *http.Client) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}

// Get token list to pair token symbol with token id.
func (g *Gecko) UpdateCoinList() {
	url := fmt.Sprintf("%s/coins/list", baseURL)
	if g.ApiKey != "" {
		url = fmt.Sprintf("%s/coins/list?x_cg_pro_api_key=%s", baseURL, g.ApiKey)
	}
	log.Println("Init Coin List......")
	resp, err := g.MakeReq(url)
	if err != nil {
		log.Println("Update Coin List Error(get list error)")
		panic(err)
	}
	err = json.Unmarshal(resp, &g.CoinList)
	if err != nil {
		log.Println("Update Coin List Error(parse list error)")
		panic(err)
	}
	log.Println("Available coins:", len(g.CoinList))
	log.Println("Coin List Update Success!")
	// for _, token := range g.CoinList {
	// 	sy := strings.ToLower(token.Symbol)
	// 	id := strings.ToLower(token.ID)
	// 	name := strings.ToLower(token.Name)
	// 	if strings.Contains(sy, ".e") || strings.Contains(id, ".e") || strings.Contains(name, ".e") {
	// 		fmt.Println("===========================================")
	// 		fmt.Println(token.ID, ">>", token.Name, ">>", token.Symbol)
	// 		fmt.Println("===========================================")
	// 	}
	// 	if strings.Contains(token.ID, "wormhole") {
	// 		fmt.Println(token.ID, ">1>", token.Name, ">1>", token.Symbol)
	// 	}
	// }
	g.ListUpdateAt = time.Now()
}
/*
??????symbol??????id
symbol???eth
id :ethereum
*/
func (g *Gecko) GetCoinId(symbol string) (string, error) {
	if g.CoinList == nil || len(g.CoinList) == 0 {
		return "", errors.New("Coin List Not Update")
	}
	// find the id
	for _, token := range g.CoinList {
		//ignore case
		if strings.EqualFold(token.Symbol, symbol) {
			// fmt.Printf("find coin==>ID:%v,Symbol:%v,Name:%v", token.ID, token.Symbol, token.Name)
			return token.ID, nil
		}
	}
	return "", fmt.Errorf("symbol not in the token list,you provided symbol:[%v]", symbol)
}
/*
??????id??????symbol
symbol???eth
id :ethereum
*/
func (g *Gecko) GetCoinSymbol(id string) (string, error) {
	if g.CoinList == nil || len(g.CoinList) == 0 {
		return "", errors.New("Coin List Not Update")
	}
	// find the id
	for _, token := range g.CoinList {
		//ignore case
		if strings.EqualFold(token.ID, id) {
			// fmt.Printf("find coin==>ID:%v,Symbol:%v,Name:%v", token.ID, token.Symbol, token.Name)
			return token.Symbol, nil
		}
	}
	return "", fmt.Errorf("id not in the token list,you provided id:[%v]", id)
}
func (g *Gecko) GetCoinMsgBySymbol(symbol string) (*Coin, error) {
	if g.CoinList == nil || len(g.CoinList) == 0 {
		return nil, errors.New("Coin List Not Update")
	}
	// find the id
	for _, token := range g.CoinList {
		//ignore case
		if strings.EqualFold(token.Symbol, symbol) {
			// fmt.Printf("find coin==>ID:%v,Symbol:%v,Name:%v", token.ID, token.Symbol, token.Name)
			fmt.Println()
			return &token, nil
		}
	}
	return nil, fmt.Errorf("symbol not in the token list,you provided symbol:[%v]", symbol)
}
func (g *Gecko) GetCoinIdByNetWork(symbol string, network string) (string, error) {
	if g.CoinList == nil || len(g.CoinList) == 0 {
		return "", errors.New("Coin List Not Update")
	}
	// ?????????????????????symbol?????????.e???????????????gecko???????????????????????????????????????
	// if network == common.AvalancheChainName {
	// 	symbol = strings.Split(symbol, ".")[0]
	// }
	// ?????????????????????: todo
	// find the id
	for _, token := range g.CoinList {
		// wormhole ?????????solana???????????????
		// wormholes ?????????id???????????????'-wormhole'??????????????????symbol??????
		if strings.Contains(token.ID, "wormhole") {
			continue
		}

		if strings.EqualFold(token.Symbol, symbol) {
			return token.ID, nil
		}
	}
	return "", fmt.Errorf("symbol not in the token list,you provided symbol:[%v]", symbol)
}
/**
???????????????
*/
func (g *Gecko) CoinPushToWechat()(string ,error) {
	params := url.Values{}
	params.Add("ids", "ethereum,apecoin,solana,aptos,near,arweave,filecoin")
	params.Add("vs_currencies", "usd")
	url := fmt.Sprintf("%s/simple/price?%s", baseURL, params.Encode())
	if g.ApiKey != "" {
		url = fmt.Sprintf("%s/simple/price?x_cg_pro_api_key=%s&%s", baseURL, g.ApiKey, params.Encode())
	}
	resp, err := g.MakeReq(url)
	if err != nil {
		return "", err
	}
	var v map[string]map[string]decimal.Decimal
	err = json.Unmarshal(resp, &v)
	if err != nil {
		return "", err
	}
	
	if len(v) == 0 {
		return "", fmt.Errorf("id or currency not existed caontained in watch list")
	}
	var msg ="watch list:\n\n"
	for i, token := range v {
		for _, price := range token {
			sym,_ :=g.GetCoinSymbol(i)
			msg += sym+"("+i+")" +" : $"+price.String()+"\n"
		}
	}
	msg += "\n(????????????9?????????10?????????)"
	return msg, nil
}
/*
example:

	price, err := g.SimpleSinglePrice("bitcoin", "usd")
*/
func (g *Gecko) GetCoinPriceById(id string, currency string) (decimal.Decimal, error) {
	idParam := []string{strings.ToLower(id)}
	vcParam := []string{strings.ToLower(currency)}
	mapp, err := g.SimplePrice(idParam, vcParam)
	curr := (mapp)[id]
	if len(curr) == 0 {
		return decimal.Zero, fmt.Errorf("id or currency not existed(id:%v)", id)
	}
	if err != nil {
		return decimal.Zero, err
	}
	for _, token := range mapp {
		for _, price := range token {
			return price, nil
		}
	}
	return decimal.Zero, nil
}

/*
example:

	ids := []string{"bitcoin", "ethereum"}
	vc := []string{"usd", "myr"}
	sp, err := g.SimplePrice(ids, vc)
*/
func (g *Gecko) SimplePrice(ids []string, vsCurrencies []string) (map[string]map[string]decimal.Decimal, error) {
	params := url.Values{}
	idsParam := strings.Join(ids[:], ",")
	vsCurrenciesParam := strings.Join(vsCurrencies[:], ",")

	params.Add("ids", idsParam)
	params.Add("vs_currencies", vsCurrenciesParam)

	url := fmt.Sprintf("%s/simple/price?%s", baseURL, params.Encode())
	if g.ApiKey != "" {
		url = fmt.Sprintf("%s/simple/price?x_cg_pro_api_key=%s&%s", baseURL, g.ApiKey, params.Encode())
	}
	resp, err := g.MakeReq(url)
	if err != nil {
		return nil, err
	}
	var v map[string]map[string]decimal.Decimal
	err = json.Unmarshal(resp, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// Some tokens like usdc.e on avalanche cannot be found by coingecko list, need special process.
func (g *Gecko) GetPriceBySymbol(symbol string, currency string) (decimal.Decimal, error) {
	if symbol == "" {
		return decimal.Zero, errors.New("symbol must not be empty")
	}
	id, err := g.GetCoinId(symbol)
	if err != nil {
		return decimal.Zero, err
	}
	return g.GetCoinPriceById(id, currency)
}

// Return token price.
//
// Some tokens like usdc.e on avalanche cannot be found by coingecko list, need special process.
func (g *Gecko) GetPriceByAddress(address string, network string, currency string, client bind.ContractBackend) (decimal.Decimal, error) {
	if address == "" {
		return decimal.Zero, errors.New("address must not be empty")
	}
	if address == "0x0000000000000000000000000000000000000000" {
		return decimal.Zero, errors.New("address must not be zero")
	}
	token, err := NewErc20(ethcom.HexToAddress(address), client)
	if err != nil {
		return decimal.Zero, err
	}
	symbol, err := token.Symbol(nil)
	if err != nil {
		return decimal.Zero, err
	}
	id, err := g.GetCoinId(symbol)
	if err != nil {
		return decimal.Zero, err
	}
	return g.GetCoinPriceById(id, currency)
}

// Return chain token price.
func (g *Gecko) GetChainTokenPrice(network string, currency string) (decimal.Decimal, error) {
	return g.GetPriceBySymbol(common.ChainTokenSymbolList[network], currency)
}

// CoinsMarket /coins/market

func (g *Gecko) SimpleCoinMarketBySymbol(vsCurrency string, symbol string) (*CoinsMarketItem, error) {
	if len(vsCurrency) == 0 {
		return nil, fmt.Errorf("vs_currency is required")
	}
	id, err := g.GetCoinId(symbol)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add("ids", id)
	// vs_currency
	params.Add("vs_currency", vsCurrency)
	// order
	order := OrderTypeObject.MarketCapDesc
	params.Add("order", order)
	// ids
	// per_page
	params.Add("per_page", Int2String(1))
	params.Add("page", Int2String(1))
	// sparkline
	params.Add("sparkline", Bool2String(false))
	// price_change_percentage
	url := fmt.Sprintf("%s/coins/markets?%s", baseURL, params.Encode())
	resp, err := g.MakeReq(url)
	if err != nil {
		return nil, err
	}
	var data *CoinsMarket
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return &(*data)[0], nil
}
func (g *Gecko) SimpleCoinMarketById(vsCurrency string, id string) (*CoinsMarketItem, error) {
	if len(vsCurrency) == 0 {
		return nil, fmt.Errorf("vs_currency is required")
	}
	id = utils.TrimAllSpace(id)
	if id == "" || id == "0" {
		return nil, errors.New("id can't be nil")
	}
	params := url.Values{}
	params.Add("ids", id)
	// vs_currency
	params.Add("vs_currency", vsCurrency)
	// order
	order := OrderTypeObject.MarketCapDesc
	params.Add("order", order)
	// ids
	// per_page
	params.Add("per_page", Int2String(1))
	params.Add("page", Int2String(1))
	// sparkline
	params.Add("sparkline", Bool2String(false))
	// price_change_percentage
	url := fmt.Sprintf("%s/coins/markets?%s", baseURL, params.Encode())
	resp, err := g.MakeReq(url)
	if err != nil {
		return nil, err
	}
	var data *CoinsMarket
	err = json.Unmarshal(resp, &data)
	if err != nil{
		return nil, err
	}
	if len(*data)==0{
		return nil, errors.New("did not get result")
	}
	return &(*data)[0], nil
}

/*
example:

	vsCurrency := "usd"
	ids := []string{"bitcoin", "ethereum", "steem"}
	perPage := 1
	page := 1
	sparkline := true
	pcp := geckoTypes.PriceChangePercentageObject
	priceChangePercentage := []string{pcp.PCP1h, pcp.PCP24h, pcp.PCP7d, pcp.PCP14d, pcp.PCP30d, pcp.PCP200d, pcp.PCP1y}
	order := geckoTypes.OrderTypeObject.MarketCapDesc
	market, err := cg.CoinsMarket(vsCurrency, ids, order, perPage, page, sparkline, priceChangePercentage)
*/
func (g *Gecko) CoinsMarket(vsCurrency string, ids []string, order string, perPage int, page int, sparkline bool, priceChangePercentage []string) (*CoinsMarket, error) {
	if len(vsCurrency) == 0 {
		return nil, fmt.Errorf("vs_currency is required")
	}
	params := url.Values{}
	// vs_currency
	params.Add("vs_currency", vsCurrency)
	// order
	if len(order) == 0 {
		order = OrderTypeObject.MarketCapDesc
	}
	params.Add("order", order)
	// ids
	if len(ids) != 0 {
		idsParam := strings.Join(ids[:], ",")
		params.Add("ids", idsParam)
	}
	// per_page
	if perPage <= 0 || perPage > 250 {
		perPage = 100
	}
	params.Add("per_page", Int2String(perPage))
	params.Add("page", Int2String(page))
	// sparkline
	params.Add("sparkline", Bool2String(sparkline))
	// price_change_percentage
	if len(priceChangePercentage) != 0 {
		priceChangePercentageParam := strings.Join(priceChangePercentage[:], ",")
		params.Add("price_change_percentage", priceChangePercentageParam)
	}
	url := fmt.Sprintf("%s/coins/markets?%s", baseURL, params.Encode())
	resp, err := g.MakeReq(url)
	if err != nil {
		return nil, err
	}
	var data *CoinsMarket
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
