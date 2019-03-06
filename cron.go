package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	ticker := time.NewTicker(time.Second * 1)
	wg.Add(1)
	go func() {
		for _ = range ticker.C {
			wg.Add(1)
			fmt.Printf("ticked at %v\n", time.Now())
			wg.Done()
		}
	}()
	wg.Wait()
	//log.Println("Starting...")
	//c := cron.New()
	//c.AddFunc("* * * * * *", func(){
	//	log.Println("Run models.CleanAllArticle...")
	//	models.CleanAllArticle()
	//})
	//c.Start()
	////
	//t1 := time.NewTimer(time.Second * 10)
	//for {
	//	select {
	//	case <-t1.C:
	//		t1.Reset(time.Second * 10)
	//	}
	//}
}
