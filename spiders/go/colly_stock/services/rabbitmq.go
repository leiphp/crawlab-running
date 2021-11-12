package services

import (
	"fmt"
	"time"
)

//执行MQ处理数据
func PushRabbitMQ(c chan int) {
	time.Sleep(10*time.Second)
	ret := <-c
	fmt.Println("接收成功", ret)
}