package eastmoney

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"log"
	"strings"
)

type article struct {
	ID                  string  `json:"id"` 					//ID
	Title          		string  `json:"title"`            	   //标题
	Source              string  `json:"source"`                //来源
	Content             string  `json:"content"`               //内容
	Reward              int     `json:"reward"`                //奖励
	ViewCount           int     `json:"view_count"`            //浏览量
	CategoryId          int     `json:"category_id"`           //分类ID
	Remark              string  `json:"remark"`                //备注
	CategoryName        string  `json:"category_name"`  	   //分类名称
	CreateTime          int64   `json:"create_time"`           //创建时间
	CreateDate          string  `json:"create_date"`           //发布日期
}

func GetArticleDetail(url string) (detail article, err error) {
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

		detail.Content = htmlquery.OutputHTML(content,true)
		detail.CreateDate = htmlquery.InnerText(create_date)
		detail.Source = htmlquery.InnerText(source)
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