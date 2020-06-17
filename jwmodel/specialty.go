package jwmodel

import (
	"github.com/Alex-zs/oucjwparser/util"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

type Specialty struct {
	name string			// 专业名称
	code string			// 专业代码
	number string		// 专业号
}

// 获取指定年份的专业列表
func GetSpecialties(year int) ([]Specialty, error){
	var specialties [] Specialty = nil
	var specialtyError error = nil
	util.SimpleDo(func(req *fasthttp.Request, resp *fasthttp.Response) {
		// 添加请求参数
		args := fasthttp.AcquireArgs()
		defer fasthttp.ReleaseArgs(args)
		args.Add("comboBoxName", "MsGrade_Specialty")
		args.Add("paramValue", "nj=" + strconv.Itoa(year))
		args.WriteTo(req.BodyWriter())
		req.Header.SetMethod(fasthttp.MethodPost)
		req.SetRequestURI(SpecialtyList)
		if err := fasthttp.DoTimeout(req, resp, util.DefaultHttpTimeout); err != nil {
			specialtyError = err
			return
		}
		// 获取专业列表大小
		size := int(gjson.GetBytes(resp.Body(), "#").Num)
		specialties = make([]Specialty, size)
		result := gjson.ParseBytes(resp.Body())
		specialtyIndex := 0
		result.ForEach(func(key, value gjson.Result) bool {
			code := gjson.Get(value.Raw, "code").Str
			name := gjson.Get(value.Raw, "name").Str
			startIndex := strings.Index(name, "[")
			endIndex := strings.Index(name, "]")
			number := name[startIndex + 1: endIndex]
			name = name[endIndex + 1:]
			specialties[specialtyIndex] = Specialty{name, code, number}
			specialtyIndex++
			return true
		})
	})
	return specialties, specialtyError
}