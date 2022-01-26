package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
	"strconv"
	"time"
)

func CreateConnection() (*redis.Pool, error) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	pwd := os.Getenv("REDIS_PWD")
	database, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	return &redis.Pool{
		MaxIdle:     10,  // 最大空闲连接数
		MaxActive:   10,  // 最大可用连接数
		IdleTimeout: 240, // 关闭超时的空闲连接
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			var c, err = redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%s", host, port),
				redis.DialDatabase(database),                           // 设置redis所用的数据库
				redis.DialReadTimeout(time.Duration(1)*time.Second),    // 读取redis数据的超时时间
				redis.DialWriteTimeout(time.Duration(1)*time.Second),   // 写入redis数据的超时时间
				redis.DialConnectTimeout(time.Duration(2)*time.Second), // 连接redis的超时时间
				redis.DialPassword(pwd),                                // 连接redis密码
			)
			if err != nil {
				log.Printf("redis连接失败%v\n", err)
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if _, err := c.Do("PING"); err != nil {
				log.Printf("redis Ping失败%v\n", err)
				return err
			}
			return nil
		},
	}, nil
}
