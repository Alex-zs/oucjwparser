package jwmodel

import (
	"bytes"
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
func (session *JwSession) GetStuCourse(year, semester int) ([]StuCourse, error) {
	stuCourses := make([]StuCourse, 0)
	var stuCourseError error = nil
	session.Do(func(req *fasthttp.Request, resp *fasthttp.Response) {
		// 添加请求参数
		req.SetRequestURI(DataTable)
		req.Header.SetMethod(fasthttp.MethodPost)
		args := fasthttp.AcquireArgs()
		defer fasthttp.ReleaseArgs(args)
		args.Add("xh", session.UserCode)
		args.Add("xn", strconv.Itoa(year))
		args.Add("xq", strconv.Itoa(semester))
		args.Add("tableId", "6093")
		args.WriteTo(req.BodyWriter())

		if err := fasthttp.DoTimeout(req, resp, util.DefaultHttpTimeout); err != nil {
			util.Log("获取选课记录失败")
			stuCourseError = err
		}
		data := util.GBK2UTF8(resp.Body())
		doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(data))
		doc.Find(".O,.E").Each(func(i int, selection *goquery.Selection) {
			course := StuCourse{}
			selection.Find("td").Each(func(i int, selection *goquery.Selection) {
				text := strings.Trim(selection.Text(), " ")
				switch i {
				case 1:
					startIndex := strings.Index(text, "[")
					endIndex := strings.Index(text, "]")
					course.classCode = text[startIndex + 1: endIndex]
					course.className = text[endIndex + 1:]
				case 2:
					course.credit = text
				case 4:
					course.classType = text
				case 5:
					course.campus = text
				case 6:
					course.classNum = text
				case 7:
					course.teacher = text
				case 8:
					course.coins = text
				case 9:
					course.buyBook = selection.Find("input").First().Is(":checked")
				case 10:
					course.retake = selection.Find("input").First().Is(":checked")
				case 11:
					course.adjust = selection.Find("input").First().Is(":checked")
				case 13:
					course.status = text
				case 14:
					course.remark = text
				}
			})
			stuCourses = append(stuCourses, course)
		})

	})
	return stuCourses, stuCourseError
}


// 获取二维html格式课程表
func (session *JwSession) GetHtmlStuCourse(year, semester int) (string, error){
	var htmlStuCourseError error = nil
	htmlStuCourse := ""
	session.Do(func(req *fasthttp.Request, resp *fasthttp.Response) {
		// 设置请求参数
		req.SetRequestURI(HtmlTimeTable)
		req.Header.SetMethod(fasthttp.MethodPost)
		args := fasthttp.AcquireArgs()
		defer fasthttp.ReleaseArgs(args)
		args.Add("xh", session.UserCode)
		args.Add("xn", strconv.Itoa(year))
		args.Add("xq", strconv.Itoa(semester))
		args.WriteTo(req.BodyWriter())

		if err := fasthttp.DoTimeout(req, resp, util.DefaultHttpTimeout); err != nil {
			util.Log("请求二维html格式课程表失败")
			htmlStuCourseError = err
			return
		}
		htmlStuCourse = string(util.GBK2UTF8(resp.Body()))
	})

	return htmlStuCourse, htmlStuCourseError
	
}

