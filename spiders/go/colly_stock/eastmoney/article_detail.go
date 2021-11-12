package eastmoney

import (
	"colly_stock/datamodels"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"log"
	"strings"
)


func GetArticleDetail(url string) (detail datamodels.Article, err error) {
	c := colly.NewCollector()

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.108 Safari/537.36"
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", "finance.eastmoney.com")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", url)
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	})
	c.OnResponse(func(resp *colly.Response) {
		doc, err := htmlquery.Parse(strings.NewReader(string(resp.Body)))
		//fmt.Println("doc:",string(resp.Body))
		if err != nil {
			log.Fatal(err)
		}
		//获取列表
		content := htmlquery.FindOne(doc, `//*[@id="ContentBody"]`)
		//fmt.Println("content:",htmlquery.InnerText(content))
		//fmt.Println("content2:",htmlquery.OutputHTML(content,true))
		create_date := htmlquery.FindOne(doc, `//*[@id="topbox"]/div[3]/div[1]/div[1]`)
		source := htmlquery.FindOne(doc, `//*[@id="topbox"]/div[3]/div[1]/div[2]`)
		if content != nil {
			detail.Content = htmlquery.OutputHTML(content,true)
		}else{
			detail.Content = ""
		}
		if create_date != nil {
			detail.CreateDate = htmlquery.InnerText(create_date)
		}else{
			detail.CreateDate = ""
		}
		if source != nil {
			detail.Source = htmlquery.InnerText(source)
		}else{
			detail.Source = ""
		}
	})


	c.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
	})

	//结束
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	err = c.Visit(url)

	return
}