package services

import (
	"colly_stock/datamodels"
	"colly_stock/initialize"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"sync"
	"time"
)

var wg sync.WaitGroup

//执行MQ处理数据
func PushRabbitMQ(c chan int) {
	articls,err := queryFinance()
	if err != nil {
		log.Println("err:",err.Error())
	}
	if len(articls) > 0 {
		for _, article := range articls{
			wg.Add(1) // 启动一个goroutine就登记+1
			go pushFinanceArt(article)
		}

	}
	wg.Wait() // 等待所有登记的goroutine都结束
	time.Sleep(3*time.Second)
	ret := <-c
	fmt.Println("接收成功", ret)
}


//推送金融文章
func pushFinanceArt(article datamodels.Article) {
	//新建channel
	visitLogCh, err := initialize.MqClient.GetChannel()
	if err != nil {
		log.Println("[推送金融文章-创建channel失败]-[%s]", err.Error())
		return
	}
	defer visitLogCh.Close()

	//绑定队列
	q, queueErr := visitLogCh.QueueDeclare(
		"lxtkj_finance_article_queue", //name visitLogQueue: "visit_log_queue"
		true,  //durable
		false, //delete when usused
		false, //exclusive
		false, //nowait
		nil,   //argments
	)
	if queueErr != nil {
		log.Println("[推送金融文章-创建queue失败]-[%s]", err.Error())
		return
	}

	//封装body
	jsonBytes, err := json.Marshal(article)
	if err != nil {
		log.Println("[推送金融文章-json格式化失败]-[%s]", err.Error())
		return
	}

	//推送信息
	pushErr := visitLogCh.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        jsonBytes,
		})

	if pushErr != nil {
		log.Println("[推送金融文章-推送失败]-[%s]", err.Error())
		return
	}
	log.Println("[推送金融文章-推送成功]")

	//推送文章成功，标记状态
	res,err := UpdateFinance(article.ID)
	if err != nil {
		log.Println("[更新金融文章失败]-[%s]", err.Error())
	}
	fmt.Println("res:",res)
	defer wg.Done() // goroutine结束就登记-1
}