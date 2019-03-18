package models

import (
	"github.com/jinzhu/gorm"
)

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

// 是否存在标签名
func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("name = ? AND deleted_on = ?", name, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

func ExistTagByID(id int) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("id = ? AND deleted_on = ?", id, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

// 添加标签
func AddTag(name string, state int, CreatedBy string) error {
	tag := Tag{
		Name:      name,
		State:     state,
		CreatedBy: CreatedBy,
	}
	if err := db.Create(&tag).Error; err != nil {
		return err
	}
	return nil
}

// 编辑标签
func EditTag(id int, data interface{}) error {
	err := db.Model(&Tag{}).Where("id = ? AND deleted_on = ?", id, 0).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

// 删除标签
func DeleteTag(id int) error {
	err := db.Where("id = ?", id).Delete(&Tag{}).Error
	if err != nil {
		return err
	}
	return nil
}

// 获取标签总数
func GetTagTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&Tag{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetTags(pageNum int, pageSize int, maps interface{}) ([]Tag, error) {
	var (
		tags []Tag
		err  error
	)
	if pageSize > 0 && pageNum > 0 {
		err = db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	} else {
		err = db.Where(maps).Find(&tags).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return tags, nil

	//db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)
}

func CleanAllTag() error {
	if err := db.Unscoped().Where("deleted_on != ?", 0).Delete(&Tag{}).Error; err != nil {
		return err
	}
	return nil
}
