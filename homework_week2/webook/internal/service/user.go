// 业务完整的处理过程
package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"my_record/webook/internal/domain"
	"my_record/webook/internal/repository"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("用户不存在或密码不对")
)

type UserServer struct {
	repo *repository.UserReopsitory
}

func NewUserServer(repo *repository.UserReopsitory) *UserServer {
	return &UserServer{
		repo: repo,
	}
}

// 注册服务 中的数据操作
// 第二个字段需要建立一个自己的数据存储格式，不能直接调要存储的数据，应该是存储的数据调用这个函数
func (svc *UserServer) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	err = svc.repo.Create(ctx, u)
	if err == ErrDuplicateEmail {
		return ErrDuplicateEmail
	}
	return err
}

func (svc *UserServer) Login(ctx context.Context, email string, password string) (domain.User, error) {
	// 判断用户的邮箱是否存在
	u, err := svc.repo.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	// 检查密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserServer) UpdateNonSentiveInfo(ctx *gin.Context, user domain.User) error {
	return svc.repo.UpdateNonZeroFields(ctx, user)

}
