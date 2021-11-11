package datamodels

type Article struct {
	ID                  string  `json:"id"` 					//ID
	Title          		string  `json:"title"`            	   //标题
	Url          		string  `json:"url"`            	   //URL
	Source              string  `json:"source"`                //来源
	Content             string  `json:"content"`               //内容
	Status              int     `json:"status"`                //状态
	ViewCount           int     `json:"view_count"`            //浏览量
	CategoryId          int     `json:"category_id"`           //分类ID
	Remark              string  `json:"remark"`                //备注
	CategoryName        string  `json:"category_name"`  	   //分类名称
	CreateTime          int64   `json:"create_time"`           //创建时间
	CreateDate          string  `json:"create_date"`           //发布日期
}
