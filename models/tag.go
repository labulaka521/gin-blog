package models

import "fmt"

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)
	return
}

// 获取标签总数
func GetTagTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)
	return
}

// 是否存在标签名
func ExistTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name = ? AND deleted_on = ?", name, 0).First(&tag)
	if tag.ID > 0 {
		return true
	}
	return false
}

// 添加标签
func AddTag(name string, state int, CreatedBy string) bool {
	db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: CreatedBy,
	})
	return true
}

func ExistTagByID(id int) bool {
	var tag Tag
	db.Select("id").Where("id = ? AND deleted_on = ?", id, 0).First(&tag)
	if tag.ID > 0 {
		return true
	}
	return false
}

func DeleteTag(id int) bool {
	err := db.Where("id = ?", id).Delete(&Tag{}).Error
	fmt.Println(err)
	if err != nil {

		return false
	}
	return true
}

func EditTag(id int, data interface{}) bool {
	db.Model(&Tag{}).Where("id = ? AND deleted_on = ?", id, 0).Updates(data)
	return true
}

func CleanAllTag() bool {
	db.Unscoped().Where("deleted_on != ?", 0).Delete(&Tag{})
	return true
}
