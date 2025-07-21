package middleware

import (
	"encoding/gob"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginJWTMiddleware struct {
	paths []string
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UserId    int64
	UserAgent string
}

var (
	JwtSecret = []byte("f2Ug8BsFjZmuUyAVi11ZA2J36Sc0RXwE")
)

const (
	JWTExpire = time.Minute * 5
)

func NewJWTLoginMiddleware() *LoginJWTMiddleware {
	return &LoginJWTMiddleware{}
}

func (l *LoginJWTMiddleware) AddIgnorePaths(paths []string) *LoginJWTMiddleware {
	l.paths = append(l.paths, paths...)
	return l
}
func (l *LoginJWTMiddleware) AddIgnorePath(path ...string) *LoginJWTMiddleware {
	l.paths = append(l.paths, path...)
	return l
}

func (l *LoginJWTMiddleware) Build() gin.HandlerFunc {
	// 注册time.Now()类型，让cookie支持存储时间数据
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenParts := strings.Split(tokenHeader, " ")
		if len(tokenParts) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := tokenParts[1]
		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.UserId == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			// ua异常
			// 可能是JWT泄露或浏览器更新
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		if time.Until(claims.ExpiresAt.Time) < (JWTExpire / 2) {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(JWTExpire))
			token = jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
			tokenStr, _ := token.SignedString(JwtSecret)
			ctx.Header("X-JWT-Token", tokenStr)
		}
		ctx.Set("claims", claims)
		ctx.Set("userId", claims.UserId)
	}
}
