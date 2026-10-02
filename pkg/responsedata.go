package biligo

/**
 * 响应数据 结构体
 */
type ApiResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Data      Data   `json:"data"`
}

type Data struct {
	AnchorInfo    AnchorInfo    `json:"anchor_info"`
	GameInfo      GameInfo      `json:"game_info"`
	WebsocketInfo WebsocketInfo `json:"websocket_info"`
}

// 主播信息
type AnchorInfo struct {
	OpenID  string `json:"open_id"`
	RoomID  int64  `json:"room_id"`
	Uface   string `json:"uface"`
	UID     int64  `json:"uid"`
	Uname   string `json:"uname"`
	UnionID string `json:"union_id"`
}

// 房间id
type GameInfo struct {
	GameID string `json:"game_id"`
}

// ws信息
type WebsocketInfo struct {
	AuthBody string   `json:"auth_body"`
	WSSLink  []string `json:"wss_link"`
}
