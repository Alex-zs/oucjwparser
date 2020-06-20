package jwmodel

import (
	"bytes"
	"github.com/Alex-zs/oucjwparser/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fasthttp"
)

type CourseScore struct {
	CourseName string		// 课程名
	Credit string			// 学分
	CourseType string		// 课程类型
	Score string			// 成绩
	semester string			// 学年学期
}

func (session *JwSession) GetScores()  ([]CourseScore, error){
	var scoreList []CourseScore
	var scoreError error = nil
	param := "ysyx=yscj&sjxz=sjxz1&ysyxS=on&sjxzS=on&userCode=" + session.UserCode

	keyParam, err := session.GetKeyParam()
	if err != nil {
		util.Log("获取成绩失败")
		return nil, err
	}
	paramMap := util.EncParamStr(param, keyParam.Key, keyParam.Time)
	session.Do(func(req *fasthttp.Request, resp *fasthttp.Response) {
		req.SetRequestURI(Score)
		for key, value := range paramMap {
			req.URI().QueryArgs().Add(key, value)
		}
		req.Header.SetMethod(fasthttp.MethodGet)
		if err := fasthttp.DoTimeout(req, resp, util.DefaultHttpTimeout); err != nil {
			scoreError = err
			return
		}
		data := util.GBK2UTF8(resp.Body())
		scoreList, scoreError = parseHtmlScore(data)
	})
	return scoreList, scoreError
}

func parseHtmlScore(data []byte) ([]CourseScore, error) {
	courseScores := make([]CourseScore, 0)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	semester := ""
	doc.Find("table").Each(func(i int, selection *goquery.Selection) {
		if i % 2 == 0 {
		 	semester =	selection.Find("td").First().Text()
		 	semester = semester[15:]
		}else {
			selection.Find("tbody").Find("tr").Each(func(i int, selection *goquery.Selection) {
				courseScore := CourseScore{semester: semester}
				selection.Find("td").Each(func(i int, selection *goquery.Selection) {
					switch i {
					case 1:
						courseScore.CourseName = selection.Text()[14:]
					case 2:
						courseScore.Credit = selection.Text()
					case 3:
						courseScore.CourseType = selection.Text()
					case 6:
						courseScore.Score = selection.Text()
					}
				})
				courseScores = append(courseScores, courseScore)
			})
		}
	})
	return courseScores, nil
}
