package service

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/zhouweixing123/blog-user-service/model"
	repo2 "github.com/zhouweixing123/blog-user-service/repo"
	"log"
	"time"
)

// 加密的盐
var key = []byte("zwx114003...")

type UserClaims struct {
	User *model.User
	jwt.StandardClaims
}

type Authable interface {
	Decode(token string) (*UserClaims, error)
	Encode(user *model.User) (string, error)
}

type TokenService struct {
	Repo repo2.Repository
}

//Decode 解析token
func (srv *TokenService) Decode(token string) (*UserClaims, error) {
	tokens, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	claims, ok := tokens.Claims.(*UserClaims)
	if ok && tokens.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (srv *TokenService) Encode(user *model.User) (string, error) {
	// 过期时间
	expireToken := time.Now().Add(time.Hour * 72).Unix()
	log.Printf("获取到的用户信息，用于加密%v", user)
	claims := UserClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "blog.user.service",
		},
	}
	log.Println(claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}
