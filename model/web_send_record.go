package model

import "union-data-analysis/lib/driver"

type WebSendRecordData struct {
	RowId          int
	CreateTime     string
	UpdateTime     string
	RewardMoney    float32
	RewardMonth    string
	RewardState    string `未结算，待结算，已结算`
	RewardSendTime string
	Mobile         string
}

type WebSendRecord struct {
	where string
	group string
}

func NewWebSendRecord() *WebSendRecord {
	return &WebSendRecord{}
}

func (ctx *WebSendRecord) Where(w map[string]string) *WebSendRecord {
	ctx.where = " where 1 "
	for k, v := range w {
		ctx.where += " and " + k + v
	}
	return ctx
}

func (ctx *WebSendRecord) Group(group string) *WebSendRecord {
	if group != "" {
		ctx.group = " group by " + group
	}
	return ctx
}

func (ctx *WebSendRecord) Insert(webSendRecordData WebSendRecordData) (int64, error) {
	n, err := driver.SQLiteDriverWeb.Insert(
		"insert into ["+WebSendRecordTableName+"]([创建时间], [更新时间],[奖励金额],[月份],[发放状态],[发放时间],[手机号])"+
			" values(?, ?, ?, ?, ?, ?, ?)",
		webSendRecordData.CreateTime, webSendRecordData.UpdateTime, webSendRecordData.RewardMoney,
		webSendRecordData.RewardMonth, webSendRecordData.RewardState, webSendRecordData.RewardSendTime,
		webSendRecordData.Mobile,
	)
	return n, err
}

func (ctx *WebSendRecord) UpdateRewardMoney(webSendRecordData WebSendRecordData) (int64, error) {
	n, err := driver.SQLiteDriverWeb.Update("UPDATE "+WebSendRecordTableName+" SET 奖励金额 = ?,更新时间 = ?"+
		" WHERE rowid = ?",
		webSendRecordData.RewardMoney, webSendRecordData.UpdateTime, webSendRecordData.RowId)
	return n, err
}

func (ctx *WebSendRecord) UpdateRewarSendTimeAndState(id, webSendRecordData WebSendRecordData) (int64, error) {
	n, err := driver.SQLiteDriverWeb.Update("UPDATE "+WebSendRecordTableName+" SET 更新时间 = ?,"+
		"发放状态 = ?,发放时间 = ? WHERE rowid = ?",
		webSendRecordData.UpdateTime, webSendRecordData.RewardState, webSendRecordData.RewardSendTime)
	return n, err
}

func (ctx *WebSendRecord) GetMobileLastRecord() (WebSendRecordData, error) {
	//取出最后一条记录
	r := driver.SQLiteDriverWeb.GetOne("select rowid,* from " + WebSendRecordTableName +
		ctx.where + " order by 更新时间 desc " + ctx.group)
	var webSendRecordData = new(WebSendRecordData)
	err := r.Scan(&webSendRecordData.RowId, &webSendRecordData.CreateTime, &webSendRecordData.UpdateTime,
		&webSendRecordData.RewardMoney, &webSendRecordData.RewardMonth, &webSendRecordData.RewardState,
		&webSendRecordData.RewardSendTime, &webSendRecordData.Mobile)
	if err != nil {
		lg.Error(err.Error())
		return *webSendRecordData, err
	}
	return *webSendRecordData, nil
}
