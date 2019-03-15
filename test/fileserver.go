package main

import (
	//_ "github.com/gomodule/redigo/redis"
	"net/http"
)

func main() {
	// This works and strip "/static/" fragment from path
	fs := http.FileServer(http.Dir("."))
	http.Handle("/go", http.StripPrefix("/go", fs))
	//http.Handle("/go", fs)

	http.ListenAndServe(":8081", nil)
}
