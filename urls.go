package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// UrlParams 将url的参数统一进行解析
type UrlParams struct {
	Limit     int  `form:"limit" json:"limit" xml:"limit"`
	Offset    int  `form:"offset" json:"offset" xml:"offset"`
	HasFilter bool `form:"has_filter" json:"has_filter" xml:"has_filter"`
	Filter    FilterRequest
	FilterMap map[string]string `form:"filter_map" json:"filter_map" xml:"filter_map"`
}

func ParserParams(c *gin.Context) UrlParams {
	params := UrlParams{}

	// 获取过滤字段
	filter, hasFilter := c.GetQueryArray("filter")

	// 获取范围值
	rangeValue, hasRange := c.GetQueryArray("range")
	logCtx := log.WithField("has_filter", hasFilter).
		WithField("hasRange", hasRange)

	params.Limit = DefaultPerPage
	params.Offset = DefaultOffset

	if len(rangeValue) == 1 && hasRange {
		logCtx.WithField("rangeValue", rangeValue[0]).
			Debug("range value has parser")
		rangeObj := make([]int, 2)

		// 解析范围
		err := json.Unmarshal([]byte(rangeValue[0]), &rangeObj)

		if err != nil {
			logCtx.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return params
		}
		logCtx.WithField("rangeFrom", rangeObj[0]).
			WithField("rangeTo", rangeObj[1]).
			Debug("range has parser")

		params.Offset = rangeObj[0]

		if rangeObj[1] > 0 {
			params.Limit = rangeObj[1]
		}
	}

	// 如果存在过滤，则将其解析到map中
	// /v1/auth/merchant/roles?filter={"roles":"642f7c008ac505a238abb4d2"}&range=[0,24]&sort=["id","DESC"]
	if hasFilter && len(filter) != 0 {

		// 将过滤器中的所有参数都解析出来供
		// 业务查询进行使用
		filterInstance := make(map[string]string, 0)

		err := json.Unmarshal([]byte(filter[0]), &filterInstance)
		if err != nil {
			// 如果是空的会解析失败
			// 暂停继续解析
			// 返回当前解析到的结果
			logCtx.Error(err)
			return params
		}

		params.HasFilter = true
		// 将解析到的map复制到参数对象中
		params.Filter = filterInstance
	}
	return params
}
