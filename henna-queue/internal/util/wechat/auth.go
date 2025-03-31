package wechat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

// Session 微信会话信息
type Session struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
}

// Code2Session 使用code获取session
func Code2Session(code string) (*Session, error) {
	appID := viper.GetString("wechat.app_id")
	appSecret := viper.GetString("wechat.app_secret")
	
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appID, appSecret, code)
	
	// 发送请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// 解析响应
	var result struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionid"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("WeChat API error: %d %s", result.ErrCode, result.ErrMsg)
	}
	
	return &Session{
		OpenID:     result.OpenID,
		SessionKey: result.SessionKey,
		UnionID:    result.UnionID,
	}, nil
}

// GetAccessToken 获取接口调用凭据
func GetAccessToken() (string, error) {
	appID := viper.GetString("wechat.app_id")
	appSecret := viper.GetString("wechat.app_secret")
	
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		appID, appSecret)
	
	// 发送请求
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	// 解析响应
	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	
	if result.ErrCode != 0 {
		return "", fmt.Errorf("WeChat API error: %d %s", result.ErrCode, result.ErrMsg)
	}
	
	return result.AccessToken, nil
} 