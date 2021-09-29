package coin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/SteinsElite/pickGinS/internal/token"
)

// This is a service response for maintain the cache of the price
const (
	endpoint = "https://api.coingecko.com/api/v3"

	price_interval = 60 * 60 // 1 hour

	trendDays     = "1"
	trendInterval = "hourly"
)

type ChartItem [2]float64
type CoinsMarketChart struct {
	ID           string       `json:"id"`
	Symbol       string       `json:"symbol"`
	Name         string       `json:"name"`
	Prices       *[]ChartItem `json:"prices"`
	MarketCaps   *[]ChartItem `json:"market_caps"`
	TotalVolumes *[]ChartItem `json:"total_volumes"`
}

type CoinInfo struct {
	Price float64   `json:"usd"`
	Rate  float64   `json:"usd_24h_change"`
	Trend []float64 `json:"trend"`
}

// client to interact with the gecko api to get price info
type CoinClient struct {
	httpClient *http.Client
	CoinCache  map[string]CoinInfo
}

func NewCoinClient() *CoinClient {
	coinCache := map[string]CoinInfo{
		token.BTCIds:  {Price: 41749.32},
		token.ETHIds:  {Price: 2898.42},
		token.USDTIds: {Price: 1.0},
		token.HTIds:   {Price: 7.73},
		token.MDXIds:  {Price: 1.12},
	}
	return &CoinClient{
		httpClient: &http.Client{},
		CoinCache:  coinCache,
	}
}

// doReq HTTP client
func doReq(req *http.Request, client *http.Client) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}

// MakeReq HTTP request helper
func (c *CoinClient) MakeReq(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}
	resp, err := doReq(req, c.httpClient)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (c *CoinClient) SimplePrice(ids string) (map[string]map[string]float64, error) {
	params := url.Values{}
	params.Add("ids", ids)
	params.Add("vs_currencies", token.Vs_currency)
	params.Add("include_24hr_change", "true")
	url := fmt.Sprintf("%s/simple/price?%s", endpoint, params.Encode())

	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}

	t := make(map[string]map[string]float64)
	err = json.Unmarshal(resp, &t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (c *CoinClient) CoinsMarketChart(ids string) (*CoinsMarketChart, error) {
	params := url.Values{}
	params.Add("vs_currency", token.Vs_currency)
	params.Add("days", trendDays)
	params.Add("interval", trendInterval)

	url := fmt.Sprintf("%s/coins/%s/market_chart?%s", endpoint, ids, params.Encode())
	resp, err := c.MakeReq(url)
	if err != nil {
		return nil, err
	}

	t := CoinsMarketChart{}
	err = json.Unmarshal(resp, &t)
	if err != nil {
		return &t, err
	}
	return &t, nil
}

func (c *CoinClient) GetLatestCoinInfo(ids string) (CoinInfo, error) {
	priceInfo, err := c.SimplePrice(ids)
	if err != nil {
		log.Println("[SimplePrice]: ", err)
		return CoinInfo{}, err
	}
	trendInfo, err := c.CoinsMarketChart(ids)
	if err != nil {
		log.Println("[CoinsMarket]: ", err)
		return CoinInfo{}, err
	}
	trend := []float64{}
	for i := range *(trendInfo.Prices) {
		trend = append(trend, (*trendInfo.Prices)[i][1])
	}
	coinInfo := CoinInfo{
		Price: priceInfo[ids][token.Vs_currency],
		Rate:  priceInfo[ids]["usd_24h_change"],
		Trend: trend,
	}
	return coinInfo, nil
}

func (c *CoinClient) updateCoinInfo() {
	for k := range c.CoinCache {
		coinInfo, err := c.GetLatestCoinInfo(k)
		if err == nil {
			c.CoinCache[k] = coinInfo
		}
	}
}

// This should run a infinite loop to maintain the coin info in a
// standlone goroutine, and it should start before the time ticker
func RunCoinInfoWatcher() {
	InitCoinClient()
	timeTicker := time.NewTicker(price_interval * time.Second)
	for {
		coinClient.updateCoinInfo()
		<-timeTicker.C
	}
}

var coinClient *CoinClient

func InitCoinClient() {
	coinClient = NewCoinClient()
	log.Println("finish init the CoinInfo watcher client")
}

// The api for other part to call
// get the specific coin price in the cache
func GetCurrentCoinPrice(coinIds string) float64 {
	return coinClient.CoinCache[coinIds].Price
}
