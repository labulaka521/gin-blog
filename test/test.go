package test

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"reflect"
)

type User struct {
	Name   string "user name" //这引号里面的就是tag
	Passwd string "user passsword"
}

//func main() {
//	user := &User{"chronos", "pass"}
//	s := reflect.TypeOf(user).Elem() //通过反射获取type定义
//	for i := 0; i < s.NumField(); i++ {
//		fmt.Println(s.Field(i).Tag) //将tag输出出来
//	}
//}

type T struct {
	f1     string "f one"
	f2     string
	f3     string "f three"
	f4, f5 int    "f four and five"
	f6     bool   `json: "test1, test2, empty"`
}

type T1 struct {
	f  string `one:"1" two:"2" blank:""`
	f2 string `one:"1" two:"2" blank:""`
}

func main_1() {
	t := reflect.TypeOf(T{})
	f1, _ := t.FieldByName("f1")
	fmt.Println(f1.Name, f1.Tag)

	f2, _ := t.FieldByName("f2")
	fmt.Println(f2.Name, f2.Tag)

	f3, _ := t.FieldByName("f3")
	fmt.Println(f2.Name, f3.Tag)

	f4, _ := t.FieldByName("f4")
	fmt.Println(f2.Name, f4.Tag)

	f6, _ := t.FieldByName("f6")
	fmt.Println(f6.Name, f6.Tag)

	t1 := reflect.TypeOf(T1{})
	f, _ := t1.FieldByName("f")
	fmt.Println(f.Tag)

	v, ok := f.Tag.Lookup("blank")
	fmt.Println(v, ok)
	fmt.Printf("%s %t\n", v, ok)

	t2 := reflect.TypeOf(T1{})
	f22, _ := t2.FieldByName("f2")
	fmt.Println(f22.Tag)

	v1, ok := f22.Tag.Lookup("one")
	fmt.Println(v1, ok)
}


type Product struct {
	Code string
	Price int
}

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("con not connect database")
	}
	defer db.Close()
	db.AutoMigrate(&Product{})

	// 创建
	db.Create(&Product{Code: "test", Price: 123})

	// 读取
	var product Product


	db.First(&product, 1)
	fmt.Println(&product)
	db.First(&product, "code = ?", "test")

	// 更新
	db.Model(&product).Update("Price", 2000)

	db.First(&product, 1)
	db.First(&product, "code = ?", "test")


}
