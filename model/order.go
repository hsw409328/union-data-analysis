package model

import (
	"union-data-analysis/lib/driver"
)

type OrderData struct {
	Id          string
	Name        string
	UnionMoney  float32
	SelfMoney   float32
	ParentMoney float32
	FirstMoney  float32
}

type Order struct {
	where string
	group string
}

func NewOrder() *Order {
	return &Order{}
}

func (ctx *Order) Where(w map[string]string) *Order {
	ctx.where = " where 1 "
	for k, v := range w {
		ctx.where += " and " + k + v
	}
	return ctx
}

func (ctx *Order) Group(group string) *Order {
	if group != "" {
		ctx.group = " group by " + group
	}
	return ctx
}

func (ctx *Order) GetAll() []OrderData {
	r, err := driver.SQLiteDriverAnalysis.GetAll("select ID,名称,sum(联盟佣金),sum(自己佣金),sum(上级佣金),sum(一级佣金) from " + OrderTableName +
		ctx.where + ctx.group)
	if err != nil {
		lg.Error(err.Error())
	}
	defer r.Close()
	var orderData = new(OrderData)
	var orderDataSlice = make([]OrderData, 0)
	for r.Next() {
		err := r.Scan(&orderData.Id, &orderData.Name, &orderData.UnionMoney, &orderData.SelfMoney, &orderData.ParentMoney, &orderData.FirstMoney)
		if err != nil {
			lg.Error(err.Error())
		}
		orderDataSlice = append(orderDataSlice, *orderData)
	}
	return orderDataSlice
}

func (ctx *Order) GetAllNoSum() []OrderData {
	r, err := driver.SQLiteDriverAnalysis.GetAll("select ID,名称,联盟佣金,自己佣金,上级佣金,一级佣金 from " + OrderTableName +
		ctx.where + ctx.group)
	if err != nil {
		lg.Error(err.Error())
	}
	defer r.Close()
	var orderData = new(OrderData)
	var orderDataSlice = make([]OrderData, 0)
	for r.Next() {
		err := r.Scan(&orderData.Id, &orderData.Name, &orderData.UnionMoney, &orderData.SelfMoney, &orderData.ParentMoney, &orderData.FirstMoney)
		if err != nil {
			lg.Error(err.Error())
		}
		orderDataSlice = append(orderDataSlice, *orderData)
	}
	return orderDataSlice
}

func (ctx *Order) GetOne(id string) OrderData {
	r := driver.SQLiteDriverAnalysis.GetOne("select  ID,名称,联盟佣金,自己佣金,上级佣金,一级佣金 from from " + OrderTableName +
		" where ID='" + id + "' limit 1 ")
	var orderData = new(OrderData)
	err := r.Scan(&orderData.Id, &orderData.Name, &orderData.UnionMoney, &orderData.SelfMoney, &orderData.ParentMoney, &orderData.FirstMoney)
	if err != nil {
		lg.Error(err.Error())
	}
	return *orderData
}
