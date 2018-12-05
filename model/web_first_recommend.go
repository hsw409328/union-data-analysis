package model

import (
	"union-data-analysis/lib/driver"
)

type WebFirstRecommendData struct {
	Mobile           string
	RecommendName    string
	TeamPersonNumber string
	TeamOrderNumber  string
	CreateTime       string
	UpdateTime       string
	RecommendId      string
}

type WebFirstRecommend struct {
	where string
	group string
}

func NewWebFirstRecommend() *WebFirstRecommend {
	return &WebFirstRecommend{}
}

func (ctx *WebFirstRecommend) Where(w map[string]string) *WebFirstRecommend {
	ctx.where = " where 1 "
	for k, v := range w {
		ctx.where += " and " + k + v
	}
	return ctx
}

func (ctx *WebFirstRecommend) Group(group string) *WebFirstRecommend {
	if group != "" {
		ctx.group = " group by " + group
	}
	return ctx
}

func (ctx *WebFirstRecommend) GetAll() []WebFirstRecommendData {
	r, err := driver.SQLiteDriverWeb.GetAll("select * from " + WebRecommendTableName +
		ctx.where + ctx.group)
	if err != nil {
		lg.Error(err.Error())
	}
	defer r.Close()
	var webFirstRecommendData = new(WebFirstRecommendData)
	var webFirstRecommendDataSlice = make([]WebFirstRecommendData, 0)
	for r.Next() {
		err := r.Scan(&webFirstRecommendData.Mobile,
			&webFirstRecommendData.RecommendName,
			&webFirstRecommendData.TeamPersonNumber,
			&webFirstRecommendData.TeamOrderNumber,
			&webFirstRecommendData.CreateTime,
			&webFirstRecommendData.UpdateTime,
			&webFirstRecommendData.RecommendId)
		if err != nil {
			lg.Error(err.Error())
		}
		webFirstRecommendDataSlice = append(webFirstRecommendDataSlice, *webFirstRecommendData)
	}
	return webFirstRecommendDataSlice
}

func (ctx *WebFirstRecommend) Insert(w WebFirstRecommendData) (int64, error) {
	n, err := driver.SQLiteDriverWeb.Insert(
		"insert into ["+WebRecommendTableName+"]([手机号], [直接推荐用户名],[团队人数],[团队订单],"+
			"[创建时间],[更新时间],[直接推荐用户ID])"+
			" values(?, ?, ?, ?, ?, ?, ?)",
		w.Mobile, w.RecommendName, w.TeamPersonNumber, w.TeamOrderNumber, w.CreateTime, w.UpdateTime, w.RecommendId,
	)
	return n, err
}

func (ctx *WebFirstRecommend) Delete(mobileName string) {
	driver.SQLiteDriverWeb.Delete("DELETE FROM "+WebRecommendTableName+" WHERE 手机号= ?", mobileName)
	return
}
