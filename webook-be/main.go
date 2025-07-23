package main

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/newton-miku/webook/webook-be/internal/config"
	"github.com/newton-miku/webook/webook-be/internal/repository"
	"github.com/newton-miku/webook/webook-be/internal/repository/dao"
	"github.com/newton-miku/webook/webook-be/internal/service"
	"github.com/newton-miku/webook/webook-be/internal/web"
	"github.com/newton-miku/webook/webook-be/internal/web/middleware"
	"github.com/newton-miku/webook/webook-be/pkg/ginx/middleware/ratelimit"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := initDB()
	server := initWebServer()

	user := initUser(db)
	user.RegisterRoutesV1(server.Group("/users"))

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})
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
		AllowHeaders:  []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"X-JWT-Token"},
		// 是否允许携带cookie
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return strings.Contains(origin, "http://localhost")
		},
		MaxAge: 12 * time.Hour,
	}))

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})

	server.Use(ratelimit.NewBuilder(redisClient, time.Minute, 50).Build())

	// store, err := redis.NewStore(16, "tcp", "localhost:6379", "", "", []byte("eY3VQBCzq8p748ME20cWWBuJs7ZVqN9W"), []byte("f2Ug8BsFjZmuUyAVi11ZA2J36Sc0RXwE"))
	// // 此处的两个key长度必须是2的整数倍，否则会报错
	// if err != nil {
	// 	panic(err)
	// }
	// // store := cookie.NewStore([]byte("secret"))
	// // store := memstore.NewStore([]byte("secret"))
	// server.Use(sessions.Sessions("mysession", store))

	server.Use(middleware.NewJWTLoginMiddleware().
		AddIgnorePath("/ping").
		AddIgnorePath("/users/login").
		AddIgnorePath("/users/signup").
		Build())
	return server
}

func initDB() *gorm.DB {
	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
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
