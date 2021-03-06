package model

import (
	"log"
	"union-data-analysis/lib/driver"
)

type WebUserData struct {
	Mobile           string
	State            string
	BankUserAccount  string
	BankUserName     string
	BankName         string
	ChildUserNumber  string
	ChildOrderNumber string
	Level            string
	Ratio            float32
	CreateTime       string
	UpdateTime       string
	UnionId          string
	UnionName        string
	UnionParentId    string
	UserNickName     string
}

type WebUsers struct {
	where string
	group string
}

func NewWebUsers() *WebUsers {
	return &WebUsers{}
}

func (ctx *WebUsers) Where(w map[string]string) *WebUsers {
	ctx.where = " where 1 "
	for k, v := range w {
		ctx.where += " and " + k + v
	}
	return ctx
}

func (ctx *WebUsers) Group(group string) *WebUsers {
	if group != "" {
		ctx.group = " group by " + group
	}
	return ctx
}

func (ctx *WebUsers) GetAll() []WebUserData {
	r, err := driver.SQLiteDriverWeb.GetAll("select 用户ID, 用户名, 上级用户ID,手机号 from " + WebUserTableName +
		ctx.where + ctx.group)
	if err != nil {
		lg.Error(err.Error())
	}
	defer r.Close()
	var webUsersData = new(WebUserData)
	var webUsersDataSlice = make([]WebUserData, 0)
	for r.Next() {
		err := r.Scan(&webUsersData.UnionId, &webUsersData.UnionName, &webUsersData.UnionParentId, &webUsersData.Mobile)
		if err != nil {
			lg.Error(err.Error())
		}
		webUsersDataSlice = append(webUsersDataSlice, *webUsersData)
	}
	return webUsersDataSlice
}

func (ctx *WebUsers) GetOne() (WebUserData, error) {
	r := driver.SQLiteDriverWeb.GetOne("select 用户ID, 用户名, 上级用户ID,等级,比例,手机号,状态 from " + WebUserTableName +
		ctx.where + ctx.group)
	var webUsersData = new(WebUserData)
	err := r.Scan(&webUsersData.UnionId, &webUsersData.UnionName, &webUsersData.UnionParentId, &webUsersData.Level,
		&webUsersData.Ratio, &webUsersData.Mobile, &webUsersData.State)
	if err != nil {
		lg.Error(err.Error())
		return *webUsersData, err
	}
	return *webUsersData, nil
}

func (ctx *WebUsers) GetOneAllField() (WebUserData, error) {
	r := driver.SQLiteDriverWeb.GetOne("select * from " + WebUserTableName +
		ctx.where + ctx.group)
	var webUsersData = new(WebUserData)
	err := r.Scan(&webUsersData.Mobile,
		&webUsersData.State,
		&webUsersData.BankUserAccount,
		&webUsersData.BankUserName,
		&webUsersData.BankName,
		&webUsersData.ChildUserNumber,
		&webUsersData.ChildOrderNumber,
		&webUsersData.Level,
		&webUsersData.Ratio,
		&webUsersData.CreateTime,
		&webUsersData.UpdateTime,
		&webUsersData.UnionId,
		&webUsersData.UnionName,
		&webUsersData.UnionParentId,
		&webUsersData.UserNickName)
	if err != nil {
		lg.Error(err.Error())
		return *webUsersData, err
	}
	return *webUsersData, nil
}

func (ctx *WebUsers) Insert(w WebUserData) (int64, error) {
	n, err := driver.SQLiteDriverWeb.Insert(
		"insert into ["+WebUserTableName+"]([手机号], [状态],[创建时间],[更新时间])"+
			" values(?, ?, ?, ?)",
		w.Mobile, w.State, w.CreateTime, w.UpdateTime,
	)
	return n, err
}

func (ctx *WebUsers) Update(w WebUserData) (int64, error) {
	n, err := driver.SQLiteDriverWeb.Update("UPDATE "+WebUserTableName+" SET 开户账号 = ?,"+
		"开户名 = ?,开户行 = ? WHERE 手机号 = ?",
		w.BankUserAccount, w.BankUserName, w.BankName,
		w.Mobile)
	if err != nil {
		log.Println(err)
	}
	return n, err
}

func (ctx *WebUsers) UpdateUserNumberAndOrderNumber(w WebUserData) (int64, error) {
	n, err := driver.SQLiteDriverWeb.Update("UPDATE "+WebUserTableName+" SET 用户量 = ?,"+
		"订单量 = ? WHERE 手机号 = ?",
		w.ChildUserNumber, w.ChildOrderNumber,
		w.Mobile)
	if err != nil {
		log.Println(err)
	}
	return n, err
}
