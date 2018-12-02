package main

import (
	"github.com/hsw409328/gofunc"
	"github.com/hsw409328/gofunc/go_hlog"
	"union-data-analysis/model"
)

var (
	lg *go_hlog.Logger
)

func init() {
	//lg = go_hlog.GetInstance("./logs/system/" + gofunc.CurrentDate()+".log")
	lg = go_hlog.GetInstance("")
}

type Achievement struct {
	Id       string               `用户ID`
	Money    float32              `业绩`
	DataTree []model.RelationUser `层级流`
}

func main() {
	lg.Info(" system start ")
	add()
	reduce()
}

func add() {

	//取出用户关系以 【下级ID】 hash形式存放
	relationModel := model.NewRelation()
	relationDataMap := relationModel.GetAllHash()

	// 计算每天的所有订单业绩
	orderDataSlice := model.NewOrder().Where(map[string]string{
		"cgtime": " between '2018-11-01 00:00:00' and '2018-12-01 23:59:59' ",
		"订单状态":   " = '订单结算'",
	}).Group("ID").GetAll()
	//每笔订单的业绩计算方法：
	//业绩=自己的佣金 -（联盟佣金*12%）
	var achievementList = make([]Achievement, 0)
	for _, v := range orderDataSlice {
		//自己的业绩
		achievementMoney := v.SelfMoney - (v.UnionMoney * 0.12)
		//找出自己的层级关系
		relationModel.ResultTreeSlice = []model.RelationUser{}
		relationModel.Tree(v.Id, relationDataMap)
		if len(relationModel.ResultTreeSlice) > 0 {
			achievementList = append(achievementList, Achievement{v.Id, achievementMoney, relationModel.ResultTreeSlice})
			continue
		}
		achievementList = append(achievementList, Achievement{Id: v.Id, Money: achievementMoney})
	}

	//遍历每一个订单的层级关系，与web领导员用户ID做对应关系。查询出获取的佣金
	webUserModel := new(model.WebUsers)
	dayRecordModel := new(model.WebDayRecord)
	for _, v := range achievementList {
		if v.DataTree != nil {
			//临时存储一条流查询到的用户等级和比例
			tmp := []model.WebUserData{}
			for _, childVal := range v.DataTree {
				w, err := webUserModel.Where(map[string]string{
					"用户ID": " = '" + childVal.Id + "'",
				}).GetOne()
				if err != nil {
					continue
				}
				tmp = append(tmp, w)
			}
			//按照比例等级进行折算
			if tmp != nil {
				var ratio float32 = 0
				var useRatio float32 = 0
				for _, webUserVal := range tmp {
					useRatio = webUserVal.Ratio - ratio
					ratio = webUserVal.Ratio
					//获取的奖励
					result := v.Money * useRatio
					lg.Asset(result)
					dayRecordModel.Insert(model.WebDayRecordData{
						webUserVal.Mobile,
						gofunc.CurrentTime(),
						gofunc.CurrentTime(),
						gofunc.CurrentDate(),
						result,
					})
				}
			}
		}
	}
}

func reduce() {

	//取出用户关系以 【下级ID】 hash形式存放
	relationModel := model.NewRelation()
	relationDataMap := relationModel.GetAllHash()

	// 计算每天的所有订单业绩
	orderDataSlice := model.NewOrder().Where(map[string]string{
		"cgtime": " between '2018-11-01 00:00:00' and '2018-12-01 23:59:59' ",
		"订单状态":   " = '失效订单'",
	}).Group("ID").GetAll()
	//每笔订单的业绩计算方法：
	//业绩=自己的佣金 -（联盟佣金*12%）
	var achievementList = make([]Achievement, 0)
	for _, v := range orderDataSlice {
		//自己的业绩
		achievementMoney := v.SelfMoney - (v.UnionMoney * 0.12)
		//找出自己的层级关系
		relationModel.ResultTreeSlice = []model.RelationUser{}
		relationModel.Tree(v.Id, relationDataMap)
		if len(relationModel.ResultTreeSlice) > 0 {
			achievementList = append(achievementList, Achievement{v.Id, achievementMoney, relationModel.ResultTreeSlice})
			continue
		}
		achievementList = append(achievementList, Achievement{Id: v.Id, Money: achievementMoney})
	}

	//遍历每一个订单的层级关系，与web领导员用户ID做对应关系。查询出获取的佣金
	webUserModel := new(model.WebUsers)
	dayRecordModel := new(model.WebDayRecord)
	for _, v := range achievementList {
		if v.DataTree != nil {
			//临时存储一条流查询到的用户等级和比例
			tmp := []model.WebUserData{}
			for _, childVal := range v.DataTree {
				w, err := webUserModel.Where(map[string]string{
					"用户ID": " = '" + childVal.Id + "'",
				}).GetOne()
				if err != nil {
					continue
				}
				tmp = append(tmp, w)
			}
			//按照比例等级进行折算
			if tmp != nil {
				var ratio float32 = 0
				var useRatio float32 = 0
				for _, webUserVal := range tmp {
					useRatio = webUserVal.Ratio - ratio
					ratio = webUserVal.Ratio
					//获取的奖励
					result := v.Money * useRatio
					lg.Asset(result)
					dayRecordModel.Insert(model.WebDayRecordData{
						webUserVal.Mobile,
						gofunc.CurrentTime(),
						gofunc.CurrentTime(),
						gofunc.CurrentDate(),
						-result,
					})
				}
			}
		}
	}
}
