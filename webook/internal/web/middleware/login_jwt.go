package middleware

import (
	"basic-go/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path ...string) *LoginJWTMiddlewareBuilder {
	for _, p := range path {
		l.paths = append(l.paths, p)
	}
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		tokenHeader := ctx.GetHeader("x-jwt-token")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := segs[1]
		claims := &web.UserClaims{}
		// ParseWithClaims 里面一定要传入指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("q0@m6)ay3(Na094ShBq9nfb=nW*D{4c"), nil
		})
		if err != nil || !token.Valid || claims.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt.Add(time.Minute)
			tokenStr, err = token.SignedString([]byte("q0@m6)ay3(Na094ShBq9nfb=nW*D{4c"))
			if err != nil {
				zap.S().Error("  jwt续约失败: " + err.Error())
			}
			ctx.Header("x-jwt-token", tokenStr)
		}

		ctx.Set("claims", claims)
	}
}
