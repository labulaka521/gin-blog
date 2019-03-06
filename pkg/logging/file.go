package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	LogSavePath = "runtime/logs/"
	LogSavename = "log"
	LogFileExt = "log"
	TimeFormat = "20060102"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s", LogSavePath)
}

func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", LogSavename, time.Now().Format(TimeFormat), LogFileExt)
	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		mkDir(getLogFilePath())
	case os.IsPermission(err):
		log.Fatalf("Permission :%v", err)
	}

	handler, err := os.OpenFile(filePath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile :%v", err)
	}
	return handler
}

func mkDir(filePath string) {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir + "/" + filePath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}