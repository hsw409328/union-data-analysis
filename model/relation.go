package model

import "union-data-analysis/lib/driver"

type RelationData struct {
	ParentId   string
	ParentName string
	NextId     string
	NextName   string
}

type RelationUser struct {
	Id   string
	Name string
}

type Relation struct {
	where           string
	group           string
	ResultTreeSlice []RelationUser
}

func NewRelation() *Relation {
	return &Relation{}
}

func (ctx *Relation) Where(w map[string]string) *Relation {
	ctx.where = " where 1 "
	for k, v := range w {
		ctx.where += " and " + k + v
	}
	return ctx
}

func (ctx *Relation) Group(group string) *Relation {
	if group != "" {
		ctx.group = " group by " + group
	}
	return ctx
}

func (ctx *Relation) GetAll() []RelationData {
	r, err := driver.SQLiteDriverAnalysis.GetAll("select 上级ID, 上级名称, 下级ID, 下级名称 from " + ProxyTableName +
		ctx.where)
	if err != nil {
		lg.Error(err.Error())
	}
	defer r.Close()
	var relationData = new(RelationData)
	var relationDataSlice = make([]RelationData, 0)
	for r.Next() {
		err := r.Scan(&relationData.ParentId, &relationData.ParentName, &relationData.NextId, &relationData.NextName)
		if err != nil {
			lg.Error(err.Error())
		}
		relationDataSlice = append(relationDataSlice, *relationData)
	}
	return relationDataSlice
}

func (ctx *Relation) GetAllHash() (result map[string]RelationUser) {
	result = map[string]RelationUser{}
	r, err := driver.SQLiteDriverAnalysis.GetAll("select 上级ID, 上级名称, 下级ID, 下级名称 from " + ProxyTableName +
		ctx.where)
	if err != nil {
		lg.Error(err.Error())
	}
	defer r.Close()
	var relationData = new(RelationData)
	for r.Next() {
		err := r.Scan(&relationData.ParentId, &relationData.ParentName, &relationData.NextId, &relationData.NextName)
		if err != nil {
			lg.Error(err.Error())
		}
		result[relationData.NextId] = RelationUser{relationData.ParentId, relationData.ParentName}
	}
	return
}

func (ctx *Relation) GetOne() RelationData {
	r := driver.SQLiteDriverAnalysis.GetOne("select 用户ID, 用户名, 上级用户ID from " + WebUserTableName +
		ctx.where + ctx.group)
	var relationData = new(RelationData)
	err := r.Scan(&relationData.ParentId, &relationData.ParentName, &relationData.NextId, &relationData.NextName)
	if err != nil {
		lg.Error(err.Error())
	}
	return *relationData
}

func (ctx *Relation) Tree(id string, data map[string]RelationUser) {
	if _, ok := data[id]; ok {
		tmp := data[id]
		ctx.ResultTreeSlice = append(ctx.ResultTreeSlice, tmp)
		ctx.Tree(tmp.Id, data)
	}
}
