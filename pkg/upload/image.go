package upload

import (
	"fmt"
	"gin-blog/pkg/file"
	"gin-blog/pkg/logging"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

// 获取图片的完整路径
func GetImageFullUrl(name string) string {
	return setting.AppSetting.ImagePrefixUrl + "/" + GetImagePath() + name
}

// 将图片名称转化为md5格式
func GetImageName(name string) string {
	ext := path.Ext(name) // 返回文件后缀名.jpeg
	// trimsuffix 返回没有提供后缀字符串的s。
	//如果s没有以后缀结尾，则返回s时不会改变。
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)
	fmt.Println(fileName, ext)
	return fileName + ext

}

// 获取图片完整路径
func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

// 获取图片路径
func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

// 检查文件后缀
func CheckImageExt(filename string) bool {
	ext := file.GetExt(filename)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}
	return false
}

// 检查图片大小
func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}
	return size <= setting.AppSetting.ImageMaxSize
}

// 检查图片保存的路径是否可写 是否存在路径 将文件写入的权限
func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}
	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := file.CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}
	return nil
}
