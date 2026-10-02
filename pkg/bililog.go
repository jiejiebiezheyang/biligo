package biligo

import (
	"log"
	"os"
)

// 默认丢弃，未初始化不输出任何东西
var biliLog = log.New(os.Stdout, "[biligo] ", log.LstdFlags)

func GetBiliLog() *log.Logger {
	return biliLog
}
