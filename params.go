package rest

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	Translate map[string]string

	Filter map[string]interface{} `json:"filter"`

	// 使用IN 方式查询
	FilterIn map[string][]string `json:"filer_in"`

	HasFilterIn bool `json:"has_filter_in"`

	MongoIDList []*primitive.ObjectID `json:"mongo_id_list"`
}

func toSkip(a, b int64) int64 {
	return a / b
}

func (p *Params) Init() *Params {
	initP := &Params{
		Page:        0,
		PerPage:     10,
		Sort:        "name",
		Order:       "ASC",
		Filter:      make(map[string]interface{}, 0),
		HasFilterIn: false,
		FilterIn:    make(map[string][]string, 0),
		MongoIDList: make([]*primitive.ObjectID, 0),
		Translate:   make(map[string]string, 0),
	}
	return initP
}

// Load 加载字符串参数
// filter={"category":"646cceafccc54408c8ccb308","category_id":"646cceafccc54408c8ccb308","status":false}&order=ASC&page=1&perPage=10&sort=id
func (p *Params) Load(payload string) *Params {
	// filter={}&range=[0,9]&sort=["id","ASC"]
	// 检测是否为空列表
	if strings.Contains(payload, "range") && !strings.Contains(payload, "perPage") {
		m, err := url.ParseQuery(payload)
		if err != nil {
			return p
		}
		queryRange := make([]int64, 2)
		if err := json.Unmarshal([]byte(m["range"][0]), &queryRange); err != nil {
			//panic(err)
			return p
		}
		// 设定范围
		// 相当于mongo limit
		perPage := queryRange[1] - queryRange[0]
		p.PerPage = perPage

		querySort := make([]string, 2)
		if err := json.Unmarshal([]byte(m["sort"][0]), &querySort); err != nil {
			//panic(err)
			return p
		}
		p.Sort = querySort[0]
		p.Order = querySort[1]

		// 设定页面
		// 相当于skip
		if queryRange[0] > perPage {
			p.Page = queryRange[0] / perPage
		}

		var filter map[string]interface{}
		if err := json.Unmarshal([]byte(m["filter"][0]), &filter); err != nil {
			//panic(err)
			return p
		}

		// 移除默认过滤
		delete(filter, "status")

		p.Filter = filter
		return p
	} else if strings.Contains(payload, "perPage") && strings.Contains(payload, "page") {
		// filter={"category":"646cceafccc54408c8ccb308","category_id":"646cceafccc54408c8ccb308","status":false}&order=ASC&page=1&perPage=10&sort=id
		m, err := url.ParseQuery(payload)
		if err != nil {
			//panic(err)
			return p
		}
		log.WithField("m", m).Warning("+++++++++++")

		pageSlice := m["page"]
		if len(pageSlice) > 0 {
			pageNumber, err := strconv.ParseInt(pageSlice[0], 10, 32)
			if err == nil {
				p.Page = pageNumber
			}
		}

		perPageSlice := m["perPage"]
		if len(perPageSlice) > 0 {
			perPageNumber, err := strconv.ParseInt(perPageSlice[0], 10, 32)
			if err == nil {
				p.PerPage = perPageNumber
			}
		}

		orderSlice := m["order"]
		if len(orderSlice) > 0 {
			p.Order = orderSlice[0]
		}

		sortSlice := m["sort"]
		if len(sortSlice) > 0 {
			p.Sort = sortSlice[0]
		}

		var filter map[string]interface{}
		b := m["filter"]
		if len(b) != 0 {
			if err := json.Unmarshal([]byte(b[0]), &filter); err != nil {
				return p
			}

			p.Filter = filter
			fmt.Println("===>", p)
			return p
		}

		log.WithField("b", b).Warning("+++++++++++")

		return p

	} else {
		// http://localhost:9988/v1/auth/merchant/rolecategory?filter={"id":["646dd319a8cd2902283f0b79"]}
		m, err := url.ParseQuery(payload)
		if err != nil {
			//panic(err)
			return p
		}

		var filterIn map[string][]string
		b := m["filter"]
		log.WithField("b", b).Warning("+++++bb++++++")
		if len(b) != 0 {

			if err := json.Unmarshal([]byte(b[0]), &filterIn); err != nil {
				log.WithField("filterIn", filterIn).
					Warning("+++++filterIn++++++")
				return p
			}

			p.FilterIn = filterIn
			p.HasFilterIn = true
			return p
		}

		objIds := make([]*primitive.ObjectID, 0)
		ids := filterIn["id"]
		log.WithField("ids", ids).Warning("+++++++++++")

		//for _, i := range ids {
		//	objID, _ := primitive.ObjectIDFromHex(i)
		//	objIds = append(objIds, &objID)
		//}
		idSlice := p.ToObjectId(ids)
		objIds = append(objIds, idSlice...)

		ofS := filterIn["of"]
		ofSlice := p.ToObjectId(ofS)
		objIds = append(objIds, ofSlice...)
		p.MongoIDList = objIds
		return p
	}
}

func (p *Params) ToObjectId(ids []string) []*primitive.ObjectID {
	objIds := make([]*primitive.ObjectID, 0)
	for _, i := range ids {
		objID, _ := primitive.ObjectIDFromHex(i)
		objIds = append(objIds, &objID)
		return objIds
	}
	return objIds
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
		// 需要翻译的key
		finalKey, ok := p.Translate[key]
		if ok {
			filter := bson.E{Key: finalKey, Value: val}
			log.WithField("filter", filter).Warning("~~~~~~~~~~~~~~")
			filters = append(filters, filter)
		} else {
			filter := bson.E{Key: key, Value: val}
			filters = append(filters, filter)
		}
	}
	return filters
}
