package main

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/newton-miku/webook/webook-be/internal/web"
)

func main() {
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
	user := web.NewUserHandler()
	user.RegisterRoutesV1(server.Group("/users"))
	server.Run(":8080")
}
