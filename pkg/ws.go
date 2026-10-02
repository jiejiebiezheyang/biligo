package biligo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// 连接 WebSocket 并完成鉴权
func (a *Config) connectWS() {
	if a.connData == nil {
		biliLog.Fatal("认证信息为空")
	}

	// 获取 ws 信息
	wsInfo := a.connData.WebsocketInfo
	if len(wsInfo.WSSLink) == 0 {
		biliLog.Fatal("未获取到 ws 链接")
	}

	for index, link := range wsInfo.WSSLink {
		biliLog.Printf("ws 链接 %d: %s\n", index+1, link)
	}

	dialer := websocket.Dialer{
		HandshakeTimeout:  10 * time.Second,
		EnableCompression: true,
	}

	var conn *websocket.Conn
	var err error

	// 1. 多 wss_link 兜底连接
	for _, link := range wsInfo.WSSLink {
		conn, _, err = dialer.Dial(link, nil)
		if err == nil {
			biliLog.Printf("ws 链接 %s 连接成功\n", link)
			break
		}
		biliLog.Printf("当前 ws 链接 %s 连接失败 err: %v\n", link, err)
	}

	if conn == nil {
		biliLog.Fatal("所有 ws链接 连接失败")
	}

	a.isReady = true
	a.wsConn = conn

	// 2. 发送 OP_AUTH(Proto 包)
	authPacket := packAuth(wsInfo.AuthBody)
	err = conn.WriteMessage(websocket.BinaryMessage, authPacket)
	if err != nil {
		biliLog.Fatal("ws 认证失败")
	}

	biliLog.Println("ws 认证发送")

	// 3. 启动读循环(核心)
	go func() {
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				biliLog.Println("ws 读错误", err)
				return
			}

			// 只处理二进制包
			if msgType != websocket.BinaryMessage {
				continue
			}

			packets, err := unpackPackets(msg)
			if err != nil {
				biliLog.Println("解包错误", err)
				continue
			}

			for _, p := range packets {
				handlePacket(p)
			}
		}
	}()

	// 4. ws 心跳 (必须)
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if a.wsConn == nil {
				return
			}

			err := a.wsConn.WriteMessage(
				websocket.BinaryMessage,
				packHeartbeat(),
			)
			if err != nil {
				biliLog.Println("ws 心跳错误", err)
				return
			}
			biliLog.Println("ws 心跳请求")
		}
	}()
	// 5. 项目心跳 (必须)
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			data := map[string]interface{}{
				"game_id": a.connData.GameInfo.GameID,
			}
			body, _ := json.Marshal(data)

			client := &http.Client{}

			req, _ := http.NewRequest("POST", biliHost+"/v2/app/heartbeat", bytes.NewReader(body))

			h := buildBiliHeader(string(body))
			h.Apply(req)

			biliLog.Println("http 心跳请求")

			resp, err := client.Do(req)

			if err != nil {
				biliLog.Println("http 心跳错误", err)
				return
			}

			var res ApiResponse
			resBody, _ := io.ReadAll(resp.Body)
			json.Unmarshal(resBody, &res)
			if res.Code == 0 {
				biliLog.Println("http 心跳应答", string(resBody))
			}
		}
	}()

	biliLog.Println("ws 完全准备好")
}
