# oucjwparser
A tool for ouc 教务系统，go语言编写。

该工具需要连接学校VPN或者连接校园网。



## 登录教务系统

```go
// 创建新的教务系统会话
session := jwmodel.JwSession{}
// 通过bool类型返回值判断是否登录成功
success := session.Login("学号", "密码")
```



## 验证码识别模块

验证码识别模块使用了第三方接口：

```http
Method: POST 
URL: https://itstudio.club/ocr/jw
body: image=<验证码图片文件>
```

