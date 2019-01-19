package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"union-data-analysis/enums"
	"union-data-analysis/web/controller"
)

func main() {
	routers := gin.Default()
	store := cookie.NewStore([]byte("xiaoqqq"))
	routers.Use(sessions.Sessions("myssssssss", store))

	routers.Static("/css", "./web/css")
	routers.Static("/images", "./web/images")
	routers.Static("/js", "./web/js")

	routers.LoadHTMLGlob("./web/views/*")
	routers.GET("/", controller.CheckLogin(), index)
	routers.GET("/login", controller.CheckHaveLogin(), login)
	routers.GET("/pay", controller.CheckHaveLogin(), pay)
	routers.GET("/person", controller.CheckLogin(), person)
	routers.GET("/apply", controller.CheckLogin(), apply)
	routers.GET("/no/apply", controller.CheckLogin(), noApply)

	routers.POST("/api/send/message", (&controller.MainController{}).SendMobile)
	routers.POST("/api/login", (&controller.MainController{}).Login)

	routers.GET("/api/get/user", controller.CheckLoginJson(), (&controller.MainController{}).GetUserInfo)
	routers.POST("/api/update/user", controller.CheckLoginJson(), (&controller.MainController{}).UpdateUserInfo)

	routers.GET("/api/get/all", controller.CheckLoginJson(), (&controller.MainController{}).GetAllMoney)
	routers.GET("/api/get/last", controller.CheckLoginJson(), (&controller.MainController{}).GetLastMoney)
	routers.GET("/api/get/current", controller.CheckLoginJson(), (&controller.MainController{}).GetCurrentMoney)
	routers.GET("/api/get/recommend", controller.CheckLoginJson(), (&controller.MainController{}).GetRecommend)

	routers.GET("/pay_callback", (&controller.MainController{}).PayCallBack)

	routers.Run()
}

func index(ctx *gin.Context) {
	tmpUser := (&controller.MainController{}).GetLoginUserInfo(ctx)
	if tmpUser.State == enums.UserStateNo {
		ctx.Redirect(http.StatusFound, "/no/apply")
	} else if tmpUser.State == enums.UserStateApply {
		ctx.Redirect(http.StatusFound, "/apply")
	}
	ctx.HTML(http.StatusOK, "data.html", gin.H{})
}
func login(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{})
}
func pay(ctx *gin.Context) {
	a := (&controller.MainController{}).GetPayCode(ctx)
	ctx.HTML(http.StatusOK, "pay.html", gin.H{
		"img_src": a,
	})
}
func person(ctx *gin.Context) {
	tmpUser := (&controller.MainController{}).GetLoginUserInfo(ctx)
	if tmpUser.State == enums.UserStateNo {
		ctx.Redirect(http.StatusFound, "/no/apply")
	} else if tmpUser.State == enums.UserStateApply {
		ctx.Redirect(http.StatusFound, "/apply")
	}
	ctx.HTML(http.StatusOK, "index.html", gin.H{})
}
func apply(ctx *gin.Context) {
	tmpUser := (&controller.MainController{}).GetLoginUserInfo(ctx)
	if tmpUser.State == enums.UserStateNo {
		ctx.Redirect(http.StatusFound, "/no/apply")
	} else if tmpUser.State == enums.UserStateYes {
		ctx.Redirect(http.StatusFound, "/")
	}
	ctx.HTML(http.StatusOK, "audit.html", gin.H{})
}

func noApply(ctx *gin.Context) {
	tmpUser := (&controller.MainController{}).GetLoginUserInfo(ctx)
	if tmpUser.State == enums.UserStateApply {
		ctx.Redirect(http.StatusFound, "/apply")
	} else if tmpUser.State == enums.UserStateYes {
		ctx.Redirect(http.StatusFound, "/")
	}
	ctx.HTML(http.StatusOK, "audit_no.html", gin.H{})
}
