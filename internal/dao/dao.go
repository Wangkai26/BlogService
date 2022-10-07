package dao

import "github.com/jinzhu/gorm"

// 第八节新增的文件
type Dao struct {
	engine *gorm.DB
}

func New(engine *gorm.DB) *Dao {
	return &Dao{engine:engine}
}