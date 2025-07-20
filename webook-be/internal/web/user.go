package web

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/newton-miku/webook/webook-be/internal/domain"
	"github.com/newton-miku/webook/webook-be/internal/service"
	"github.com/newton-miku/webook/webook-be/internal/web/middleware"
)

type UserHandler struct {
	svc         *service.UserService
	EmailReg    *regexp2.Regexp
	PasswordReg *regexp2.Regexp
	DateReg     *regexp2.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		// 邮箱正则（邮箱用户名部分，可以包含字母、数字、点、下划线、百分号、加号和减号）
		EmailRegPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		// EmailRegPattern = `^\w+(-+.\w+)*@\w+(-.\w+)*.\w+(-.\w+)*$`
		// 密码正则（长度6位以上，包含字母和数字）
		PasswordRegPattern = `^(?=.*[a-zA-Z])(?=.*\d).{6,}$`
		// 日期正则
		DateRegPattern = `^\d{4}-\d{2}-\d{2}$`
	)
	emailReg := regexp2.MustCompile(EmailRegPattern, regexp2.None)
	pwdReg := regexp2.MustCompile(PasswordRegPattern, regexp2.None)
	dateReg := regexp2.MustCompile(DateRegPattern, regexp2.None)
	return &UserHandler{
		svc:         svc,
		EmailReg:    emailReg,
		PasswordReg: pwdReg,
		DateReg:     dateReg,
	}
}

type Msg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/users/signup", u.SignUp)
	server.POST("/users/login", u.Login)
	server.POST("/users/logout", u.Logout)
	server.GET("/users/profile", u.Profile)
	server.POST("/users/edit", u.Edit)
}
func (u *UserHandler) RegisterRoutesV1(ug *gin.RouterGroup) {
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/logout", u.Logout)
	ug.GET("/profile", u.Profile)
	ug.POST("/edit", u.Edit)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	// Bind 方法会根据请求的 Content-Type 来选择绑定器
	// 如果绑定失败，会返回 400
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := u.EmailReg.MatchString(req.Email)
	if err != nil {
		fmt.Println("邮箱校验时系统内部出错,err:", err)
		ctx.JSON(http.StatusOK, Msg{
			Code: 500,
			Msg:  "系统内部出错,请稍后再试",
		})
		return
	}
	if !ok {
		data := Msg{Code: 400, Msg: "邮箱格式有误"}
		ctx.JSON(http.StatusOK, data)
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusOK, Msg{
			Code: 400,
			Msg:  "两次输入的密码不一致",
		})
		return
	}

	ok, err = u.PasswordReg.MatchString(req.Password)
	if err != nil {
		fmt.Println("密码校验时系统内部出错,err:", err)
		ctx.JSON(http.StatusOK, Msg{
			Code: 500,
			Msg:  "系统内部出错,请稍后再试",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Msg{
			Code: 400,
			Msg:  "密码格式有误，至少包含字母、数字，且长度不低于6位",
		})
		return
	}
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: []byte(req.Password),
	})
	if err != nil {
		if errors.Is(err, service.ErrUserDuplicateEmail) {
			ctx.JSON(http.StatusOK, Msg{
				Code: 500,
				Msg:  "该邮箱已被注册",
			})
		} else {
			ctx.JSON(http.StatusOK, Msg{
				Code: 500,
				Msg:  "注册失败，系统内部错误",
			})
			fmt.Println("注册失败,err:", err)
		}
		return
	}

	ctx.JSON(http.StatusOK, Msg{
		Code: 200,
		Msg:  "注册成功",
	})
	// fmt.Printf("req: %v\n", req)

}
func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: []byte(req.Password),
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserOrPassword) {
			ctx.JSON(http.StatusOK, Msg{
				Code: 400,
				Msg:  "邮箱或者密码不正确",
			})
		} else {
			ctx.JSON(http.StatusOK, Msg{
				Code: 500,
				Msg:  "登录时系统发生错误",
			})
		}
		return
	}
	claims := middleware.JWTClaims{
		UserId: user.Id,
	}
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(middleware.JWTExpire))
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(middleware.JwtSecret)
	if err != nil {
		fmt.Println("JWT 签名错误,err:", err)
		ctx.JSON(http.StatusOK, Msg{
			Code: http.StatusInternalServerError,
			Msg:  "登录时系统发生错误",
		})
		return
	}
	ctx.Header("X-JWT-Token", tokenStr)
	print(user.Id)
	ctx.JSON(http.StatusOK, Msg{
		Code: 0,
		Msg:  "登录成功",
	})
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Clear()
	sess.Save()
	ctx.JSON(http.StatusOK, Msg{
		Code: 0,
		Msg:  "登出成功",
	})
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	c, _ := ctx.Get("claims")
	claims, ok := c.(*middleware.JWTClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Msg{
			Code: http.StatusInternalServerError,
			Msg:  "内部错误",
		})
	}
	id := claims.UserId
	user, err := u.svc.Profile(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrProfileNotFound) {
		}
	}
	ctx.JSON(http.StatusOK, user)
}
func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Birthday string `json:"birthday"`
		Nickname string `json:"nickname"`
		Summary  string `json:"aboutMe"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 校验日期格式是否为2023-05-05
	ok, err := u.DateReg.MatchString(req.Birthday)
	if err != nil {
		fmt.Println("生日格式校验时系统内部出错,err:", err)
		ctx.JSON(http.StatusOK, Msg{
			Code: 500,
			Msg:  "系统内部出错,请稍后再试",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Msg{
			Code: 400,
			Msg:  "生日格式有误",
		})
		return
	}
	// 将字符串转换为 time.Time 类型进行日期比较
	reqBirthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		ctx.JSON(http.StatusOK, Msg{
			Code: 400,
			Msg:  "生日解析失败",
		})
		return
	}

	// 获取当前时间并校验请求中的日期是否超过当前时间
	now := time.Now()
	if reqBirthday.After(now) {
		ctx.JSON(http.StatusOK, Msg{
			Code: 400,
			Msg:  "生日不能超过当前日期",
		})
		return
	}

	// 设置允许的最小日期，例如：1900-01-01
	minDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.Now().Location())
	if reqBirthday.Before(minDate) {
		ctx.JSON(http.StatusOK, Msg{
			Code: 400,
			Msg:  "生日不能早于1900年1月1日",
		})
		return
	}

	sess := sessions.Default(ctx)
	id := sess.Get("userId")
	err = u.svc.UpdateProfile(ctx, domain.UserProfile{
		UID:      id.(int64),
		Birthday: req.Birthday,
		Nickname: req.Nickname,
		Summary:  req.Summary,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, fmt.Sprint(err))
		return
	}
	ctx.JSON(http.StatusOK, Msg{Code: 0, Msg: "更新成功"})

}
