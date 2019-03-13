package main

import (
	"fmt"
	"gin-blog/pkg/setting"
	"gin-blog/routers"
	"log"
	"net/http"
)

//func main() {
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
//}

func main() {

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HttpPort),
		Handler:        routers.InitRouter(),
		ReadTimeout:    setting.ReadTimeOut,
		WriteTimeout:   setting.WriteTimeOut,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Printf("Listen: %s\n", err)
	}
}
