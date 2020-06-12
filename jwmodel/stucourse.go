package jwmodel

import (
	"fmt"
	"github.com/Alex-zs/oucjwparser/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

type StuCourse struct {
	classCode string		// 课程代码
	classNum string			// 选课号
	className string		// 课程名
	credit string			// 学分
	classType string		// 课程类型
	campus string			// 校区
	teacher string			// 教师
	coins string			// 投币
	buyBook	bool			// 买书
	retake	bool			// 重修
	adjust bool				// 调剂
	status string			// 选课状态
	remark string			// 备注
}

// 根据学年、学期获得该学年、学期的选课记录
func (session *JwSession) GetStuCourse(year, semester int) *[]StuCourse {
	var stuCourses []StuCourse = nil
	session.Do(func(req *fasthttp.Request, resp *fasthttp.Response) {
		// 添加请求参数
		req.SetRequestURI(DataTable)
		req.Header.SetMethod(fasthttp.MethodPost)
		args := &fasthttp.Args{}
		args.Add("xh", session.UserCode)
		args.Add("xn", strconv.Itoa(year))
		args.Add("xq", strconv.Itoa(semester))
		args.Add("tableId", "6093")
		args.WriteTo(req.BodyWriter())

		fasthttp.DoTimeout(req, resp, util.DefaultHttpTimeout)

		data := util.GBKBytes2UTF8(resp.Body())
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(data))
		slector := doc.Find("tbody")
		fmt.Println(slector)
	})
	return &stuCourses
}
