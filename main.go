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

var urls = []string{
	/*	bitflyer
		API Doc: https://lightning.bitflyer.jp/docs?lang=ja
	*/
	"https://api.bitflyer.jp/v1/ticker?product_code=BTC_JPY",
	"https://api.bitflyer.jp/v1/board?product_code=BTC_JPY",
	"https://api.bitflyer.jp/v1/ticker?product_code=FX_BTC_JPY",
	"https://api.bitflyer.jp/v1/board?product_code=FX_BTC_JPY",
	"https://api.bitflyer.jp/v1/ticker?product_code=ETH_BTC",
	"https://api.bitflyer.jp/v1/board?product_code=ETH_BTC",
	"https://api.bitflyer.jp/v1/ticker?product_code=BCH_BTC",
	"https://api.bitflyer.jp/v1/board?product_code=BCH_BTC",
	"https://api.bitflyer.jp/v1/ticker?product_code=BTCJPY_MAT1WK",
	"https://api.bitflyer.jp/v1/board?product_code=BTCJPY_MAT1WK",
	"https://api.bitflyer.jp/v1/ticker?product_code=BTCJPY_MAT2WK",
	"https://api.bitflyer.jp/v1/board?product_code=BTCJPY_MAT2WK",

	/*	bitbank
		API Doc: https://docs.bitbank.cc/
	*/
	"https://public.bitbank.cc/btc_jpy/ticker",
	"https://public.bitbank.cc/btc_jpy/depth",
	"https://public.bitbank.cc/xrp_jpy/ticker",
	"https://public.bitbank.cc/xrp_jpy/depth",
	"https://public.bitbank.cc/mona_jpy/ticker",
	"https://public.bitbank.cc/mona_jpy/depth",
	"https://public.bitbank.cc/bcc_jpy/ticker",
	"https://public.bitbank.cc/bcc_jpy/depth",

	/*	BTCBOX
		API Doc: https://www.btcbox.co.jp/help/asm
	*/
	"https://www.btcbox.co.jp/api/v1/ticker/",
	"https://www.btcbox.co.jp/api/v1/depth/",

	/*	coincheck
		API Doc: https://coincheck.com/ja/documents/exchange/api
	*/
	"https://coincheck.com/api/ticker",
	"https://coincheck.com/api/order_books",

	/*	Zaif
		API Doc: https://corp.zaif.jp/api-docs/
	*/
	"https://api.zaif.jp/api/1/ticker/btc_jpy",
	"https://api.zaif.jp/api/1/depth/btc_jpy",
	"https://api.zaif.jp/api/1/ticker/xem_jpy",
	"https://api.zaif.jp/api/1/depth/xem_jpy",
	"https://api.zaif.jp/api/1/ticker/mona_jpy",
	"https://api.zaif.jp/api/1/depth/mona_jpy",

	/*	Kraken
		API Doc: https://www.kraken.com/help/api
	*/
	"https://api.kraken.com/0/public/Ticker?pair=XETHZGBP",
	"https://api.kraken.com/0/public/Depth?pair=XETHZGBP",
	"https://api.kraken.com/0/public/Ticker?pair=GNOEUR",
	"https://api.kraken.com/0/public/Depth?pair=GNOEUR",
	"https://api.kraken.com/0/public/Ticker?pair=GNOXBT",
	"https://api.kraken.com/0/public/Depth?pair=GNOXBT",
	"https://api.kraken.com/0/public/Ticker?pair=XICNXXBT",
	"https://api.kraken.com/0/public/Depth?pair=XICNXXBT",
	"https://api.kraken.com/0/public/Ticker?pair=XXLMZEUR",
	"https://api.kraken.com/0/public/Depth?pair=XXLMZEUR",
	"https://api.kraken.com/0/public/Ticker?pair=XXLMZUSD",
	"https://api.kraken.com/0/public/Depth?pair=XXLMZUSD",
	"https://api.kraken.com/0/public/Ticker?pair=XETCXETH",
	"https://api.kraken.com/0/public/Depth?pair=XETCXETH",
	"https://api.kraken.com/0/public/Ticker?pair=XXBTZGBP",
	"https://api.kraken.com/0/public/Depth?pair=XXBTZGBP",
	"https://api.kraken.com/0/public/Ticker?pair=XXMRXXBT",
	"https://api.kraken.com/0/public/Depth?pair=XXMRXXBT",
	"https://api.kraken.com/0/public/Ticker?pair=XXMRZUSD",
	"https://api.kraken.com/0/public/Depth?pair=XXMRZUSD",
	"https://api.kraken.com/0/public/Ticker?pair=DASHEUR",
	"https://api.kraken.com/0/public/Depth?pair=DASHEUR",
	"https://api.kraken.com/0/public/Ticker?pair=XXBTZJPY",
	"https://api.kraken.com/0/public/Depth?pair=XXBTZJPY",
	"https://api.kraken.com/0/public/Ticker?pair=XXRPZEUR",
	"https://api.kraken.com/0/public/Depth?pair=XXRPZEUR",
	"https://api.kraken.com/0/public/Ticker?pair=XMLNXETH",
	"https://api.kraken.com/0/public/Depth?pair=XMLNXETH",
	"https://api.kraken.com/0/public/Ticker?pair=GNOETH",
	"https://api.kraken.com/0/public/Depth?pair=GNOETH",
	"https://api.kraken.com/0/public/Ticker?pair=XETHZEUR",
	"https://api.kraken.com/0/public/Depth?pair=XETHZEUR",
	"https://api.kraken.com/0/public/Ticker?pair=XETHZJPY",
	"https://api.kraken.com/0/public/Depth?pair=XETHZJPY",
	"https://api.kraken.com/0/public/Ticker?pair=XETHZUSD",
	"https://api.kraken.com/0/public/Depth?pair=XETHZUSD",
	"https://api.kraken.com/0/public/Ticker?pair=XZECXXBT",
	"https://api.kraken.com/0/public/Depth?pair=XZECXXBT",
	"https://api.kraken.com/0/public/Ticker?pair=XICNXETH",
	"https://api.kraken.com/0/public/Depth?pair=XICNXETH",
	"https://api.kraken.com/0/public/Ticker?pair=XLTCXXBT",
	"https://api.kraken.com/0/public/Depth?pair=XLTCXXBT",
	"https://api.kraken.com/0/public/Ticker?pair=XLTCZEUR",
	"https://api.kraken.com/0/public/Depth?pair=XLTCZEUR",
	"https://api.kraken.com/0/public/Ticker?pair=XMLNXXBT",
	"https://api.kraken.com/0/public/Depth?pair=XMLNXXBT",
	"https://api.kraken.com/0/public/Ticker?pair=XREPZEUR",
	"https://api.kraken.com/0/public/Depth?pair=XREPZEUR",
	"https://api.kraken.com/0/public/Ticker?pair=XXLMXXBT",
	"https://api.kraken.com/0/public/Depth?pair=XXLMXXBT",
	"https://api.kraken.com/0/public/Ticker?pair=XXRPZCAD",
	"https://api.kraken.com/0/public/Depth?pair=XXRPZCAD",
	"https://api.kraken.com/0/public/Ticker?pair=BCHUSD",
	"https://api.kraken.com/0/public/Depth?pair=BCHUSD",
	"https://api.kraken.com/0/public/Ticker?pair=XREPXETH",
	"https://api.kraken.com/0/public/Depth?pair=XREPXETH",
	"https://api.kraken.com/0/public/Ticker?pair=XREPXXBT",
	"https://api.kraken.com/0/public/Depth?pair=XREPXXBT",
	"https://api.kraken.com/0/public/Ticker?pair=BCHXBT",
	"https://api.kraken.com/0/public/Depth?pair=BCHXBT",
	"https://api.kraken.com/0/public/Ticker?pair=XXRPZJPY",
	"https://api.kraken.com/0/public/Depth?pair=XXRPZJPY",
	"https://api.kraken.com/0/public/Ticker?pair=EOSETH",
	"https://api.kraken.com/0/public/Depth?pair=EOSETH",
	"https://api.kraken.com/0/public/Ticker?pair=XETCXXBT",
	"https://api.kraken.com/0/public/Depth?pair=XETCXXBT",
	"https://api.kraken.com/0/public/Ticker?pair=XETCZUSD",
	"https://api.kraken.com/0/public/Depth?pair=XETCZUSD",
	"https://api.kraken.com/0/public/Ticker?pair=XXRPXXBT",
	"https://api.kraken.com/0/public/Depth?pair=XXRPXXBT",
	"https://api.kraken.com/0/public/Ticker?pair=EOSUSD",
	"https://api.kraken.com/0/public/Depth?pair=EOSUSD",
	"https://api.kraken.com/0/public/Ticker?pair=USDTZUSD",
	"https://api.kraken.com/0/public/Depth?pair=USDTZUSD",
	"https://api.kraken.com/0/public/Ticker?pair=XLTCZUSD",
	"https://api.kraken.com/0/public/Depth?pair=XLTCZUSD",
	"https://api.kraken.com/0/public/Ticker?pair=XREPZUSD",
	"https://api.kraken.com/0/public/Depth?pair=XREPZUSD",
	"https://api.kraken.com/0/public/Ticker?pair=XXBTZCAD",
	"https://api.kraken.com/0/public/Depth?pair=XXBTZCAD",
	"https://api.kraken.com/0/public/Ticker?pair=XXMRZEUR",
	"https://api.kraken.com/0/public/Depth?pair=XXMRZEUR",
	"https://api.kraken.com/0/public/Ticker?pair=XXRPZUSD",
	"https://api.kraken.com/0/public/Depth?pair=XXRPZUSD",
	"https://api.kraken.com/0/public/Ticker?pair=EOSEUR",
	"https://api.kraken.com/0/public/Depth?pair=EOSEUR",
	"https://api.kraken.com/0/public/Ticker?pair=XETHZCAD",
	"https://api.kraken.com/0/public/Depth?pair=XETHZCAD",
	"https://api.kraken.com/0/public/Ticker?pair=XXBTZEUR",
	"https://api.kraken.com/0/public/Depth?pair=XXBTZEUR",
	"https://api.kraken.com/0/public/Ticker?pair=GNOUSD",
	"https://api.kraken.com/0/public/Depth?pair=GNOUSD",
	"https://api.kraken.com/0/public/Ticker?pair=XETCZEUR",
	"https://api.kraken.com/0/public/Depth?pair=XETCZEUR",
	"https://api.kraken.com/0/public/Ticker?pair=XETHXXBT",
	"https://api.kraken.com/0/public/Depth?pair=XETHXXBT",
	"https://api.kraken.com/0/public/Ticker?pair=XXBTZUSD",
	"https://api.kraken.com/0/public/Depth?pair=XXBTZUSD",
	"https://api.kraken.com/0/public/Ticker?pair=BCHEUR",
	"https://api.kraken.com/0/public/Depth?pair=BCHEUR",
	"https://api.kraken.com/0/public/Ticker?pair=DASHUSD",
	"https://api.kraken.com/0/public/Depth?pair=DASHUSD",
	"https://api.kraken.com/0/public/Ticker?pair=EOSXBT",
	"https://api.kraken.com/0/public/Depth?pair=EOSXBT",
	"https://api.kraken.com/0/public/Ticker?pair=XZECZEUR",
	"https://api.kraken.com/0/public/Depth?pair=XZECZEUR",
	"https://api.kraken.com/0/public/Ticker?pair=XXDGXXBT",
	"https://api.kraken.com/0/public/Depth?pair=XXDGXXBT",
	"https://api.kraken.com/0/public/Ticker?pair=XZECZUSD",
	"https://api.kraken.com/0/public/Depth?pair=XZECZUSD",

	/*	Fisco
		API Doc: https://fcce.jp/api-docs
	*/
	"https://api.fcce.jp/api/1/ticker/btc_jpy",
	"https://api.fcce.jp/api/1/depth/btc_jpy",
	"https://api.fcce.jp/api/1/ticker/mona_jpy",
	"https://api.fcce.jp/api/1/depth/mona_jpy",

	/*	FIREX
		API Doc: https://firex.jp/api-docs
	*/
	"https://api.firex.jp/api/1/ticker/btc_jpy",
	"https://api.firex.jp/api/1/depth/btc_jpy",

	/*	Bitstamp
		API Doc: https://www.bitstamp.net/api/
	*/
	"https://www.bitstamp.net/api/v2/ticker/btcusd",
	"https://www.bitstamp.net/api/v2/order_book/btcusd",
	"https://www.bitstamp.net/api/v2/ticker/btceur",
	"https://www.bitstamp.net/api/v2/order_book/btceur",
	"https://www.bitstamp.net/api/v2/ticker/eurusd",
	"https://www.bitstamp.net/api/v2/order_book/eurusd",
	"https://www.bitstamp.net/api/v2/ticker/xrpusd",
	"https://www.bitstamp.net/api/v2/order_book/xrpusd",
	"https://www.bitstamp.net/api/v2/ticker/xrpeur",
	"https://www.bitstamp.net/api/v2/order_book/xrpeur",
	"https://www.bitstamp.net/api/v2/ticker/ltcusd",
	"https://www.bitstamp.net/api/v2/order_book/ltcusd",
	"https://www.bitstamp.net/api/v2/ticker/ltceur",
	"https://www.bitstamp.net/api/v2/order_book/ltceur",
	"https://www.bitstamp.net/api/v2/ticker/ethusd",
	"https://www.bitstamp.net/api/v2/order_book/ethusd",
	"https://www.bitstamp.net/api/v2/ticker/etheur",
	"https://www.bitstamp.net/api/v2/order_book/etheur",

	/*	BTCC
		API Doc: https://www.btcc.com/apidocs/usd-spot-exchange-market-data-rest-api
	*/
	"https://spotusd-data.btcc.com/data/pro/ticker?symbol=btcusd",
	"https://spotusd-data.btcc.com/data/pro/orderbook?symbol=btcusd",

	/*	OKCoin
		https://www.okcoin.com/rest_api.html
	*/
	"https://www.okcoin.com/api/v1/ticker.do?symbol=btc_usd",
	"https://www.okcoin.com/api/v1/depth.do?symbol=btc_usd",
	"https://www.okcoin.com/api/v1/ticker.do?symbol=ltc_usd",
	"https://www.okcoin.com/api/v1/depth.do?symbol=ltc_usd",
	"https://www.okcoin.com/api/v1/ticker.do?symbol=eth_usd",
	"https://www.okcoin.com/api/v1/depth.do?symbol=eth_usd",
	"https://www.okcoin.com/api/v1/ticker.do?symbol=etc_usd",
	"https://www.okcoin.com/api/v1/depth.do?symbol=etc_usd",
	"https://www.okcoin.com/api/v1/ticker.do?symbol=bcc_usd",
	"https://www.okcoin.com/api/v1/depth.do?symbol=bcc_usd",

	/*	huobi
		API Doc: https://github.com/huobiapi/API_Docs_en/wiki
	*/
	"http://api.huobi.com/staticmarket/ticker_btc_json.js",
	"http://api.huobi.com/staticmarket/depth_btc_json.js",
	"http://api.huobi.com/staticmarket/ticker_ltc_json.js",
	"http://api.huobi.com/staticmarket/depth_ltc_json.js",
	"http://api.huobi.com/usdmarket/ticker_btc_json.js",
	"http://api.huobi.com/usdmarket/depth_btc_json.js",

	/*	Bitfinex
		API Doc: https://docs.bitfinex.com/v1/reference#rest-public-ticker
	*/
	"https://api.bitfinex.com/v1/pubticker/btcusd",
	"https://api.bitfinex.com/v1/book/btcusd",
	"https://api.bitfinex.com/v1/pubticker/ltcusd",
	"https://api.bitfinex.com/v1/book/ltcusd",
	"https://api.bitfinex.com/v1/pubticker/ethusd",
	"https://api.bitfinex.com/v1/book/ethusd",
	"https://api.bitfinex.com/v1/pubticker/etcusd",
	"https://api.bitfinex.com/v1/book/etcusd",
	"https://api.bitfinex.com/v1/pubticker/rrtusd",
	"https://api.bitfinex.com/v1/book/rrtusd",
	"https://api.bitfinex.com/v1/pubticker/zecusd",
	"https://api.bitfinex.com/v1/book/zecusd",
	"https://api.bitfinex.com/v1/pubticker/xmrusd",
	"https://api.bitfinex.com/v1/book/xmrusd",
	"https://api.bitfinex.com/v1/pubticker/dshusd",
	"https://api.bitfinex.com/v1/book/dshusd",
	"https://api.bitfinex.com/v1/pubticker/bccusd",
	"https://api.bitfinex.com/v1/book/bccusd",
	"https://api.bitfinex.com/v1/pubticker/bcuusd",
	"https://api.bitfinex.com/v1/book/bcuusd",
	"https://api.bitfinex.com/v1/pubticker/xrpusd",
	"https://api.bitfinex.com/v1/book/xrpusd",
	"https://api.bitfinex.com/v1/pubticker/iotusd",
	"https://api.bitfinex.com/v1/book/iotusd",
	"https://api.bitfinex.com/v1/pubticker/eosusd",
	"https://api.bitfinex.com/v1/book/eosusd",
	"https://api.bitfinex.com/v1/pubticker/sanusd",
	"https://api.bitfinex.com/v1/book/sanusd",
	"https://api.bitfinex.com/v1/pubticker/omgusd",
	"https://api.bitfinex.com/v1/book/omgusd",
	"https://api.bitfinex.com/v1/pubticker/bchusd",
	"https://api.bitfinex.com/v1/book/bchusd",

	/*
		poloniex
		AIP Doc: https://poloniex.com/support/api/
	*/
	"https://poloniex.com/public?command=returnTicker",
	"https://poloniex.com/public?command=returnOrderBook&currencyPair=USDT_REP",
	"https://poloniex.com/public?command=returnOrderBook&currencyPair=USDT_ZEC",
	"https://poloniex.com/public?command=returnOrderBook&currencyPair=USDT_ETH",
	"https://poloniex.com/public?command=returnOrderBook&currencyPair=USDT_BTC",
	"https://poloniex.com/public?command=returnOrderBook&currencyPair=USDT_ETC",
	"https://poloniex.com/public?command=returnOrderBook&currencyPair=USDT_BCH",
	"https://poloniex.com/public?command=returnOrderBook&currencyPair=USDT_DASH",
	"https://poloniex.com/public?command=returnOrderBook&currencyPair=USDT_NXT",
	"https://poloniex.com/public?command=returnOrderBook&currencyPair=USDT_LTC",
	"https://poloniex.com/public?command=returnOrderBook&currencyPair=USDT_XMR",
	"https://poloniex.com/public?command=returnOrderBook&currencyPair=USDT_XRP",
	"https://poloniex.com/public?command=returnOrderBook&currencyPair=USDT_STR",

	/*
		bithumb
		API Doc: https://www.bithumb.com/u1/US127
	*/
	"https://api.bithumb.com/public/ticker/BTC",
	"https://api.bithumb.com/public/orderbook/BTC",
	"https://api.bithumb.com/public/ticker/ETH",
	"https://api.bithumb.com/public/orderbook/ETH",
	"https://api.bithumb.com/public/ticker/DASH",
	"https://api.bithumb.com/public/orderbook/DASH",
	"https://api.bithumb.com/public/ticker/LTC",
	"https://api.bithumb.com/public/orderbook/LTC",
	"https://api.bithumb.com/public/ticker/ETC",
	"https://api.bithumb.com/public/orderbook/ETC",
	"https://api.bithumb.com/public/ticker/XRP",
	"https://api.bithumb.com/public/orderbook/XRP",
	"https://api.bithumb.com/public/ticker/BCH",
	"https://api.bithumb.com/public/orderbook/BCH",

	/*
		Bittrex
		API Doc: https://bittrex.com/home/api
	*/
	"https://bittrex.com/api/v1.1/public/getticker?market=USDT-BCC",
	"https://bittrex.com/api/v1.1/public/getorderbook?market=USDT-BCC&type=both",
	"https://bittrex.com/api/v1.1/public/getticker?market=USDT-BTC",
	"https://bittrex.com/api/v1.1/public/getorderbook?market=USDT-BTC&type=both",
	"https://bittrex.com/api/v1.1/public/getticker?market=USDT-DASH",
	"https://bittrex.com/api/v1.1/public/getorderbook?market=USDT-DASH&type=both",
	"https://bittrex.com/api/v1.1/public/getticker?market=USDT-ETC",
	"https://bittrex.com/api/v1.1/public/getorderbook?market=USDT-ETC&type=both",
	"https://bittrex.com/api/v1.1/public/getticker?market=USDT-ETH",
	"https://bittrex.com/api/v1.1/public/getorderbook?market=USDT-ETH&type=both",
	"https://bittrex.com/api/v1.1/public/getticker?market=USDT-LTC",
	"https://bittrex.com/api/v1.1/public/getorderbook?market=USDT-LTC&type=both",
	"https://bittrex.com/api/v1.1/public/getticker?market=USDT-NEO",
	"https://bittrex.com/api/v1.1/public/getorderbook?market=USDT-NEO&type=both",
	"https://bittrex.com/api/v1.1/public/getticker?market=USDT-XMR",
	"https://bittrex.com/api/v1.1/public/getorderbook?market=USDT-XMR&type=both",
	"https://bittrex.com/api/v1.1/public/getticker?market=USDT-XRP",
	"https://bittrex.com/api/v1.1/public/getorderbook?market=USDT-XRP&type=both",
	"https://bittrex.com/api/v1.1/public/getticker?market=USDT-ZEC",
	"https://bittrex.com/api/v1.1/public/getorderbook?market=USDT-ZEC&type=both",

	/*
		HitBTC
		API Doc: https://hitbtc.com/api
	*/
	"https://api.hitbtc.com/api/1/public/BTCUSD/ticker",
	"https://api.hitbtc.com/api/1/public/BTCUSD/orderbook",
	"https://api.hitbtc.com/api/1/public/BTCEUR/ticker",
	"https://api.hitbtc.com/api/1/public/BTCEUR/orderbook",
	"https://api.hitbtc.com/api/1/public/LTCUSD/ticker",
	"https://api.hitbtc.com/api/1/public/LTCUSD/orderbook",
	"https://api.hitbtc.com/api/1/public/LTCEUR/ticker",
	"https://api.hitbtc.com/api/1/public/LTCEUR/orderbook",
	"https://api.hitbtc.com/api/1/public/ETHEUR/ticker",
	"https://api.hitbtc.com/api/1/public/ETHEUR/orderbook",
	"https://api.hitbtc.com/api/1/public/LSKEUR/ticker",
	"https://api.hitbtc.com/api/1/public/LSKEUR/orderbook",
	"https://api.hitbtc.com/api/1/public/STEEMEUR/ticker",
	"https://api.hitbtc.com/api/1/public/STEEMEUR/orderbook",

	/*
		Gemini
		API Doc: https://docs.gemini.com/rest-api/
	*/
	"https://api.gemini.com/v1/pubticker/btcusd",
	"https://api.gemini.com/v1/book/btcusd",
	"https://api.gemini.com/v1/pubticker/ethusd",
	"https://api.gemini.com/v1/book/ethusd",
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

func saveFile(data []byte, id int, dir string) error {
	//idStr := fmt.Sprintf("%04d", id)
	//dir := fmt.Sprintf("contents/%s", idStr)
	os.MkdirAll(dir, os.FileMode(0700))
	filePath := path.Join(dir, fmt.Sprintf("%04d.json", id))
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
	if len(os.Args) != 2 {
		fmt.Println("Zipファイルの出力先を指定してください。")
		os.Exit(1)
	}
	output := os.Args[1]

	now := time.Now().UTC()
	fmt.Println(now)
	date := now.Format("20060102-1504")
	dir := path.Join(output, date)

	errors := make(chan error)
	threads := make(chan struct{}, 64)

	for i, url := range urls {
		id := i + 1
		//fmt.Printf("%04d: %s\n", id, url)
		go func(id int, url string, dir string) {
			threads <- struct{}{}
			defer func() { <-threads }()

			content, err := getUrlContent(url)
			if err == nil {
				errors <- saveFile(content, id, dir)
			} else {
				errors <- err
			}
		}(id, url, dir)
	}

	successes := 0
	failures := 0
	for range urls {
		error := <-errors
		if error == nil {
			successes++
		} else {
			failures++
			fmt.Println(error)
		}
	}
	fmt.Printf("Total:%d Successes:%d Failures:%d\n", len(urls), successes, failures)

	zipPath := fmt.Sprintf("%s.zip", dir)
	zipit(dir, zipPath)
	os.RemoveAll(dir)
	fmt.Printf("Create: %s", zipPath)
}
