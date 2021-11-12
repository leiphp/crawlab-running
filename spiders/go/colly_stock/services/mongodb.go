package services

import (
	"colly_stock/datamodels"
	"colly_stock/initialize"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)


type CurlInfo struct {
	DNS float64 `json:"NAMELOOKUP_TIME"` //NAMELOOKUP_TIME
	TCP float64 `json:"CONNECT_TIME"`    //CONNECT_TIME - DNS
	SSL float64 `json:"APPCONNECT_TIME"` //APPCONNECT_TIME - CONNECT_TIME
}

type ConnectData struct {
	Latency  float64  `json:"latency"`
	RespCode int      `json:"respCode"`
	Url      string   `json:"url"`
	Detail   CurlInfo `json:"details"`
}

type Sensor struct {
	ISP       string
	Clientutc int64
	DataByAPP map[string]ConnectData
}

//增加:使用collection.InsertOne()来插入一条Document记录：
func InsertFinance(record datamodels.Article) (insertID primitive.ObjectID) {
	insertRest, err := initialize.Collection.InsertOne(context.TODO(), record)
	if err != nil {
		fmt.Println("InsertOne err:",err)
		return
	}
	insertID = insertRest.InsertedID.(primitive.ObjectID)
	return insertID
}

//查询:这里引入一个filter来匹配MongoDB数据库中的Document记录，使用bson.D类型来构建filter。
//使用collection.FindOne()来查询单个Document记录。这个方法返回一个可以解码为值的结果。
func queryFinance() ([]datamodels.Article,error){
	//查询一条记录
	//filter := bson.D{
	//	{"status", 0},
	//	{"id", "202111112178447978"},
	//}
	//
	//var article datamodels.Article
	//
	//err := initialize.Collection.FindOne(context.TODO(), filter).Decode(&article)
	//if err != nil {
	//	fmt.Printf("查询数据报错:%s\n", err.Error())
	//	return nil, err
	//}
	//fmt.Printf("Found a single document: %+v\n", article)
	//articles := make([]datamodels.Article,0)
	//articles = append(articles,article)
	//return articles,nil

	// 查询多个
	// 将选项传递给Find()
	findOptions := options.Find()
	findOptions.SetLimit(200)

	filter := bson.D{
		{"status", 0},
	}

	// 定义一个切片用来存储查询结果
	var results []datamodels.Article

	// 把bson.D{{}}作为一个filter来匹配所有文档
	cur, err := initialize.Collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// 查找多个文档返回一个光标
	// 遍历游标允许我们一次解码一个文档
	for cur.Next(context.TODO()) {
		// 创建一个值，将单个文档解码为该值
		var elem datamodels.Article
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	// 完成后关闭游标
	cur.Close(context.TODO())
	fmt.Printf("Found multiple documents num: %#v\n", len(results))
	return results,nil
}

//更新:使用collection.UpdateOne()更新单个Document记录。
func UpdateFinance(id string) (int64,error){
	//修改一条数据
	filter := bson.D{
		{"id", id},
	}

	//更新内容
	update := bson.M{
		"$set": bson.M{
			"status": 1,
		},
	}

	updateResult, err := initialize.Collection.UpdateOne(context.TODO(), filter, update)
	fmt.Println("updateResult:",updateResult)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return updateResult.ModifiedCount,nil
}

//删除:使用collection.DeleteOne()来删除一条记录，仍然使用刚才的filter。
func DeleteSensor(collection *mongo.Collection, ispBefore string, ispAfter string) {
	//筛选数据
	timestamp := time.Now().UTC().Unix()
	start := timestamp - 1800
	end := timestamp
	filter := bson.D{
		{"isp", ispBefore},
		{"$and", bson.A{
			bson.D{{"clientutc", bson.M{"$gte": start}}},
			bson.D{{"clientutc", bson.M{"$lte": end}}},
		}},
	}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	fmt.Println("deleteResult:",deleteResult)
}



