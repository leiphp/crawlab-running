package main

import (
	"colly_stock/eastmoney"
	"fmt"
	"log"
	"strconv"
	"time"
)


func main() {

	var err error
	//err = sse.GetStockListA("e:\\sseA.csv")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = sse.GetStockListB("e:\\sseB.csv")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//获取资讯列表
	articleList, err := eastmoney.GetArticleList("https://finance.eastmoney.com/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("articleList:",articleList,len(articleList))
	for k,v := range articleList {
		article,err := eastmoney.GetArticleDetail(v["url"])
		if err != nil {
			log.Fatal(err)
			break
		}
		v["content"]= article.Content
		v["create_date"]= article.CreateDate
		v["create_time"]=  strconv.FormatInt(time.Now().Unix(),10)
		v["source"]=  article.Source
		articleList[k] = v

	}
	fmt.Println("articleList2:",articleList,len(articleList))
}
