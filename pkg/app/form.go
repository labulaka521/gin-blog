package app

import (
	"fmt"
	"gin-blog/pkg/e"
	"gin-blog/pkg/logging"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 检验数据 及 绑定表单请求结构
func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)

	if err != nil {
		fmt.Println(err)
		return http.StatusBadRequest, e.INVALID_PARAMS
	}
	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		logging.Warn(err)
		return http.StatusInternalServerError, e.ERROR
	}
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, e.INVALID_PARAMS
	}
	return http.StatusOK, e.SUCCESS
}
