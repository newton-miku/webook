package main

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/newton-miku/webook/webook-be/internal/repository"
	"github.com/newton-miku/webook/webook-be/internal/repository/dao"
	"github.com/newton-miku/webook/webook-be/internal/service"
	"github.com/newton-miku/webook/webook-be/internal/web"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := initDB()
	server := initWebServer()

	user := initUser(db)
	user.RegisterRoutesV1(server.Group("/users"))

	server.Run(":8080")
}

func initUser(db *gorm.DB) *web.UserHandler {
	userDao := dao.NewUserDAO(db)
	resp := repository.NewUserRepository(userDao)
	svc := service.NewUserService(resp)
	user := web.NewUserHandler(svc)
	return user
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	// 处理跨域插件
	server.Use(cors.New(cors.Config{
		// AllowOrigins: []string{"http://localhost:3000"},
		// AllowMethods: []string{"PUT", "PATCH", "POST"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 是否允许携带cookie
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return strings.Contains(origin, "http://localhost")
		},
		MaxAge: 12 * time.Hour,
	}))
	return server
}

func initDB() *gorm.DB {
	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13306)/webook"))
	if err != nil {
		panic(err)
	}

	// 建表
	err = dao.InitTable(db)
	if err != nil {
		// 如果不成功则panic
		panic(err)
	}
	return db
}
