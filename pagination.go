package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 参考文档
// https://universalglue.dev/posts/gin-pagination/

const (
	// 默认每页展示
	DefaultPerPage = 15
	// 默认展示最大记录每页
	DefaultMaxPerPage = 99
	// 默认起始索引
	DefaultOffset = 0
)

const (
	// reactjs-admin 分页范围模版
	// fmt.Sprintf("items %d-%d/%d", params.Offset, params.Limit, parser.OnePage)
	RectJsAdminPageTpl = "items %d-%d/%d"
)

// GetContentRange 获取分页范围信息
// offset 与数据库的offset意义相同
// perPage 与数据库limit 意义相同
// totalPage 与数据库中counter意义相同
func GetContentRange(tpl string, offset uint, perPage uint, totalPage uint) string {
	return fmt.Sprintf(tpl, offset, perPage, totalPage)
}

// RenderList 列表展示
func RenderList(c *gin.Context, contentRange string, counter int64, obj any) {
	// 写入头部信息，用于reactjs-admin 进行识别从而完成分页
	c.Header("Content-Range", contentRange)
	c.Header("X-Total-Count", fmt.Sprintf("%d", counter))
	// 返回数据部分
	c.JSON(http.StatusOK, obj)
}