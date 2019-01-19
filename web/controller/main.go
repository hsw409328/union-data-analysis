package controller

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-with/wxpay"
	"github.com/hsw409328/gofunc"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
	"union-data-analysis/enums"
	"union-data-analysis/model"
)

var (
	webUserModel      = new(model.WebUsers)
	webSendRecord     = new(model.WebSendRecord)
	webRecommendModel = new(model.WebFirstRecommend)
	webPayOrderModel  = new(model.WebPayOrder)
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

func (ctx *MainController) GetLoginUserInfo(c *gin.Context) model.WebUserData {
	s := sessions.Default(c)
	u := s.Get("loginUser")
	var tmpUser model.WebUserData
	err := json.Unmarshal([]byte(gofunc.InterfaceToString(u)), &tmpUser)
	if err != nil {
		log.Println(err)
	}
	w, err := webUserModel.Where(map[string]string{
		"手机号": " ='" + tmpUser.Mobile + "' ",
	}).GetOneAllField()
	if err != nil {
		log.Println(err)
	}
	return w
}

func (ctx *MainController) GetUserInfo(c *gin.Context) {
	u := ctx.GetLoginUserInfo(c)
	if u.Mobile == "" {
		c.JSON(http.StatusOK, gin.H{"data": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": u})
	return
}

func (ctx *MainController) GetAllMoney(c *gin.Context) {
	u := ctx.GetLoginUserInfo(c)
	r := webSendRecord.Where(map[string]string{
		"手机号": " = '" + u.Mobile + "'",
	}).GetAll()
	if len(r) <= 0 {
		c.JSON(http.StatusOK, gin.H{"data": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
	return
}

func (ctx *MainController) UpdateUserInfo(c *gin.Context) {
	bankAccount := template.HTMLEscapeString(c.PostForm("ba"))
	bankUser := template.HTMLEscapeString(c.PostForm("bu"))
	bankName := template.HTMLEscapeString(c.PostForm("bn"))
	u := ctx.GetLoginUserInfo(c)
	_, err := webUserModel.Update(model.WebUserData{
		BankName:        bankName,
		BankUserName:    bankUser,
		BankUserAccount: bankAccount,
		Mobile:          u.Mobile,
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"data": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "更新成功"})
	return
}

func (ctx *MainController) GetLastMoney(c *gin.Context) {
	u := ctx.GetLoginUserInfo(c)
	a := gofunc.TimeUnixIntToStringCustom(gofunc.LastTime("m", -1), "2006-01")
	aSlice := strings.Split(a, "-")
	r := webSendRecord.Where(map[string]string{
		"手机号": " = '" + u.Mobile + "'",
		"年份":  " = '" + aSlice[0] + "'",
		"月份":  " = '" + aSlice[1] + "'",
	}).GetOne()
	if r.Mobile == "" {
		c.JSON(http.StatusOK, gin.H{"data": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
}

func (ctx *MainController) GetCurrentMoney(c *gin.Context) {
	u := ctx.GetLoginUserInfo(c)
	a := time.Now().Format("2006-01")
	aSlice := strings.Split(a, "-")
	r := webSendRecord.Where(map[string]string{
		"手机号": " = '" + u.Mobile + "'",
		"年份":  " = '" + aSlice[0] + "'",
		"月份":  " = '" + aSlice[1] + "'",
	}).GetOne()
	if r.Mobile == "" {
		c.JSON(http.StatusOK, gin.H{"data": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
}

func (ctx *MainController) GetRecommend(c *gin.Context) {
	u := ctx.GetLoginUserInfo(c)
	r := webRecommendModel.Where(map[string]string{
		"手机号": " = '" + u.Mobile + "'",
	}).GetAll()
	if len(r) <= 0 {
		c.JSON(http.StatusOK, gin.H{"data": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
}

func (ctx *MainController) GetPayCode(c *gin.Context) string {
	u := ctx.GetLoginUserInfo(c)
	orderId := gofunc.GetGuuid()
	webPayOrderModel.Insert(model.WebPayOrderData{
		OrderId:    orderId,
		Mobile:     u.Mobile,
		PayFee:     99,
		OpenId:     "",
		CreateTime: gofunc.CurrentTime(),
		UpdateTime: "",
		State:      `未支付`,
	})
	wx := wxpay.NewClient("wx7b6a4b52b3472bc0", "1339710501", "1kue34jiueiuieuqeoiquweoqiuweoiq")
	params := make(wxpay.Params)
	params.SetString("appid", wx.AppId)
	params.SetString("mch_id", wx.MchId)
	params.SetString("body", "支付")
	params.SetString("out_trade_no", orderId)
	params.SetString("total_fee", "1")
	params.SetString("spbill_create_ip", gofunc.GetLocalIp())
	params.SetString("notify_url", "http://partner.51xiaoq.com/pay_callback")
	params.SetString("nonce_str", gofunc.GetGuuid())
	params.SetString("trade_type", "NATIVE")
	params.SetString("sign", wx.Sign(params))

	url := "https://api.mch.weixin.qq.com/pay/unifiedorder"

	// 发送查询企业付款请求
	ret, err := wx.Post(url, params, false)
	if err != nil {
		log.Println(err)
		return ""
	}
	//map[appid:wx7b6a4b52b3472bc0 sign:76647F5354801D41D31A0B7DE6D4D1AA prepay_id:wx19112946883652944d18326a2956806708 code_url:weixin://wxpay/bizpayurl?pr=sHjpmfL return_code:SUCCESS return_msg:OK mch_id:1339710501 nonce_str:DMUh1puEw5g4G5po result_code:SUCCESS trade_type:NATIVE]
	if ret.GetString("return_code") == "SUCCESS" {
		return ret.GetString("code_url")
	}
	return ""
}

func (ctx *MainController) PayCallBack(c *gin.Context) {
	a := webPayOrderModel.WxpayCallback(c.Request)
	c.XML(http.StatusOK, a)
}
