package services

import (
	"fmt"
	"time"
)

func PushRabbitMQ(c chan int) {
	time.Sleep(10*time.Second)
	ret := <-c
	fmt.Println("接收成功", ret)
}