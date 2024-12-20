package common

type Result struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type PageResult struct {
	Total   int64       `josn:total`   // 总记录数
	Records interface{} `josn:records` // 当前页数据集合
}
