package jwmodel

import (
	"github.com/Alex-zs/oucjwparser/util"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
)

// 判断学生是否可以重修某门课程
// @param classCode 课程代码
func (session *JwSession) CanRetake(classCode string) (bool, error) {
	var retakeError error = nil
	retake := false
	session.Do(func(req *fasthttp.Request, resp *fasthttp.Response) {
		req.SetRequestURI(RepeatCourse)
		req.Header.SetMethod(fasthttp.MethodGet)
		args := fasthttp.AcquireArgs()
		defer fasthttp.ReleaseArgs(args)
		args.Add("xh", session.UserCode)
		args.Add("kcdm", classCode)
		args.WriteTo(req.BodyWriter())

		if err := fasthttp.DoTimeout(req, resp, util.DefaultHttpTimeout); err != nil {
			retakeError = err
			return
		}

		resultList := gjson.GetManyBytes(resp.Body(), "result", "status")

		if len(resultList) != 2 {
			retakeError = &util.JwError{Msg: "获取重修参数错误"}
		}
		retake = (resultList[0].Str == "1") && (resultList[1].Str == "200")
	})

	return retake, retakeError
}

