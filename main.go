package main

import (
	"context"
	"fmt"
	"gin-blog/pkg/setting"
	"gin-blog/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"syscall"
)

//func main() {
//	//r := gin.Default()
//	//r.GET("/ping", func(c *gin.Context) {
//	//	c.JSON(200, gin.H{
//	//		"message": "pong",
//	//	})
//	//})
//	endless.DefaultReadTimeOut = setting.ReadTimeOut
//	endless.DefaultWriteTimeOut = setting.WriteTimeOut
//	endless.DefaultMaxHeaderBytes = 1 << 20
//	endPoint := fmt.Sprintf(":%d", setting.HttpPort)
//
//	server := endless.NewServer(endPoint, routers.InitRouter())
//	server.BeforeBegin = func(add string) {
//		log.Printf("Actual pid is %d", syscall.Getpid())
//	}
//	err := server.ListenAndServe()
//	if err != nil {
//		log.Printf("Server err: %v\n", err)
//	}
//	//
//	//r := routers.InitRouter()
//	//s := &http.Server{
//	//	Addr:           fmt.Sprintf(":%d", setting.HttpPort),
//	//	Handler:        r,
//	//	ReadTimeout:    setting.ReadTimeOut,
//	//	WriteTimeout:   setting.WriteTimeOut,
//	//	MaxHeaderBytes: 1 << 20,
//	//}
//	//
//	//s.ListenAndServe()
//
//}

func main() {
	router := routers.InitRouter()
	s := &http.Server{
		Addr:	fmt.Sprintf(":%d", setting.HttpPort),
		Handler:	router,
		ReadTimeout: setting.ReadTimeOut,
		WriteTimeout: setting.WriteTimeOut,
		MaxHeaderBytes: 1 << 20,
	}
	//log.Println(syscall.Getpid())
	go func() {
		log.Println(syscall.Getpid())
		if err := s.ListenAndServe(); err != nil {
			log.Printf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<- quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}