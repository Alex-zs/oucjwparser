package jwmodel

import (
	"github.com/Alex-zs/oucjwparser/util"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
	"strings"
)

type Student struct {
	UserCode      string // 学号
	Grade         string // 年级
	Specialty     string // 专业名称
	SpecialtyCode string // 专业代码
}

// 获取个人年级、专业信息
func (session *JwSession) GetInfo() (*Student, error)  {
	var infoError error = nil
	resultStr := ""
	session.Do(func(req *fasthttp.Request, resp *fasthttp.Response) {
		req.SetRequestURI(GradeSpeciality)
		req.Header.SetMethod(fasthttp.MethodPost)
		if err := fasthttp.DoTimeout(req, resp, util.DefaultHttpTimeout); err != nil {
			util.Log("请求获取个人年级、专业信息失败")
			infoError = err
			return
		}
		resultStr = gjson.GetBytes(resp.Body(), "result").Str
	})
	if infoError != nil{
		return nil, infoError
	}
	resultStr = strings.ReplaceAll(resultStr, `\`, "")
	resultList := gjson.GetMany(resultStr, "nj", "zymc", "zydm")
	return &Student{
		session.UserCode, resultList[0].Str,
		resultList[1].Str, resultList[2].Str}, nil
}
