package main

import (
	"flag"
	"github.com/hsw409328/gofunc"
	"github.com/hsw409328/gofunc/go_hlog"
	"time"
	"union-data-analysis/lib/driver"
	"union-data-analysis/model"
)

var (
	lg                     *go_hlog.Logger
	webUserModel           = new(model.WebUsers)
	proxyModel             = new(model.Relation)
	webFirstRecommendModel = new(model.WebFirstRecommend)
	oModel                 = new(model.Order)
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
	lastDate := gofunc.TimeUnixIntToStringCustom(gofunc.LastTime("d", -1), "2006-01-02")
	lastStartTime := lastDate + " 00:00:00"
	lastEndTime := lastDate + " 23:59:59"
	//第一个参数，为参数名称，第二个参数为默认值，第三个参数是说明
	dbName := flag.String("name", "", "输入要分析的数据库路径名称，如：./db/2018-12-02-003925.db")
	flag.Parse()
	if *dbName != "" {
		driver.SQLiteDriverAnalysis = driver.NewSQLite(*dbName)
	}
	//统计每天增加的金额
	add(lastDate, lastStartTime, lastEndTime)
	//统计每天减少的金额
	reduce(lastDate, lastStartTime, lastEndTime)
	//更新发放记录数据
	updateSendRecord(lastDate)
	//计算每个WEB用户的一级推荐人
	webUserFirstRecommend()
}

func add(lastDate string, lastStartTime string, lastEndTime string) {

	//取出用户关系以 【下级ID】 hash形式存放
	relationModel := model.NewRelation()
	relationDataMap := relationModel.GetAllHash()

	// 计算每天的所有订单业绩
	orderDataSlice := model.NewOrder().Where(map[string]string{
		//"cgtime": " between '" + lastStartTime + "' and '" + lastEndTime + "' ",
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
					dayRecordModel.Insert(model.WebDayRecordData{
						webUserVal.Mobile,
						gofunc.CurrentTime(),
						gofunc.CurrentTime(),
						lastDate,
						result,
					})
				}
			}
		}
	}
}

func reduce(lastDate string, lastStartTime string, lastEndTime string) {

	//取出用户关系以 【下级ID】 hash形式存放
	relationModel := model.NewRelation()
	relationDataMap := relationModel.GetAllHash()

	// 计算每天的所有订单业绩
	orderDataSlice := model.NewOrder().Where(map[string]string{
		//"cgtime": " between '" + lastStartTime + "' and '" + lastEndTime + "' ",
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
					dayRecordModel.Insert(model.WebDayRecordData{
						webUserVal.Mobile,
						gofunc.CurrentTime(),
						gofunc.CurrentTime(),
						lastDate,
						-result,
					})
				}
			}
		}
	}
}

func updateSendRecord(lastDate string) {
	//获取昨天汇总之后所有的记录表
	webDayRecordModel := new(model.WebDayRecord)
	r := webDayRecordModel.Where(map[string]string{
		"奖励日期": " = '" + lastDate + "'",
	}).GetAllSum()
	for _, v := range r {
		//获取未结算的单子
		webSendRecordModel := new(model.WebSendRecord)
		tmp, _ := webSendRecordModel.Where(map[string]string{
			"发放状态": " ='未结算'",
			"年份":   " ='" + time.Now().Format("2006") + "'",
			"月份":   " ='" + time.Now().Format("01") + "'",
			"手机号":  "= '" + v.Mobile + "' ",
		}).GetMobileLastRecord()
		if tmp.RowId == 0 {
			//执行插入操作
			_, err := webSendRecordModel.Insert(model.WebSendRecordData{
				0,
				gofunc.CurrentTime(),
				gofunc.CurrentTime(),
				v.RewardMoney,
				time.Now().Format("01"),
				"未结算",
				"",
				v.Mobile,
				time.Now().Format("2006"),
			})
			if err != nil {
				lg.Error(err.Error())
				continue
			}
		} else {
			//执行更新操作 如果当前执行的月份大于旧月份，将原数据状态更新为“待结算”
			n := time.Now().Format("200601")
			old := tmp.RewardYear + tmp.RewardMonth
			if n > old {
				tmp.RewardState = "待结算"
			}
			tmp.RewardMoney += v.RewardMoney
			tmp.UpdateTime = gofunc.CurrentTime()
			webSendRecordModel.UpdateRewarSendTimeAndState(tmp)
			_, err := webSendRecordModel.UpdateRewardMoney(tmp)
			if err != nil {
				lg.Error(err.Error())
				continue
			}
		}
	}
}

func webUserFirstRecommend() {
	//读取web用户表，取出绑定的ID
	r := webUserModel.GetAll()
	//去代理表，通过上级ID，查出自己的下级
	for _, v := range r {
		//用户量和订单量
		var userNumber, orderNumber int
		tmp := proxyModel.Where(map[string]string{
			"上级ID": " = '" + v.UnionId + "'",
		}).GetAll()
		userNumber += len(tmp)
		//删旧的一级推荐
		webFirstRecommendModel.Delete(v.Mobile)
		//将新数据放入一级推荐表
		for _, vv := range tmp {
			//统计该人的订单
			o := oModel.Where(map[string]string{
				"ID": " ='" + vv.NextId + "' ",
			}).GetAll()
			orderNumber += len(o)
			//统计该人的一级推荐
			child := proxyModel.Where(map[string]string{
				"上级ID": " = '" + vv.NextId + "'",
			}).GetAll()
			webFirstRecommendModel.Insert(model.WebFirstRecommendData{
				Mobile:           v.Mobile,
				RecommendName:    vv.NextName,
				RecommendId:      vv.NextId,
				TeamOrderNumber:  gofunc.InterfaceToString(len(o)),
				TeamPersonNumber: gofunc.InterfaceToString(len(child)),
				CreateTime:       gofunc.CurrentTime(),
				UpdateTime:       gofunc.CurrentTime(),
			})
		}
		//更新用户表的 用户量和订单量
		webUserModel.UpdateUserNumberAndOrderNumber(model.WebUserData{
			Mobile:           v.Mobile,
			ChildOrderNumber: gofunc.InterfaceToString(orderNumber),
			ChildUserNumber:  gofunc.InterfaceToString(userNumber),
		})
	}
}
