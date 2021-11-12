package main

import (
	"colly_stock/eastmoney"
	"colly_stock/initialize"
	"colly_stock/services"
	"fmt"
	"log"
	"strconv"
	"time"
)


func main() {
	initialize.Init()
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

		article.ID = v["id"]
		article.Title = v["title"]
		article.Url = v["url"]
		article.CreateTime = time.Now().Unix()
		article.Status = 0
		id := services.InsertFinance(article)
		fmt.Println("return id is:",id)

	}
	//fmt.Println("articleList2:",articleList,len(articleList))
	//关闭连接
	//defer func() {
	//	err = initialize.Client.Disconnect(context.TODO())
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println("Connection to MongoDB closed.")
	//}()

	//异步MQ推送数据到线上
	ch := make(chan int)
	go services.PushRabbitMQ(ch)
	ch <- 10
	fmt.Println("发送成功")

	fmt.Println("程序执行完成！")
	//time.Sleep(60*time.Second)

	defer services.Disconnect(initialize.Client)
}
