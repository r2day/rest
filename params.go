package rest

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
	"strconv"
	"strings"
)

// Params 解析URL 参数
type Params struct {
	// Page 页数
	Page int64 `json:"page"`
	// PerPage 每页展示记录数量
	PerPage int64 `json:"per_page"`
	// Sort 排序
	Sort string `json:"sort"`
	// 排序方式 DESC/AES
	Order string `json:"order"`

	Filter map[string]interface{} `json:"filter"`
}

func toSkip(a, b int64) int64 {
	return a / b
}

// Load 加载字符串参数
// filter={"category":"646cceafccc54408c8ccb308","category_id":"646cceafccc54408c8ccb308","status":false}&order=ASC&page=1&perPage=10&sort=id
func (p *Params) Load(payload string) *Params {

	// filter={}&range=[0,9]&sort=["id","ASC"]
	// 检测是否为空列表
	if strings.Contains(payload, "range") && !strings.Contains(payload, "perPage") {
		m, err := url.ParseQuery(payload)
		if err != nil {
			panic(err)
		}
		p = &Params{
			Page:    0,
			PerPage: 10,
			Sort:    "name",
			Order:   "ASC",
			Filter:  make(map[string]interface{}, 0),
		}
		queryRange := make([]int64, 2)
		if err := json.Unmarshal([]byte(m["range"][0]), &queryRange); err != nil {
			panic(err)
		}
		// 设定范围
		// 相当于mongo limit
		perPage := queryRange[1] - queryRange[0]
		p.PerPage = perPage

		querySort := make([]string, 2)
		if err := json.Unmarshal([]byte(m["sort"][0]), &querySort); err != nil {
			panic(err)
		}
		p.Sort = querySort[0]
		p.Order = querySort[1]

		// 设定页面
		// 相当于skip
		p.Page = queryRange[0] / perPage

		var filter map[string]interface{}
		if err := json.Unmarshal([]byte(m["filter"][0]), &filter); err != nil {
			panic(err)
		}

		// 移除默认过滤
		delete(filter, "status")

		p.Filter = filter
		return p
	} else {
		// filter={"category":"646cceafccc54408c8ccb308","category_id":"646cceafccc54408c8ccb308","status":false}&order=ASC&page=1&perPage=10&sort=id
		m, err := url.ParseQuery(payload)
		if err != nil {
			panic(err)
		}

		pageNumber, err := strconv.ParseInt(m["page"][0], 10, 32)
		if err == nil {
			p.Page = pageNumber
		}

		perPageNumber, err := strconv.ParseInt(m["perPage"][0], 10, 32)
		if err == nil {
			p.PerPage = perPageNumber
		}

		p.Order = m["order"][0]

		p.Sort = m["sort"][0]

		var filter map[string]interface{}
		if err := json.Unmarshal([]byte(m["filter"][0]), &filter); err != nil {
			panic(err)
		}

		p.Filter = filter
		fmt.Println("===>", p)
		return p
	}
}

// ToMongoOptions 加载字符串参数
// filter={"category":"646cceafccc54408c8ccb308","category_id":"646cceafccc54408c8ccb308","status":false}&order=ASC&page=1&perPage=10&sort=id
func (p *Params) ToMongoOptions() *options.FindOptions {
	// 进行必要分页处理
	opt := options.Find()
	sortValue := -1
	if p.Order == "ASC" {
		sortValue = 1
	}
	// 设定排序
	opt.SetSort(bson.M{p.Sort: sortValue})
	// 设定分页
	opt.SetSkip(p.Page)
	opt.SetLimit(p.PerPage)
	return opt
}

// ToMongoFilter 加载字符串参数
// filter={"category":"646cceafccc54408c8ccb308","category_id":"646cceafccc54408c8ccb308","status":false}&order=ASC&page=1&perPage=10&sort=id
func (p *Params) ToMongoFilter(merchantID string, accessLevel uint) bson.D {
	filters := bson.D{{Key: "meta.merchant_id", Value: merchantID},
		{"meta.access_level", bson.D{{"$lte", accessLevel}}}}

	for key, val := range p.Filter {
		filter := bson.E{Key: key, Value: val}
		filters = append(filters, filter)
	}
	return filters
}
