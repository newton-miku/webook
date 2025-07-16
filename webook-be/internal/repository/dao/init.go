package dao

import "gorm.io/gorm"

// 初始化表结构
func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
