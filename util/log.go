package util

import (
	"log"
	"strings"
)

// 是否打印日志
var PrintLog = true
// 是否第一次记录日志，用来初始化日志
var firstLog = false

// 打印日志
func Log(msg ...string) {
	if PrintLog {
		if firstLog {
			log.SetFlags(log.Ldate|log.Ltime|log.Llongfile)
			firstLog = false
		}
		log.Println(strings.Join(msg, " "))
	}
}
