package main

import (
	"gin-blog/models"
	"time"
)

func main() {

	ticker := time.NewTicker(time.Second * 2)
	for i := range ticker.C {
		_ = i
		models.CleanAllTag()
		models.CleanAllArticle()
	}
}
