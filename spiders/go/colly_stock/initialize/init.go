package initialize

import (
	"context"
	"fmt"
	//"go.mongodb.org/mongo-driver/bson"    //BOSN解析包
	"go.mongodb.org/mongo-driver/mongo"    //MongoDB的Go驱动包
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	Client     *mongo.Client     //数据库客户端
	Collection *mongo.Collection //客户端
)

//	提供系统初始化，全局变量
func Init() {
	var err error
	// 设置mongoDB客户端连接信息
	param := fmt.Sprintf("mongodb://10.250.200.223:27017")
	clientOptions := options.Client().ApplyURI(param)

	// 建立客户端连接
	Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}

	// 检查连接情况
	err = Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}
	fmt.Println("Connected to MongoDB!")

	//指定要操作的数据集
	Collection = Client.Database("crawlab_test").Collection("results_eastmoney")
	//执行增删改查操作

	//// 断开客户端连接
	//err = client.Disconnect(context.TODO())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Connection to MongoDB closed.")
}
