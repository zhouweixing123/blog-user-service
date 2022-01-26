package main

import (
	"fmt"
	"github.com/micro/go-micro"
	db2 "github.com/zhouweixing123/blog-user-service/db"
	"github.com/zhouweixing123/blog-user-service/handler"
	blog_user_service "github.com/zhouweixing123/blog-user-service/proto/user"
	redis2 "github.com/zhouweixing123/blog-user-service/redis"
	repo2 "github.com/zhouweixing123/blog-user-service/repo"
	"github.com/zhouweixing123/blog-user-service/service"
	"log"
)

func main() {
	db, err := db2.CreateConnection()
	if err != nil {
		log.Fatalf("数据库连接失败, err:%v", err)
		return
	}
	defer db.Close()
	redisDB, err := redis2.CreateConnection()
	if err != nil {
		log.Fatalf("redis链接失败, err: %v\n", err)
		return
	}
	defer redisDB.Close()
	db.AutoMigrate(blog_user_service.User{})
	repo := &repo2.UserRepository{
		DB: db,
	}
	resetRepo := &repo2.PwdResetRepo{
		Redis: redisDB,
	}
	token := &service.TokenService{
		Repo: repo,
	}
	srv := micro.NewService(
		micro.Name("blog.service.user"),
		micro.Version("latest"),
	)
	pubsub := srv.Server().Options().Broker
	srv.Init()
	blog_user_service.RegisterUserServiceHandler(srv.Server(), &handler.UserService{
		Repo:      repo,
		Token:     token,
		ResetRepo: resetRepo,
		PubSub:    pubsub,
	})
	if err := srv.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
