package blog_service_user

import (
	"github.com/jinzhu/gorm"
	uuid2 "github.com/satori/go.uuid"
)

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid2.NewV4()
	return scope.SetColumn("Id", uuid.String())
}
