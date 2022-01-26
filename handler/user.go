package handler

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/broker"
	"github.com/zhouweixing123/blog-user-service/model"
	blog_user_service "github.com/zhouweixing123/blog-user-service/proto/user"
	"github.com/zhouweixing123/blog-user-service/repo"
	"github.com/zhouweixing123/blog-user-service/service"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"log"
	"os"
)

type UserService struct {
	Repo      repo.Repository
	Token     service.Authable
	ResetRepo repo.PwdInterface
	PubSub    broker.Broker
}

//Register 用户信息注册
func (srv *UserService) Register(ctx context.Context, user *blog_user_service.User, rsp *blog_user_service.Response) (err error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(user.Pwd), bcrypt.DefaultCost)
	log.Printf("注册时密码加密的错误$:%v\n", err)
	if err != nil {
		return err
	}
	user.Pwd = string(pwd)
	userModel := &model.User{}
	userOrm, _ := userModel.ToORM(user)
	err = srv.Repo.Register(userOrm)
	log.Printf("注册的错误信息:%v\n", err)
	if err != nil {
		return err
	}
	rsp.User, _ = userOrm.ToProtobuf()
	return nil
}

//GetAllUser 获取所有的用户信息
func (srv *UserService) GetAllUser(ctx context.Context, req *blog_user_service.Request, rsp *blog_user_service.Response) error {
	users, err := srv.Repo.GetAllUser()
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}
	userItems := make([]*blog_user_service.User, len(users))
	for index, userItem := range users {
		userItemProtobuf, _ := userItem.ToProtobuf()
		userItems[index] = userItemProtobuf
	}
	rsp.Users = userItems
	return nil
}

//GetByAccount 通过账号获取用户信息
func (srv *UserService) GetByAccount(ctx context.Context, req *blog_user_service.User, rsp *blog_user_service.Response) error {
	user, err := srv.Repo.GetByAccount(req.Account)
	log.Printf("通过账号获取用户信息:%v\n", err)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if user != nil {
		rsp.User, _ = user.ToProtobuf()
	}
	log.Printf("通过账号获取的用户信息:%v\n", err)
	return nil
}

//GetByEmail 通过邮箱获取用户信息
func (src *UserService) GetByEmail(ctx context.Context, req *blog_user_service.User, rsp *blog_user_service.Response) error {
	user, err := src.Repo.GetByEmail(req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if user != nil {
		rsp.User, _ = user.ToProtobuf()
	}
	return nil
}

//GetByPhone 通过手机号获取用户信息
func (srv *UserService) GetByPhone(ctx context.Context, req *blog_user_service.User, rsp *blog_user_service.Response) error {
	user, err := srv.Repo.GetByPhone(req.Phone)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if user != nil {
		rsp.User, _ = user.ToProtobuf()
	}
	return nil
}

//GetById 通过用户的主键ID获取用户信息
func (srv *UserService) GetById(ctx context.Context, req *blog_user_service.User, response *blog_user_service.Response) error {
	user, err := srv.Repo.GetById(req.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.New("该用户不存在")
	}
	if user != nil {
		response.User, _ = user.ToProtobuf()
	}
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

//UpdateUser 修改用户信息
func (srv *UserService) UpdateUser(ctx context.Context, req *blog_user_service.User, rsp *blog_user_service.Response) error {
	if req.Id == "" {
		return errors.New("用户ID不能为空")
	}
	if req.Pwd != "" {
		pwd, err := bcrypt.GenerateFromPassword([]byte(req.Pwd), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("用户密码加密失败%v\n", err)
			return err
		}
		req.Pwd = string(pwd)
	}
	user, _ := srv.Repo.GetById(req.Id)
	err := srv.Repo.UpdateUser(user)
	if err != nil {
		log.Printf("用户信息修改失败%v\n", err)
		return err
	}
	rsp.User = req
	return nil
}

//CreatePwdResetToken 存储重置密码的Token
func (srv *UserService) CreatePwdResetToken(ctx context.Context, req *blog_user_service.PwdReset, rsp *blog_user_service.PwdResetResponse) error {
	if req.Email == "" {
		log.Println("邮箱不能为空")
		return errors.New("邮箱不能为空")
	}
	userInfo, err := srv.Repo.GetByEmail(req.Email)
	if err != nil {
		log.Printf("通过邮箱获取用户信息: %v\n", err)
		return errors.New("用户不存在")
	}
	err = srv.ResetRepo.Create(req)
	if err != nil {
		log.Printf("创建修改密码token失败%v\n", err)
		return err
	}
	if userInfo != nil {
		err = srv.publishEvent(req)
		if err != nil {
			return err
		}
		rsp.PwdReset = req
	}
	return nil
}

//ValidatePwdResetToken 校验重置密码的Token
func (srv *UserService) ValidatePwdResetToken(ctx context.Context, req *blog_user_service.Token, rsp *blog_user_service.Token) error {
	rsp.Valid = false
	if req.Token == "" {
		log.Printf("token不能为空")
		return errors.New("非法请求")
	}
	_, err := srv.ResetRepo.GetByToken(req.Token)
	if err != nil {
		log.Printf("token获取失败%v\n", err)
		return errors.New("链接已过期")
	}
	rsp.Valid = true
	return nil
}

//DelPwdResetToken 删除重置密码的Token
func (srv *UserService) DelPwdResetToken(ctx context.Context, req *blog_user_service.PwdReset, rsp *blog_user_service.PwdResetResponse) error {
	err := srv.ResetRepo.Delete(req)
	if err != nil {
		log.Printf("重置密码的Token删除失败%v\n", err)
		return err
	}
	rsp.PwdReset = nil
	return err
}
func (srv *UserService) publishEvent(reset *blog_user_service.PwdReset) error {
	body, err := json.Marshal(reset)
	if err != nil {
		return err
	}
	msg := &broker.Message{
		Header: map[string]string{
			"email": reset.Email,
			"token": reset.Token,
		},
		Body: body,
	}
	topic := os.Getenv("SEND_EMAIL_TOPIC")
	err = srv.PubSub.Publish(topic, msg)
	if err != nil {
		log.Printf("[pub]消息发送失败:%v\n", err)
		return err
	}
	return nil
}
