package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RenderUpdate 更新成功返回
func RenderUpdate(c *gin.Context, id string, message string) {
	// 写入头部信息，用于reactjs-admin 进行识别从而完成分页
	// 返回数据部分
	rsp := SimpleResponse{
		Id:      id,
		Message: message,
	}
	c.JSON(http.StatusOK, rsp)
}
