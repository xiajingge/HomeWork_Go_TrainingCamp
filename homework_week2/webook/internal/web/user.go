// 包含用户的调用的方法，利于注册、登录等单独的方法
package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"my_record/webook/internal/domain"
	"my_record/webook/internal/service"
	"net/http"
	"time"
)

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,72}$`
	bizLogin             = "login"
)

type UserHandler struct {
	emailRegexExp *regexp.Regexp
	passwordRegex *regexp.Regexp
	svc           *service.UserServer
}

func NewUserHandler(svc *service.UserServer) *UserHandler {
	return &UserHandler{
		emailRegexExp: regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegex: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:           svc,
	}
}

// 用户注册时的路由响应
func (c *UserHandler) Signup(ctx *gin.Context) {
	// 创建接收请求数据的结构：最小化原则
	type SignupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	// 使用bind方法来根据前端数据格式来自动填充数据结构
	var req SignupReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 验证邮箱格式
	isEmail, err := c.emailRegexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}
	// 验证密码格式 英文字母、数字、特殊符号
	isPassword, err := c.passwordRegex.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含字母、数字、特殊符号，并且长度不能小于8位")
		return
	}

	// 验证两次密码是否相同
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次密码输入不相同")
		return
	}

	// gorm 数据库部分功能实现
	// 数据存储
	err = c.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	switch err {
	case nil:
		ctx.String(http.StatusOK, "注册成功")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "邮箱已被注册")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

// 用户登录时的路由响应
func (c *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := c.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			MaxAge: 15 * 60,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (c *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Nickname string `json:"nickname"` // 昵称
		Birthday string `json:"birthday"` // 生日 YYYY-MM-DD
		AboutMe  string `json:"aboutMe"`  //个人简介
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	sess := sessions.Default(ctx)
	userId := sess.Get("userId")
	if userId == nil { // 检查userId是否存在
		return
	}
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "输入的生日格式不对，请修改为YYYY-MM-DD格式,例如1999-01-01")
	}
	c.svc.UpdateNonSentiveInfo(ctx, domain.User{
		Id:       userId.(int64),
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})
}

func (c *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "你正在查看内容")
}

// 用户需使用到的路由
func (c *UserHandler) RigisterRoutes(server *gin.Engine) {
	//server.POST("/users/signup", c.Signup)
	//server.POST("/users/login", c.Login)
	//server.POST("/users/edit", c.Edit)
	//server.GET("/users/profile", c.Profile)

	// 优化：利用gin的分组路由功能
	s := server.Group("/users")
	s.POST("/signup", c.Signup)
	s.POST("/login", c.Login)
	s.POST("/edit", c.Edit)
	s.GET("/profile", c.Profile)
}
