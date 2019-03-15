package v1

import (
	"gin-blog/app"
	"gin-blog/models"
	"gin-blog/pkg/e"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 获取单个文章
func GetArticle(c *gin.Context) {
	appG := app.Gin{c}
	id, _ := com.StrTo(c.Param("id")).Int()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID 必须大于0")

	// 提前判断错误并处理错误
	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, code, nil)
		return
	}
	// 代码执行到这里一定正确
	var data interface{}
	if !models.ExistArticleByID(id) {
		code = e.ERROR_NOT_EXIST_ARTICLE
		appG.Response(http.StatusOK, code, data)
		return
	}
	data = models.GetArticle(id)
	code = e.SUCCESS
	appG.Response(http.StatusOK, code, data)
	//code = e.INVALID_PARAMS
	//var data interface{}
	//if !valid.HasErrors() {
	//	if models.ExistArticleByID(id) {
	//		data = models.GetArticle(id)
	//		code = e.SUCCESS
	//	} else {
	//		code = e.ERROR_NOT_EXIST_ARTICLE
	//	}
	//
	//} else {
	//	for _, err := range valid.Errors {
	//		//log.Println(err.Key, err.Message)
	//		logging.Info(err.Key, err.Message)
	//	}
	//}
	//c.JSON(http.StatusOK, gin.H{
	//	"code": code,
	//	"msg":  e.GetMsg(code),
	//	"data": data,
	//})
}

// 获取多个文章
func GetArticles(c *gin.Context) {
	appG := app.Gin{c}
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state, _ = com.StrTo(arg).Int()
		maps["state"] = state
		maps["deleted_on"] = 0
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId, _ = com.StrTo(arg).Int()
		maps["tag_id"] = tagId

		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}
	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, code, data)
		return
	}
	// 正确执行
	code = e.SUCCESS

	data["lists"] = models.GetArticles(util.GetPage(c), setting.AppSetting.PageSize, maps)
	data["total"] = models.GetArticleTotal(maps)

	appG.Response(http.StatusOK, code, data)

	//code := e.INVALID_PARAMS
	//if !valid.HasErrors() {
	//	code = e.SUCCESS
	//
	//	data["lists"] = models.GetArticles(util.GetPage(c), setting.AppSetting.PageSize, maps)
	//	data["total"] = models.GetArticleTotal(maps)
	//} else {
	//	for _, err := range valid.Errors {
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

// 新增文章
func AddArticle(c *gin.Context) {
	appG := app.Gin{c}
	tagId, _ := com.StrTo(c.Query("tag_id")).Int()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	CoverImageUrl := c.Query("cover_image_url")
	state, _ := com.StrTo(c.DefaultQuery("state", "0")).Int()

	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Required(CoverImageUrl, "cover_image_url").Message("图片URL为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	data := make(map[string]interface{})
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, code, data)
		return
	}
	if !models.ExistTagByID(tagId) {
		code = e.ERROR_NOT_EXIST_TAG
		appG.Response(http.StatusOK, code, data)
		return
	}
	code = e.SUCCESS
	data = make(map[string]interface{})
	data["tag_id"] = tagId
	data["title"] = title
	data["desc"] = desc
	data["content"] = content
	data["created_by"] = createdBy
	data["state"] = state
	data["cover_image_url"] = CoverImageUrl
	models.AddArticle(data)

	appG.Response(http.StatusOK, code, data)
	return

	//code := e.INVALID_PARAMS
	//if !valid.HasErrors() {
	//	if models.ExistTagByID(tagId) {
	//		data := make(map[string]interface{})
	//		data["tag_id"] = tagId
	//		data["title"] = title
	//		data["desc"] = desc
	//		data["content"] = content
	//		data["created_by"] = createdBy
	//		data["state"] = state
	//		data["cover_image_url"] = CoverImageUrl
	//		models.AddArticle(data)
	//		code = e.SUCCESS
	//	} else {
	//		code = e.ERROR_NOT_EXIST_TAG
	//	}
	//} else {
	//	for _, err := range valid.Errors {
	//		logging.Info(err.Key, err.Message)
	//	}
	//}
	//fmt.Println(code)
	//c.JSON(http.StatusOK, gin.H{
	//	"code": code,
	//	"msg":  e.GetMsg(code),
	//	"data": make(map[string]interface{}),
	//})
}

//修改文章
func EditArticle(c *gin.Context) {
	appG := app.Gin{c}
	valid := validation.Validation{}

	id, _ := com.StrTo(c.Param("id")).Int()
	tagId, _ := com.StrTo(c.Query("tag_id")).Int()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state, _ = com.StrTo(arg).Int()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	valid.MaxSize(content, 65536, "content").Message("内容最长为65535字符")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")

	code := e.INVALID_PARAMS
	data := make(map[string]interface{})
	// 请求信息有错误
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, code, data)
		return
	}
	// 文章不存在
	if !models.ExistArticleByID(id) {
		code = e.ERROR_NOT_EXIST_ARTICLE
		appG.Response(http.StatusOK, code, data)
		return
	}
	// 标签不存在
	if !models.ExistTagByID(tagId) {
		code = e.ERROR_NOT_EXIST_TAG
		appG.Response(http.StatusOK, code, data)
		return
	}
	// 错误检查完毕
	code = e.SUCCESS
	if tagId > 0 {
		data["tag_id"] = tagId
	}
	if title != "" {
		data["title"] = title
	}
	if desc != "" {
		data["desc"] = desc
	}
	if content != "" {
		data["content"] = content
	}

	data["modified_by"] = modifiedBy
	models.EditArticle(id, data)

	appG.Response(http.StatusOK, code, data)
	return

	//code := e.INVALID_PARAMS
	//if !valid.HasErrors() {
	//	if models.ExistArticleByID(id) {
	//		if models.ExistTagByID(tagId) {
	//			data := make(map[string]interface{})
	//			if tagId > 0 {
	//				data["tag_id"] = tagId
	//			}
	//			if title != "" {
	//				data["title"] = title
	//			}
	//			if desc != "" {
	//				data["desc"] = desc
	//			}
	//			if content != "" {
	//				data["content"] = content
	//			}
	//
	//			data["modified_by"] = modifiedBy
	//			models.EditArticle(id, data)
	//			code = e.SUCCESS
	//		} else {
	//			code = e.ERROR_NOT_EXIST_TAG
	//		}
	//	} else {
	//		code = e.ERROR_NOT_EXIST_ARTICLE
	//	}
	//} else {
	//	for _, err := range valid.Errors {
	//		//log.Println(err.Key, err.Message)
	//		logging.Info(err.Key, err.Message)
	//	}
	//}
	//c.JSON(http.StatusOK, gin.H{
	//	"code": code,
	//	"msg":  e.GetMsg(code),
	//	"data": make(map[string]string),
	//})
}

// 删除文章
func DeleteArticle(c *gin.Context) {
	appG := app.Gin{c}
	id, _ := com.StrTo(c.Param("id")).Int()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	data := make(map[string]interface{})
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, code, data)
		return
	}
	if !models.ExistArticleByID(id) {
		code = e.ERROR_NOT_EXIST_ARTICLE
		appG.Response(http.StatusOK, code, data)
		return
	}
	models.DeleteArticle(id)
	code = e.SUCCESS

	appG.Response(http.StatusOK, code, data)
	return

	//code := e.INVALID_PARAMS
	//if !valid.HasErrors() {
	//	if models.ExistArticleByID(id) {
	//		models.DeleteArticle(id)
	//		code = e.SUCCESS
	//	} else {
	//		code = e.ERROR_NOT_EXIST_ARTICLE
	//	}
	//} else {
	//	for _, err := range valid.Errors {
	//		//log.Println(err.Key, err.Message)
	//		logging.Info(err.Key, err.Message)
	//	}
	//}
	//c.JSON(http.StatusOK, gin.H{
	//	"code": code,
	//	"msg":  e.GetMsg(code),
	//	"data": make(map[string]string),
	//})
}
