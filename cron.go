package main

import (
	"gin-blog/models"
	"time"
)

func main() {

	ticker := time.NewTicker(time.Second * 2)
	for _ := range ticker.C {
		models.CleanAllTag()
		models.CleanAllArticle()
	}
}
