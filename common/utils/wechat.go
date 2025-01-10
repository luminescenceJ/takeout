package utils

import (
	"encoding/json"
	"takeout/global"
)

// GetOpenID 调用微信接口服务，获取微信用户的openid
func GetOpenID(code string) string {
	var (
		resultJSON string
		err        error
	)
	const WxLogin = "https://api.weixin.qq.com/sns/jscode2session"
	// 准备微信接口服务参数
	values := map[string]string{
		"appid":      global.Config.Wechat.AppId,
		"secret":     global.Config.Wechat.AppSecret,
		"js_code":    code,
		"grant_type": "authorization_code",
	}
	// HTTP调用微信接口获取用户数据
	if resultJSON, err = DoGET(WxLogin, values); err != nil {
		global.Log.Error("GetOpenID 调用微信接口服务失败：", err.Error())
		return ""
	}
	wxLoginResponse := struct {
		SessionKey string `json:"session_key"`
		OpenID     string `json:"openid"`
	}{}
	if err = json.Unmarshal([]byte(resultJSON), &wxLoginResponse); err != nil {
		global.Log.Error("wxLoginResponse 不匹配：", err.Error())
		return ""
	}
	return wxLoginResponse.OpenID
}
