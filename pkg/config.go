package biligo

import "github.com/gorilla/websocket"

const biliHost = "https://live-open.biliapi.com"

// 配置信息
type Config struct {
	code      string          // 身份码
	connData  *Data           // 主播信息
	wsConn    *websocket.Conn // ws 信息
	isReady   bool            // 是否准备好
	appId     int             // 应用id
	appkeyId  string          // 账号id
	appSecret string          // 账号密钥
}

// 获取主播信息
func GetAnchorInfo() AnchorInfo {
	return cfg.connData.AnchorInfo
}

// 获取身份码
func GetCode() string {
	return cfg.code
}

// 是否准备好
func IsReady() bool {
	if cfg == nil {
		return false
	}
	return cfg.isReady
}
