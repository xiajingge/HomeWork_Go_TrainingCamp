// data access object 数据访问对象，代表数据库操作
package dao

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

var (
	ErrDuplicateEmail = errors.New("邮箱冲突，该邮箱已被注册")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

// 向数据库中插入数据，也可以叫存储数据
func (dao *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now // 创建时间
	u.Utime = now // 更新时间

	err := dao.db.WithContext(ctx).Create(&u).Error // 写入数据库
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 用户冲突，邮箱冲突
			return ErrDuplicateEmail
		}
	}
	return err
}

// 更新数据库中相应数据
func (dao *UserDao) UpdateById(ctx *gin.Context, u User) error {
	return dao.db.WithContext(ctx).Model(&u).Where("id = ?", u.Id).
		Updates(map[string]any{
			"utime":    time.Now().UnixMilli(),
			"nickname": u.Nickname,
			"birthday": u.Birthday,
			"about_me": u.AboutMe,
		}).Error
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

// 要往数据库中写的数据，domain的数据不能代表数据库
type User struct {
	// 增加自增组件
	Id int64 `gorm:"primary_key,autoIncrement"`
	// 邮箱应该是唯一属性，每个用户的Email肯定是不相同的
	Email    string `gorm:"unique"`
	Password string
	Nickname string
	Birthday time.Time
	AboutMe  string
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}
