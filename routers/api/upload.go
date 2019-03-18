package api

import (
	"gin-blog/pkg/app"
	"gin-blog/pkg/e"
	"gin-blog/pkg/logging"
	"gin-blog/pkg/upload"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadImage(c *gin.Context) {
	appG := app.Gin{c}
	code := e.INVALID_PARAMS
	data := make(map[string]string)

	file, image, err := c.Request.FormFile("image")

	// 错误处理
	if err != nil {
		logging.Warn(err)
		code = e.ERROR
		appG.Response(http.StatusOK, code, data)
		return
	}

	// 获取上传的文件错误
	if image == nil {
		appG.Response(http.StatusOK, code, data)
		return
	}

	// 获取图片名并且转换为相应的md5格式
	imageName := upload.GetImageName(image.Filename)
	// 获取图片的保存的全路径
	fullPath := upload.GetImageFullPath()
	// 获取图片的保存路径
	savePath := upload.GetImagePath()

	// 文件保存的路径
	src := fullPath + imageName

	// 检查文件的格式和大小
	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
		code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		appG.Response(http.StatusOK, code, data)
		return
	}

	// 检查文件的保存路径
	err = upload.CheckImage(src)
	if err != nil {
		logging.Warn(err)
		code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
		appG.Response(http.StatusOK, code, data)
		return
	}

	// 保存文件
	if err = appG.C.SaveUploadedFile(image, src); err != nil {
		logging.Warn(err)
		code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
		appG.Response(http.StatusOK, code, data)
		return
	}
	// 保存成功
	code = e.SUCCESS
	data["image_url"] = upload.GetImageFullUrl(imageName)
	data["image_save_url"] = savePath + imageName
	appG.Response(http.StatusOK, code, data)
	return

	//if image == nil {
	//	code = e.INVALID_PARAMS
	//
	//} else {
	//	imageName := upload.GetImageName(image.Filename)
	//	fullPath := upload.GetImageFullPath()
	//	savepath := upload.GetImagePath()
	//
	//	src := fullPath + imageName
	//	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
	//		code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
	//	} else {
	//		err := upload.CheckImage(fullPath)
	//		if err != nil {
	//			logging.Warn(err)
	//			code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
	//		} else if err = c.SaveUploadedFile(image, src); err != nil {
	//			logging.Warn(err)
	//			code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
	//		} else {
	//			data["image_url"] = upload.GetImageFullUrl(imageName)
	//			data["image_save_url"] = savepath + imageName
	//		}
	//	}
	//
	//}
	//c.JSON(http.StatusOK, gin.H{
	//	"code": code,
	//	"msg":  e.GetMsg(code),
	//	"data": data,
	//})
}
