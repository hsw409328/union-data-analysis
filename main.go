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
	relationMapAll         = make(map[string][]model.RelationData)
	userNumber,
	orderNumber,
	tmpOnceUserNumber,
	tmpOnceOrderNumber int
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

type WebUserDataSlice []model.WebUserData

// 获取此 slice 的长度
func (w WebUserDataSlice) Len() int { return len(w) }

// 根据元素的年龄降序排序 （此处按照自己的业务逻辑写）
func (w WebUserDataSlice) Less(i, j int) bool {
	return w[i].Level > w[j].Level
}

// 交换数据
func (w WebUserDataSlice) Swap(i, j int) { w[i], w[j] = w[j], w[i] }

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
		"cgtime": " between '" + lastStartTime + "' and '" + lastEndTime + "' ",
		//"cgtime": " between '2018-11-01 00:00:00' and '2018-12-01 23:59:59' ",
		//"cgtime": " between '2018-01-01 00:00:00' and '" + lastEndTime + "'  ",
		"订单状态": " = '订单结算'",
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
					if useRatio < 0 {
						useRatio = 0
					}
					ratio = webUserVal.Ratio
					//获取的奖励
					result := v.Money * useRatio
					dayRecordModel.Insert(model.WebDayRecordData{
						webUserVal.Mobile,
						gofunc.CurrentTime(),
						gofunc.CurrentTime(),
						lastDate,
						result,
						gofunc.InterfaceToString(useRatio),
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
		"cgtime": " between '" + lastStartTime + "' and '" + lastEndTime + "' ",
		//"cgtime": " between '2018-11-01 00:00:00' and '2018-12-01 23:59:59' ",
		//"cgtime": " between '2018-01-01 00:00:00' and '" + lastEndTime + "' ",
		"订单状态": " = '失效订单'",
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
						gofunc.InterfaceToString(useRatio),
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
	}).Group("手机号").GetAllSum()
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
	//取出所有代理关系表数据
	relationAll := proxyModel.GetAll()
	for _, tmp := range relationAll {
		relationMapAll[tmp.ParentId] = append(relationMapAll[tmp.ParentId], tmp)
	}
	//绑定的ID匹配代理关系表的上级ID
	for _, v := range r {
		orderNumber, userNumber = 0, 0
		//删旧的一级推荐
		webFirstRecommendModel.Delete(v.Mobile)

		//再用匹配到的子级ID接着匹配父级ID
		if tmpSlice, ok := relationMapAll[v.UnionId]; ok {
			for _, tmp := range tmpSlice {
				tmpOnceUserNumber, tmpOnceOrderNumber = 0, 0
				//统计当前人的订单 如果不为空，则算有效用户
				//统计该人的订单
				tmpOrderNumber := oModel.Where(map[string]string{
					"ID": " ='" + tmp.NextId + "' ",
				}).GetAllNoSum()
				if len(tmpOrderNumber) > 0 {
					userNumber += 1
					orderNumber += len(tmpOrderNumber)
					tmpOnceUserNumber += 1
					tmpOnceOrderNumber += len(tmpOrderNumber)
				}
				//重复 再用匹配到的子级ID接着匹配父级ID
				recursionChild(tmp.NextId)
				webFirstRecommendModel.Insert(model.WebFirstRecommendData{
					Mobile:           v.Mobile,
					RecommendName:    tmp.NextName,
					RecommendId:      tmp.NextId,
					TeamOrderNumber:  gofunc.InterfaceToString(tmpOnceOrderNumber),
					TeamPersonNumber: gofunc.InterfaceToString(tmpOnceUserNumber),
					CreateTime:       gofunc.CurrentTime(),
					UpdateTime:       gofunc.CurrentTime(),
				})
			}
		}

		//更新用户表的 用户量和订单量
		webUserModel.UpdateUserNumberAndOrderNumber(model.WebUserData{
			Mobile:           v.Mobile,
			ChildOrderNumber: gofunc.InterfaceToString(orderNumber),
			ChildUserNumber:  gofunc.InterfaceToString(userNumber),
		})
	}
}

func recursionChild(pid string) {
	if tmpSlice, ok := relationMapAll[pid]; ok {
		for _, tmp := range tmpSlice {
			//统计当前人的订单 如果不为空，则算有效用户
			//统计该人的订单
			tmpOrderNumber := oModel.Where(map[string]string{
				"ID": " ='" + tmp.NextId + "' ",
			}).GetAllNoSum()
			if len(tmpOrderNumber) > 0 {
				userNumber += 1
				orderNumber += len(tmpOrderNumber)
				tmpOnceUserNumber += 1
				tmpOnceOrderNumber += len(tmpOrderNumber)
			}
			//重复 再用匹配到的子级ID接着匹配父级ID
			recursionChild(tmp.NextId)
		}
	}

	return
}
