package biligo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// 项目配置信息
var cfg *Config

// 项目启动
func Start(code string, appId int, appSecret string, appkeyId string, eventBus *Bus) {
	if eventBus == nil {
		biliLog.Fatalf("eventBus 未初始化")
	}
	if cfg == nil {
		// 初始化配置
		cfg = &Config{
			code:      code,      // 身份码
			appId:     appId,     // 应用id
			appkeyId:  appkeyId,  // 账号id
			appSecret: appSecret, // 账号密钥
		}
		// 启动客户端
		cfg.startClient()
	}
}

// 项目关闭
func Shutdown() {
	if cfg != nil && cfg.isReady {
		// 关闭链接
		cfg.endClient()
	}
}

// 开启客户端
func (a *Config) startClient() error {
	if a.isReady {
		return errors.New("请勿重复启动")
	}

	if a.code == "" {
		panic("身份码为空")
	}

	// 应用启动请求数据
	data := map[string]interface{}{
		"code":   a.code,
		"app_id": a.appId,
	}
	body, _ := json.Marshal(data)

	client := &http.Client{}

	// 构建请求
	req, _ := http.NewRequest("POST", biliHost+"/v2/app/start", bytes.NewReader(body))
	// 构建请求头参数
	h := buildBiliHeader(string(body))
	// 设置请求头
	h.Apply(req)

	// 发送连接请求
	resp, _ := client.Do(req)
	var res ApiResponse
	resBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(resBody, &res)
	a.connData = &res.Data
	// 认证成功
	if res.Code == 0 {
		a.connData = &res.Data
		// 主播头像转 base64
		a.connData.AnchorInfo.Uface = ImageToBase64(a.connData.AnchorInfo.Uface)
		// 连接 ws
		a.connectWS()
	} else {
		msg := getBiliErrorMessage(res.Code)
		// 失败提示
		return errors.New(msg)
	}
	return nil
}

// 关闭客户端
func (a *Config) endClient() {
	if !a.isReady || a.code == "" || a.connData.GameInfo.GameID == "" {
		return
	}
	data := map[string]interface{}{
		"app_id":  a.appId,
		"game_id": a.connData.GameInfo.GameID,
	}
	body, _ := json.Marshal(data)

	client := &http.Client{}

	// 构建请求
	req, _ := http.NewRequest("POST", biliHost+"/v2/app/end", bytes.NewReader(body))
	h := buildBiliHeader(string(body))
	// 设置请求头
	h.Apply(req)
	// 发送请求
	resp, _ := client.Do(req)
	var res ApiResponse
	resBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(resBody, &res)
	if res.Code == 0 {
		a.isReady = false
		biliLog.Println("应用关闭 code", res.Code, "message", res.Message)
	}
	a.connData.GameInfo.GameID = ""
}
