package routers

import (
	"gin-blog/pkg/upload"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"

	_ "gin-blog/docs"
	"gin-blog/pkg/setting"
	"gin-blog/routers/api"
	"gin-blog/routers/api/v1"
)

func InitRouter() *gin.Engine {
	r := gin.New() // 返回一个新的引擎

	r.Use(gin.Logger(), gin.Recovery()) // 绑定日志和恢复中间件

	gin.SetMode(setting.ServerSetting.RunMode)

	//r.GET("/test", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "test",
	//	})
	//})

	// http.StringPrefix主要是从请求的URL的路径中删除给定的前缀，最终返回一个Handler
	// 通常 http.FileServer 要与 http.StripPrefix 相结合使用，否则当你运行：
	//
	//http.Handle("/upload/images", http.FileServer(http.Dir("upload/images")))
	//会无法正确的访问到文件目录，因为 /upload/images 也包含在了 URL 路径中，必须使用：
	//http.Handle("/upload/images", http.StripPrefix("upload/images", http.FileServer(http.Dir("upload/images"))))
	r.StaticFS("upload/images", http.Dir(upload.GetImageFullPath()))
	r.GET("/auth", api.GetAuth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload", api.UploadImage)

	apiv1 := r.Group("/api/v1/")
	//apiv1.Use(middleware.JWT())
	{
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		///新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		// 获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		// 获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		// 新建文章
		apiv1.POST("/articles", v1.AddArticle)
		// 更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		// 删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}
	return r
}
