package main

import (
	"fmt"
	"gin-blog/models"
	"gin-blog/pkg/logging"
	"gin-blog/pkg/setting"
	"gin-blog/routers"
	"log"
	"net/http"
)

func init() {
	setting.Setup()
	models.Setup()
	logging.Setup()

}

func main() {

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        routers.InitRouter(),
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	logging.Info("Sevrer start ")
	if err := s.ListenAndServe(); err != nil {
		log.Printf("Listen: %s\n", err)
	}

	//	endless.DefaultReadTimeOut = setting.ReadTimeOut
	//	endless.DefaultHammerTime = setting.WriteTimeOut
	//	endless.DefaultMaxHeaderBytes = 1 << 20
	//	endPoint := fmt.Sprintf(":%d", setting.HttpPort)
	//	server := endless.NewServer(endPoint, routers.InitRouter())
	//	server.BeforeBegin = func(add string) {
	//		log.Printf("Actual pid is %d", syscall.Getpid())
	//	}
	//	if err := server.ListenAndServe(); err != nil {
	//		log.Printf("Server err: %v", err)
	//	}
}
