package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type ApiType int

func (a ApiType) Label() string {
	return API_LABELS[a]
}

const (
	TICKER ApiType = iota
	ORDERBOOK
	TRADES
	API_COUNT
)

var API_LABELS = [API_COUNT]string{
	"ticker",
	"orderbook",
	"trades",
}

type ApiUrls [API_COUNT]string

type Url struct {
	Service  string
	ApiType  ApiType
	Currency string
	Url      string
}

type Service struct {
	Name       string
	Currencies []string
	ApiUrls    ApiUrls
}

func (s Service) Urls() (urls []Url) {
	for _, currency := range s.Currencies {
		for api, apiUrl := range s.ApiUrls {
			if len(apiUrl) == 0 {
				continue
			}
			url := apiUrl
			if len(currency) > 0 {
				url = fmt.Sprintf(apiUrl, currency)
			}
			urls = append(urls, Url{
				Service:  s.Name,
				ApiType:  ApiType(api),
				Currency: currency,
				Url:      url,
			})
		}
	}
	return urls
}

type Services []Service

func (ss Services) validate() error {
	errorStrings := []string{}
	// Check overlap service's currency
	for _, service := range services {
		currencies := map[string]int{}
		for _, currency := range service.Currencies {
			if _, exist := currencies[currency]; exist {
				errorStrings = append(errorStrings, fmt.Sprintf(
					"ERROR: %s overlap in %s.\n", currency, service.Name))
			} else {
				currencies[currency] = 0
			}
		}
	}
	if len(errorStrings) > 0 {
		return errors.New(strings.Join(errorStrings, ""))
	}
	return nil
}

var services = Services{
	// https://lightning.bitflyer.jp/docs?lang=ja
	Service{
		Name: "bitflyer",
		Currencies: []string{
			"BTC_JPY",
			"FX_BTC_JPY",
			"ETH_BTC",
			"BCH_BTC",
			"BTCJPY_MAT1WK",
			"BTCJPY_MAT2WK"},
		ApiUrls: ApiUrls{
			"https://api.bitflyer.jp/v1/ticker?product_code=%s",
			"https://api.bitflyer.jp/v1/board?product_code=%s"},
	},
	// https://docs.bitbank.cc/
	Service{
		Name: "bitbank",
		Currencies: []string{
			"btc_jpy",
			"xrp_jpy",
			"mona_jpy",
			"bcc_jpy"},
		ApiUrls: ApiUrls{
			"https://public.bitbank.cc/%s/ticker",
			"https://public.bitbank.cc/%s/ticker"},
	},
	// https://www.btcbox.co.jp/help/asm
	Service{
		Name: "BTCBOX",
		Currencies: []string{
			""},
		ApiUrls: ApiUrls{
			"https://www.btcbox.co.jp/api/v1/ticker/",
			"https://www.btcbox.co.jp/api/v1/depth/"},
	},
	// https://coincheck.com/ja/documents/exchange/api
	Service{
		Name: "coincheck",
		Currencies: []string{
			""},
		ApiUrls: ApiUrls{
			"https://coincheck.com/api/ticker",
			"https://coincheck.com/api/order_books"},
	},
	// https://corp.zaif.jp/api-docs/
	Service{
		Name: "zaif",
		Currencies: []string{
			"btc_jpy",
			"xem_jpy",
			"mona_jpy"},
		ApiUrls: ApiUrls{
			"https://api.zaif.jp/api/1/ticker/%s",
			"https://api.zaif.jp/api/1/depth/%s"},
	},
	// https://www.kraken.com/help/api
	Service{
		Name: "Kraken",
		Currencies: []string{
			"XETHZGBP",
			"GNOEUR",
			"GNOXBT",
			"XICNXXBT",
			"XXLMZEUR",
			"XXLMZUSD",
			"XETCXETH",
			"XXBTZGBP",
			"XXMRXXBT",
			"XXMRZUSD",
			"DASHEUR",
			"XXBTZJPY",
			"XXRPZEUR",
			"XMLNXETH",
			"GNOETH",
			"XETHZEUR",
			"XETHZJPY",
			"XETHZUSD",
			"XZECXXBT",
			"XICNXETH",
			"XLTCXXBT",
			"XLTCZEUR",
			"XMLNXXBT",
			"XREPZEUR",
			"XXLMXXBT",
			"XXRPZCAD",
			"BCHUSD",
			"XREPXETH",
			"XREPXXBT",
			"BCHXBT",
			"XXRPZJPY",
			"EOSETH",
			"XETCXXBT",
			"XETCZUSD",
			"XXRPXXBT",
			"EOSUSD",
			"USDTZUSD",
			"XLTCZUSD",
			"XREPZUSD",
			"XXBTZCAD",
			"XXMRZEUR",
			"XXRPZUSD",
			"EOSEUR",
			"XETHZCAD",
			"XXBTZEUR",
			"GNOUSD",
			"XETCZEUR",
			"XETHXXBT",
			"XXBTZUSD",
			"BCHEUR",
			"DASHUSD",
			"EOSXBT",
			"XZECZEUR",
			"XXDGXXBT",
			"XZECZUSD"},
		ApiUrls: ApiUrls{
			"https://api.kraken.com/0/public/Ticker?pair=%s",
			"https://api.kraken.com/0/public/Depth?pair=%s"},
	},
	// https://fcce.jp/api-docs
	Service{
		Name: "Fisco",
		Currencies: []string{
			"btc_jpy",
			"mona_jpy"},
		ApiUrls: ApiUrls{
			"",
			"https://api.fcce.jp/api/1/ticker/%s",
			"https://api.fcce.jp/api/1/depth/%s"},
	},
	// https://firex.jp/api-docs
	Service{
		Name: "FIREX",
		Currencies: []string{
			"btc_jpy"},
		ApiUrls: ApiUrls{
			"https://api.firex.jp/api/1/ticker/%s",
			"https://api.firex.jp/api/1/depth/%s"},
	},
	// https://www.bitstamp.net/api/
	Service{
		Name: "Bitstamp",
		Currencies: []string{
			"btcusd",
			"btceur",
			"eurusd",
			"xrpusd",
			"xrpeur",
			"ltcusd",
			"ltceur",
			"ethusd",
			"etheur"},
		ApiUrls: ApiUrls{
			"https://www.bitstamp.net/api/v2/ticker/%s",
			"https://www.bitstamp.net/api/v2/order_book/%s"},
	},
	// https://www.btcc.com/apidocs/usd-spot-exchange-market-data-rest-api
	Service{
		Name: "BTCC",
		Currencies: []string{
			"btcusd"},
		ApiUrls: ApiUrls{
			"https://spotusd-data.btcc.com/data/pro/ticker?symbol=%s",
			"https://spotusd-data.btcc.com/data/pro/orderbook?symbol=%s"},
	},
	// https://www.okcoin.com/rest_api.html
	Service{
		Name: "OKCoin",
		Currencies: []string{
			"btc_usd",
			"ltc_usd",
			"eth_usd",
			"etc_usd",
			"bcc_usd"},
		ApiUrls: ApiUrls{
			"https://www.okcoin.com/api/v1/ticker.do?symbol=%s",
			"https://www.okcoin.com/api/v1/depth.do?symbol=%s"},
	},
	// https://github.com/huobiapi/API_Docs_en/wiki
	Service{
		Name: "huobistatic",
		Currencies: []string{
			"btc",
			"ltc"},
		ApiUrls: ApiUrls{
			"http://api.huobi.com/staticmarket/ticker_%s_json.js",
			"http://api.huobi.com/staticmarket/depth_%s_json.js"},
	},
	Service{
		Name: "huobiusd",
		Currencies: []string{
			"btc"},
		ApiUrls: ApiUrls{
			"http://api.huobi.com/usdmarket/ticker_%s_json.js",
			"http://api.huobi.com/usdmarket/depth_%s_json.js"},
	},
	// https://docs.bitfinex.com/v1/reference#rest-public-ticker
	Service{
		Name: "Bitfinex",
		Currencies: []string{
			"btcusd",
			"ltcusd",
			"ethusd",
			"etcusd",
			"rrtusd",
			"zecusd",
			"xmrusd",
			"dshusd",
			"bccusd",
			"bcuusd",
			"xrpusd",
			"iotusd",
			"eosusd",
			"sanusd",
			"omgusd",
			"bchusd"},
		ApiUrls: ApiUrls{
			"https://api.bitfinex.com/v1/pubticker/%s",
			"https://api.bitfinex.com/v1/book/%s"},
	},
	// https://poloniex.com/support/api/
	Service{
		Name: "poloniexTicker",
		Currencies: []string{
			"",
		},
		ApiUrls: ApiUrls{
			"https://poloniex.com/public?command=returnTicker",
		},
	},
	Service{
		Name: "poloniex",
		Currencies: []string{
			"USDT_REP",
			"USDT_ZEC",
			"USDT_ETH",
			"USDT_BTC",
			"USDT_ETC",
			"USDT_BCH",
			"USDT_DASH",
			"USDT_NXT",
			"USDT_LTC",
			"USDT_XMR",
			"USDT_XRP",
			"USDT_STR"},
		ApiUrls: ApiUrls{
			"",
			"https://poloniex.com/public?command=returnOrderBook&currencyPair=%s"},
	},
	// https://www.bithumb.com/u1/US127
	Service{
		Name: "bithumb",
		Currencies: []string{
			"BTC",
			"ETH",
			"DASH",
			"LTC",
			"ETC",
			"XRP",
			"BCH"},
		ApiUrls: ApiUrls{
			"https://api.bithumb.com/public/ticker/%s",
			"https://api.bithumb.com/public/orderbook/%s"},
	},
	// https://bittrex.com/home/api
	Service{
		Name: "Bittrex",
		Currencies: []string{
			"USDT-BCC",
			"USDT-BTC",
			"USDT-DASH",
			"USDT-ETC",
			"USDT-ETH",
			"USDT-LTC",
			"USDT-NEO",
			"USDT-XMR",
			"USDT-XRP",
			"USDT-ZEC"},
		ApiUrls: ApiUrls{
			"https://bittrex.com/api/v1.1/public/getticker?market=%s",
			"https://bittrex.com/api/v1.1/public/getorderbook?market=%s&type=both"},
	},
	// https://hitbtc.com/api
	Service{
		Name: "HitBTC",
		Currencies: []string{
			"BTCUSD",
			"BTCEUR",
			"LTCUSD",
			"LTCEUR",
			"ETHEUR",
			"LSKEUR",
			"STEEMEUR"},
		ApiUrls: ApiUrls{
			"https://api.hitbtc.com/api/1/public/%s/ticker",
			"https://api.hitbtc.com/api/1/public/%s/orderbook"},
	},
	// https://docs.gemini.com/rest-api/
	Service{
		Name: "Gemini",
		Currencies: []string{
			"btcusd",
			"ethusd"},
		ApiUrls: ApiUrls{
			"https://api.gemini.com/v1/pubticker/%s",
			"https://api.gemini.com/v1/book/%s"},
	},
}

func getUrlContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	if resp.StatusCode != http.StatusOK {
		errorString := fmt.Sprintf("Not StatusOK. Code:%d URL:%s", resp.StatusCode, url)
		return []byte{}, errors.New(errorString)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func saveFile(data []byte, dir, fileName string) error {
	os.MkdirAll(dir, os.FileMode(0700))
	filePath := path.Join(dir, fileName)
	err := ioutil.WriteFile(filePath, data, os.FileMode(0600))
	return err
}

func zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

func main() {
	err := services.validate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(os.Args) != 2 {
		fmt.Println("Zipファイルの出力先を指定してください。")
		os.Exit(1)
	}
	output := os.Args[1]

	now := time.Now().UTC()
	fmt.Println(now)
	date := now.Format("20060102_1504")
	dir := path.Join(output, date)

	errors := make(chan error)
	threads := make(chan struct{}, 64)
	urlCount := 0

	for _, service := range services {
		for _, url := range service.Urls() {
			urlCount++
			go func(url Url, dir, date string) {
				threads <- struct{}{}
				defer func() { <-threads }()

				content, err := getUrlContent(url.Url)
				if err == nil {
					fileName := fmt.Sprintf(
						"%s_%s_%s_%s.json",
						date,
						url.Service,
						url.Currency,
						url.ApiType.Label())
					errors <- saveFile(content, dir, fileName)
				} else {
					errors <- err
				}
			}(url, dir, date)
		}
	}

	successes := 0
	failures := 0
	for i := 0; i < urlCount; i++ {
		error := <-errors
		if error == nil {
			successes++
		} else {
			failures++
			fmt.Println(error)
		}
	}
	fmt.Printf("Total:%d Successes:%d Failures:%d\n", urlCount, successes, failures)

	zipPath := fmt.Sprintf("%s.zip", dir)
	zipit(dir, zipPath)
	os.RemoveAll(dir)
	fmt.Println("Create:", zipPath)
}
