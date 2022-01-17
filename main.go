package main

import (
	"fmt"
	"github.com/micro/go-micro"
	db2 "github.com/zhouweixing123/blog-user-service/db"
	"github.com/zhouweixing123/blog-user-service/handler"
	blog_user_service "github.com/zhouweixing123/blog-user-service/proto/user"
	repo2 "github.com/zhouweixing123/blog-user-service/repo"
	"github.com/zhouweixing123/blog-user-service/service"
	"log"
)

func main() {
	db, err := db2.CreateConnection()
	if err != nil {
		log.Fatalf("数据库连接失败, err:%v", err)
	}
	db.AutoMigrate(blog_user_service.User{})
	repo := &repo2.UserRepository{
		DB: db,
	}
	token := &service.TokenService{
		Repo: repo,
	}
	srv := micro.NewService(
		micro.Name("blog.user.service"),
		micro.Version("latest"),
	)
	srv.Init()
	blog_user_service.RegisterUserServiceHandler(srv.Server(), &handler.UserService{
		Repo:  repo,
		Token: token,
	})
	if err := srv.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
