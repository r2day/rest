package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type SortTypeEnum int

const (
	DESC SortTypeEnum = iota
	AES
)

// ReqRange 分页范围
type ReqRange struct {
	Limit  int `form:"limit" json:"limit" xml:"limit"`
	Offset int `form:"offset" json:"offset" xml:"offset"`
}

// 排序方式
type ReqSort struct {
	// 需要排序的字段
	Key string `form:"limit" json:"limit" xml:"limit"`
	// 排序方式 0/1
	SortType SortTypeEnum `form:"offset" json:"offset" xml:"offset"`
}

// UrlParams 将url的参数统一进行解析
type UrlParams struct {
	Range ReqRange `form:"range" json:"range" xml:"range"`
	Sort  ReqSort  `form:"sort" json:"sort" xml:"sort"`

	HasFilter bool `form:"has_filter" json:"has_filter" xml:"has_filter"`
	// Filter    FilterRequest
	FilterMap map[string][]string `form:"filter_map" json:"filter_map" xml:"filter_map"`
	// Filter    FilterRequest
	FilterCommon CommonFilter `form:"filter_common" json:"filter_common" xml:"filter_common"`
}

// ParserParams 解析url请求参数
func ParserParams(c *gin.Context) UrlParams {
	params := UrlParams{}

	// 获取过滤字段
	filter, hasFilter := c.GetQueryArray("filter")

	// 获取范围值
	rangeValue, hasRange := c.GetQuery("range")

	// 获取排序
	sort, hasSort := c.GetQuery("sort")
	logCtx := log.WithField("has_filter", hasFilter).
		WithField("hasRange", hasRange).WithField("range", rangeValue).
		WithField("hasSort", hasSort).WithField("sort", sort)

	logCtx.Info("==========================")
	// 赋新的值
	reqRange := ReqRange{
		Offset: DefaultPerPage,
		Limit:  DefaultOffset,
	}

	// 初始化
	params.Range = reqRange

	if hasRange {
		logCtx.Debug("ready to parser range value")
		rangeObj := make([]int, 2)

		// 解析范围
		err := json.Unmarshal([]byte(rangeValue), &rangeObj)

		if err != nil {
			logCtx.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return params
		}
		logCtx.WithField("rangeFrom", rangeObj[0]).
			WithField("rangeTo", rangeObj[1]).
			Debug("range has parser")

		// 防止用户恶意输入
		limitValue := DefaultPerPage
		if rangeObj[1] > 0 {
			limitValue = rangeObj[1]
		}

		// 赋新的值
		reqRange := ReqRange{
			Offset: rangeObj[0],
			Limit:  limitValue,
		}
		params.Range = reqRange
	}

	if hasSort {
		logCtx.Debug("ready to parser range value")
		sortObj := make([]string, 2)

		// 解析范围
		err := json.Unmarshal([]byte(sort), &sortObj)

		if err != nil {
			logCtx.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return params
		}
		logCtx.WithField("sortKey", sortObj[0]).
			WithField("sortType", sortObj[1]).
			Debug("range has parser")

		// 防止用户恶意输入
		sortValue := AES
		if sortObj[1] == "DESC" {
			sortValue = DESC
		}

		// 赋新的值
		reqSort := ReqSort{
			Key:      sortObj[0],
			SortType: sortValue,
		}
		params.Sort = reqSort
	}
	// 如果存在过滤，则将其解析到map中
	// /v1/auth/merchant/roles?filter={"roles":"642f7c008ac505a238abb4d2"}&range=[0,24]&sort=["id","DESC"]
	if hasFilter {

		// 将过滤器中的所有参数都解析出来供
		// 业务查询进行使用
		// 第一次尝试解析如下格式
		// ?filter={"roles":["643319a80e352fc415f598e1","64331bcba09fc7395567ba6c"]}
		filterInstance := make(map[string][]string, 0)
		parseByMap2SliceIsFailed := false

		err := json.Unmarshal([]byte(filter[0]), &filterInstance)
		if err != nil {
			// 如果是空的会解析失败
			// 暂停继续解析
			// 返回当前解析到的结果
			logCtx.Error(err)
			// return params
			parseByMap2SliceIsFailed = true
		}

		if parseByMap2SliceIsFailed {
			filterInstance2 := CommonFilter{}
			err := json.Unmarshal([]byte(filter[0]), &filterInstance2)
			if err != nil {
				// 如果是空的会解析失败
				// 暂停继续解析
				// 返回当前解析到的结果
				logCtx.Error(err)
				// 将解析到的map复制到参数对象中
				params.FilterCommon = filterInstance2
			}
		} else {
			// 将解析到的map复制到参数对象中
			params.FilterMap = filterInstance
		}

		params.HasFilter = true

	}

	logCtx.WithField("filterMap", params.FilterMap).Info("+++++++++++")
	return params
}
