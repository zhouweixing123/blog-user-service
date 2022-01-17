package handler

import (
	"errors"
	blog_user_service "github.com/zhouweixing123/blog-user-service/proto/user"
	"github.com/zhouweixing123/blog-user-service/repo"
	"github.com/zhouweixing123/blog-user-service/service"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"log"
)

type UserService struct {
	Repo  repo.Repository
	Token service.Authable
}

func (srv *UserService) Get(ctx context.Context, req *blog_user_service.User, rsp *blog_user_service.Response) error {
	user, err := srv.Repo.Get(req.Id)
	if err != nil {
		return err
	}
	rsp.User = user
	return nil
}

func (srv *UserService) GetAll(ctx context.Context, req *blog_user_service.Request, rsp *blog_user_service.Response) error {
	users, err := srv.Repo.GetAll()
	if err != nil {
		return err
	}
	rsp.Users = users
	return nil
}

func (srv *UserService) Create(ctx context.Context, req *blog_user_service.User, rsp *blog_user_service.Response) error {
	pwd, err := bcrypt.GenerateFromPassword([]byte(req.Pwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Pwd = string(pwd)
	err = srv.Repo.Create(req)
	if err != nil {
		return err
	}
	rsp.User = req
	return nil
}

func (srv *UserService) Auth(ctx context.Context, req *blog_user_service.User, rsp *blog_user_service.Token) error {
	log.Println("登录的信息:", req.Phone, req.Pwd)
	user, err := srv.Repo.GetByEmail(req.Phone)
	log.Println("登录的用户信息", user)
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
