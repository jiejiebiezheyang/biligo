package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"

	biligo "github.com/jiejiebiezheyang/biligo/pkg"
)

func TestBilibili(t *testing.T) {
	// 测试代码
	code := os.Getenv("BILIBILI_CODE")
	appId := os.Getenv("BILIBILI_APP_ID")
	AppSecret := os.Getenv("BILIBILI_APP_SECRET")
	AppkeyId := os.Getenv("BILIBILI_APP_KEY_ID")

	biliLog := biligo.GetBiliLog()

	if code == "" || appId == "" || AppSecret == "" || AppkeyId == "" {
		t.Skip("环境变量未配置")
	}

	appIdInt, err := strconv.Atoi(appId)
	if err != nil {
		biliLog.Fatalf("BILIBILI_APP_ID 解析失败: %v", err)
	}

	// 初始化事件总线
	eventBus := biligo.InitAndGetBiliBus()

	// const (
	// 	LIVE_OPEN_PLATFORM_DM              = "LIVE_OPEN_PLATFORM_DM"              // 弹幕信息
	// 	LIVE_OPEN_PLATFORM_DM_MIRROR       = "LIVE_OPEN_PLATFORM_DM_MIRROR"       // 跨房弹幕信息
	// 	LIVE_OPEN_PLATFORM_SEND_GIFT       = "LIVE_OPEN_PLATFORM_SEND_GIFT"       // 礼物信息
	// 	LIVE_OPEN_PLATFORM_SUPER_CHAT      = "LIVE_OPEN_PLATFORM_SUPER_CHAT"      // 付费留言
	// 	LIVE_OPEN_PLATFORM_SUPER_CHAT_DEL  = "LIVE_OPEN_PLATFORM_SUPER_CHAT_DEL"  // 付费留言下线
	// 	LIVE_OPEN_PLATFORM_GUARD           = "LIVE_OPEN_PLATFORM_GUARD"           // 付费大航海
	// 	LIVE_OPEN_PLATFORM_LIKE            = "LIVE_OPEN_PLATFORM_LIKE"            // 点赞信息
	// 	LIVE_OPEN_PLATFORM_LIVE_ROOM_ENTER = "LIVE_OPEN_PLATFORM_LIVE_ROOM_ENTER" // 进入房间
	// 	LIVE_OPEN_PLATFORM_LIVE_START      = "LIVE_OPEN_PLATFORM_LIVE_START"      // 开始直播
	// 	LIVE_OPEN_PLATFORM_LIVE_END        = "LIVE_OPEN_PLATFORM_LIVE_END"        // 结束直播
	// 	LIVE_OPEN_PLATFORM_INTERACTION_END = "LIVE_OPEN_PLATFORM_INTERACTION_END" // 消息推送结束通知
	// )
	// 分别注册所有事件处理函数
	eventBus.Subscribe(biligo.LIVE_OPEN_PLATFORM_DM, func(e biligo.Event) {
		biliLog.Printf("弹幕信息:%s", e.DataPacket.Data.Msg)
	})
	eventBus.Subscribe(biligo.LIVE_OPEN_PLATFORM_DM_MIRROR, func(e biligo.Event) {
		biliLog.Printf("跨房弹幕信息:%s", e.DataPacket.Data.Msg)
	})
	eventBus.Subscribe(biligo.LIVE_OPEN_PLATFORM_SEND_GIFT, func(e biligo.Event) {
		biliLog.Printf("礼物信息:%s", e.DataPacket.Data.Msg)
	})
	eventBus.Subscribe(biligo.LIVE_OPEN_PLATFORM_SUPER_CHAT, func(e biligo.Event) {
		biliLog.Printf("付费留言信息:%s", e.DataPacket.Data.Msg)
	})
	eventBus.Subscribe(biligo.LIVE_OPEN_PLATFORM_SUPER_CHAT_DEL, func(e biligo.Event) {
		biliLog.Printf("付费留言下线:%s", e.DataPacket.Data.Msg)
	})
	eventBus.Subscribe(biligo.LIVE_OPEN_PLATFORM_GUARD, func(e biligo.Event) {
		biliLog.Printf("付费大航海:%s", e.DataPacket.Data.Msg)
	})
	eventBus.Subscribe(biligo.LIVE_OPEN_PLATFORM_LIKE, func(e biligo.Event) {
		biliLog.Printf("点赞信息:%s", e.DataPacket.Data.Msg)
	})
	eventBus.Subscribe(biligo.LIVE_OPEN_PLATFORM_LIVE_ROOM_ENTER, func(e biligo.Event) {
		biliLog.Printf("进入房间:%s", e.DataPacket.Data.Msg)
	})
	eventBus.Subscribe(biligo.LIVE_OPEN_PLATFORM_LIVE_START, func(e biligo.Event) {
		biliLog.Printf("开始直播:%s", e.DataPacket.Data.Msg)
	})
	eventBus.Subscribe(biligo.LIVE_OPEN_PLATFORM_LIVE_END, func(e biligo.Event) {
		biliLog.Printf("结束直播:%s", e.DataPacket.Data.Msg)
	})
	eventBus.Subscribe(biligo.LIVE_OPEN_PLATFORM_INTERACTION_END, func(e biligo.Event) {
		biliLog.Printf("推送结束:%s", e.DataPacket.Data.Msg)
	})

	biligo.Start(code, appIdInt, AppSecret, AppkeyId)

	waitShutdown()

	biligo.Shutdown()
}

func waitShutdown() {
	biliLog := biligo.GetBiliLog()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	biliLog.Println("收到停止信号, 程序准备退出")
}
