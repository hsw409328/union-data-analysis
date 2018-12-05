package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		v := session.Get("loginUser")
		if v == nil {
			ctx.Redirect(http.StatusFound, "/login")
		}
	}
}

func CheckLoginJson() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		v := session.Get("loginUser")
		if v == nil {
			ctx.JSON(http.StatusOK, gin.H{"m": "未登录"})
			return
		}
	}
}

func CheckHaveLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		v := session.Get("loginUser")
		if v != nil {
			ctx.Redirect(http.StatusFound, "/")
		}
	}
}
