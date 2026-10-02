# BiliGo

[![Go Version](https://img.shields.io/badge/Go-1.26.2-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

[中文](README.zh-CN.md) | **English**

A Bilibili Open Platform Live SDK built on the WebSocket protocol, providing an event bus mechanism to subscribe to various live room events.

## Quick Start

```go
package main

import (
	"os"
	"os/signal"
	"syscall"

	biligo "github.com/jiejiebiezheyang/biligo/pkg"
)

func main() {
	// Initialize the event bus and subscribe to danmaku events
	bus := biligo.InitAndGetBiliBus()
	bus.Subscribe(biligo.LIVE_OPEN_PLATFORM_DM, func(e biligo.Event) {
		biligo.GetBiliLog().Printf("danmaku: %s", e.DataPacket.Data.Msg)
	})

	// Start the client (credentials come from the Bilibili Open Platform)
	biligo.Start(code, appId, appSecret, appkeyId)

	// Wait for an exit signal, then shut down
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	biligo.Shutdown()
}
```

## Supported Events

| Constant                             | Description           |
| ------------------------------------ | --------------------- |
| `LIVE_OPEN_PLATFORM_DM`              | Danmaku               |
| `LIVE_OPEN_PLATFORM_DM_MIRROR`       | Cross-room danmaku    |
| `LIVE_OPEN_PLATFORM_SEND_GIFT`       | Gift                  |
| `LIVE_OPEN_PLATFORM_SUPER_CHAT`      | Paid message          |
| `LIVE_OPEN_PLATFORM_SUPER_CHAT_DEL`  | Paid message removed  |
| `LIVE_OPEN_PLATFORM_GUARD`           | Guard purchase        |
| `LIVE_OPEN_PLATFORM_LIKE`            | Like                  |
| `LIVE_OPEN_PLATFORM_LIVE_ROOM_ENTER` | Room enter            |
| `LIVE_OPEN_PLATFORM_LIVE_START`      | Live start            |
| `LIVE_OPEN_PLATFORM_LIVE_END`        | Live end              |
| `LIVE_OPEN_PLATFORM_INTERACTION_END` | Push end notification |

## Event Data Structure

Every event callback receives a `biligo.Event`, whose `DataPacket.Data` is a `LiveRoomData` struct (see [pkg/eventdata.go](pkg/eventdata.go) for the full definition), containing:

```go
type LiveRoomData struct {
	RoomID    int64  // Room ID
	OpenID    string // User OpenID
	Uname     string // Username
	Timestamp int64  // Timestamp

	// Danmaku
	Msg      string // Danmaku content
	DmType   int64  // Danmaku type

	// Gift
	GiftName  string // Gift name
	GiftNum   int64  // Gift count
	Price     int64  // Price

	// Paid message
	Message   string // Message content
	RMB       int64  // Amount

	// Live start/end
	Title    string // Live title
	AreaName string // Area name
}
```

## Project Structure

```
biligo/
├── biligo_test.go  # Integration test / usage example
└── pkg/
    ├── bilibili.go     # Start/stop client
    ├── bililog.go      # Logger
    ├── config.go       # Configuration management
    ├── event.go        # Event bus
    ├── eventdata.go    # Event data structures
    ├── proto.go        # WebSocket protocol encoding/decoding
    ├── request.go      # API request signing
    ├── responsedata.go # API response structs
    ├── utils.go        # Utility functions
    └── ws.go           # WebSocket connection and heartbeat
```

## License

Licensed under the [MIT License](LICENSE).

Copyright (c) 2026 jiejiebiezheyang
