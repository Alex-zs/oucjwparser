# oucjwparser
一个为方便获取教务系统信息的小工具，需要连接学校VPN或者连接校园网使用。



### 登录教务系统

```go
// 创建新的教务系统会话
session := jwmodel.NewSession()
// 输入教务系统账号密码，返回bool，error
success, err := session.Login("学号", "密码")
if err != nil {
  // 打印错误
  fmt.Println(err.Error())
  return
}
if success {
  fmt.Println("登录成功")
}
```



### 获取个人专业、年级信息

```go
// 创建会话
session := jwmodel.NewSession()
// 登录
if success, _ := session.Login("学号", "密码"); success {
  // 获取信息
  student, err := session.GetInfo()
  if err != nil{
    fmt.Printf(err.Error())
    return
  }
  // 打印个人的专业、年级信息
  fmt.Printf(student.Specialty, student.Grade)
}
```



### 获取指定年级的专业列表

```go
// 获取2019年的所有专业
specialties, err := jwmodel.GetSpecialties(2019)
if err != nil {
  // 打印错误
  fmt.Println(err.Error())
  return
}
// 打印专业
fmt.Println(specialties)
```



### 获取指定学年学期的学生课表

```go
// 创建会话
session := jwmodel.NewSession()
// 登录
if success, _ := session.Login("学号", "密码"); success {
  // 获取课表
  stuCourses, err := session.GetStuCourse(2019, 2)
  if err != nil {
    // 处理错误
    fmt.Println(err.Error())
    return
  }
  // 打印课表
  fmt.Println(stuCourses)
}
```



### 指定学年学期的学生二维html课表

```go
// 创建会话
session := jwmodel.NewSession()
// 登录
if success, _ := session.Login("学号", "密码"); success {
  // 获取二维形式课表
  stuCourses, err := session.GetHtmlStuCourse(2019, 2)
  if err != nil {
    // 处理错误
    fmt.Println(err.Error())
    return
  }
  // 打印二维html课表
  fmt.Println(stuCourses)
}
```



### 获取指定学年学期的选课币

```go
// 创建会话
session := jwmodel.NewSession()
// 登录
if success, _ := session.Login("学号", "密码"); success {
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

```



### 查看是否可以重修某一门课

```go
// 创建会话
session := jwmodel.NewSession()
// 登录
if success, _ := session.Login("学号", "密码"); success {
  // 课程代码是12位数字，一般出现在课程名前面，非选课号
  retake, err := session.CanRetake("课程代码")
  if err != nil {
    fmt.Println(err.Error())
    return
  }
  fmt.Println(retake)
}
```



### 获取教务系统加密参数

部分请求需要加密参数，例如获取成绩、选课

```go
// 创建会话
session := jwmodel.NewSession()
// 登录
if success, _ := session.Login("学号", "密码"); success {
  // 获取加密参数
  keyParam, err := session.GetKeyParam()
  if err != nil {
    fmt.Println(err.Error())
    return
  }
  // 打印加密参数x
  fmt.Println(keyParam)
}
```



### 获取个人成绩

```go
// 创建会话
session := jwmodel.NewSession()
// 登录
if success, _ := session.Login("学号", "密码"); success {
  // 获取成绩
  scoreList, err := session.GetScores()
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(scoreList)
}
```



### 获取各个模块要求学分

```go
// 创建会话
session := jwmodel.NewSession()
// 登录
if success, _ := session.Login("学号", "密码"); success {
  // 获取模块学分
	creditList, err := session.GetCreditRequire()
  if err != nil {
  	fmt.Println(err)
  	return
  }
  fmt.Println(creditList)
}
```





###  验证码识别模块

验证码识别模块使用了第三方接口：

```http
Method: POST 
URL: https://itstudio.club/ocr/jw
body: image=<验证码图片文件>
```

