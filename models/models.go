package models

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"gin-blog/pkg/setting"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

//func init() {
//	var (
//		err                                               error
//		dbType, dbName, user, password, host, tablePrefix string
//	)
//
//	sec, err := setting.Cfg.GetSection("database")
//	if err != nil {
//		log.Fatal(2, "Fail to get section 'databse': %v", err)
//	}
//
//	dbType = sec.Key("TYPE").String()
//	dbName = sec.Key("NAME").String()
//	user = sec.Key("USER").String()
//	password = sec.Key("PASSWORD").String()
//	host = sec.Key("HOST").String()
//	tablePrefix = sec.Key("TABLE_PREFIX").String()
//
//	db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
//		user,
//		password,
//		host,
//		dbName,
//	))
//
//	if err != nil {
//		log.Println(err)
//	}
//
//	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
//	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
//	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
//
//	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
//		return tablePrefix + defaultTableName
//	}
//	// SingularTable use singular table by default
//	db.SingularTable(true)
//	db.DB().SetMaxIdleConns(10)
//	db.DB().SetMaxOpenConns(100)
//}

func Setup() {
	var err error
	// 连接数据库
	db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name,
	))
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}
	// SingularTable use singular table by default
	db.SingularTable(true)
	// gorm 回调函数替换
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func CloseDB() {
	defer db.Close()
}

// 创建时调用回调函数更新创建时间和修改时候
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				_ = createTimeField.Set(nowTime)
			}
		}
		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				_ = modifyTimeField.Set(nowTime)
			}
		}
	}
}

// 更新操作时调用回调函数修改 字段
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); ok {
		_ = scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

//
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok { // 检查是否手动指定了delete_option
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn") //  获取指定的删除字段，若存在则UPDATE软删除，若不存在则DELETE硬删除

		if !scope.Search.Unscoped && hasDeletedOnField { //
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(), // 返回引用的表名
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()), // 返回组合好的sql
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
