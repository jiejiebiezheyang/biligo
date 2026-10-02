# BiliGo

[![Go Version](https://img.shields.io/badge/Go-1.26.2-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**中文** | [English](README.md)

Bilibili 开放平台直播 SDK，基于 WebSocket 协议，提供事件总线机制订阅直播间各类事件。

## 快速开始

```go
package main

import (
	"os"
	"os/signal"
	"syscall"

	biligo "github.com/jiejiebiezheyang/biligo/pkg"
)

func main() {
	// 初始化事件总线并订阅弹幕事件
	bus := biligo.InitAndGetBiliBus()
	bus.Subscribe(biligo.LIVE_OPEN_PLATFORM_DM, func(e biligo.Event) {
		biligo.GetBiliLog().Printf("弹幕: %s", e.DataPacket.Data.Msg)
	})

	// 启动客户端（参数来自 B 站开放平台）
	biligo.Start(code, appId, appSecret, appkeyId, bus)

	// 等待退出信号后关闭
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	biligo.Shutdown()
}
```

## 支持的事件

| 常量                                 | 说明         |
| ------------------------------------ | ------------ |
| `LIVE_OPEN_PLATFORM_DM`              | 弹幕         |
| `LIVE_OPEN_PLATFORM_DM_MIRROR`       | 跨房弹幕     |
| `LIVE_OPEN_PLATFORM_SEND_GIFT`       | 礼物         |
| `LIVE_OPEN_PLATFORM_SUPER_CHAT`      | 付费留言     |
| `LIVE_OPEN_PLATFORM_SUPER_CHAT_DEL`  | 付费留言下线 |
| `LIVE_OPEN_PLATFORM_GUARD`           | 大航海       |
| `LIVE_OPEN_PLATFORM_LIKE`            | 点赞         |
| `LIVE_OPEN_PLATFORM_LIVE_ROOM_ENTER` | 进入房间     |
| `LIVE_OPEN_PLATFORM_LIVE_START`      | 开始直播     |
| `LIVE_OPEN_PLATFORM_LIVE_END`        | 结束直播     |
| `LIVE_OPEN_PLATFORM_INTERACTION_END` | 推送结束通知 |

## 事件数据结构

所有事件的回调函数接收 `biligo.Event`，其中 `DataPacket.Data` 为 `LiveRoomData` 结构体（完整字段定义见 [pkg/eventdata.go](pkg/eventdata.go)），主要包括：

```go
type LiveRoomData struct {
	RoomID    int64  // 房间号
	OpenID    string // 用户 OpenID
	Uname     string // 用户名
	Timestamp int64  // 时间戳

	// 弹幕
	Msg      string // 弹幕内容
	DmType   int64  // 弹幕类型

	// 礼物
	GiftName  string // 礼物名称
	GiftNum   int64  // 礼物数量
	Price     int64  // 价格

	// 付费留言
	Message   string // 留言内容
	RMB       int64  // 金额

	// 直播开始/结束
	Title    string // 直播标题
	AreaName string // 分区名称
}
```

## 项目结构

```
biligo/
├── biligo_test.go  # 集成测试/使用示例
└── pkg/
    ├── bilibili.go     # 启动/关闭客户端
    ├── bililog.go      # 日志记录
    ├── config.go       # 配置管理
    ├── event.go        # 事件总线
    ├── eventdata.go    # 事件数据结构
    ├── proto.go        # WebSocket 协议编解码
    ├── request.go      # API 请求签名
    ├── responsedata.go # API 响应结构体
    ├── utils.go        # 工具函数
    └── ws.go           # WebSocket 连接与心跳
```

## 开源协议

本项目基于 [MIT](LICENSE) 协议开源。

Copyright (c) 2026 jiejiebiezheyang
