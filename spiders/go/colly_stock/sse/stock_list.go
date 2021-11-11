package sse

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

/*GetStockListA 获取上海证券交易所股票列表
A股
*/
func GetStockListA(saveFile string) (err error) {

	stocks, err := getStockList("http://query.sse.com.cn/security/stock/downloadStockListFile.do?csrcCode=&stockCode=&areaName=&stockType=1")
	if err != nil {
		return err
	}

	err = saveStockList2CSV(stocks, saveFile)
	return
}

/*GetStockListB 获取上海证券交易所股票列表
B股
*/
func GetStockListB(saveFile string) (err error) {
	stocks, err := getStockList("http://query.sse.com.cn/security/stock/downloadStockListFile.do?csrcCode=&stockCode=&areaName=&stockType=2")
	fmt.Println("stockB:\n",stocks)
	if err != nil {
		return err
	}
	err = saveStockList2CSV(stocks, saveFile)
	return
}

func getStockList(url string) (stockList string, err error) {

	//GET http://query.sse.com.cn/security/stock/downloadStockListFile.do?csrcCode=&stockCode=&areaName=&stockType=1 HTTP/1.1
	//Host: query.sse.com.cn
	//Connection: keep-alive
	//Accept: */*
	//Origin: http://www.sse.com.cn
	//User-Agent: Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.108 Safari/537.36
	//Referer: http://www.sse.com.cn/assortment/stock/list/share/
	//Accept-Encoding: gzip, deflate
	//Accept-Language: zh-CN,zh;q=0.9`

	c := colly.NewCollector()

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.108 Safari/537.36"
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", "query.sse.com.cn")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", "http://www.sse.com.cn")
		r.Headers.Set("Referer", "http://www.sse.com.cn/assortment/stock/list/share/") //关键头 如果没有 则返回 错误
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	})
	c.OnResponse(func(resp *colly.Response) {
		stockList = string(resp.Body)
		fmt.Println("stockList:",stockList)
	})

	c.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
	})

	err = c.Visit(url)

	return
}


func saveStockList2CSV(stockList string, file string) (err error) {

	vals := strings.Split(stockList, "\n")

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	fw := csv.NewWriter(f)

	for _, row := range vals {

		rSplits := strings.Split(row, "\t")

		rSplitsRslt := make([]string, 0)
		for _, sp := range rSplits {
			trimSp := strings.Trim(sp, " ")
			if len(trimSp) > 0 {
				rSplitsRslt = append(rSplitsRslt, trimSp)
			}
		}
		if len(rSplitsRslt) > 0 {
			err = fw.Write(rSplitsRslt)
			if err != nil {
				return err
			}
		}
	}
	fw.Flush()

	return
}
