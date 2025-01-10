package response

type WxUserVO struct {
	Id     int64  `json:"id"`     // 用户id
	Openid string `json:"openid"` // 微信用户openid
	Token  string `json:"token"`  // jwt令牌
}
