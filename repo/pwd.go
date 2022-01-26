package repo

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	blog_service_user "github.com/zhouweixing123/blog-user-service/proto/user"
	"log"
)

type PwdInterface interface {
	Create(reset *blog_service_user.PwdReset) error
	GetByToken(token string) (*blog_service_user.PwdReset, error)
	Delete(reset *blog_service_user.PwdReset) error
}
type PwdResetRepo struct {
	Redis *redis.Pool
}

//Create 创建重置密码的token
func (rep *PwdResetRepo) Create(reset *blog_service_user.PwdReset) error {
	redis := rep.Redis.Get()
	_, err := redis.Do("SET", reset.Token, reset.Email)
	if err != nil {
		log.Printf("token添加失败%v\n", err)
		return err
	}
	_, err = redis.Do("EXPIRE", reset.Token, 60)
	log.Printf("设置token的过期时间%v\n", err)
	if err != nil {
		return err
	}
	return nil
}

//GetByToken 通过token获取邮箱
func (repo *PwdResetRepo) GetByToken(token string) (*blog_service_user.PwdReset, error) {
	reset := &blog_service_user.PwdReset{}
	conn := repo.Redis.Get()
	token, err := redis.String(conn.Do("GET", token))
	if err != nil {
		log.Printf("查询失败%v\n", err)
		return nil, err
	}
	if token != "" {
		reset.Token = token
	} else {
		return nil, errors.New("查询失败")
	}
	return reset, err
}

//Delete 通过token删除已修改完成密码的密码重置token
func (repo *PwdResetRepo) Delete(reset *blog_service_user.PwdReset) error {
	conn := repo.Redis.Get()
	_, err := conn.Do("DEL", reset.Token)
	if err != nil {
		log.Printf("删除失败%v\n", err)
		return err
	}
	return err
}
