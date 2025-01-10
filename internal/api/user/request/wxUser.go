package request

type WxUserLoginDTO struct {
	Code string `json:"code" binding:"required` // 微信授权码
}
