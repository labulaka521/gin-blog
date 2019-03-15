package app

import (
	"gin-blog/pkg/e"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

// 统一数据响应入口
func (g *Gin) Response(httpCode int, errCode int, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"code": httpCode,
		"msg":  e.GetMsg(httpCode),
		"data": data,
	})
	return
}
