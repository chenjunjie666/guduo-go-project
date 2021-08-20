// 数据库常用字段定义
package model

import "gorm.io/plugin/soft_delete"

type Id struct {
	ID PrimaryKey `json:"id" gorm:"primaryKey"`
}

type AutoTime struct {
	CreatedAt SecondTimeStamp `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt SecondTimeStamp `json:"updated_at" gorm:"autoUpdateTime"`
}

type DeleteAt struct {
	DeletedAt soft_delete.DeletedAt `json:"delete_at"`
}

// 一些常用的数据库定义组合

// 常用组合，包含 id，created_at，updated_at
type Fields struct {
	Id
	AutoTime
}

// 常用组合，包含 id，created_at，updated_at，deleted_at
type FieldsWithSoftDelete struct {
	Id
	AutoTime
	DeleteAt
}
