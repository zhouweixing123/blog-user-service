package blog_service_user

import (
	"github.com/jinzhu/gorm"
	uuid2 "github.com/satori/go.uuid"
	"time"
)

//BeforeCreate 在注册用户前添加用户的主键ID
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid2.NewV4()
	_ = scope.SetColumn("CreatedAt", time.Now().Format(time.RFC3339))
	_ = scope.SetColumn("RegisterTime", time.Now().Format(time.RFC3339))
	return scope.SetColumn("Id", uuid.String())
}

//BeforeSave 用户修改自动添加修改时间的字段
func (u *User) BeforeSave(scope *gorm.Scope) error {
	_ = scope.SetColumn("UpdateAt", time.Now().Format(time.RFC3339))
	return nil
}
