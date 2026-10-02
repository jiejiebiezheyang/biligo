package biligo

import (
	"encoding/json"
	"sync"
)

// 事件类型常量
const (
	LIVE_OPEN_PLATFORM_DM              = "LIVE_OPEN_PLATFORM_DM"              // 弹幕信息
	LIVE_OPEN_PLATFORM_DM_MIRROR       = "LIVE_OPEN_PLATFORM_DM_MIRROR"       // 跨房弹幕信息
	LIVE_OPEN_PLATFORM_SEND_GIFT       = "LIVE_OPEN_PLATFORM_SEND_GIFT"       // 礼物信息
	LIVE_OPEN_PLATFORM_SUPER_CHAT      = "LIVE_OPEN_PLATFORM_SUPER_CHAT"      // 付费留言
	LIVE_OPEN_PLATFORM_SUPER_CHAT_DEL  = "LIVE_OPEN_PLATFORM_SUPER_CHAT_DEL"  // 付费留言下线
	LIVE_OPEN_PLATFORM_GUARD           = "LIVE_OPEN_PLATFORM_GUARD"           // 付费大航海
	LIVE_OPEN_PLATFORM_LIKE            = "LIVE_OPEN_PLATFORM_LIKE"            // 点赞信息
	LIVE_OPEN_PLATFORM_LIVE_ROOM_ENTER = "LIVE_OPEN_PLATFORM_LIVE_ROOM_ENTER" // 进入房间
	LIVE_OPEN_PLATFORM_LIVE_START      = "LIVE_OPEN_PLATFORM_LIVE_START"      // 开始直播
	LIVE_OPEN_PLATFORM_LIVE_END        = "LIVE_OPEN_PLATFORM_LIVE_END"        // 结束直播
	LIVE_OPEN_PLATFORM_INTERACTION_END = "LIVE_OPEN_PLATFORM_INTERACTION_END" // 消息推送结束通知
)

// 事件结构体
type Event struct {
	// eventType 事件类型
	EventType string
	// Data 是事件携带的数据
	DataPacket DanmuMessage
}

// Handler 是事件处理函数类型
//
// 当事件触发时, EventBus 会调用对应的 Handler
type Handler func(Event)

// Bus 是事件总线
//
// 1. 保存所有事件监听者
// 2. 接收事件发布
// 3. 将事件分发给对应监听者
//
// 结构:
//
//	         EventBus
//	            |
//	    +-------+-------+
//	    |       |       |
//	Handler Handler Handler
//
// WS收到消息后发送到 Bus
// Bus 再通知其他业务模块
type Bus struct {

	// mu 用来保护 handlers.
	mu sync.RWMutex

	// handlers 保存所有事件监听函数.
	//
	// key:
	//     事件名称
	//
	// value:
	//     监听这个事件的所有处理函数
	handlers map[string][]Handler
}

var (
	biliBus *Bus      // 事件总线
	busOnce sync.Once // 事件总线初始化
)

// InitAndGetBiliBus 初始化并获取事件总线 进程内单例
func InitAndGetBiliBus() *Bus {
	busOnce.Do(func() {
		biliBus = &Bus{
			// 创建可以使用的 map
			handlers: make(map[string][]Handler),
		}
	})
	return biliBus
}

// Subscribe 注册事件监听.
//
// 参数:
// eventType:
//
//	要监听的事件
//
// handler:
//
//	收到事件后执行的函数
func (b *Bus) Subscribe(eventType string, handler Handler) {

	// 加写锁
	b.mu.Lock()

	// 函数结束时自动释放锁
	defer b.mu.Unlock()

	// 添加监听函数
	b.handlers[eventType] = append(b.handlers[eventType], handler)

}

// Publish 发布事件
func (b *Bus) publish(e Event) {

	// 加读锁
	b.mu.RLock()

	// 获取对应事件的所有监听函数
	handlers := b.handlers[e.EventType]

	// 释放读锁
	b.mu.RUnlock()

	// 依次通知所有监听者
	for _, handler := range handlers {
		// 使用异步执行
		go handler(e)
	}
}

// 处理消息
func dispatchBusinessMessage(body []byte) {
	var base BaseMessage
	if err := json.Unmarshal(body, &base); err != nil {
		return
	}
	var m DanmuMessage
	json.Unmarshal(body, &m)
	// 发布事件
	biliBus.publish(Event{
		EventType:  base.Cmd,
		DataPacket: m,
	})
}
