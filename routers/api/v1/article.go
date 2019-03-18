package v1

import (
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
	"gin-blog/pkg/logging"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
	"gin-blog/service/article_service"
	"gin-blog/service/tag_service"
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

	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, article)

}

// 获取多个文章
func GetArticles(c *gin.Context) {
	appG := app.Gin{c}
	valid := validation.Validation{}

	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state")
	}

	tagId := -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "tag_id")
	}
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{
		TagID:    tagId,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	total, err := articleService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	articles, err := articleService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["list"] = articles
	data["total"] = total

	appG.Response(http.StatusOK, e.SUCCESS, data)

}

type AddArticleForm struct {
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	CreatedBy     string `form:"created_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

// 新增文章
// @Produce  json
// @Param tag_id body int true "TagID"
// @Param title body string true "Title"
// @Param desc body string true "Desc"
// @Param content body string true "Content"
// @Param created_by body string true "CreatedBy"
// @Param state body int true "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	var (
		appG = app.Gin{c}
		form AddArticleForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)

	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err := tagService.ExistByID()
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := article_service.Article{
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
	}
	err = articleService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

type EditArticleForm struct {
	ID            int    `form:"id" valid:"Required;Min(1)"`
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	ModifiedBy    string `form:"modified_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

//修改文章
// @Produce  json
// @Param id path int true "ID"
// @Param tag_id body string false "TagID"
// @Param title body string false "Title"
// @Param desc body string false "Desc"
// @Param content body string false "Content"
// @Param modified_by body string true "ModifiedBy"
// @Param state body int false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	var (
		appG = app.Gin{c}
		form = EditArticleForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	articleService := article_service.Article{
		ID:            form.ID,
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		ModifiedBy:    form.ModifiedBy,
		State:         form.State,
	}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err = tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = articleService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

	//appG := app.Gin{c}
	//valid := validation.Validation{}
	//
	//id, _ := com.StrTo(c.Param("id")).Int()
	//tagId, _ := com.StrTo(c.Query("tag_id")).Int()
	//title := c.Query("title")
	//desc := c.Query("desc")
	//content := c.Query("content")
	//modifiedBy := c.Query("modified_by")
	//
	//var state int = -1
	//if arg := c.Query("state"); arg != "" {
	//	state, _ = com.StrTo(arg).Int()
	//	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	//}
	//valid.Min(id, 1, "id").Message("ID必须大于0")
	//valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	//valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	//valid.MaxSize(content, 65536, "content").Message("内容最长为65535字符")
	//valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	//valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	//
	//code := e.INVALID_PARAMS
	//data := make(map[string]interface{})
	//// 请求信息有错误
	//if valid.HasErrors() {
	//	app.MarkErrors(valid.Errors)
	//	appG.Response(http.StatusOK, code, data)
	//	return
	//}
	//// 文章不存在
	//if !models.ExistArticleByID(id) {
	//	code = e.ERROR_NOT_EXIST_ARTICLE
	//	appG.Response(http.StatusOK, code, data)
	//	return
	//}
	//// 标签不存在
	//if !models.ExistTagByID(tagId) {
	//	code = e.ERROR_NOT_EXIST_TAG
	//	appG.Response(http.StatusOK, code, data)
	//	return
	//}
	//// 错误检查完毕
	//code = e.SUCCESS
	//if tagId > 0 {
	//	data["tag_id"] = tagId
	//}
	//if title != "" {
	//	data["title"] = title
	//}
	//if desc != "" {
	//	data["desc"] = desc
	//}
	//if content != "" {
	//	data["content"] = content
	//}
	//
	//data["modified_by"] = modifiedBy
	//models.EditArticle(id, data)
	//
	//appG.Response(http.StatusOK, code, data)
	//return

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
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles/{id} [delete]
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

	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	err = articleService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

	//if !models.ExistArticleByID(id) {
	//	code = e.ERROR_NOT_EXIST_ARTICLE
	//	appG.Response(http.StatusOK, code, data)
	//	return
	//}
	//models.DeleteArticle(id)
	//code = e.SUCCESS
	//
	//appG.Response(http.StatusOK, code, data)
	//return

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
