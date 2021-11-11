package eastmoney

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"log"
	"strings"
)

func GetArticleList(url string) (articleList []map[string]string, err error) {
	//articleList = make([]map[string]string,0)
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
		r.Headers.Set("Host", "finance.eastmoney.com")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", "https://finance.eastmoney.com/")
		//r.Headers.Set("Referer", "http://www.sse.com.cn/assortment/stock/list/share/") //关键头 如果没有 则返回 错误
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	})
	c.OnResponse(func(resp *colly.Response) {
		//stockList = string(resp.Body)
		//fmt.Println("stockList:",stockList)
		doc, err := htmlquery.Parse(strings.NewReader(string(resp.Body)))
		//fmt.Println("doc:",string(resp.Body))
		if err != nil {
			log.Fatal(err)
		}
		//获取列表
		nodes := htmlquery.Find(doc, `//div[@class="left"]`)
		if len(nodes) >1 {
			for _, node := range nodes {
				articleMap := make(map[string]string)
				url := htmlquery.FindOne(node, "./a/@href")
				title := htmlquery.FindOne(node, "./a/@title")
				fmt.Println("url:",htmlquery.InnerText(url),"title:",htmlquery.InnerText(title))
				//title := htmlquery.FindOne(node, `.//span[@class="title"]/text()`)
				//log.Println(strings.Split(htmlquery.InnerText(url), "/")[4], htmlquery.InnerText(title))
				articleMap["id"] = strings.Split(strings.Split(htmlquery.InnerText(url), "/")[4],".")[0]
				articleMap["url"] = htmlquery.InnerText(url)
				articleMap["title"] = htmlquery.InnerText(title)
				articleList = append(articleList, articleMap)
			}
		}
		//detail := htmlquery.Find(doc, `//div[@class='txtinfos']`)
		//fmt.Println("detail2:",detail)
		//fmt.Println("detail:",htmlquery.InnerText(detail[0]))

	})

	c.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
	})


	//因为最大深度设置2，
	//当前第一级 html里的 每个a标签都会回调访问
	//c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	//	link := e.Attr("href")
	//	//fmt.Println("link:",link)
	//	// 查找行首以 ?start=0&filter= 的字符串（非贪婪模式）
	//	reg := regexp.MustCompile(`(?U)^http://finance.eastmoney.com/a/(\d+).html`)
	//	regMatch := reg.FindAllString(link, -1)
	//	//如果找的到的话
	//	if(len(regMatch) > 0){
	//
	//		link = regMatch[0]
	//		fmt.Println("link2:",link)
	//		//访问该链接
	//		e.Request.Visit(link)
	//	}
	//
	//	// Visit link found on page
	//})

	//结束
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	err = c.Visit(url)

	return
}