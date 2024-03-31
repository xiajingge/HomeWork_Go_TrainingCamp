// 与业务强相关的
package repository

import (
	"context"
	"github.com/gin-gonic/gin"
	"my_record/webook/internal/domain"
	"my_record/webook/internal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserReopsitory struct {
	dao *dao.UserDao
}

func NewUserReopsitory(dao *dao.UserDao) *UserReopsitory {
	return &UserReopsitory{
		dao: dao,
	}
}

func (repo *UserReopsitory) Create(ctx context.Context, u domain.User) error {
	err := repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	if err == dao.ErrDuplicateEmail {
		return ErrDuplicateEmail
	}
	return err
}

func (repo *UserReopsitory) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDoain(u), err
}

func (repo *UserReopsitory) toDoain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}
}

func (repo *UserReopsitory) UpdateNonZeroFields(ctx *gin.Context, user domain.User) error {
	return repo.dao.UpdateById(ctx, dao.User{
		Id:       user.Id,
		Nickname: user.Nickname,
		Birthday: user.Birthday,
		AboutMe:  user.AboutMe,
	})
}
