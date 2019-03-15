package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/gomodule/redigo/redis"
)

type Message struct {
	Name string `json:"name"`
	Body string `json:"body"`
	Time int64  `json:"time"`
}

func main() {
	m := Message{"Alice", "Hello", 1294706395881547000}
	b, _ := json.Marshal(&m)
	fmt.Println(string(b)) //{
}
