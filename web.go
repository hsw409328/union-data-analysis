package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
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
	routers.GET("/person", controller.CheckLogin(), person)
	routers.GET("/apply", controller.CheckLogin(), apply)
	routers.GET("/no/apply", controller.CheckLogin(), noApply)

	routers.POST("/api/send/message", (&controller.MainController{}).SendMobile)
	routers.POST("/api/login", (&controller.MainController{}).Login)

	routers.Run()
}

func index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{})
}
func login(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{})
}
func person(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "data.html", gin.H{})
}
func apply(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "audit.html", gin.H{})
}
func noApply(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "audit_no.html", gin.H{})
}
