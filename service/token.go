package service

import (
	"github.com/dgrijalva/jwt-go"
	blog_user_service "github.com/zhouweixing123/blog-user-service/proto/user"
	repo2 "github.com/zhouweixing123/blog-user-service/repo"
	"time"
)

// 加密的盐
var key = []byte("zwx114003...")

type UserClaims struct {
	User *blog_user_service.User
	jwt.StandardClaims
}

type Authable interface {
	Decode(token string) (*UserClaims, error)
	Encode(user *blog_user_service.User) (string, error)
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

func (srv *TokenService) Encode(user *blog_user_service.User) (string, error) {
	// 过期时间
	expireToken := time.Now().Add(time.Hour * 72).Unix()
	claims := UserClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "laracom.user.service",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(key)
}
