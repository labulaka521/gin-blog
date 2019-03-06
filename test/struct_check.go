package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const tagName = "validate"

// eamil匹配规则
var mailRe = regexp.MustCompile(`\A[\w+\-.]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+\z`)

type Validator interface {
	Validate(interface{}) (bool, error)
}

// 默认校验规则
type DefaultValidator struct {
}


func (v DefaultValidator) Validate(val interface{}) (bool, error) {
	return true, nil
}

type NumberValidator struct {
	Min int
	Max int
}

// 验证数字
func (v NumberValidator) Validate(val interface{}) (bool, error) {
	num := val.(int)
	if num < v.Min {
		return false, fmt.Errorf("should be greater than %v", v.Min)
	}
	if v.Max > v.Min && num > v.Max {
		return false, fmt.Errorf("should be less than %v", v.Max)
	}
	return true, nil
}

type StringValidator struct {
	Min int
	Max int
}

// 字符串验证
func (v StringValidator) Validate(val interface{}) (bool, error) {
	l := len(val.(string))

	// 字符串长度为0
	if l == 0 {
		return false, fmt.Errorf("connot be blank")
	}

	// 字符串长度小雨定义的最小长度
	if l < v.Min {
		return false, fmt.Errorf("should be at least %v chars long", v.Min)
	}

	if v.Max >= v.Min && l > v.Max {
		return false, fmt.Errorf("should be less than %v chars long", v.Max)
	}
	return true, nil
}

type EmailValidator struct {

}

// 验证邮箱
func (v EmailValidator) Validate(val interface{}) (bool, error) {
	if !mailRe.MatchString(val.(string)) {
		return false, fmt.Errorf("is not a valid email address")
	}
	return true, nil
}

// 返回对应的结构体
func getValidatorFromTag(tag string) Validator {
	args := strings.Split(tag, ",")

	switch args[0] {
	case "number":
		validator := NumberValidator{}
		// 将对象的值赋值到validator中将连续的空格分隔值存储到由格式确定的连续参数中。
		fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)

		return validator
	case "string":
		validator := StringValidator{}
		fmt.Sscanf(strings.Join(args[1:],","), "min=%d,max=%d", &validator.Min, &validator.Max)
		return validator
	case "email":
		return EmailValidator{}
	}
	return DefaultValidator{}
}

func validateStruct(s interface{}) []error {
	errs := []error{}

	v := reflect.ValueOf(s)

	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get(tagName)
		//fmt.Println(tag)
		if tag == "" || tag == "-" {
			continue
		}

		validator := getValidatorFromTag(tag)

		valid, err := validator.Validate(v.Field(i).Interface())

		if !valid && err != nil {
			errs = append(errs, fmt.Errorf("%s %s", v.Type().Field(i).Name, err.Error()))
		}

	}
	return errs
}

type User struct {
	Id int `validate:"number,min=1,max=1000"`
	Name string `validate:"string,min=2,max=10"`
	Bio string `validate:"string"`
	Email string `validate:"email"`
}

func main() {
	user := User{
		Id: 0,
		Name: "superlongsstring",
		Bio: "",
		Email: "foobar@53.com",
	}
	for i, err := range validateStruct(user) {
		fmt.Printf("\t%d %s\n", i+1, err.Error())
	}
}


