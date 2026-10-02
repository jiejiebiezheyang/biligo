package biligo

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// 下载图片转 base64 防止前端跨域问题
func ImageToBase64(url string) string {

	// 下载图片
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取字节
	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 转 Base64
	base64Str := base64.StdEncoding.EncodeToString(imgBytes)

	// 获取 Content-Type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		// 如果没返回类型，可以根据链接后缀简单判断
		if strings.HasSuffix(url, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(url, ".jpg") || strings.HasSuffix(url, ".jpeg") {
			contentType = "image/jpeg"
		} else {
			contentType = "application/octet-stream"
		}
	}

	// 拼成可直接用的 <img> src
	return fmt.Sprintf("data:%s;base64,%s", contentType, base64Str)
}
