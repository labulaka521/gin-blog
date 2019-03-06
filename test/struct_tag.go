package main

import (
	"fmt"
	"reflect"
)

const tagName = "validate"

type User struct {
	Id int `validate:"-"`
	Name string `validate:"presence,min=2,max=23"`
	Email string `validate:"email,required"`
}
type User1 struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

func main() {
	user := User1{
		Id: 1 ,
		Name: "wang",
		Email: "wang@li.tao",
	}
	fmt.Println(user)
	t := reflect.TypeOf(user)

	fmt.Println(t.Kind(), t.Name(), t.NumField())
	for i := 0;i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)
		fmt.Println(field.Name ,field.Type.Name(), tag)
	}
}
