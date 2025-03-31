package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

// SendSubscribeMessage 发送订阅消息
func SendSubscribeMessage(openID, title, content string) error {
	// 获取access_token
	token, err := GetAccessToken()
	if err != nil {
		return err
	}
	
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s", token)
	templateID := viper.GetString("wechat.template_id")
	
	// 构建消息数据
	data := map[string]interface{}{
		"touser": openID,
		"template_id": templateID,
		"page": "pages/queue/queue",
		"data": map[string]interface{}{
			"thing1": map[string]string{
				"value": title,
			},
			"thing2": map[string]string{
				"value": content,
			},
			"time3": map[string]string{
				"value": "马上",
			},
		},
	}
	
	// 序列化数据
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	// 发送请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	// 解析响应
	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	
	if result.ErrCode != 0 {
		return fmt.Errorf("WeChat API error: %d %s", result.ErrCode, result.ErrMsg)
	}
	
	return nil
} 