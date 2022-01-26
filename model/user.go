package model

import (
	blog_service_user "github.com/zhouweixing123/blog-user-service/proto/user"
	"strconv"
	"time"
)

//gorm.Model           // 包括ID/CreatedAt/UpdatedAt/DeletedAt

type User struct {
	Id               string    `gorm:"type:varchar(100);primary_key;not null"`
	Account          string    `gorm:"type:varchar(100);unique_index:u_account_email_phone;not null;default:''"`
	Email            string    `gorm:"type:varchar(100);unique_index:u_account_email_phone;not null;default:''"`
	Phone            string    `gorm:"type:char(20);unique_index:u_account_email_phone;not null;default:''"`
	ProfilePhoto     string    `gorm:"type:varchar(255)"`
	Pwd              string    `gorm:"type:varchar(255);not null;default:''"`
	Nickname         string    `gorm:"type:varchar(20);not null;default:''"`
	Ip               uint64    `gorm:"type:int;default:0"`
	RegistrationTime time.Time `gorm:"type:timestamp;null"`
	Birthday         time.Time `gorm:"type:date;not null"`
	CreatedAt        time.Time `gorm:"type:timestamp;null"`
	UpdatedAt        time.Time `gorm:"type:timestamp;null"`
}

func (model *User) ToORM(req *blog_service_user.User) (*User, error) {
	if req.Id != "" {
		model.Id = req.Id
	}
	if req.Account != "" {
		model.Account = req.Account
	}
	if req.Email != "" {
		model.Email = req.Email
	}
	if req.Phone != "" {
		model.Phone = req.Phone
	}
	if req.Pwd != "" {
		model.Pwd = req.Pwd
	}
	if req.Nickname != "" {
		model.Nickname = req.Nickname
	}
	if strconv.Itoa(int(req.Ip)) != "" {
		ip, _ := strconv.ParseUint(strconv.FormatInt(req.Ip, 10), 10, 64)
		model.Ip = uint64(ip)
	}
	if req.RegistrationTime != "" {
		model.RegistrationTime, _ = time.Parse("2006-01-02", req.RegistrationTime)
	}
	if req.Birthday != "" {
		model.Birthday, _ = time.Parse("2006-01-02", req.Birthday)
	}
	if req.CreatedAt != "" {
		model.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", req.CreatedAt)
	}
	if req.UpdatedAt != "" {
		model.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", req.UpdatedAt)
	}
	return model, nil
}

func (model *User) ToProtobuf() (*blog_service_user.User, error) {
	var user = &blog_service_user.User{}
	user.Id = model.Id
	user.Account = model.Account
	user.Email = model.Email
	user.Phone = model.Phone
	user.ProfilePhoto = model.ProfilePhoto
	user.Nickname = model.Nickname
	user.Ip = int64(model.Ip)
	user.RegistrationTime = model.RegistrationTime.Format("2006-01-02")
	user.Birthday = model.Birthday.Format("2006-01-02")
	user.CreatedAt = model.CreatedAt.Format("2006-01-02 15:04:05")
	user.UpdatedAt = model.UpdatedAt.Format("2006-01-02 15:04:05")
	return user, nil
}
