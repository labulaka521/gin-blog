package test

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

type User struct {
	User_id int `gorm:"primary_key"`
	Name string
	Pwd string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	// 创建连接
	db, err :=  gorm.Open("sqlite3", "test1.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&User{})

	user := User{}

	// 添加
	db.Create(&User{Name: "wang", Pwd: "cyusbcj"})

	// 查询ID为1
	db.First(&user, 1)
	//输出查询的结果
	fmt.Println(user.Name, user.Pwd, &user)

	db.Save(&user)
	// 提交

	//db.Where()

}