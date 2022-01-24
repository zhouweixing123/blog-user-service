package handler

import (
	"errors"
	"github.com/jinzhu/gorm"
	blog_user_service "github.com/zhouweixing123/blog-user-service/proto/user"
	"github.com/zhouweixing123/blog-user-service/repo"
	"github.com/zhouweixing123/blog-user-service/service"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type UserService struct {
	Repo  repo.Repository
	Token service.Authable
}

//Register 用户信息注册
func (srv *UserService) Register(ctx context.Context, user *blog_user_service.User, rsp *blog_user_service.Response) (err error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(user.Pwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Pwd = string(pwd)
	err = srv.Repo.Register(user)
	if err != nil {
		return err
	}
	rsp.User = user
	return nil
}

//GetAllUser 获取所有的用户信息
func (srv *UserService) GetAllUser(ctx context.Context, req *blog_user_service.Request, rsp *blog_user_service.Response) error {
	users, err := srv.Repo.GetAllUser()
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}
	rsp.Users = users
	return nil
}

//GetByAccount 通过账号获取用户信息
func (srv *UserService) GetByAccount(ctx context.Context, req *blog_user_service.User, rsp *blog_user_service.Response) error {
	user, err := srv.Repo.GetByAccount(req.Account)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	rsp.User = user
	return nil
}

//GetByEmail 通过邮箱获取用户信息
func (src *UserService) GetByEmail(ctx context.Context, req *blog_user_service.User, rsp *blog_user_service.Response) error {
	user, err := src.Repo.GetByEmail(req.Email)
	if err != nil {
		return err
	}
	rsp.User = user
	return nil
}

//GetByPhone 通过手机号获取用户信息
func (srv *UserService) GetByPhone(ctx context.Context, req *blog_user_service.User, rsp *blog_user_service.Response) error {
	user, err := srv.Repo.GetByPhone(req.Phone)
	if err != nil {
		return err
	}
	rsp.User = user
	return nil
}

//Auth 登录
func (srv *UserService) Auth(ctx context.Context, req *blog_user_service.User, rsp *blog_user_service.Token) error {
	// req.Account可以传递账号、邮箱、手机号
	user, err := srv.Repo.Auth(req.Account)
	if err != nil {
		return err
	}
	// 校验密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(user.Pwd), []byte(req.Pwd))
	if err != nil {
		return err
	}
	// 生成token
	token, err := srv.Token.Encode(user)
	if err != nil {
		return err
	}
	rsp.Token = token
	return nil
}

//ValidateToken 校验用户传递的token
func (srv *UserService) ValidateToken(ctx context.Context, req *blog_user_service.Token, rsp *blog_user_service.Token) error {
	claims, err := srv.Token.Decode(req.Token)
	if err != nil {
		return err
	}
	if claims.User.Id == "" {
		return errors.New("无效的用户信息")
	}
	rsp.Valid = true
	return nil
}
