package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RenderDelete 删除成功返回
func RenderDelete(c *gin.Context, id string, message string) {
	// 写入头部信息，用于reactjs-admin 进行识别从而完成分页
	// 返回数据部分
	rsp := CreateResponse{
		Id:      id,
		Message: message,
	}
	c.JSON(http.StatusNoContent, rsp)
}
