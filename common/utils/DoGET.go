package utils

import (
	"io"
	"net/http"
	"net/url"
)

// DoGET 发送GET方式请求
func DoGET(targetUrl string, values map[string]string) (string, error) {
	var (
		err    error
		resp   *http.Response
		urlVal = url.Values{} // query参数
		body   []byte
	)
	// 将字符串url解析为URL结构
	u, _ := url.ParseRequestURI(targetUrl)
	// 遍历传入的query参数
	for k, v := range values {
		urlVal.Set(k, v)
	}
	// 为URL添加query参数
	u.RawQuery = urlVal.Encode()
	// 发送GET请求，获取响应
	resp, err = http.Get(u.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// 读取响应体
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// 返回响应数据
	return string(body), nil
}
