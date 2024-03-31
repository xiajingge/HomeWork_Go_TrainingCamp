package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 静态路由
	server.GET("/login", func(c *gin.Context) {
		c.String(http.StatusOK, "这是静态路由显示内容")
	})

	// 参数路由
	server.GET("/login:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "name : "+name)
	})

	server.GET("/login2", func(c *gin.Context) {
		name := c.Query("name")
		c.String(http.StatusOK, "name ：%s", name)
	})

	// 通配符路由
	server.GET("/login/*.html", func(c *gin.Context) {
		name := c.Param(".html")
		c.String(http.StatusOK, "name : %s", name)
	})
}

func main() {
	// 初始化web 引擎
	server := gin.Default()

	// 设置路由
	// 静态路由
	server.GET("/login", func(c *gin.Context) {
		c.String(http.StatusOK, "这是静态路由显示内容")
	})

	// 参数路由
	server.GET("/login:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "name : "+name)
	})

	server.GET("/login2", func(c *gin.Context) {
		name := c.Query("name")
		c.String(http.StatusOK, "name ：%s", name)
	})

	// 通配符路由
	server.GET("/login/*.html", func(c *gin.Context) {
		name := c.Param(".html")
		c.String(http.StatusOK, "name : %s", name)
	})

	// 运行web服务,如果不写参数，默认是8080端口，注意参数前有冒号
	go server.Run(":8080")

	server2 := gin.Default()
	server2.POST("/log", func(c *gin.Context) {
		c.String(http.StatusOK, "这是监听8081端口的POST服务")
	})
	server2.Run(":8081")

}
