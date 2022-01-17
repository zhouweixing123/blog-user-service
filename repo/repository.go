package repo

import (
	"github.com/jinzhu/gorm"
	blog_user_service "user-service/proto/user"
)

type Repository interface {
	Create(user *blog_user_service.User) error
	Get(id string) (*blog_user_service.User, error)
	GetAll() ([]*blog_user_service.User, error)
	GetByEmail(email string) (*blog_user_service.User, error)
}
type UserRepository struct {
	DB *gorm.DB
}

func (repo *UserRepository) Create(user *blog_user_service.User) error {
	err := repo.DB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetByEmail(email string) (*blog_user_service.User, error) {
	user := &blog_user_service.User{}
	err := repo.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) Get(id string) (*blog_user_service.User, error) {
	user := &blog_user_service.User{}
	user.Id = id
	err := repo.DB.First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) GetAll() ([]*blog_user_service.User, error) {
	var users []*blog_user_service.User
	err := repo.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
