package response

// TurnoverReportVO 营业额统计返回数据模型
type TurnoverReportVO struct {
	DateList     string `json:"dateList"`
	TurnoverList string `json:"turnoverList"`
}

// UserReportVO 用户统计返回数据模型
type UserReportVO struct {
	DateList      string `json:"dateList"`
	TotalUserList string `json:"totalUserList"`
	NewUserList   string `json:"newUserList"`
}

// OrderReportVO 订单统计返回数据模型
type OrderReportVO struct {
	DateList            string  `json:"dateList"`
	OrderCountList      string  `json:"orderCountList"`
	ValidOrderCountList string  `json:"validOrderCountList"`
	TotalOrderCount     int     `json:"totalOrderCount"`
	ValidOrderCount     int     `json:"validOrderCount"`
	OrderCompletionRate float64 `json:"orderCompletionRate"`
}

// SalesTop10ReportVO 销量top10返回数据模型
type SalesTop10ReportVO struct {
	NameList   string `json:"nameList"`
	NumberList string `json:"numberList"`
}
