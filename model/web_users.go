package model

import "union-data-analysis/lib/driver"

type WebUserData struct {
	Id               int
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
	r, err := driver.SQLiteDriverWeb.GetAll("select 用户ID, 用户名, 上级用户ID from " + WebUserTableName +
		ctx.where + ctx.group)
	if err != nil {
		lg.Error(err.Error())
	}
	defer r.Close()
	var webUsersData = new(WebUserData)
	var webUsersDataSlice = make([]WebUserData, 0)
	for r.Next() {
		err := r.Scan(&webUsersData.UnionId, &webUsersData.UnionName, &webUsersData.UnionParentId)
		if err != nil {
			lg.Error(err.Error())
		}
		webUsersDataSlice = append(webUsersDataSlice, *webUsersData)
	}
	return webUsersDataSlice
}

func (ctx *WebUsers) GetOne() (WebUserData, error) {
	r := driver.SQLiteDriverWeb.GetOne("select 用户ID, 用户名, 上级用户ID,等级,比例,手机号 from " + WebUserTableName +
		ctx.where + ctx.group)
	var webUsersData = new(WebUserData)
	err := r.Scan(&webUsersData.UnionId, &webUsersData.UnionName, &webUsersData.UnionParentId, &webUsersData.Level,
		&webUsersData.Ratio, &webUsersData.Mobile)
	if err != nil {
		lg.Error(err.Error())
		return *webUsersData, err
	}
	return *webUsersData, nil
}
