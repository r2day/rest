package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UrlParams 将url的参数统一进行解析
type UrlParams struct {
	Limit     int  `form:"limit" json:"limit" xml:"limit"`
	Offset    int  `form:"offset" json:"offset" xml:"offset"`
	HasFilter bool `form:"has_filter" json:"has_filter" xml:"has_filter"`
	Filter    FilterRequest
}

func ParserParams(c *gin.Context) UrlParams {
	params := UrlParams{}

	// 获取过滤字段
	filter, isFilterOk := c.GetQueryArray("filter")

	// 获取范围值
	rangeValue, isRangeOk := c.GetQueryArray("range")

	params.Limit = DefaultPerPage
	params.Offset = DefaultOffset

	if len(rangeValue) == 1 && isRangeOk {
		println("rangeValue-->", rangeValue[0], isRangeOk)
		rangeObj := make([]int, 2)

		// 解析范围
		err := json.Unmarshal([]byte(rangeValue[0]), &rangeObj)

		if err != nil {
			fmt.Println("json.Unmarshal failed-->", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return params
		}
		params.Offset = rangeObj[0]

		if rangeObj[1] > 0 {
			params.Limit = rangeObj[1]
		}
	}

	if isFilterOk && len(filter) != 0 {

		// 将过滤器中的所有参数都解析出来供
		// 业务查询进行使用
		filterInstance := FilterRequest{}

		err := json.Unmarshal([]byte(filter[0]), &filterInstance)
		if err != nil {
			// 如果是空的会解析失败
			return params
		}

		// 检查如果所有过滤字段都没有被解析到那么
		// 直接返回
		if filterInstance.Status == "" &&
			filterInstance.CategoryId == "" &&
			filterInstance.ProductId == "" &&
			filterInstance.BrandId == "" {
			return params
		}
		params.HasFilter = true
		params.Filter = filterInstance
	}
	return params
}
