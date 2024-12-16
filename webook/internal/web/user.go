package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern        = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	ByteKey              string = "q0@m6)ay3(Na094ShBq9nfb=nW*D{4c"
)

type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	// REST 风格
	//server.POST("/user", h.SignUp)
	//server.PUT("/user", h.SignUp)
	//server.GET("/users/:username", h.Profile)
	ug := server.Group("/users")
	// POST /users/signup
	ug.POST("/signup", h.SignUp)
	// POST /users/login
	ug.POST("/login", h.LoginJWT)
	// POST /users/edit
	ug.POST("/edit", h.Edit)
	// GET /users/profile
	ug.GET("/profile", h.Profile)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不对")
		return
	}

	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含字母、数字、特殊字符，并且不少于八位")
		return
	}

	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, gin.H{"msg": "注册成功"})
	case errors.Is(err, service.ErrDuplicateEmail):
		ctx.JSON(http.StatusOK, gin.H{"msg": "邮箱冲突，请换一个"})
	default:
		ctx.JSON(http.StatusOK, gin.H{"msg": "系统错误"})
	}
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "参数错误",
		})
		return
	}

	//拿到数据库Find返回的信息
	user, err := h.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "用户名或密码错误",
		})
	}
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "系统错误",
		})
	}

	//cookie
	sess := sessions.Default(ctx)
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		MaxAge: 60,
	})
	sess.Save()

	// jwt携带token，方便后续Profiel从中解析uid然后去数据库查表
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
		Uid: user.Id,
	})
	tokenStr, err := token.SignedString([]byte(ByteKey))

	ctx.Header("x-jwt-token", tokenStr)
	ctx.String(http.StatusOK, "登录成功")
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "登录成功",
	})

}

func (h *UserHandler) Edit(ctx *gin.Context) {

}

func (h *UserHandler) ProfileJWT(ctx *gin.Context) {
	claims, ok := ctx.Get("claims")
	if !ok {
		zap.S().Error(" Profile 没有拿到claims")
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "系统错误",
		})
		return
	}
	claims, ok = claims.(*UserClaims)
	if !ok {
		zap.S().Error(" Profile claims错误")
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "系统错误",
		})
		return
	}

	ctx.String(http.StatusOK, "这是 profile")
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "这是 profile")
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid int64
}
