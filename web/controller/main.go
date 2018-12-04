package controller

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gofunc"
	"net/http"
	"union-data-analysis/enums"
	"union-data-analysis/model"
)

var (
	webUserModel = new(model.WebUsers)
)

type MainController struct {
	*gin.Context
}

func (ctx *MainController) SendMobile(c *gin.Context) {
	m := c.PostForm("mobile")
	if m == "" {
		c.JSON(http.StatusInternalServerError, "手机号不能为空")
		return
	}
	session := sessions.Default(c)
	checkIsSend := session.Get(m)
	if checkIsSend != nil {
		c.JSON(http.StatusInternalServerError, "不要重复发送")
		return
	}
	authCode := GetRandomSalt()
	err := SendMobile(m, authCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "发送失败，请检查手机号是否正确")
		return
	}
	session.Set(m, authCode)
	session.Save()
	c.JSON(http.StatusOK, "发送成功")
	return
}

func (ctx *MainController) Login(c *gin.Context) {
	m := c.PostForm("mobile")
	a := c.PostForm("authCode")
	session := sessions.Default(c)
	oldAuthCode := session.Get(m)
	if oldAuthCode == nil {
		c.JSON(http.StatusInternalServerError, "未发送短信")
		return
	}
	if a != oldAuthCode {
		c.JSON(http.StatusInternalServerError, "验证失败")
		return
	}
	u, err := webUserModel.Where(map[string]string{
		"手机号": " ='" + m + "'",
	}).GetOne()
	if err != nil {
		if u.Mobile == "" {
			w := model.WebUserData{
				Mobile:     m,
				State:      enums.UserStateApply,
				CreateTime: gofunc.CurrentTime(),
				UpdateTime: gofunc.CurrentTime(),
			}
			webUserModel.Insert(w)
			by, err := json.Marshal(w)
			if err != nil {
				lg.Error(err)
			}
			session.Set("loginUser", string(by))
			session.Save()
			c.JSON(http.StatusOK, enums.UserStateApply)
			return
		}
	}
	by, err := json.Marshal(u)
	if err != nil {
		lg.Error(err)
	}
	session.Set("loginUser", string(by))
	session.Save()
	if u.State == enums.UserStateApply {
		c.JSON(http.StatusOK, enums.UserStateApply)
		return
	} else if u.State == enums.UserStateNo {
		c.JSON(http.StatusOK, enums.UserStateNo)
		return
	} else {
		c.JSON(http.StatusOK, enums.UserStateYes)
		return
	}
}
