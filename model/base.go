package model

import "github.com/hsw409328/gofunc/go_hlog"

var (
	lg *go_hlog.Logger
)

const (
	OrderTableName = "订单管理"
	ProxyTableName = "代理管理"

	WebUserTableName       = "用户表"
	WebPayOrderTableName   = "用户支付表"
	WebDayRecordTableName  = "每日记录表"
	WebSendRecordTableName = "发放记录表"
	WebRecommendTableName  = "一级推荐表"
)

func init() {
	lg = go_hlog.GetInstance("")
}
