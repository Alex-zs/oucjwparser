package main

import (
	"fmt"
	"github.com/Alex-zs/oucjwparser/jwmodel"
)

func main() {
	// 创建会话
	session := jwmodel.NewSession()
	// 登录
	if success, _ := session.Login("17020031002", "chen1234"); success {
		// 获取选课币状态
		coinStatus, err := session.GetCoinStatus(2019, 2)
		// 处理错误
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// 打印选课币
		fmt.Println(coinStatus)
	}
}
