package model

import (
	"encoding/xml"
	"github.com/hsw409328/gofunc"
	"io/ioutil"
	"log"
	"net/http"
	"union-data-analysis/lib/driver"
)

type WXPayNotifyReq struct {
	Return_code    string `xml:"return_code"`
	Return_msg     string `xml:"return_msg"`
	Appid          string `xml:"appid"`
	Mch_id         string `xml:"mch_id"`
	Nonce          string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	Result_code    string `xml:"result_code"`
	Openid         string `xml:"openid"`
	Is_subscribe   string `xml:"is_subscribe"`
	Trade_type     string `xml:"trade_type"`
	Bank_type      string `xml:"bank_type"`
	Total_fee      int    `xml:"total_fee"`
	Fee_type       string `xml:"fee_type"`
	Cash_fee       int    `xml:"cash_fee"`
	Cash_fee_Type  string `xml:"cash_fee_type"`
	Transaction_id string `xml:"transaction_id"`
	Out_trade_no   string `xml:"out_trade_no"`
	Attach         string `xml:"attach"`
	Time_end       string `xml:"time_end"`
}

type WXPayNotifyResp struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
}

type WebPayOrderData struct {
	RowId      int
	OrderId    string
	Mobile     string
	PayFee     int
	OpenId     string
	CreateTime string
	UpdateTime string
	State      string `未支付 已支付`
}

type WebPayOrder struct {
	where string
	group string
}

func NewWebPayOrder() *WebPayOrder {
	return &WebPayOrder{}
}

func (ctx *WebPayOrder) Where(w map[string]string) *WebPayOrder {
	ctx.where = " where 1 "
	for k, v := range w {
		ctx.where += " and " + k + v
	}
	return ctx
}

func (ctx *WebPayOrder) Group(group string) *WebPayOrder {
	if group != "" {
		ctx.group = " group by " + group
	}
	return ctx
}

func (ctx *WebPayOrder) GetOne() (WebPayOrderData, error) {
	r := driver.SQLiteDriverWeb.GetOne("select 订单ID, 手机号, 支付金额,OPENID,创建时间,更新时间,状态 from " + WebPayOrderTableName +
		ctx.where + ctx.group)
	var WebPayOrderData = new(WebPayOrderData)
	err := r.Scan(&WebPayOrderData.OpenId, &WebPayOrderData.Mobile, &WebPayOrderData.PayFee, &WebPayOrderData.OpenId,
		&WebPayOrderData.CreateTime, &WebPayOrderData.UpdateTime, &WebPayOrderData.State)
	if err != nil {
		lg.Error(err.Error())
		return *WebPayOrderData, err
	}
	return *WebPayOrderData, nil
}

func (ctx *WebPayOrder) Insert(w WebPayOrderData) (int64, error) {
	n, err := driver.SQLiteDriverWeb.Insert(
		"insert into ["+WebPayOrderTableName+"]([订单ID],[手机号],[支付金额],[OPENID],[创建时间],[更新时间],[状态])"+
			" values(?, ?, ?, ?, ?, ?, ?)",
		w.OrderId, w.Mobile, w.PayFee, "", w.CreateTime, w.UpdateTime, w.State,
	)
	return n, err
}

func (ctx *WebPayOrder) Update(w WebPayOrderData) (int64, error) {
	n, err := driver.SQLiteDriverWeb.Update("UPDATE "+WebUserTableName+" SET OPENID = ?,"+
		"更新时间 = ?,状态 = ? WHERE 订单ID = ?",
		w.OpenId, w.UpdateTime, w.State,
		w.OrderId)
	if err != nil {
		log.Println(err)
	}
	return n, err
}

//具体的微信支付回调函数的范例
func (ctx *WebPayOrder) WxpayCallback(r *http.Request) WXPayNotifyResp {
	var resp WXPayNotifyResp
	resp.Return_code = "FAIL"
	resp.Return_msg = "failed to verify sign, please retry!"
	// body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("读取http body失败，原因!", err)
		return resp
	}
	defer r.Body.Close()

	var mr WXPayNotifyReq
	err = xml.Unmarshal(body, &mr)
	if err != nil {
		log.Println("解析HTTP Body格式到xml失败，原因!", err)
		return resp
	}

	var reqMap map[string]interface{}
	reqMap = make(map[string]interface{}, 0)

	reqMap["return_code"] = mr.Return_code
	reqMap["return_msg"] = mr.Return_msg
	reqMap["appid"] = mr.Appid
	reqMap["mch_id"] = mr.Mch_id
	reqMap["nonce_str"] = mr.Nonce
	reqMap["result_code"] = mr.Result_code
	reqMap["openid"] = mr.Openid
	reqMap["is_subscribe"] = mr.Is_subscribe
	reqMap["trade_type"] = mr.Trade_type
	reqMap["bank_type"] = mr.Bank_type
	reqMap["total_fee"] = mr.Total_fee
	reqMap["fee_type"] = mr.Fee_type
	reqMap["cash_fee"] = mr.Cash_fee
	reqMap["cash_fee_type"] = mr.Cash_fee_Type
	reqMap["transaction_id"] = mr.Transaction_id
	reqMap["out_trade_no"] = mr.Out_trade_no
	reqMap["attach"] = mr.Attach
	reqMap["time_end"] = mr.Time_end

	if mr.Return_code == "SUCCESS" {
		//这里就可以更新我们的后台数据库了，其他业务逻辑同理。
		ctx.Update(WebPayOrderData{
			OpenId: mr.Openid, UpdateTime: gofunc.CurrentTime(), State: "已支付",
			OrderId: mr.Out_trade_no,
		})
		resp.Return_code = "SUCCESS"
		resp.Return_msg = "OK"
	}

	return resp
}
