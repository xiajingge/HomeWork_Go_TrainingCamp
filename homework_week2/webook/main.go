package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"my_record/webook/internal/repository"
	"my_record/webook/internal/repository/dao"
	"my_record/webook/internal/service"
	"my_record/webook/internal/web"
	"my_record/webook/internal/web/middleware"
	"strings"
	"time"
)

func main() {
	db := initDB()
	//db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	//if err != nil {
	//	panic(err)
	//}
	//err = dao.InitTables(db)

	server := initWebServer()
	//// web引擎创建
	//server := gin.Default()
	//
	//// 跨域问题解决
	//server.Use(cors.New(cors.Config{
	//	//AllowAllOrigins: true,
	//	//AllowOrigins:     []string{"http://localhost:3000"},
	//	AllowCredentials: true, // 是否允许用户带上用户认证信息，比如cookie
	//
	//	//AllowHeaders: []string{"content-type"},
	//	AllowHeaders: []string{"Content-Type", "Authorization"}, // 业务请求中允许带的头
	//	// 这个是允许前端访问你的后端响应中带的头部
	//	ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
	//	//AllowMethods: []string{"POST"}, // 允许的路由，一般设置为全部
	//	AllowOriginFunc: func(origin string) bool {
	//		if strings.HasPrefix(origin, "http://localhost") {
	//			//if strings.Contains(origin, "localhost") {
	//			return true
	//		}
	//		return strings.Contains(origin, "your_company.com")
	//	}, // 哪些来源是允许的
	//	MaxAge: 12 * time.Hour,
	//}))

	// 用户要使用的方法调用
	hdl := initUserHdl(db)
	hdl.RigisterRoutes(server)
	// 启动引擎
	server.Run(":8080")
}

func initUserHdl(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDao(db)
	ur := repository.NewUserReopsitory(ud)
	us := service.NewUserServer(ur)
	hdl := web.NewUserHandler(us)
	//var hdl = &web.UserHandler{}
	return hdl
}

// 初始化建立数据库
func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

// 初始化建立server
func initWebServer() *gin.Engine {
	// web引擎创建
	server := gin.Default()

	// 跨域问题解决
	server.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		//AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true, // 是否允许用户带上用户认证信息，比如cookie

		//AllowHeaders: []string{"content-type"},
		AllowHeaders: []string{"Content-Type", "Authorization"}, // 业务请求中允许带的头
		// 这个是允许前端访问你的后端响应中带的头部
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
		//AllowMethods: []string{"POST"}, // 允许的路由，一般设置为全部
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				//if strings.Contains(origin, "localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		}, // 哪些来源是允许的
		MaxAge: 12 * time.Hour,
	}))

	// 用于提取session
	// 存储数据的，也就是我的 userId 存哪里
	store := cookie.NewStore([]byte("secret")) // 初学者 直接存cookie里面
	server.Use(sessions.Sessions("ssid", store))

	// 执行登录校验
	login := &middleware.LoginMiddlewareBuilder{}
	server.Use(login.CheckLogin())
	return server
}
