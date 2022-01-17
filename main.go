package main

import (
	"fmt"
	"github.com/micro/go-micro"
	"log"
	db2 "user-service/db"
	"user-service/handler"
	blog_user_service "user-service/proto/user"
	repo2 "user-service/repo"
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
	srv := micro.NewService(
		micro.Name("blog.user.service"),
		micro.Version("latest"),
	)
	srv.Init()
	blog_user_service.RegisterUserServiceHandler(srv.Server(), &handler.UserService{
		Repo: repo,
	})
	if err := srv.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
