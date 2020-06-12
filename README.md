# oucjwparser
A tool for ouc 教务系统，go语言编写。

该工具需要连接学校VPN或者连接校园网。



### 登录教务系统

```go
// 创建新的教务系统会话
session := jwmodel.JwSession{}
// 通过bool类型返回值判断是否登录成功
success := session.Login("学号", "密码")
```



### 获取个人专业、年级信息

```go
// 创建会话
session := jwmodel.JwSession{}
// 登录
if session.Login("学号", "密码") {
  	// 获取信息
    studentInfo := session.GetInfo()
    fmt.Println(studentInfo)
}
```



### 获取指定年级的专业列表

```go
// 打印2019年的所有专业
fmt.Println(jwmodel.GetSpecialties(2019))		
```





###  验证码识别模块

验证码识别模块使用了第三方接口：

```http
Method: POST 
URL: https://itstudio.club/ocr/jw
body: image=<验证码图片文件>
```

