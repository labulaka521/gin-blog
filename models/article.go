package models

import (
	"gin-blog/pkg/logging"
	"github.com/jinzhu/gorm"
)

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"` // 声明这个字段为索引
	Tag   Tag `json:"tag"`                 // 嵌套的struct，利用TagID和Tag模型互相关联，在执行查询的时候能够达到Article、Tag关联查询的功能

	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	State         int    `json:"state"`
}

//func (article *Article) BeforeCreate(scope *gorm.Scope) error {
//	scope.SetColumn("CreatedOn", time.Now().Unix())
//	return nil
//}
//
//func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
//	scope.SetColumn("ModifiedOn", time.Now().Unix())
//	return nil
//}

func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ? AND deleted_on = ?", id, 0).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if article.ID > 0 {
		return true, nil
	}
	return false, nil
}

func GetArticleTotal(maps interface{}) (int, error) {
	var count int
	err := db.Model(&Article{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetArticles(pageNum int, pageSize int, maps interface{}) ([]*Article, error) {
	var (
		articles []*Article
		err      error
	)
	if pageSize > 0 && pageNum > 0 {
		//err = db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
		err = db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	} else {
		err = db.Preload("Tag").Where(maps).Find(&articles).Error
		//err = db.Where(maps).Find(&tags).Error
	}
	//err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Warn(err)
		return nil, err
	}
	return articles, nil
}

//Article有一个结构体成员是TagID，就是外键。gorm会通过类名+ID的方式去找到这两个类之间的关联关系
//Article有一个结构体成员是Tag，就是我们嵌套在Article里的Tag结构体，我们可以通过Related进行关联查询
func GetArticle(id int) (*Article, error) {
	var article Article
	err := db.Where("id = ? AND deleted_on = ? ", id, 0).First(&article).Related(&article.Tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &article, nil
}

func EditArticle(id int, data interface{}) error {
	err := db.Model(&Article{}).Where("id = ? AND deleted_on = ?", id, 0).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

func AddArticle(data map[string]interface{}) error {
	article := Article{
		TagID:         data["tag_id"].(int),
		Title:         data["title"].(string),
		Desc:          data["desc"].(string),
		Content:       data["content"].(string),
		CreatedBy:     data["created_by"].(string),
		State:         data["state"].(int),
		CoverImageUrl: data["cover_image_url"].(string),
	}
	err := db.Create(&article).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteArticle(id int) error {
	err := db.Where("id = ?", id).Delete(Article{}).Error
	if err != nil {
		return err
	}
	return nil
}

func CleanAllArticle() error {
	err := db.Unscoped().Where("deleted_on != ?", 0).Delete(&Article{}).Error
	if err != nil {
		return err
	}
	return nil
}
