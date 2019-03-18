package api

import (
	"gin-blog/models"
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
	"gin-blog/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}
	a := auth{Username: username, Password: password}

	appG := app.Gin{c}

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	// 检查字段
	ok, _ := valid.Valid(&a)
	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, code, data)
		return
	}
	// 检查用户名密码是否正确
	if !models.CheckAuth(username, password) {
		code = e.ERROR_AUTH
		appG.Response(http.StatusOK, code, data)
		return
	}
	// 生成token
	token, err := util.GenerateToken(username, password)
	if err != nil {
		code = e.ERROR_AUTH_TOKEN
		appG.Response(http.StatusOK, code, data)
		return
	}
	data["token"] = token
	code = e.SUCCESS
	appG.Response(http.StatusOK, code, data)

	//if ok {
	//	isExist := models.CheckAuth(username, password)
	//	if isExist {
	//		token, err := util.GenerateToken(username, password)
	//		if err != nil {
	//			code = e.ERROR_AUTH_TOKEN
	//		} else {
	//			data["token"] = token
	//			code = e.SUCCESS
	//		}
	//	} else {
	//		code = e.ERROR_AUTH
	//	}
	//
	//} else {
	//	for _, err := range valid.Errors {
	//		//log.Println(err.Key, err.Message)
	//		logging.Info(err.Key, err.Message)
	//	}
	//}
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"code": code,
	//	"msg":  e.GetMsg(code),
	//	"data": data,
	//})
}
