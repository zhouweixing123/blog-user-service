package repo

import (
	"github.com/jinzhu/gorm"
	"github.com/zhouweixing123/blog-user-service/model"
	"log"
)

type Repository interface {
	//Register 用户注册
	Register(user *model.User) error
	//GetAllUser 获取所有的用户信息
	GetAllUser() ([]*model.User, error)
	//GetByAccount 通过账号获取用户信息
	GetByAccount(account string) (*model.User, error)
	//GetByEmail 通过邮箱获取用户信息
	GetByEmail(email string) (*model.User, error)
	//GetByPhone 通过手机号获取用户信息
	GetByPhone(phone string) (*model.User, error)
	//GetById 通过主键ID获取用户信息
	GetById(id string) (*model.User, error)
	//Auth 登录 account可接收邮箱/手机号/账号
	Auth(account string) (*model.User, error)
	//UpdateUser 修改用户信息
	UpdateUser(user *model.User) error
}

type UserRepository struct {
	DB *gorm.DB
}

//Register 用户信息注册
func (repo *UserRepository) Register(user *model.User) error {
	err := repo.DB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

//GetAllUser 获取所有的用户信息
func (repo *UserRepository) GetAllUser() ([]*model.User, error) {
	var users []*model.User
	err := repo.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

//GetByAccount 通过用户账号获取用户信息
func (repo *UserRepository) GetByAccount(account string) (*model.User, error) {
	user := &model.User{}
	err := repo.DB.Where("account = ?", account).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

//GetByEmail 通过邮箱获取用户信息
func (repo *UserRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := repo.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

//GetByPhone 通过手机号获取用户信息
func (repo *UserRepository) GetByPhone(phone string) (*model.User, error) {
	user := &model.User{}
	err := repo.DB.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

//GetById 通过用户的主键ID获取用户信息
func (repo *UserRepository) GetById(id string) (*model.User, error) {
	user := &model.User{}
	err := repo.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

//Auth 登录
func (repo *UserRepository) Auth(account string) (*model.User, error) {
	user := &model.User{}
	err := repo.DB.Where("account = ? or email = ? or phone = ?", account, account, account).First(&user).Error
	log.Println("%v\n", err)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//UpdateUser 修改用户信息
func (repo *UserRepository) UpdateUser(user *model.User) (err error) {
	err = repo.DB.Save(user).Error
	if err != nil {
		return err
	}
	return nil
}
