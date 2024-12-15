package web

import "github.com/gin-gonic/gin"

type UserHandler struct {
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/users")
	group.POST("/login", u.Login)
	group.POST("/signUp", u.SignUp)
	group.POST("/edit", u.Edit)
	group.POST("/profile", u.Profile)
}

func (u *UserHandler) Login(ctx *gin.Context) {}

func (u *UserHandler) SignUp(ctx *gin.Context) {}

func (u *UserHandler) Edit(ctx *gin.Context) {}

func (u *UserHandler) Profile(ctx *gin.Context) {}
