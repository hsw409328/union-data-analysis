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

func (ctx *WebDayRecord) Insert(webDayRecordData WebDayRecordData) (int64, error) {
	n, err := driver.SQLiteDriverWeb.Insert(
		"insert into [每日记录表]([创建时间], [更新时间], [奖励金额], [手机号], [奖励日期]) values(?, ?, ?, ?, ?)",
		webDayRecordData.CreateTime, webDayRecordData.UpdateTime, webDayRecordData.RewardMoney,
		webDayRecordData.Mobile, webDayRecordData.RewardDate,
	)
	return n, err
}
