package controller

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hsw409328/gofunc"
	"log"
	"net/http"
	"union-data-analysis/model"
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

func CheckPay() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		v := session.Get("loginUser")
		if v == nil {
			ctx.Redirect(http.StatusFound, "/login")
		}
		u := session.Get("loginUser")
		var tmpUser model.WebUserData
		err := json.Unmarshal([]byte(gofunc.InterfaceToString(u)), &tmpUser)
		if err != nil {
			log.Println(err)
			ctx.Redirect(http.StatusFound, "/login")
		}
		//查询是否已经支付
		w, _ := new(model.WebPayOrder).Where(map[string]string{"手机号": "=" + tmpUser.Mobile, "状态": "已支付"}).GetOne()
		if w.OrderId == "" {
			ctx.Redirect(http.StatusFound, "/pay")
		}
		ctx.Redirect(http.StatusFound, "/")
	}
}
