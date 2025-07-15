package web

import (
	"fmt"
	"net/http"

	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	EmailReg    *regexp2.Regexp
	PasswordReg *regexp2.Regexp
}

func NewUserHandler() *UserHandler {
	const (
		// 邮箱正则（较为简单，仅验证@和.）
		EmailRegPattern = `^\w+(-+.\w+)*@\w+(-.\w+)*.\w+(-.\w+)*$`
		// 密码正则（长度6位以上，包含字母和数字）
		PasswordRegPattern = `^(?=.*[a-zA-Z])(?=.*\d).{6,}$`
	)
	emailReg := regexp2.MustCompile(EmailRegPattern, regexp2.None)
	pwdReg := regexp2.MustCompile(PasswordRegPattern, regexp2.None)
	return &UserHandler{
		EmailReg:    emailReg,
		PasswordReg: pwdReg,
	}
}

type Msg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/users/signup", u.SignUp)
	server.POST("/users/login", u.Login)
	server.GET("/users/profile", u.Profile)
	server.POST("/users/edit", u.Edit)
}
func (u *UserHandler) RegisterRoutesV1(ug *gin.RouterGroup) {
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
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
		data := Msg{Code: 500, Msg: "系统内部出错"}
		fmt.Println("邮箱校验时系统内部出错,err:", err)
		ctx.JSON(http.StatusOK, data)
		return
	}
	if !ok {
		data := Msg{Code: 400, Msg: "邮箱格式有误"}
		ctx.JSON(http.StatusOK, data)
		return
	}
	if req.Password != req.ConfirmPassword {
		data := Msg{Code: 400, Msg: "两次输入的密码不一致"}
		ctx.JSON(http.StatusOK, data)
		return
	}

	ok, err = u.PasswordReg.MatchString(req.Password)
	if err != nil {
		data := Msg{Code: 500, Msg: "系统内部出错"}
		fmt.Println("密码校验时系统内部出错,err:", err)
		ctx.JSON(http.StatusOK, data)
		return
	}
	if !ok {
		data := Msg{Code: 400, Msg: "密码格式有误，至少包含字母、数字，且长度不低于6位"}
		ctx.JSON(http.StatusOK, data)
		return
	}

	data := Msg{Code: 200, Msg: "注册成功"}
	ctx.JSON(http.StatusOK, data)
	fmt.Printf("req: %v\n", req)

}
func (u *UserHandler) Login(ctx *gin.Context) {

}
func (u *UserHandler) Profile(ctx *gin.Context) {

}
func (u *UserHandler) Edit(ctx *gin.Context) {

}
