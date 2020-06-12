package jwmodel

import (
	"github.com/Alex-zs/oucjwparser/util"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
	"strings"
)

type Student struct {
	userCode string			// 学号
	grade string			// 年级
	specialty string		// 专业名称
	specialtyCode string	// 专业代码
}

// 获取个人年级、专业信息
func (session *JwSession) GetInfo() *Student  {
	resultStr := ""
	session.Do(func(req *fasthttp.Request, resp *fasthttp.Response) {
		req.SetRequestURI(GradeSpeciality)
		req.Header.SetMethod(fasthttp.MethodPost)
		fasthttp.DoTimeout(req, resp, util.DefaultHttpTimeout)
		resultStr = gjson.GetBytes(resp.Body(), "result").Str
	})
	resultStr = strings.ReplaceAll(resultStr, `\`, "")
	resultList := gjson.GetMany(resultStr, "nj", "zymc", "zydm")
	if len(resultList) == 3 {
		return &Student{
			session.UserCode, resultList[0].Str,
			resultList[1].Str, resultList[2].Str}
	} else{
		return &Student{}
	}
}
