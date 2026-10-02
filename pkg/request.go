package biligo

// 请求相关

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// 统一请求头结构体
type RequestHeader struct {
	Accept           string
	ContentType      string
	ContentMD5       string
	Timestamp        string
	SignatureMethod  string
	SignatureNonce   string
	AccessKeyID      string
	SignatureVersion string
	Authorization    string
}

// 初始化
func newRequestHeader() *RequestHeader {

	return &RequestHeader{
		Accept:           "application/json",
		ContentType:      "application/json",
		Timestamp:        strconv.FormatInt(time.Now().Unix(), 10),
		SignatureMethod:  "HMAC-SHA256",
		SignatureNonce:   uuid.New().String(),
		AccessKeyID:      cfg.appkeyId,
		SignatureVersion: "1.0",
	}

}

// 签名
func (h *RequestHeader) sign(body string) {
	// 请求体 哈希
	sum := md5.Sum([]byte(body))
	h.ContentMD5 = hex.EncodeToString(sum[:])
	// 拼接字符串
	canonical := strings.Join([]string{
		"x-bili-accesskeyid:" + h.AccessKeyID,
		"x-bili-content-md5:" + h.ContentMD5,
		"x-bili-signature-method:" + h.SignatureMethod,
		"x-bili-signature-nonce:" + h.SignatureNonce,
		"x-bili-signature-version:" + h.SignatureVersion,
		"x-bili-timestamp:" + h.Timestamp,
	}, "\n")
	// 加密签名
	hmacHash := hmac.New(sha256.New, []byte(cfg.appSecret))
	hmacHash.Write([]byte(canonical))
	h.Authorization = hex.EncodeToString(hmacHash.Sum(nil))
}

// 一键构建
func buildBiliHeader(body string) *RequestHeader {
	r := newRequestHeader()
	r.sign(body)
	return r
}

// 设置请求头
func (h *RequestHeader) Apply(req *http.Request) {
	if h.Accept != "" {
		req.Header.Set("Accept", h.Accept)
	}
	if h.ContentType != "" {
		req.Header.Set("Content-Type", h.ContentType)
	}
	if h.ContentMD5 != "" {
		req.Header.Set("x-bili-content-md5", h.ContentMD5)
	}
	if h.Timestamp != "" {
		req.Header.Set("x-bili-timestamp", h.Timestamp)
	}
	if h.SignatureMethod != "" {
		req.Header.Set("x-bili-signature-method", h.SignatureMethod)
	}
	if h.SignatureNonce != "" {
		req.Header.Set("x-bili-signature-nonce", h.SignatureNonce)
	}
	if h.AccessKeyID != "" {
		req.Header.Set("x-bili-accesskeyid", h.AccessKeyID)
	}
	if h.SignatureVersion != "" {
		req.Header.Set("x-bili-signature-version", h.SignatureVersion)
	}
	if h.Authorization != "" {
		req.Header.Set("Authorization", h.Authorization)
	}
}

// 获取错误信息
func getBiliErrorMessage(code int) string {
	// 错误码与错误信息映射表
	errorMap := map[int]string{
		// 4000: "参数错误,请检查必填参数,参数大小限制",
		// 4001: "应用无效,请检查header的x-bili-accesskeyid是否为空,或者有效",
		// 4002: "签名异常,请检查header的Authorization",
		// 4003: "请求过期,请检查header的x-bili-timestamp",
		// 4004: "重复请求,请检查header的x-bili-nonce",
		// 4005: "签名method异常,请检查header的x-bili-signature-method",
		// 4006: "版本异常,请检查header的x-bili-version",
		// 4007: "IP白名单限制,请确认请求服务器是否在报备的白名单内",
		// 4008: "权限异常,请确认接口权限",
		// 4009: "接口访问限制,请确认接口权限及请求频率",
		// 4010: "接口不存在,请确认请求接口url",
		// 4011: "Content-Type不为application/json,请检查header的Content-Type",
		// 4012: "MD5校验失败,请检查header的x-bili-content-md5",
		// 4013: "Accept不为application/json,请检查header的Accept",

		5000: "服务异常,请联系B站对接同学",
		5001: "请求超时",
		5002: "内部错误,请联系B站对接同学",
		5003: "配置错误,请联系B站对接同学",
		5004: "房间白名单限制,未上架应用仅能连接开发者自己的直播间",
		5005: "房间黑名单限制,请联系B站对接同学",
		5011: "应用权限限制,请联系B站对接同学",

		6000: "验证码错误,验证码校验失败",
		6001: "手机号码错误",
		6002: "验证码已过期",
		6003: "验证码频率限制",

		6010: "房间号不能为空",
		6011: "没有查询到房间",
		6012: "主播信息为空",
		6013: "互玩游戏关闭失败",
		6014: "插件关闭失败",
		6015: "直播工具关闭失败",

		7000: "不在游戏内,当前房间未进行互动游戏",
		7001: "请求冷却期,建议10秒后重试",
		7002: "房间重复游戏,当前房间正在进行游戏",
		7003: "心跳过期,game_id错误或互动游戏已关闭",
		7004: "批量心跳超过最大值,单次最大200",
		7005: "批量心跳ID重复",
		7007: "身份码错误",
		7008: "插件重复开启",
		7009: "无道具投放权限",
		7010: "超过上限,同一应用单个直播间最多5个连接",

		// 8002: "项目无权限访问,确认项目ID是否正确",
	}

	// 查表返回
	if msg, ok := errorMap[code]; ok {
		return msg
	}

	// 兜底
	return "未知错误"
}
