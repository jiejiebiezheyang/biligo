package biligo

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
)

// 长链消息包结构体
type Packet struct {
	PacketLen uint32 // 整包长度
	HeaderLen uint16 // 头长度, 固定 16
	Version   uint16 // 协议版本
	Operation uint32 // 操作码
	Sequence  uint32 // 序号
	Body      []byte // 包体
}

// 消息类型
const (
	OP_HEARTBEAT       = 2
	OP_HEARTBEAT_REPLY = 3
	OP_MESSAGE         = 5
	OP_AUTH            = 7
	OP_AUTH_REPLY      = 8
)

// 包装数据
func packPacket(op uint32, body []byte) []byte {
	headerLen := uint16(16)
	packetLen := uint32(headerLen) + uint32(len(body))

	buf := make([]byte, packetLen)

	binary.BigEndian.PutUint32(buf[0:4], packetLen)
	binary.BigEndian.PutUint16(buf[4:6], headerLen)
	binary.BigEndian.PutUint16(buf[6:8], 1) // version = 1
	binary.BigEndian.PutUint32(buf[8:12], op)
	binary.BigEndian.PutUint32(buf[12:16], 1)

	copy(buf[16:], body)
	return buf
}

// 包装认证数据
func packAuth(authBody string) []byte {
	return packPacket(OP_AUTH, []byte(authBody))
}

// 包装心跳数据
func packHeartbeat() []byte {
	return packPacket(OP_HEARTBEAT, nil)
}

// 解包
func unpackPackets(data []byte) ([]Packet, error) {
	var packets []Packet
	offset := 0

	for {
		if len(data[offset:]) < 16 {
			break
		}

		packetLen := int(binary.BigEndian.Uint32(data[offset : offset+4]))
		headerLen := int(binary.BigEndian.Uint16(data[offset+4 : offset+6]))
		version := binary.BigEndian.Uint16(data[offset+6 : offset+8])
		operation := binary.BigEndian.Uint32(data[offset+8 : offset+12])
		sequence := binary.BigEndian.Uint32(data[offset+12 : offset+16])

		if offset+packetLen > len(data) {
			break
		}

		body := data[offset+headerLen : offset+packetLen]

		packets = append(packets, Packet{
			PacketLen: uint32(packetLen),
			HeaderLen: uint16(headerLen),
			Version:   version,
			Operation: operation,
			Sequence:  sequence,
			Body:      body,
		})

		offset += packetLen
	}

	return packets, nil
}

// 处理包数据
func handlePacket(p Packet) {
	switch p.Operation {

	case OP_AUTH_REPLY:
		biliLog.Println("ws 认证回复:", string(p.Body))

	case OP_HEARTBEAT_REPLY:
		biliLog.Println("ws 心跳应答")

	case OP_MESSAGE:
		handleMessage(p)

	default:
		biliLog.Println("未知操作:", p.Operation)
	}
}

func inflateZlib(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

// 处理 OP_MESSAGE
func handleMessage(p Packet) {
	switch p.Version {

	case 0, 1:
		// 如果Version=0，Body中就是实际发送的数据
		biliLog.Println("收到消息:", string(p.Body))
		dispatchBusinessMessage(p.Body)

	case 2:
		// 如果Version=2，Body中是经过压缩后的数据，请使用zlib解压
		// Version=2时，zlib压缩后的body格式可能包含多个完整的proto包
		raw, err := inflateZlib(p.Body)
		if err != nil {
			biliLog.Println("zlib 错误:", err)
			return
		}

		subPackets, _ := unpackPackets(raw)
		// 递归处理
		for _, sp := range subPackets {
			handlePacket(sp)
		}

	default:
		biliLog.Println("未知版本:", p.Version)
	}
}
