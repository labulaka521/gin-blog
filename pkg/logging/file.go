package logging

import (
	"fmt"
	"gin-blog/pkg/file"
	"gin-blog/pkg/setting"
	"os"
	"time"
)

//
//var (
//	LogSavePath = "runtime/logs/"
//	LogSavename = "log"
//	LogFileExt = "log"
//	TimeFormat = "20060102"
//)

func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s", setting.AppSetting.LogSavename,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
	)
}

//func getLogFileFullPath() string {
//	prefixPath := getLogFilePath()
//	suffixPath := fmt.Sprintf("%s%s.%s", LogSavename, time.Now().Format(TimeFormat), LogFileExt)
//	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
//}

//func openLogFile(filePath string) *os.File {
//	_, err := os.Stat(filePath)
//	switch {
//	case os.IsNotExist(err):
//		mkDir(getLogFilePath())
//	case os.IsPermission(err):
//		log.Fatalf("Permission :%v", err)
//	}
//
//	handler, err := os.OpenFile(filePath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
//	if err != nil {
//		log.Fatalf("Fail to OpenFile :%v", err)
//	}
//	return handler
//}
//
//func mkDir(filePath string) {
//	dir, _ := os.Getwd()
//	err := os.MkdirAll(dir + "/" + filePath, os.ModePerm)
//	if err != nil {
//		panic(err)
//	}
//}

func openLogFile(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}
	src := dir + "/" + filePath
	perm := file.CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", err)
	}

	err = file.IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := file.Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}
	return f, nil

}
