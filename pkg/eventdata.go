package biligo

// 事件数据

type BaseMessage struct {
	Cmd string `json:"cmd"`
}

type DanmuMessage struct {
	Cmd  string       `json:"cmd"`
	Data LiveRoomData `json:"data"`
}

type LiveRoomData struct {
	// --- 通用字段 (出现在多个事件中) ---
	RoomID    int64  `json:"room_id,omitempty"`
	OpenID    string `json:"open_id,omitempty"`
	UnionID   string `json:"union_id,omitempty"`
	UID       int64  `json:"uid,omitempty"` // 已废弃，固定为0
	Uname     string `json:"uname,omitempty"`
	UFace     string `json:"uface,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	MsgID     string `json:"msg_id,omitempty"`

	// --- 弹幕 (LIVE_OPEN_PLATFORM_DM) ---
	Msg                    string `json:"msg,omitempty"`
	FansMedalLevel         int64  `json:"fans_medal_level,omitempty"`
	FansMedalName          string `json:"fans_medal_name,omitempty"`
	FansMedalWearingStatus bool   `json:"fans_medal_wearing_status,omitempty"`
	GuardLevel             int64  `json:"guard_level,omitempty"`
	EmojiImgURL            string `json:"emoji_img_url,omitempty"`
	DmType                 int64  `json:"dm_type,omitempty"`
	GloryLevel             int    `json:"glory_level,omitempty"`
	ReplyOpenID            string `json:"reply_open_id,omitempty"`
	ReplyUname             string `json:"reply_uname,omitempty"`
	IsAdmin                int    `json:"is_admin,omitempty"`

	// --- 礼物 (LIVE_OPEN_PLATFORM_SEND_GIFT) ---
	GiftID     int64       `json:"gift_id,omitempty"`
	GiftName   string      `json:"gift_name,omitempty"`
	GiftNum    int64       `json:"gift_num,omitempty"`
	Price      int64       `json:"price,omitempty"`
	RPrice     int64       `json:"r_price,omitempty"`
	Paid       bool        `json:"paid,omitempty"`
	GiftIcon   string      `json:"gift_icon,omitempty"`
	ComboGift  bool        `json:"combo_gift,omitempty"`
	AnchorInfo *AnchorInfo `json:"anchor_info,omitempty"`
	ComboInfo  *ComboInfo  `json:"combo_info,omitempty"`
	BlindGift  *BlindGift  `json:"blind_gift,omitempty"`

	// --- 付费留言 (LIVE_OPEN_PLATFORM_SUPER_CHAT) ---
	MessageID int64  `json:"message_id,omitempty"`
	Message   string `json:"message,omitempty"`
	RMB       int64  `json:"rmb,omitempty"`
	StartTime int64  `json:"start_time,omitempty"`
	EndTime   int64  `json:"end_time,omitempty"`

	// --- 付费留言下线 (LIVE_OPEN_PLATFORM_SUPER_CHAT_DEL) ---
	MessageIDs []int64 `json:"message_ids,omitempty"`

	// --- 大航海 (LIVE_OPEN_PLATFORM_GUARD) ---
	UserInfo  *UserInfo `json:"user_info,omitempty"`
	GuardNum  int64     `json:"guard_num,omitempty"`
	GuardUnit string    `json:"guard_unit,omitempty"`

	// --- 点赞 (LIVE_OPEN_PLATFORM_LIKE) ---
	LikeText  string `json:"like_text,omitempty"`
	LikeCount int64  `json:"like_count,omitempty"`

	// --- 开始/结束直播 (LIVE_OPEN_PLATFORM_LIVE_START / END) ---
	AreaName string `json:"area_name,omitempty"`
	Title    string `json:"title,omitempty"`

	// --- 推送结束 (LIVE_OPEN_PLATFORM_INTERACTION_END) ---
	GameID string `json:"game_id,omitempty"`
}

// ComboInfo 连击信息 (用于礼物事件)
type ComboInfo struct {
	ComboBaseNum int64  `json:"combo_base_num,omitempty"`
	ComboCount   int64  `json:"combo_count,omitempty"`
	ComboID      string `json:"combo_id,omitempty"`
	ComboTimeout int64  `json:"combo_timeout,omitempty"`
}

// BlindGift 盲盒信息 (用于礼物事件)
type BlindGift struct {
	BlindGiftID int64 `json:"blind_gift_id,omitempty"`
	Status      bool  `json:"status,omitempty"`
}

// UserInfo 用户信息 (用于大航海事件)
type UserInfo struct {
	UID     int64  `json:"uid,omitempty"`
	OpenID  string `json:"open_id,omitempty"`
	UnionID string `json:"union_id,omitempty"`
	Uname   string `json:"uname,omitempty"`
	UFace   string `json:"uface,omitempty"`
}
