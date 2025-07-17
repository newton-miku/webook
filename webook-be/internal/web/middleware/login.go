package middleware

import (
	"net/http"

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
	}
}
