package middleware

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddleware struct {
	paths []string
}

func NewLoginMiddleware() *LoginMiddleware {
	return &LoginMiddleware{}
}

func (l *LoginMiddleware) AddIgnorePaths(paths []string) *LoginMiddleware {
	l.paths = append(l.paths, paths...)
	return l
}
func (l *LoginMiddleware) AddIgnorePath(path ...string) *LoginMiddleware {
	l.paths = append(l.paths, path...)
	return l
}

func (l *LoginMiddleware) Build() gin.HandlerFunc {
	// 注册time.Now()类型，让cookie支持存储时间数据
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()
		val := sess.Get("lastUpdateTime")
		lastUpdateTime, ok := val.(time.Time)
		if val == nil || !ok || now.Sub(lastUpdateTime) > 30*time.Second {
			sess.Set("lastUpdateTime", now)
			sess.Set("userId", id)
			sess.Options(sessions.Options{
				MaxAge: 60,
			})
			sess.Save()
		}
	}
}
