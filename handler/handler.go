package handler

import (
	blog_user_service "github.com/zhouweixing123/blog-user-service/proto/user"
	"github.com/zhouweixing123/blog-user-service/repo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type UserService struct {
	Repo repo.Repository
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
