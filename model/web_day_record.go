package model

import "union-data-analysis/lib/driver"

type WebDayRecordData struct {
	Mobile      string
	CreateTime  string
	UpdateTime  string
	RewardDate  string
	RewardMoney float32
}

type WebDayRecord struct {
	where string
	group string
}

func NewWebDayRecord() *WebDayRecord {
	return &WebDayRecord{}
}

func (ctx *WebDayRecord) Where(w map[string]string) *WebDayRecord {
	ctx.where = " where 1 "
	for k, v := range w {
		ctx.where += " and " + k + v
	}
	return ctx
}

func (ctx *WebDayRecord) Group(group string) *WebDayRecord {
	if group != "" {
		ctx.group = " group by " + group
	}
	return ctx
}

func (ctx *WebDayRecord) GetAll() []WebDayRecordData {
	r, err := driver.SQLiteDriverWeb.GetAll("select  奖励金额,手机号,奖励日期 from " + WebDayRecordTableName +
		ctx.where + ctx.group)
	if err != nil {
		lg.Error(err.Error())
	}
	defer r.Close()
	var webDayRecordData = new(WebDayRecordData)
	var webDayRecordDataSlice = make([]WebDayRecordData, 0)
	for r.Next() {
		err := r.Scan(&webDayRecordData.RewardMoney, &webDayRecordData.Mobile, &webDayRecordData.RewardDate)
		if err != nil {
			lg.Error(err.Error())
		}
		webDayRecordDataSlice = append(webDayRecordDataSlice, *webDayRecordData)
	}
	return webDayRecordDataSlice
}

func (ctx *WebDayRecord) GetAllSum() []WebDayRecordData {
	r, err := driver.SQLiteDriverWeb.GetAll("select  SUM(奖励金额),手机号,奖励日期 from " + WebDayRecordTableName +
		ctx.where + ctx.group)
	if err != nil {
		lg.Error(err.Error())
	}
	defer r.Close()
	var webDayRecordData = new(WebDayRecordData)
	var webDayRecordDataSlice = make([]WebDayRecordData, 0)
	for r.Next() {
		err := r.Scan(&webDayRecordData.RewardMoney, &webDayRecordData.Mobile, &webDayRecordData.RewardDate)
		if err != nil {
			lg.Error(err.Error())
		}
		webDayRecordDataSlice = append(webDayRecordDataSlice, *webDayRecordData)
	}
	return webDayRecordDataSlice
}

func (ctx *WebDayRecord) Insert(webDayRecordData WebDayRecordData) (int64, error) {
	n, err := driver.SQLiteDriverWeb.Insert(
		"insert into ["+WebDayRecordTableName+"]([创建时间], [更新时间], [奖励金额], [手机号], [奖励日期])"+
			" values(?, ?, ?, ?, ?)",
		webDayRecordData.CreateTime, webDayRecordData.UpdateTime, webDayRecordData.RewardMoney,
		webDayRecordData.Mobile, webDayRecordData.RewardDate,
	)
	return n, err
}
